require 'shellwords'
require 'forwardable'
require 'delegate'

module Hub
  # Methods for inspecting the environment, such as reading git config,
  # repository info, and other.
  module Context
    extend Forwardable

    NULL = defined?(File::NULL) ? File::NULL : File.exist?('/dev/null') ? '/dev/null' : 'NUL'

    # Shells out to git to get output of its commands
    class GitReader
      attr_reader :executable

      def initialize(executable = nil, &read_proc)
        @executable = executable || 'git'
        # caches output when shelling out to git
        read_proc ||= lambda { |cache, cmd|
          str = command_to_string(cmd)
          result = silence_stderr { %x{#{str}}.chomp }
          cache[cmd] = $?.success? && !result.empty? ? result : nil
        }
        @cache = Hash.new(&read_proc)
      end

      def add_exec_flags(flags)
        @executable = Array(executable).concat(flags)
      end

      def read_config(cmd, all = false)
        config_cmd = ['config', (all ? '--get-all' : '--get'), *cmd]
        config_cmd = config_cmd.join(' ') unless cmd.respond_to? :join
        read config_cmd
      end

      def read(cmd)
        @cache[cmd]
      end

      def stub_config_value(key, value, get = '--get')
        stub_command_output "config #{get} #{key}", value
      end

      def stub_command_output(cmd, value)
        @cache[cmd] = value.nil? ? nil : value.to_s
      end

      def stub!(values)
        @cache.update values
      end

      private

      def to_exec(args)
        args = Shellwords.shellwords(args) if args.respond_to? :to_str
        Array(executable) + Array(args)
      end

      def command_to_string(cmd)
        full_cmd = to_exec(cmd)
        full_cmd.respond_to?(:shelljoin) ? full_cmd.shelljoin : full_cmd.join(' ')
      end

      def silence_stderr
        oldio = STDERR.dup
        STDERR.reopen(NULL)
        yield
      ensure
        STDERR.reopen(oldio)
      end
    end

    module GitReaderMethods
      extend Forwardable

      def_delegator :git_reader, :read_config, :git_config
      def_delegator :git_reader, :read, :git_command

      def self.extended(base)
        base.extend Forwardable
        base.def_delegators :'self.class', :git_config, :git_command
      end
    end

    class Error < RuntimeError; end
    class FatalError < Error; end

    private

    def git_reader
      @git_reader ||= GitReader.new ENV['GIT']
    end

    include GitReaderMethods
    private :git_config, :git_command

    def local_repo(fatal = true)
      return nil if defined?(@local_repo) && @local_repo == false
      @local_repo =
        if git_dir = git_command('rev-parse -q --git-dir')
          LocalRepo.new(git_reader, current_dir, git_dir)
        elsif fatal
          raise FatalError, "Not a git repository"
        else
          false
        end
    end

    repo_methods = [
      :current_branch, :git_dir,
      :remote_branch_and_project,
      :repo_owner, :repo_host,
      :remotes, :remotes_group, :origin_remote
    ]
    def_delegator :local_repo, :name, :repo_name
    def_delegators :local_repo, *repo_methods
    private :repo_name, *repo_methods

    def master_branch
      if local_repo(false)
        local_repo.master_branch
      else
        # FIXME: duplicates functionality of LocalRepo#master_branch
        Branch.new nil, 'refs/heads/master'
      end
    end

    class LocalRepo < Struct.new(:git_reader, :dir, :git_dir)
      include GitReaderMethods

      def name
        if project = main_project
          project.name
        else
          File.basename(dir)
        end
      end

      def repo_owner
        if project = main_project
          project.owner
        end
      end

      def repo_host
        project = main_project and project.host
      end

      def main_project
        remote = origin_remote and remote.project
      end

      def remote_branch_and_project(username_fetcher, prefer_upstream = false)
        project = main_project
        if project and branch = current_branch
          branch = branch.push_target(username_fetcher.call(project.host), prefer_upstream)
          project = remote_by_name(branch.remote_name).project if branch && branch.remote?
        end
        [branch, project]
      end

      def current_branch
        @current_branch ||= branch_at_ref('HEAD')
      end

      def branch_at_ref(*parts)
        begin
          head = file_read(*parts)
        rescue Errno::ENOENT
          return nil
        else
          Branch.new(self, head.rstrip) if head.sub!('ref: ', '')
        end
      end

      def file_read(*parts)
        File.read(File.join(git_dir, *parts))
      end

      def file_exist?(*parts)
        File.exist?(File.join(git_dir, *parts))
      end

      def master_branch
        if remote = origin_remote
          default_branch = branch_at_ref("refs/remotes/#{remote}/HEAD")
        end
        default_branch || Branch.new(self, 'refs/heads/master')
      end

      ORIGIN_NAMES = %w[ upstream github origin ]

      def remotes
        @remotes ||= begin
          names = []
          url_memo = Hash.new {|h,k| names << k; h[k]=[] }
          git_command('remote -v').to_s.split("\n").map do |line|
            next if line !~ /^(.+?)\t(.+) \(/
            name, url = $1, $2
            url_memo[name] << url
          end
          ((ORIGIN_NAMES + names) & names).map do |name|
            urls = url_memo[name].uniq
            Remote.new(self, name, urls)
          end
        end
      end

      def remotes_for_publish(owner_name)
        list = ORIGIN_NAMES.map {|n| remote_by_name(n) }
        list << remotes.find {|r| p = r.project and p.owner == owner_name }
        list.compact.uniq.reverse
      end

      def remotes_group(name)
        git_config "remotes.#{name}"
      end

      def origin_remote
        remotes.detect {|r| r.urls.any? }
      end

      def remote_by_name(remote_name)
        remotes.find {|r| r.name == remote_name }
      end

      def known_host?(host)
        default = default_host
        default == host || "ssh.#{default}" == host ||
          git_config('hub.host', :all).to_s.split("\n").include?(host)
      end

      def self.default_host
        ENV['GITHUB_HOST'] || main_host
      end

      def self.main_host
        'github.com'
      end

      extend Forwardable
      def_delegators :'self.class', :default_host, :main_host

      def ssh_config
        @ssh_config ||= SshConfig.new
      end
    end

    class GithubProject < Struct.new(:local_repo, :owner, :name, :host)
      def self.from_url(url, local_repo)
        if local_repo.known_host?(url.host)
          _, owner, name = url.path.split('/', 4)
          GithubProject.new(local_repo, owner, name.sub(/\.git$/, ''), url.host)
        end
      end

      attr_accessor :repo_data

      def initialize(*args)
        super
        self.name = self.name.tr(' ', '-')
        self.host ||= (local_repo || LocalRepo).default_host
        self.host = host.sub(/^ssh\./i, '') if 'ssh.github.com' == host.downcase
      end

      def private?
        repo_data ? repo_data.fetch('private') :
          host != (local_repo || LocalRepo).main_host
      end

      def owned_by(new_owner)
        new_project = dup
        new_project.owner = new_owner
        new_project
      end

      def name_with_owner
        "#{owner}/#{name}"
      end

      def ==(other)
        name_with_owner == other.name_with_owner
      end

      def remote
        local_repo.remotes.find { |r| r.project == self }
      end

      def web_url(path = nil, protocol_config = nil)
        project_name = name_with_owner
        if project_name.sub!(/\.wiki$/, '')
          unless '/wiki' == path
            path = if path =~ %r{^/commits/} then '/_history'
                   else path.to_s.sub(/\w+/, '_\0')
                   end
            path = '/wiki' + path
          end
        end
        '%s://%s/%s' % [
          protocol_config ? protocol_config.call(host) : 'https',
          host,
          project_name + path.to_s
        ]
      end

      def git_url(options = {})
        if options[:https] then "https://#{host}/"
        elsif options[:private] or private? then "git@#{host}:"
        else "git://#{host}/"
        end + name_with_owner + '.git'
      end
    end

    class GithubURL < DelegateClass(URI::HTTP)
      extend Forwardable

      attr_reader :project
      def_delegator :project, :name, :project_name
      def_delegator :project, :owner, :project_owner

      def self.resolve(url, local_repo)
        u = URI(url)
        if %[http https].include? u.scheme and project = GithubProject.from_url(u, local_repo)
          self.new(u, project)
        end
      rescue URI::InvalidURIError
        nil
      end

      def initialize(uri, project)
        @project = project
        super(uri)
      end

      # segment of path after the project owner and name
      def project_path
        path.split('/', 4)[3]
      end
    end

    class Branch < Struct.new(:local_repo, :name)
      alias to_s name

      def short_name
        name.sub(%r{^refs/(remotes/)?.+?/}, '')
      end

      def master?
        master_name = if local_repo then local_repo.master_branch.short_name
        else 'master'
        end
        short_name == master_name
      end

      def upstream
        if branch = local_repo.git_command("rev-parse --symbolic-full-name #{short_name}@{upstream}")
          Branch.new local_repo, branch
        end
      end

      def push_target(owner_name, prefer_upstream = false)
        push_default = local_repo.git_config('push.default')
        if %w[upstream tracking].include?(push_default)
          upstream
        else
          short = short_name
          refs = local_repo.remotes_for_publish(owner_name).map { |remote|
            "refs/remotes/#{remote}/#{short}"
          }
          refs.reverse! if prefer_upstream
          if branch = refs.detect {|ref| local_repo.file_exist?(ref) }
            Branch.new(local_repo, branch)
          end
        end
      end

      def remote?
        name.index('refs/remotes/') == 0
      end

      def remote_name
        name =~ %r{^refs/remotes/([^/]+)} and $1 or
          raise Error, "can't get remote name from #{name.inspect}"
      end
    end

    class Remote < Struct.new(:local_repo, :name, :raw_urls)
      alias to_s name

      def ==(other)
        other.respond_to?(:to_str) ? name == other.to_str : super
      end

      def project
        urls.each { |url|
          if valid = GithubProject.from_url(url, local_repo)
            return valid
          end
        }
        nil
      end

      def github_url
        urls.detect {|url| local_repo.known_host?(url.host) }
      end

      def urls
        @urls ||= raw_urls.map do |url|
          with_normalized_url(url) do |normalized|
            begin
              uri_parse(normalized)
            rescue URI::InvalidURIError
            end
          end
        end.compact
      end

      def with_normalized_url(url)
        if url =~ %r{^[\w-]+://} || url =~ %r{^([^/]+?):}
          url = "ssh://#{$1}/#{$'}" if $1
          yield url
        end
      end

      def uri_parse uri
        uri = URI.parse uri
        if uri.host != local_repo.default_host
          ssh = local_repo.ssh_config
          uri.host = ssh.get_value(uri.host, :HostName) { uri.host }
        end
        uri
      end
    end

    ## helper methods for local repo, GH projects

    def github_project(name, owner = nil)
      if owner and owner.index('/')
        owner, name = owner.split('/', 2)
      elsif name and name.index('/')
        owner, name = name.split('/', 2)
      else
        name ||= repo_name
        owner ||= github_user
      end

      if local_repo(false) and main_project = local_repo.main_project
        project = main_project.dup
        project.owner = owner
        project.name = name
        project
      else
        GithubProject.new(local_repo(false), owner, name)
      end
    end

    def git_url(owner = nil, name = nil, options = {})
      project = github_project(name, owner)
      project.git_url({:https => https_protocol?}.update(options))
    end

    def resolve_github_url(url)
      GithubURL.resolve(url, local_repo) if url =~ /^https?:/
    end

    # legacy setting
    def http_clone?
      git_config('--bool hub.http-clone') == 'true'
    end

    def https_protocol?
      git_config('hub.protocol') == 'https' or http_clone?
    end

    def git_alias_for(name)
      git_config "alias.#{name}"
    end

    def rev_list(a, b)
      git_command("rev-list --cherry-pick --right-only --no-merges #{a}...#{b}")
    end

    PWD = Dir.pwd

    def current_dir
      PWD
    end

    def is_repo?
      !!local_repo(false)
    end

    def git_editor
      # possible: ~/bin/vi, $SOME_ENVIRONMENT_VARIABLE, "C:\Program Files\Vim\gvim.exe" --nofork
      editor = git_command 'var GIT_EDITOR'
      editor.gsub!(/\$(\w+|\{\w+\})/) { ENV[$1.tr('{}', '')] }
      editor = ENV[$1] if editor =~ /^\$(\w+)$/
      editor = File.expand_path editor if (editor =~ /^[~.]/ or editor.index('/')) and editor !~ /["']/
      # avoid shellsplitting "C:\Program Files"
      if File.exist? editor then [editor]
      else editor.shellsplit
      end
    end

    module System
      # Cross-platform web browser command; respects the value set in $BROWSER.
      # 
      # Returns an array, e.g.: ['open']
      def browser_launcher
        browser = ENV['BROWSER'] || (
          osx? ? 'open' : windows? ? %w[cmd /c start] :
          %w[xdg-open cygstart x-www-browser firefox opera mozilla netscape].find { |comm| which comm }
        )

        abort "Please set $BROWSER to a web launcher to use this command." unless browser
        Array(browser)
      end

      def osx?
        require 'rbconfig'
        RbConfig::CONFIG['host_os'].to_s.include?('darwin')
      end

      def windows?
        require 'rbconfig'
        RbConfig::CONFIG['host_os'] =~ /msdos|mswin|djgpp|mingw|windows/
      end

      def unix?
        require 'rbconfig'
        RbConfig::CONFIG['host_os'] =~ /(aix|darwin|linux|(net|free|open)bsd|cygwin|solaris|irix|hpux)/i
      end

      # Cross-platform way of finding an executable in the $PATH.
      #
      #   which('ruby') #=> /usr/bin/ruby
      def which(cmd)
        exts = ENV['PATHEXT'] ? ENV['PATHEXT'].split(';') : ['']
        ENV['PATH'].split(File::PATH_SEPARATOR).each do |path|
          exts.each { |ext|
            exe = "#{path}/#{cmd}#{ext}"
            return exe if File.executable? exe
          }
        end
        return nil
      end

      # Checks whether a command exists on this system in the $PATH.
      #
      # name - The String name of the command to check for.
      #
      # Returns a Boolean.
      def command?(name)
        !which(name).nil?
      end

      def tmp_dir
        ENV['TMPDIR'] || ENV['TEMP'] || '/tmp'
      end

      def terminal_width
        if unix?
          width = %x{stty size 2>#{NULL}}.split[1].to_i
          width = %x{tput cols 2>#{NULL}}.to_i if width.zero?
        else
          width = 0
        end
        width < 10 ? 78 : width
      end
    end

    include System
    extend System
  end
end
