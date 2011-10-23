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

    # Caches output when shelling out to git
    GIT_CONFIG = ShellOutCache.new(ENV['GIT'] || 'git') do |cache, cmd|
      full_cmd = cache.to_exec(cmd)
      cmd_string = full_cmd.respond_to?(:shelljoin) ? full_cmd.shelljoin : full_cmd.join(' ')
      result = %x{#{cmd_string}}.chomp
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
      REMOTES[default_remote][:user]
    end

    def repo_user
      REMOTES[current_remote][:user]
    end

    def repo_name
      REMOTES[default_remote][:repo] || current_dirname
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

    def tracked_branch
      branch = current_branch && tracked_for(current_branch)
      normalize_branch(branch) if branch
    end

    def remotes
      list = GIT_CONFIG['remote'].to_s.split("\n")
      main = list.delete('origin') and list.unshift(main)
      list
    end

    def remotes_group(name)
      GIT_CONFIG["config remotes.#{name}"]
    end

    def current_remote
      return if remotes.empty?
      (current_branch && remote_for(current_branch)) || default_remote
    end

    def default_remote
      remotes.first
    end

    def normalize_branch(branch)
      branch.sub('refs/heads/', '')
    end

    def remote_for(branch)
      GIT_CONFIG['config branch.%s.remote' % normalize_branch(branch)]
    end

    def tracked_for(branch)
      GIT_CONFIG['config branch.%s.merge' % normalize_branch(branch)]
    end

    def http_clone?
      GIT_CONFIG['config --bool hub.http-clone'] == 'true'
    end

    def git_alias_for(name)
      GIT_CONFIG["config alias.#{name}"]
    end

    def github_url(options = {})
      repo = options[:repo]
      user, repo = repo.split('/') if repo && repo.index('/')
      user ||= options[:user] || github_user
      repo ||= repo_name
      secure = options[:private]

      if options[:web]
        scheme = secure ? 'https:' : 'http:'
        path = options[:web] == true ? '' : options[:web].to_s
        if repo =~ /\.wiki$/
          repo = repo.sub(/\.wiki$/, '')
          unless '/wiki' == path
            path = '/wiki%s' % if path =~ %r{^/commits/} then '/_history'
              else path.sub(/\w+/, '_\0')
              end
          end
        end
        '%s//github.com/%s/%s%s' % [scheme, user, repo, path]
      else
        if secure
          url = 'git@github.com:%s/%s.git'
        elsif http_clone?
          url = 'http://github.com/%s/%s.git'
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

    # Cross-platform web browser command; respects the value set in $BROWSER.
    # 
    # Returns an array, e.g.: ['open']
    def browser_launcher
      require 'rbconfig'
      browser = ENV['BROWSER'] ||
        (RbConfig::CONFIG['host_os'].include?('darwin') && 'open') ||
        (RbConfig::CONFIG['host_os'] =~ /msdos|mswin|djgpp|mingw|windows/ && 'start') ||
        %w[xdg-open cygstart x-www-browser firefox opera mozilla netscape].find { |comm| which comm }

      abort "Please set $BROWSER to a web launcher to use this command." unless browser
      Array(browser)
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
