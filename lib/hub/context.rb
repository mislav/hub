require 'shellwords'
require 'forwardable'

module Hub
  # Provides methods for inspecting the environment, such as GitHub user/token
  # settings, repository info, and similar.
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
          result = %x{#{command_to_string(cmd)} 2>#{NULL}}.chomp
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

    private

    def git_reader
      @git_reader ||= GitReader.new ENV['GIT']
    end

    include GitReaderMethods
    private :git_config, :git_command

    def local_repo
      @local_repo ||= begin
        LocalRepo.new git_reader, current_dir if is_repo?
      end
    end

    repo_methods = [
      :current_branch, :master_branch,
      :current_project, :upstream_project,
      :repo_owner,
      :remotes, :remotes_group, :origin_remote
    ]
    def_delegator :local_repo, :name, :repo_name
    def_delegators :local_repo, *repo_methods
    private :repo_name, *repo_methods

    class LocalRepo < Struct.new(:git_reader, :dir)
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

      def main_project
        remote = origin_remote and remote.project
      end

      def upstream_project
        if upstream = current_branch.upstream
          remote = remote_by_name upstream.remote_name
          remote.project
        end
      end

      def current_project
        upstream_project || main_project
      end

      def current_branch
        if branch = git_command('symbolic-ref -q HEAD')
          Branch.new self, branch
        end
      end

      def master_branch
        Branch.new self, 'refs/heads/master'
      end

      def remotes
        @remotes ||= begin
          # TODO: is there a plumbing command to get a list of remotes?
          list = git_command('remote').to_s.split("\n")
          # force "origin" to be first in the list
          main = list.delete('origin') and list.unshift(main)
          list.map { |name| Remote.new self, name }
        end
      end

      def remotes_group(name)
        git_config "remotes.#{name}"
      end

      def origin_remote
        remotes.first
      end

      def remote_by_name(remote_name)
        remotes.find {|r| r.name == remote_name }
      end
    end

    class GithubProject < Struct.new(:local_repo, :owner, :name)
      def name_with_owner
        "#{owner}/#{name}"
      end

      def ==(other)
        name_with_owner == other.name_with_owner
      end

      def remote
        local_repo.remotes.find { |r| r.project == self }
      end

      def web_url(path = nil)
        project_name = name_with_owner
        if project_name.sub!(/\.wiki$/, '')
          unless '/wiki' == path
            path = if path =~ %r{^/commits/} then '/_history'
                   else path.to_s.sub(/\w+/, '_\0')
                   end
            path = '/wiki' + path
          end
        end
        'https://github.com/' + project_name + path.to_s
      end

      def git_url(options = {})
        if options[:https] then 'https://github.com/'
        elsif options[:private] then 'git@github.com:'
        else 'git://github.com/'
        end + name_with_owner + '.git'
      end
    end

    class Branch < Struct.new(:local_repo, :name)
      alias to_s name

      def short_name
        name.split('/').last
      end

      def master?
        short_name == 'master'
      end

      def upstream
        if branch = local_repo.git_command("rev-parse --symbolic-full-name #{short_name}@{upstream}")
          Branch.new local_repo, branch
        end
      end

      def remote?
        name.index('refs/remotes/') == 0
      end

      def remote_name
        name =~ %r{^refs/remotes/([^/]+)} and $1 or
          raise "can't get remote name from #{name.inspect}"
      end
    end

    class Remote < Struct.new(:local_repo, :name)
      alias to_s name

      def ==(other)
        other.respond_to?(:to_str) ? name == other.to_str : super
      end

      def project
        if urls.find { |u| u =~ %r{\bgithub\.com[:/](.+)/(.+).git$} }
          GithubProject.new local_repo, $1, $2
        end
      end

      def urls
        @urls ||= local_repo.git_config("remote.#{name}.url", :all).to_s.split("\n")
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

      GithubProject.new local_repo, owner, name
    end

    def git_url(owner = nil, name = nil, options = {})
      project = github_project(name, owner)
      project.git_url({:https => https_protocol?}.update(options))
    end

    LGHCONF = "http://help.github.com/git-email-settings/"

    # Either returns the GitHub user as set by git-config(1) or aborts
    # with an error message.
    def github_user(fatal = true)
      if user = ENV['GITHUB_USER'] || git_config('github.user')
        user
      elsif fatal
        abort("** No GitHub user set. See #{LGHCONF}")
      end
    end

    def github_token(fatal = true)
      if token = ENV['GITHUB_TOKEN'] || git_config('github.token')
        token
      elsif fatal
        abort("** No GitHub token set. See #{LGHCONF}")
      end
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

    PWD = Dir.pwd

    def current_dir
      PWD
    end

    def git_dir
      git_command 'rev-parse -q --git-dir'
    end

    def is_repo?
      !!git_dir
    end

    def git_editor
      # possible: ~/bin/vi, $SOME_ENVIRONMENT_VARIABLE, "C:\Program Files\Vim\gvim.exe" --nofork
      editor = git_command 'var GIT_EDITOR'
      editor = ENV[$1] if editor =~ /^\$(\w+)$/
      editor = File.expand_path editor if (editor =~ /^[~.]/ or editor.index('/')) and editor !~ /["']/
      editor.shellsplit
    end

    # Cross-platform web browser command; respects the value set in $BROWSER.
    # 
    # Returns an array, e.g.: ['open']
    def browser_launcher
      browser = ENV['BROWSER'] || (
        osx? ? 'open' : windows? ? 'start' :
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
  end
end
