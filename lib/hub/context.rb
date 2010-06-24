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

      if url && url.to_s =~ %r{\bgithub\.com[:/](.+)/(.+).git$}
        cache[remote] = { :user => $1, :repo => $2 }
      else
        cache[remote] = { }
      end
    end

    LGHCONF = "http://github.com/guides/local-github-config"

    def repo_owner
      REMOTES[default_remote][:user]
    end

    def repo_user
      REMOTES[current_remote][:user]
    end

    def repo_name
      REMOTES[default_remote][:repo] || File.basename(Dir.pwd)
    end

    # Either returns the GitHub user as set by git-config(1) or aborts
    # with an error message.
    def github_user(fatal = true)
      if user = GIT_CONFIG['config github.user']
        user
      elsif fatal
        abort("** No GitHub user set. See #{LGHCONF}")
      end
    end

    def github_token(fatal = true)
      if token = GIT_CONFIG['config github.token']
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
      list = GIT_CONFIG['remote'].split("\n")
      main = list.delete('origin') and list.unshift(main)
      list
    end

    def remotes_group(name)
      GIT_CONFIG["config remotes.#{name}"]
    end

    def current_remote
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

    def github_url(options = {})
      repo = options[:repo]
      user, repo = repo.split('/') if repo && repo.index('/')
      user ||= options[:user] || github_user
      repo ||= repo_name
      secure = options[:private]

      if options[:web] == 'wiki'
        scheme = secure ? 'https:' : 'http:'
        '%s//wiki.github.com/%s/%s/' % [scheme, user, repo]
      elsif options[:web]
        scheme = secure ? 'https:' : 'http:'
        path = options[:web] == true ? '' : options[:web].to_s
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
  end
end
