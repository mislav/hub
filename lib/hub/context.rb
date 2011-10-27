require 'shellwords'

module Hub
  # Provides methods for inspecting the environment, such as GitHub user/token
  # settings, repository info, and similar.
  module Context
    private

    class ShellOutCache < Hash
      attr_accessor :executable

      def initialize(executable = nil, &block)
        super(&block)
        @executable = executable
      end

      def to_exec(args)
        args = Shellwords.shellwords(args) if args.respond_to? :to_str
        Array(executable) + Array(args)
      end
    end

    NULL = defined?(File::NULL) ? File::NULL : File.exist?('/dev/null') ? '/dev/null' : 'NUL'

    # Caches output when shelling out to git
    GIT_CONFIG = ShellOutCache.new(ENV['GIT'] || 'git') do |cache, cmd|
      full_cmd = cache.to_exec(cmd)
      cmd_string = full_cmd.respond_to?(:shelljoin) ? full_cmd.shelljoin : full_cmd.join(' ')
      result = %x{#{cmd_string} 2>#{NULL}}.chomp
      cache[cmd] = $?.success? && !result.empty? ? result : nil
    end

    # Parses URLs for git remotes and stores info
    REMOTES = Hash.new do |cache, remote|
      if remote
        urls = GIT_CONFIG["config --get-all remote.#{remote}.url"].to_s.split("\n")

        if urls.find { |u| u =~ %r{\bgithub\.com[:/](.+)/(.+).git$} } 
          cache[remote] = { :user => $1, :repo => $2 }
        else
          cache[remote] = { }
        end
      else
        cache[remote] = { }
      end
    end

    LGHCONF = "http://help.github.com/git-email-settings/"

    def repo_owner
      if repo = origin_repo
        repo.owner
      end
    end

    def repo_user
      if repo = current_repo
        repo.owner
      end
    end

    def repo_name
      if repo = origin_repo
        repo.name
      else
        current_dirname
      end
    end

    def origin_repo
      if remote = default_remote
        remote.repo
      end
    end

    def current_repo
      if remote = current_remote
        remote.repo
      end
    end

    class Remote < Struct.new(:name)
      def self.all
        @remotes ||= begin
          list = GIT_CONFIG['remote'].to_s.split("\n")
          main = list.delete('origin') and list.unshift(main)
          list.map { |name| new(name) }
        end
      end

      def self.clear!
        remove_instance_variable '@remotes' if defined? @remotes
      end

      def self.origin
        all.first
      end

      def self.by_name(remote_name)
        all.find {|r| r.name == remote_name }
      end

      def self.for_repo(repo)
        all.find {|r| r.repo == repo }
      end

      alias to_s name

      def ==(other)
        other.respond_to?(:to_str) ? name == other.to_str : super
      end

      def repo
        return @repo if defined? @repo
        @repo = if urls.find { |u| u =~ %r{\bgithub\.com[:/](.+)/(.+).git$} } 
          Repo.new($1, $2)
        end
      end

      def urls
        @urls ||= GIT_CONFIG["config --get-all remote.#{name}.url"].to_s.split("\n")
      end
    end

    class Repo < Struct.new(:owner, :name)
      def self.from_string(str, context_repo = nil)
        if str.index('/')
          new(*str.split('/', 2))
        else
          new(str, context_repo.name)
        end
      end

      def from_owner(another_owner, another_name = nil)
        self.class.new(another_owner, another_name || self.name)
      end

      def ref_for(branch)
        Ref.new(self, branch)
      end
    end

    class Ref < Struct.new(:repo, :branch, :remote)
      def self.from_github_ref(branch, context_repo)
        if branch.index(':')
          owner_with_name, branch = branch.split(':', 2)
          context_repo = context_repo.from_owner(*owner_with_name.split('/', 2))
        end
        new(context_repo, branch)
      end

      def self.upstream_for(branch, repo = nil)
        if GIT_CONFIG["name-rev #{branch}@{upstream} --name-only --refs='refs/remotes/*' --no-undefined"] =~ %r{^remotes/(.+)}
          remote_name, branch = $1.split('/', 2)
          new(repo, branch, Remote.by_name(remote_name))
        end
      end

      def initialize(*args)
        super
        if remote.nil?
          self.remote = Remote.for_repo(repo)
        elsif repo.nil?
          self.repo = remote.repo
        end
      end

      def to_local_ref
        "#{remote}/#{branch}"
      end

      def to_github_ref
        "#{repo.owner}:#{branch}"
      end
    end

    # Either returns the GitHub user as set by git-config(1) or aborts
    # with an error message.
    def github_user(fatal = true)
      if user = ENV['GITHUB_USER'] || GIT_CONFIG['config github.user']
        user
      elsif fatal
        abort("** No GitHub user set. See #{LGHCONF}")
      end
    end

    def github_token(fatal = true)
      if token = ENV['GITHUB_TOKEN'] || GIT_CONFIG['config github.token']
        token
      elsif fatal
        abort("** No GitHub token set. See #{LGHCONF}")
      end
    end

    def current_branch
      GIT_CONFIG['symbolic-ref -q HEAD']
    end

    def upstream_ref(branch, repo = nil)
      Ref.upstream_for(normalize_branch(branch), repo)
    end

    def tracked_branch
      current_branch and ref = upstream_ref(current_branch) and ref.branch
    end

    def remotes
      Remote.all
    end

    def remotes_group(name)
      GIT_CONFIG["config remotes.#{name}"]
    end

    def current_remote
      return if remotes.empty?
      (ref = current_branch && upstream_ref(current_branch)) and ref.remote or default_remote
    end

    def default_remote
      Remote.origin
    end

    def normalize_branch(branch)
      branch.sub('refs/heads/', '')
    end

    # legacy setting
    def http_clone?
      GIT_CONFIG['config --bool hub.http-clone'] == 'true'
    end

    def https_protocol?
      GIT_CONFIG['config hub.protocol'] == 'https' or http_clone?
    end

    def git_alias_for(name)
      GIT_CONFIG["config alias.#{name}"]
    end

    def github_url(options = {})
      repo = options[:repo]
      user, repo = repo.split('/') if repo && repo.index('/')
      user ||= options[:user] || github_user
      repo ||= repo_name

      if options[:web]
        path = options[:web] == true ? '' : options[:web].to_s
        if repo =~ /\.wiki$/
          repo = repo.sub(/\.wiki$/, '')
          unless '/wiki' == path
            path = '/wiki%s' % if path =~ %r{^/commits/} then '/_history'
              else path.sub(/\w+/, '_\0')
              end
          end
        end
        'https://github.com/%s/%s%s' % [user, repo, path]
      else
        if https_protocol?
          url = 'https://github.com/%s/%s.git'
        elsif options[:private]
          url = 'git@github.com:%s/%s.git'
        else
          url = 'git://github.com/%s/%s.git'
        end

        url % [user, repo]
      end
    end

    DIRNAME = File.basename(Dir.pwd)

    def current_dirname
      DIRNAME
    end

    def git_dir
      GIT_CONFIG['rev-parse --git-dir']
    end

    def is_repo?
      !!git_dir
    end

    def git_editor
      # possible: ~/bin/vi, $SOME_ENVIRONMENT_VARIABLE, "C:\Program Files\Vim\gvim.exe" --nofork
      editor = GIT_CONFIG['var GIT_EDITOR']
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
