module Hub
  # The Commands module houses the git commands that hub
  # lovingly wraps. If a method exists here, it is expected to have a
  # corresponding git command which either gets run before or after
  # the method executes.
  #
  # The typical flow is as follows:
  #
  # 1. hub is invoked from the command line:
  #    $ hub clone rtomayko/tilt
  #
  # 2. The Runner class is initialized:
  #    >> Hub::Runner.new('clone', 'rtomayko/tilt')
  #
  # 3. The method representing the git subcommand is executed with the
  #    full args:
  #    >> Hub::Commands.clone(['clone', 'rtomayko/tilt'])
  #
  # 4. That method rewrites the args as it sees fit:
  #    >> args[1] = "git://github.com/" + args[1] + ".git"
  #    => "git://github.com/rtomayko/tilt.git"
  #
  # 5. The new args are used to run `git`:
  #    >> exec "git", "clone", "git://github.com/rtomayko/tilt.git"
  #
  # An optional `after` callback can be set. If so, it is run after
  # step 5 (which then performs a `system` call rather than an
  # `exec`). See `Hub::Args` for more information on the `after` callback.
  module Commands
    # We are a blank slate.
    instance_methods.each { |m| undef_method(m) unless m =~ /(^__|send|to\?$)/ }
    extend self

    # provides git interrogation methods
    extend Context

    NAME_RE = /[\w.][\w.-]*/
    OWNER_RE = /[a-zA-Z0-9][a-zA-Z0-9-]*/
    NAME_WITH_OWNER_RE = /^(?:#{NAME_RE}|#{OWNER_RE}\/#{NAME_RE})$/

    CUSTOM_COMMANDS = %w[alias create browse compare fork pull-request]

    def run(args)
      slurp_global_flags(args)

      # Hack to emulate git-style
      args.unshift 'help' if args.empty?

      cmd = args[0]
      if expanded_args = expand_alias(cmd)
        cmd = expanded_args[0]
        expanded_args.concat args[1..-1]
      end

      respect_help_flags(expanded_args || args) if custom_command? cmd

      # git commands can have dashes
      cmd = cmd.gsub(/(\w)-/, '\1_')
      if method_defined?(cmd) and cmd != 'run'
        args.replace expanded_args if expanded_args
        send(cmd, args)
      end
    rescue Errno::ENOENT
      if $!.message.include? "No such file or directory - git"
        abort "Error: `git` command not found"
      else
        raise
      end
    rescue Context::FatalError => err
      abort "fatal: #{err.message}"
    end

    # $ hub pull-request
    # $ hub pull-request "My humble contribution"
    # $ hub pull-request -i 92
    # $ hub pull-request https://github.com/rtomayko/tilt/issues/92
    def pull_request(args)
      args.shift
      options = { }
      force = explicit_owner = false
      base_project = local_repo.main_project
      head_project = local_repo.current_project

      unless current_branch
        abort "Aborted: not currently on any branch."
      end

      unless base_project
        abort "Aborted: the origin remote doesn't point to a GitHub repository."
      end

      from_github_ref = lambda do |ref, context_project|
        if ref.index(':')
          owner, ref = ref.split(':', 2)
          project = github_project(context_project.name, owner)
        end
        [project || context_project, ref]
      end

      while arg = args.shift
        case arg
        when '-f'
          force = true
        when '-b'
          base_project, options[:base] = from_github_ref.call(args.shift, base_project)
        when '-h'
          head = args.shift
          explicit_owner = !!head.index(':')
          head_project, options[:head] = from_github_ref.call(head, head_project)
        when '-i'
          options[:issue] = args.shift
        else
          if url = resolve_github_url(arg) and url.project_path =~ /^issues\/(\d+)/
            options[:issue] = $1
            base_project = url.project
          elsif !options[:title] then options[:title] = arg
          else
            abort "invalid argument: #{arg}"
          end
        end
      end

      options[:project] = base_project
      options[:base] ||= master_branch.short_name

      if tracked_branch = options[:head].nil? && current_branch.upstream
        if !tracked_branch.remote?
          # The current branch is tracking another local branch. Pretend there is
          # no upstream configuration at all.
          tracked_branch = nil
        elsif base_project == head_project and tracked_branch.short_name == options[:base]
          $stderr.puts "Aborted: head branch is the same as base (#{options[:base].inspect})"
          warn "(use `-h <branch>` to specify an explicit pull request head)"
          abort
        end
      end
      options[:head] ||= (tracked_branch || current_branch).short_name

      # when no tracking, assume remote branch is published under active user's fork
      user = github_user(head_project.host)
      if head_project.owner != user and !tracked_branch and !explicit_owner
        head_project = head_project.owned_by(user)
      end

      remote_branch = "#{head_project.remote}/#{options[:head]}"
      options[:head] = "#{head_project.owner}:#{options[:head]}"

      if !force and tracked_branch and local_commits = rev_list(remote_branch, nil)
        $stderr.puts "Aborted: #{local_commits.split("\n").size} commits are not yet pushed to #{remote_branch}"
        warn "(use `-f` to force submit a pull request anyway)"
        abort
      end

      if args.noop?
        puts "Would request a pull to #{base_project.owner}:#{options[:base]} from #{options[:head]}"
        exit
      end

      unless options[:title] or options[:issue]
        base_branch = "#{base_project.remote}/#{options[:base]}"
        commits = rev_list(base_branch, remote_branch).to_s.split("\n")

        case commits.size
        when 0
          default_message = commit_summary = nil
        when 1
          format = '%w(78,0,0)%s%n%+b'
          default_message = git_command "show -s --format='#{format}' #{commits.first}"
          commit_summary = nil
        else
          format = '%h (%aN, %ar)%n%w(78,3,3)%s%n%+b'
          default_message = nil
          commit_summary = git_command "log --no-color --format='%s' --cherry %s...%s" %
            [format, base_branch, remote_branch]
        end

        options[:title], options[:body] = pullrequest_editmsg(commit_summary) { |msg|
          msg.puts default_message if default_message
          msg.puts ""
          msg.puts "# Requesting a pull to #{base_project.owner}:#{options[:base]} from #{options[:head]}"
          msg.puts "#"
          msg.puts "# Write a message for this pull request. The first block"
          msg.puts "# of text is the title and the rest is description."
        }
      end

      pull = api_client.create_pullrequest(options)

      args.executable = 'echo'
      args.replace [pull['html_url']]
    rescue GitHubAPI::Exceptions
      response = $!.response
      display_api_exception("creating pull request", response)
      if 404 == response.status
        base_url = base_project.web_url.split('://', 2).last
        warn "Are you sure that #{base_url} exists?"
      end
      exit 1
    end

    # $ hub clone rtomayko/tilt
    # > git clone git://github.com/rtomayko/tilt.
    #
    # $ hub clone -p kneath/hemingway
    # > git clone git@github.com:kneath/hemingway.git
    #
    # $ hub clone tilt
    # > git clone git://github.com/YOUR_LOGIN/tilt.
    #
    # $ hub clone -p github
    # > git clone git@github.com:YOUR_LOGIN/hemingway.git
    def clone(args)
      ssh = args.delete('-p')
      has_values = /^(--(upload-pack|template|depth|origin|branch|reference|name)|-[ubo])$/

      idx = 1
      while idx < args.length
        arg = args[idx]
        if arg.index('-') == 0
          idx += 1 if arg =~ has_values
        else
          # $ hub clone rtomayko/tilt
          # $ hub clone tilt
          if arg =~ NAME_WITH_OWNER_RE and !File.directory?(arg)
            name, owner = arg, nil
            owner, name = name.split('/', 2) if name.index('/')
            project = github_project(name, owner || github_user)
            ssh ||= args[0] != 'submodule' && project.owner == github_user(project.host) { }
            args[idx] = project.git_url(:private => ssh, :https => https_protocol?)
          end
          break
        end
        idx += 1
      end
    end

    # $ hub submodule add wycats/bundler vendor/bundler
    # > git submodule add git://github.com/wycats/bundler.git vendor/bundler
    #
    # $ hub submodule add -p wycats/bundler vendor/bundler
    # > git submodule add git@github.com:wycats/bundler.git vendor/bundler
    #
    # $ hub submodule add -b ryppl --name pip ryppl/pip vendor/pip
    # > git submodule add -b ryppl --name pip git://github.com/ryppl/pip.git vendor/pip
    def submodule(args)
      return unless index = args.index('add')
      args.delete_at index

      clone(args)
      args.insert index, 'add'
    end

    # $ hub remote add pjhyett
    # > git remote add pjhyett git://github.com/pjhyett/THIS_REPO.git
    #
    # $ hub remote add -p mojombo
    # > git remote add mojombo git@github.com:mojombo/THIS_REPO.git
    #
    # $ hub remote add origin
    # > git remote add origin git://github.com/YOUR_LOGIN/THIS_REPO.git
    def remote(args)
      if %w[add set-url].include?(args[1])
        name = args.last
        if name =~ /^(#{OWNER_RE})$/ || name =~ /^(#{OWNER_RE})\/(#{NAME_RE})$/
          user, repo = $1, $2 || repo_name
        end
      end
      return unless user # do not touch arguments

      ssh = args.delete('-p')

      if args.words[2] == 'origin' && args.words[3].nil?
        # Origin special case triggers default user/repo
        user, repo = github_user, repo_name
      elsif args.words[-2] == args.words[1]
        # rtomayko/tilt => rtomayko
        # Make sure you dance around flags.
        idx = args.index( args.words[-1] )
        args[idx] = user
      else
        # They're specifying the remote name manually (e.g.
        # git remote add blah rtomayko/tilt), so just drop the last
        # argument.
        args.pop
      end

      args << git_url(user, repo, :private => ssh)
    end

    # $ hub fetch mislav
    # > git remote add mislav git://github.com/mislav/REPO.git
    # > git fetch mislav
    #
    # $ hub fetch --multiple mislav xoebus
    # > git remote add mislav ...
    # > git remote add xoebus ...
    # > git fetch --multiple mislav xoebus
    def fetch(args)
      # $ hub fetch --multiple <name1>, <name2>, ...
      if args.include?('--multiple')
        names = args.words[1..-1]
      # $ hub fetch <name>
      elsif remote_name = args.words[1]
        # $ hub fetch <name1>,<name2>,...
        if remote_name =~ /^\w+(,\w+)+$/
          index = args.index(remote_name)
          args.delete(remote_name)
          names = remote_name.split(',')
          args.insert(index, *names)
          args.insert(index, '--multiple')
        else
          names = [remote_name]
        end
      else
        names = []
      end

      projects = names.map { |name|
        unless name !~ /^#{OWNER_RE}$/ or remotes.include?(name) or remotes_group(name)
          project = github_project(nil, name)
          repo_info = api_client.repo_info(project)
          if repo_info.success?
            project.repo_data = repo_info.data
            project
          end
        end
      }.compact

      if projects.any?
        projects.each do |project|
          args.before ['remote', 'add', project.owner, project.git_url(:https => https_protocol?)]
        end
      end
    end

    # $ git checkout https://github.com/defunkt/hub/pull/73
    # > git remote add -f -t feature git://github:com/mislav/hub.git
    # > git checkout --track -B mislav-feature mislav/feature
    def checkout(args)
      _, url_arg, new_branch_name = args.words
      if url = resolve_github_url(url_arg) and url.project_path =~ /^pull\/(\d+)/
        pull_id = $1
        pull_data = api_client.pullrequest_info(url.project, pull_id)

        args.delete new_branch_name
        user, branch = pull_data['head']['label'].split(':', 2)
        abort "Error: #{user}'s fork is not available anymore" unless pull_data['head']['repo']
        new_branch_name ||= "#{user}-#{branch}"

        if remotes.include? user
          args.before ['remote', 'set-branches', '--add', user, branch]
          args.before ['fetch', user, "+refs/heads/#{branch}:refs/remotes/#{user}/#{branch}"]
        else
          url = github_project(url.project_name, user).git_url(:private => pull_data['head']['repo']['private'],
                                                               :https => https_protocol?)
          args.before ['remote', 'add', '-f', '-t', branch, user, url]
        end
        idx = args.index url_arg
        args.delete_at idx
        args.insert idx, '--track', '-B', new_branch_name, "#{user}/#{branch}"
      end
    end

    # $ git merge https://github.com/defunkt/hub/pull/73
    # > git fetch git://github.com/mislav/hub.git +refs/heads/feature:refs/remotes/mislav/feature
    # > git merge mislav/feature --no-ff -m 'Merge pull request #73 from mislav/feature...'
    def merge(args)
      _, url_arg = args.words
      if url = resolve_github_url(url_arg) and url.project_path =~ /^pull\/(\d+)/
        pull_id = $1
        pull_data = api_client.pullrequest_info(url.project, pull_id)

        user, branch = pull_data['head']['label'].split(':', 2)
        abort "Error: #{user}'s fork is not available anymore" unless pull_data['head']['repo']

        url = github_project(url.project_name, user).git_url(:private => pull_data['head']['repo']['private'],
                                                             :https => https_protocol?)

        merge_head = "#{user}/#{branch}"
        args.before ['fetch', url, "+refs/heads/#{branch}:refs/remotes/#{merge_head}"]

        idx = args.index url_arg
        args.delete_at idx
        args.insert idx, merge_head, '--no-ff', '-m',
                    "Merge pull request ##{pull_id} from #{merge_head}\n\n#{pull_data['title']}"
      end
    end

    # $ git cherry-pick http://github.com/mislav/hub/commit/a319d88#comments
    # > git remote add -f mislav git://github.com/mislav/hub.git
    # > git cherry-pick a319d88
    #
    # $ git cherry-pick mislav@a319d88
    # > git remote add -f mislav git://github.com/mislav/hub.git
    # > git cherry-pick a319d88
    #
    # $ git cherry-pick mislav@SHA
    # > git fetch mislav
    # > git cherry-pick SHA
    def cherry_pick(args)
      unless args.include?('-m') or args.include?('--mainline')
        ref = args.words.last
        if url = resolve_github_url(ref) and url.project_path =~ /^commit\/([a-f0-9]{7,40})/
          sha = $1
          project = url.project
        elsif ref =~ /^(#{OWNER_RE})@([a-f0-9]{7,40})$/
          owner, sha = $1, $2
          project = local_repo.main_project.owned_by(owner)
        end

        if project
          args[args.index(ref)] = sha

          if remote = project.remote and remotes.include? remote
            args.before ['fetch', remote.to_s]
          else
            args.before ['remote', 'add', '-f', project.owner, project.git_url(:https => https_protocol?)]
          end
        end
      end
    end

    # $ hub am https://github.com/defunkt/hub/pull/55
    # > curl https://github.com/defunkt/hub/pull/55.patch -o /tmp/55.patch
    # > git am /tmp/55.patch
    def am(args)
      if url = args.find { |a| a =~ %r{^https?://(gist\.)?github\.com/} }
        idx = args.index(url)
        gist = $1 == 'gist.'
        # strip the fragment part of the url
        url = url.sub(/#.+/, '')
        # strip extra path from "pull/42/files", "pull/42/commits"
        url = url.sub(%r{(/pull/\d+)/\w*$}, '\1') unless gist
        ext = gist ? '.txt' : '.patch'
        url += ext unless File.extname(url) == ext
        patch_file = File.join(tmp_dir, "#{gist ? 'gist-' : ''}#{File.basename(url)}")
        # TODO: remove dependency on curl
        args.before 'curl', ['-#LA', "hub #{Hub::Version}", url, '-o', patch_file]
        args[idx] = patch_file
      end
    end

    # $ hub apply https://github.com/defunkt/hub/pull/55
    # > curl https://github.com/defunkt/hub/pull/55.patch -o /tmp/55.patch
    # > git apply /tmp/55.patch
    alias_method :apply, :am

    # $ hub init -g
    # > git init
    # > git remote add origin git@github.com:USER/REPO.git
    def init(args)
      if args.delete('-g')
        project = github_project(File.basename(current_dir))
        url = project.git_url(:private => true, :https => https_protocol?)
        args.after ['remote', 'add', 'origin', url]
      end
    end

    # $ hub fork
    # ... hardcore forking action ...
    # > git remote add -f YOUR_USER git@github.com:YOUR_USER/CURRENT_REPO.git
    def fork(args)
      unless project = local_repo.main_project
        abort "Error: repository under 'origin' remote is not a GitHub project"
      end
      forked_project = project.owned_by(github_user(project.host))

      existing_repo = api_client.repo_info(forked_project)
      if existing_repo.success?
        parent_data = existing_repo.data['parent']
        parent_url  = parent_data && resolve_github_url(parent_data['html_url'])
        if !parent_url or parent_url.project != project
          abort "Error creating fork: %s already exists on %s" %
            [ forked_project.name_with_owner, forked_project.host ]
        end
      else
        api_client.fork_repo(project) unless args.noop?
      end

      if args.include?('--no-remote')
        exit
      else
        url = forked_project.git_url(:private => true, :https => https_protocol?)
        args.replace %W"remote add -f #{forked_project.owner} #{url}"
        args.after 'echo', ['new remote:', forked_project.owner]
      end
    rescue GitHubAPI::Exceptions
      display_api_exception("creating fork", $!.response)
      exit 1
    end

    # $ hub create
    # ... create repo on github ...
    # > git remote add -f origin git@github.com:YOUR_USER/CURRENT_REPO.git
    def create(args)
      if !is_repo?
        abort "'create' must be run from inside a git repository"
      else
        owner = github_user
        args.shift
        options = {}
        options[:private] = true if args.delete('-p')
        new_repo_name = nil

        until args.empty?
          case arg = args.shift
          when '-d'
            options[:description] = args.shift
          when '-h'
            options[:homepage] = args.shift
          else
            if arg =~ /^[^-]/ and new_repo_name.nil?
              new_repo_name = arg
              owner, new_repo_name = new_repo_name.split('/', 2) if new_repo_name.index('/')
            else
              abort "invalid argument: #{arg}"
            end
          end
        end
        new_repo_name ||= repo_name
        new_project = github_project(new_repo_name, owner)

        if api_client.repo_exists?(new_project)
          warn "#{new_project.name_with_owner} already exists on #{new_project.host}"
          action = "set remote origin"
        else
          action = "created repository"
          unless args.noop?
            repo_data = api_client.create_repo(new_project, options)
            new_project = github_project(repo_data['full_name'])
          end
        end

        url = new_project.git_url(:private => true, :https => https_protocol?)

        if remotes.first != 'origin'
          args.replace %W"remote add -f origin #{url}"
        else
          args.replace %W"remote -v"
        end

        args.after 'echo', ["#{action}:", new_project.name_with_owner]
      end
    rescue GitHubAPI::Exceptions
      display_api_exception("creating repository", $!.response)
      exit 1
    end

    # $ hub push origin,staging cool-feature
    # > git push origin cool-feature
    # > git push staging cool-feature
    def push(args)
      return if args[1].nil? || !args[1].index(',')

      refs    = args.words[2..-1]
      remotes = args[1].split(',')
      args[1] = remotes.shift

      if refs.empty?
        # add current branch as explicit ref when there are no refs specified
        refs = [current_branch.short_name]
        args.concat refs
      end

      remotes.each do |name|
        args.after ['push', name, *refs]
      end
    end

    # $ hub browse
    # > open https://github.com/CURRENT_REPO
    #
    # $ hub browse -- issues
    # > open https://github.com/CURRENT_REPO/issues
    #
    # $ hub browse pjhyett/github-services
    # > open https://github.com/pjhyett/github-services
    #
    # $ hub browse github-services
    # > open https://github.com/YOUR_LOGIN/github-services
    #
    # $ hub browse github-services wiki
    # > open https://github.com/YOUR_LOGIN/github-services/wiki
    def browse(args)
      args.shift
      browse_command(args) do
        dest = args.shift
        dest = nil if dest == '--'

        if dest
          # $ hub browse pjhyett/github-services
          # $ hub browse github-services
          project = github_project dest
          branch = master_branch
        else
          # $ hub browse
          project = current_project
          branch = current_branch && current_branch.upstream || master_branch
        end

        abort "Usage: hub browse [<USER>/]<REPOSITORY>" unless project

        # $ hub browse -- wiki
        path = case subpage = args.shift
        when 'commits'
          "/commits/#{branch.short_name}"
        when 'tree', NilClass
          "/tree/#{branch.short_name}" if branch and !branch.master?
        else
          "/#{subpage}"
        end

        project.web_url(path)
      end
    end

    # $ hub compare 1.0..fix
    # > open https://github.com/CURRENT_REPO/compare/1.0...fix
    # $ hub compare refactor
    # > open https://github.com/CURRENT_REPO/compare/refactor
    # $ hub compare myfork feature
    # > open https://github.com/myfork/REPO/compare/feature
    # $ hub compare -u 1.0...2.0
    # "https://github.com/CURRENT_REPO/compare/1.0...2.0"
    def compare(args)
      args.shift
      browse_command(args) do
        if args.empty?
          branch = current_branch.upstream
          if branch and not branch.master?
            range = branch.short_name
            project = current_project
          else
            abort "Usage: hub compare [USER] [<START>...]<END>"
          end
        else
          sha_or_tag = /((?:#{OWNER_RE}:)?\w[\w.-]+\w)/
          # replaces two dots with three: "sha1...sha2"
          range = args.pop.sub(/^#{sha_or_tag}\.\.#{sha_or_tag}$/, '\1...\2')
          project = if owner = args.pop then github_project(nil, owner)
                    else current_project
                    end
        end

        project.web_url "/compare/#{range}"
      end
    end

    # $ hub hub standalone
    # Prints the "standalone" version of hub for an easy, memorable
    # installation sequence:
    #
    # $ gem install hub
    # $ hub hub standalone > ~/bin/hub && chmod 755 ~/bin/hub
    # $ gem uninstall hub
    def hub(args)
      return help(args) unless args[1] == 'standalone'
      require 'hub/standalone'
      Hub::Standalone.build $stdout
      exit
    rescue LoadError
      abort "hub is already running in standalone mode."
    rescue Errno::EPIPE
      exit # ignore broken pipe
    end

    def alias(args)
      shells = %w[bash zsh sh ksh csh fish]

      script = !!args.delete('-s')
      shell = args[1] || ENV['SHELL']
      abort "hub alias: unknown shell" if shell.nil? or shell.empty?
      shell = File.basename shell

      unless shells.include? shell
        $stderr.puts "hub alias: unsupported shell"
        warn "supported shells: #{shells.join(' ')}"
        abort
      end

      if script
        puts "alias git=hub"
        if 'zsh' == shell
          puts "if type compdef >/dev/null; then"
          puts "   compdef hub=git"
          puts "fi"
        end
      else
        profile = case shell
          when 'bash' then '~/.bash_profile'
          when 'zsh'  then '~/.zshrc'
          when 'ksh'  then '~/.profile'
          else
            'your profile'
          end

        puts "# Wrap git automatically by adding the following to #{profile}:"
        puts
        puts 'eval "$(hub alias -s)"'
      end

      exit
    end

    # $ hub version
    # > git version
    # (print hub version)
    def version(args)
      args.after 'echo', ['hub version', Version]
    end
    alias_method "--version", :version

    # $ hub help
    # (print improved help text)
    def help(args)
      command = args.words[1]

      if command == 'hub'
        puts hub_manpage
        exit
      elsif command.nil? && !args.has_flag?('-a', '--all')
        ENV['GIT_PAGER'] = '' unless args.has_flag?('-p', '--paginate') # Use `cat`.
        puts improved_help_text
        exit
      end
    end
    alias_method "--help", :help

  private
    #
    # Helper methods are private so they cannot be invoked
    # from the command line.
    #

    def api_client
      @api_client ||= begin
        config_file = ENV['HUB_CONFIG'] || '~/.config/hub'
        file_store = GitHubAPI::FileStore.new File.expand_path(config_file)
        file_config = GitHubAPI::Configuration.new file_store
        GitHubAPI.new file_config, :app_url => 'http://defunkt.io/hub/'
      end
    end

    def github_user host = nil, &block
      host ||= (local_repo(false) || Context::LocalRepo).default_host
      api_client.config.username(host, &block)
    end

    def custom_command? cmd
      CUSTOM_COMMANDS.include? cmd
    end

    # Show short usage help for `-h` flag, and open man page for `--help`
    def respect_help_flags args
      return if args.size > 2
      case args[1]
      when '-h'
        pattern = /(git|hub) #{Regexp.escape args[0].gsub('-', '\-')}/
        hub_raw_manpage.each_line { |line|
          if line =~ pattern
            $stderr.print "Usage: "
            $stderr.puts line.gsub(/\\f./, '').gsub('\-', '-')
            abort
          end
        }
        abort "Error: couldn't find usage help for #{args[0]}"
      when '--help'
        puts hub_manpage
        exit
      end
    end

    # The text print when `hub help` is run, kept in its own method
    # for the convenience of the author.
    def improved_help_text
      <<-help
usage: git [--version] [--exec-path[=<path>]] [--html-path] [--man-path] [--info-path]
           [-p|--paginate|--no-pager] [--no-replace-objects] [--bare]
           [--git-dir=<path>] [--work-tree=<path>] [--namespace=<name>]
           [-c name=value] [--help]
           <command> [<args>]

Basic Commands:
   init       Create an empty git repository or reinitialize an existing one
   add        Add new or modified files to the staging area
   rm         Remove files from the working directory and staging area
   mv         Move or rename a file, a directory, or a symlink
   status     Show the status of the working directory and staging area
   commit     Record changes to the repository

History Commands:
   log        Show the commit history log
   diff       Show changes between commits, commit and working tree, etc
   show       Show information about commits, tags or files

Branching Commands:
   branch     List, create, or delete branches
   checkout   Switch the active branch to another branch
   merge      Join two or more development histories (branches) together
   tag        Create, list, delete, sign or verify a tag object

Remote Commands:
   clone      Clone a remote repository into a new directory
   fetch      Download data, tags and branches from a remote repository
   pull       Fetch from and merge with another repository or a local branch
   push       Upload data, tags and branches to a remote repository
   remote     View and manage a set of remote repositories

Advanced Commands:
   reset      Reset your staging area or working directory to another point
   rebase     Re-apply a series of patches in one branch onto another
   bisect     Find by binary search the change that introduced a bug
   grep       Print files with lines matching a pattern in your codebase

GitHub Commands:
   pull-request   Open a pull request on GitHub
   fork           Make a fork of a remote repository on GitHub and add as remote
   create         Create this repository on GitHub and add GitHub as origin
   browse         Open a GitHub page in the default browser
   compare        Open a compare page on GitHub

See 'git help <command>' for more information on a specific command.
help
    end

    # Extract global flags from the front of the arguments list.
    # Makes sure important ones are supplied for calls to subcommands.
    #
    # Known flags are:
    #   --version --exec-path=<path> --html-path
    #   -p|--paginate|--no-pager --no-replace-objects
    #   --bare --git-dir=<path> --work-tree=<path>
    #   -c name=value --help
    #
    # Special: `--version`, `--help` are replaced with "version" and "help".
    # Ignored: `--exec-path`, `--html-path` are kept in args list untouched.
    def slurp_global_flags(args)
      flags = %w[ --noop -c -p --paginate --no-pager --no-replace-objects --bare --version --help ]
      flags2 = %w[ --exec-path= --git-dir= --work-tree= ]

      # flags that should be present in subcommands, too
      globals = []
      # flags that apply only to main command
      locals = []

      while args[0] && (flags.include?(args[0]) || flags2.any? {|f| args[0].index(f) == 0 })
        flag = args.shift
        case flag
        when '--noop'
          args.noop!
        when '--version', '--help'
          args.unshift flag.sub('--', '')
        when '-c'
          # slurp one additional argument
          config_pair = args.shift
          # add configuration to our local cache
          key, value = config_pair.split('=', 2)
          git_reader.stub_config_value(key, value)

          globals << flag << config_pair
        when '-p', '--paginate', '--no-pager'
          locals << flag
        else
          globals << flag
        end
      end

      git_reader.add_exec_flags(globals)
      args.add_exec_flags(globals)
      args.add_exec_flags(locals)
    end

    # Handles common functionality of browser commands like `browse`
    # and `compare`. Yields a block that returns params for `github_url`.
    def browse_command(args)
      url_only = args.delete('-u')
      warn "Warning: the `-p` flag has no effect anymore" if args.delete('-p')
      url = yield

      args.executable = url_only ? 'echo' : browser_launcher
      args.push url
    end

    # Returns the terminal-formatted manpage, ready to be printed to
    # the screen.
    def hub_manpage
      abort "** Can't find groff(1)" unless command?('groff')

      require 'open3'
      out = nil
      Open3.popen3(groff_command) do |stdin, stdout, _|
        stdin.puts hub_raw_manpage
        stdin.close
        out = stdout.read.strip
      end
      out
    end

    # The groff command complete with crazy arguments we need to run
    # in order to turn our raw roff (manpage markup) into something
    # readable on the terminal.
    def groff_command
      "groff -Wall -mtty-char -mandoc -Tascii"
    end

    # Returns the raw hub manpage. If we're not running in standalone
    # mode, it's a file sitting at the root under the `man`
    # directory.
    #
    # If we are running in standalone mode the manpage will be
    # included after the __END__ of the file so we can grab it using
    # DATA.
    def hub_raw_manpage
      if File.exists? file = File.dirname(__FILE__) + '/../../man/hub.1'
        File.read(file)
      else
        DATA.read
      end
    end

    # All calls to `puts` in after hooks or commands are paged,
    # git-style.
    def puts(*args)
      page_stdout
      super
    end

    # http://nex-3.com/posts/73-git-style-automatic-paging-in-ruby
    def page_stdout
      return if not $stdout.tty? or windows?

      read, write = IO.pipe

      if Kernel.fork
        # Parent process, become pager
        $stdin.reopen(read)
        read.close
        write.close

        # Don't page if the input is short enough
        ENV['LESS'] = 'FSRX'

        # Wait until we have input before we start the pager
        Kernel.select [STDIN]

        pager = ENV['GIT_PAGER'] ||
          `git config --get-all core.pager`.split.first || ENV['PAGER'] ||
          'less -isr'

        pager = 'cat' if pager.empty?

        exec pager rescue exec "/bin/sh", "-c", pager
      else
        # Child process
        $stdout.reopen(write)
        $stderr.reopen(write) if $stderr.tty?
        read.close
        write.close
      end
    rescue NotImplementedError
      # fork might not available, such as in JRuby
    end

    def pullrequest_editmsg(changes)
      message_file = File.join(git_dir, 'PULLREQ_EDITMSG')
      File.open(message_file, 'w') { |msg|
        yield msg
        if changes
          msg.puts "#\n# Changes:\n#"
          msg.puts changes.gsub(/^/, '# ').gsub(/ +$/, '')
        end
      }
      edit_cmd = Array(git_editor).dup
      edit_cmd << '-c' << 'set ft=gitcommit' if edit_cmd[0] =~ /^[mg]?vim$/
      edit_cmd << message_file
      system(*edit_cmd)
      abort "can't open text editor for pull request message" unless $?.success?
      title, body = read_editmsg(message_file)
      abort "Aborting due to empty pull request title" unless title
      [title, body]
    end

    def read_editmsg(file)
      title, body = '', ''
      File.open(file, 'r') { |msg|
        msg.each_line do |line|
          next if line.index('#') == 0
          ((body.empty? and line =~ /\S/) ? title : body) << line
        end
      }
      title.tr!("\n", ' ')
      title.strip!
      body.strip!

      [title =~ /\S/ ? title : nil, body =~ /\S/ ? body : nil]
    end

    def expand_alias(cmd)
      if expanded = git_alias_for(cmd)
        if expanded.index('!') != 0
          require 'shellwords' unless defined?(::Shellwords)
          Shellwords.shellwords(expanded)
        end
      end
    end
    
    def display_api_exception(action, response)
      $stderr.puts "Error #{action}: #{response.message.strip} (HTTP #{response.status})"
      if 422 == response.status and response.error_message?
        # display validation errors
        msg = response.error_message
        msg = msg.join("\n") if msg.respond_to? :join
        warn msg
      end
    end

  end
end
