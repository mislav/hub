module Hub
  # Provides methods for inspecting the environment, such as GitHub user/token
  # settings, repository info, and similar.
  module Context
    # Caches output when shelling out to git
    GIT_CONFIG = Hash.new do |cache, cmd|
      result = %x{git #{cmd}}.chomp
      cache[cmd] = $?.success? && !result.empty? ? result : nil
    end

    # Parses URLs for git remotes and stores info
    REMOTES = Hash.new do |cache, remote|
      url = GIT_CONFIG["config remote.#{remote}.url"]

      if url and url.to_s =~ %r{\bgithub\.com[:/](.+)/(.+).git$}
        cache[remote] = { :user => $1, :repo => $2 }
      else
        cache[remote] = { }
      end
    end

    LGHCONF = "http://github.com/guides/local-github-config"

    private

    def repo_owner
      REMOTES['origin'][:user]
    end
    alias repo_user repo_owner

    def repo_name
      REMOTES['origin'][:repo] || File.basename(Dir.pwd)
    end

    # Either returns the GitHub user as set by git-config(1) or aborts
    # with an error message.
    def github_user(fatal = true)
      GIT_CONFIG['config github.user'] or
        fatal ? abort("** No GitHub user set. See #{LGHCONF}") : nil
    end

    def github_token(fatal = true)
      GIT_CONFIG['config github.token'] or
        fatal ? abort("** No GitHub token set. See #{LGHCONF}") : nil
    end

    def http_clone?
      GIT_CONFIG['config --bool hub.http-clone'] == 'true'
    end

    def github_url(options = {})
      user, repo = options[:user], options[:repo]
      user, repo = repo.split('/') if user.nil? and repo and repo.index('/')
      user ||= github_user
      repo ||= repo_name
      secure = options[:private]

      if options[:web]
        scheme = secure ? 'https:' : 'http:'
        path = options[:web] == true ? '' : options[:web].to_s
        '%s//github.com/%s/%s%s' % [scheme, user, repo, path]
      else
        if secure
          'git@github.com:%s/%s.git'
        elsif http_clone?
          'http://github.com/%s/%s.git'
        else
          'git://github.com/%s/%s.git'
        end % [user, repo]
      end
    end
  end
end
