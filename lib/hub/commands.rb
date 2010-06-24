module Hub
  # See context.rb
  module Context; end

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
  # 2. The Hub class is initialized:
  #    >> hub = Hub.new('clone', 'rtomayko/tilt')
  #
  # 3. The method representing the git subcommand is executed with the
  #    full args:
  #    >> Commands.clone('clone', 'rtomayko/tilt')
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

    # Provides `github_url` and various inspection methods
    extend Context

    API_REPO = 'http://github.com/api/v2/yaml/repos/show/%s/%s'
    API_FORK = 'http://github.com/api/v2/yaml/repos/fork/%s/%s'

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

      last_args = args[1..-1].reject { |arg| arg == "--" }.last(3)
      last_args.each do |arg|
        if arg =~ /^-/
          # Skip mandatory arguments.
          last_args.shift if arg =~ /^(--(ref|o|br|u|t|d)[^=]+|-(o|b|u|d))$/
          next
        end

        if arg =~ %r{.+?://|.+?@} || File.directory?(arg)
          # Bail out early for URLs and local paths.
          break
        elsif arg.scan('/').size <= 1 && !arg.include?(':')
          # $ hub clone rtomayko/tilt
          # $ hub clone tilt
          args[args.index(arg)] = github_url(:repo => arg, :private => ssh)
          break
        end
      end
    end

    # $ hub submodule add wycats/bundler vendor/bundler
    # > git submodule add git://github.com/wycats/bundler.git vendor/bundler
    #
    # $ hub submodule add -p wycats/bundler vendor/bundler
    # > git submodule add git@github.com:wycats/bundler.git vendor/bundler
    #
    # $ hub submodule add -b ryppl ryppl/pip vendor/bundler
    # > git submodule add -b ryppl git://github.com/ryppl/pip.git vendor/pip
    def submodule(args)
      return unless index = args.index('add')
      args.delete_at index

      branch = args.index('-b') || args.index('--branch')
      if branch
        args.delete_at branch
        branch_name = args.delete_at branch
      end

      clone(args)

      if branch_name
        args.insert branch, '-b', branch_name
      end
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
      return if args[1] != 'add' || args.last =~ %r{.+?://|.+?@|^[./]}

      ssh = args.delete('-p')

      # user/repo
      args.last =~ /\b(.+?)(?:\/(.+))?$/
      user, repo = $1, $2

      if args.words[2] == 'origin' && args.words[3].nil?
        # Origin special case triggers default user/repo
        user = repo = nil
      elsif args.words[-2] == args.words[1]
        # rtomayko/tilt => rtomayko
        # Make sure you dance around flags.
        idx = args.index( args.words[-1] )
        args[idx] = user
      else
        # They're specifying the remote name manually (e.g.
        # git remote add blah rtomayko/tilt), so just drop the last
        # argument.
        args.replace args[0...-1]
      end

      args << github_url(:user => user, :repo => repo, :private => ssh)
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
      elsif name = args.words[1]
        # $ hub fetch <name1>,<name2>,...
        if name =~ /^\w+(,\w+)+$/
          index = args.index(name)
          args.delete(name)
          names = name.split(',')
          args.insert(index, *names)
          args.insert(index, '--multiple')
        else
          names = [name]
        end
      else
        names = []
      end

      names.reject! { |name|
        name =~ /\W/ or remotes.include?(name) or
          remotes_group(name) or not repo_exists?(name)
      }

      if names.any?
        commands = names.map { |name| "git remote add #{name} #{github_url(:user => name)}" }
        commands << args.to_exec.join(' ')
        args.replace commands.shift.split(' ')
        args.shift # don't want "git"
        args.after commands.join('; ')
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
        case ref = args.words.last
        when %r{^(https?:)//github.com/(.+?)/(.+?)/commit/([a-f1-9]{7,40})}
          scheme, user, repo, sha = $1, $2, $3, $4
          args[args.index(ref)] = sha
        when /^(\w+)@([a-f1-9]{7,40})$/
          scheme, user, repo, sha = nil, $1, nil, $2
          args[args.index(ref)] = sha
        else
          user = nil
        end

        if user
          # cherry-pick comes after the fetch
          args.after args.to_exec.join(' ')

          if user == repo_owner
            # fetch from origin if the repo belongs to the user
            args.replace ['fetch', default_remote]
          elsif remotes.include?(user)
            args.replace ['fetch', user]
          else
            secure = scheme == 'https:'
            remote_url = github_url(:user => user, :repo => repo, :private => secure)
            args.replace ['remote', 'add', '-f', user, remote_url]
          end
        end
      end
    end

    # $ hub init -g
    # > git init
    # > git remote add origin git@github.com:USER/REPO.git
    def init(args)
      if args.delete('-g')
        url = github_url(:private => true, :repo => File.basename(Dir.pwd))
        args.after "git remote add origin #{url}"
      end
    end

    # $ hub fork
    # ... hardcore forking action ...
    # > git remote add -f YOUR_USER git@github.com:YOUR_USER/CURRENT_REPO.git
    def fork(args)
      # can't do anything without token and original owner name
      if github_user && github_token && repo_owner
        if repo_exists?(github_user)
          puts "#{github_user}/#{repo_name} already exists on GitHub"
        else
          fork_repo
        end

        if args.include?('--no-remote')
          exit
        else
          url = github_url(:private => true)
          args.replace %W"remote add -f #{github_user} #{url}"
          args.after { puts "new remote: #{github_user}" }
        end
      end
    end

    # $ hub push origin,staging cool-feature
    # > git push origin cool-feature
    # > git push staging cool-feature
    def push(args)
      return unless args[1] =~ /,/

      branch  = args[2]
      remotes = args[1].split(',')
      args[1] = remotes.shift

      after = "git push #{remotes.shift} #{branch}"

      while remotes.length > 0
        after += "; git push #{remotes.shift} #{branch}"
      end

      args.after after
    end

    # $ hub browse
    # > open http://github.com/CURRENT_REPO
    #
    # $ hub browse -- issues
    # > open http://github.com/CURRENT_REPO/issues
    #
    # $ hub browse pjhyett/github-services
    # > open http://github.com/pjhyett/github-services
    #
    # $ hub browse -p pjhyett/github-fi
    # > open https://github.com/pjhyett/github-fi
    #
    # $ hub browse github-services
    # > open http://github.com/YOUR_LOGIN/github-services
    #
    # $ hub browse github-services wiki
    # > open http://wiki.github.com/YOUR_LOGIN/github-services
    #
    # $ hub browse -p github-fi
    # > open https://github.com/YOUR_LOGIN/github-fi
    def browse(args)
      args.shift
      browse_command(args) do
        user = repo = nil
        dest = args.shift
        dest = nil if dest == '--'

        if dest
          # $ hub browse pjhyett/github-services
          # $ hub browse github-services
          repo = dest
        elsif repo_user
          # $ hub browse
          user = repo_user
        else
          abort "Usage: hub browse [<USER>/]<REPOSITORY>"
        end

        params = { :user => user, :repo => repo }

        # $ hub browse -- wiki
        case subpage = args.shift
        when 'wiki'
          params[:web] = 'wiki'
        when 'commits'
          branch = (!dest && tracked_branch) || 'master'
          params[:web] = "/commits/#{branch}"
        when 'tree', NilClass
          branch = !dest && tracked_branch
          params[:web] = "/tree/#{branch}" if branch && branch != 'master'
        else
          params[:web] = "/#{subpage}"
        end

        params
      end
    end

    # $ hub compare 1.0...fix
    # > open http://github.com/CURRENT_REPO/compare/1.0...fix
    # $ hub compare refactor
    # > open http://github.com/CURRENT_REPO/compare/refactor
    # $ hub compare myfork feature
    # > open http://github.com/myfork/REPO/compare/feature
    # $ hub compare -p myfork topsecret
    # > open https://github.com/myfork/REPO/compare/topsecret
    # $ hub compare -u 1.0...2.0
    # prints "http://github.com/CURRENT_REPO/compare/1.0...2.0"
    def compare(args)
      args.shift
      browse_command(args) do
        if args.empty?
          branch = tracked_branch
          if branch && branch != 'master'
            range, user = branch, repo_user
          else
            abort "Usage: hub compare [USER] [<START>...]<END>"
          end
        else
          range = args.pop
          user = args.pop || repo_user
        end
        { :user => user, :web => "/compare/#{range}" }
      end
    end

    # $ hub hub standalone
    # Prints the "standalone" version of hub for an easy, memorable
    # installation sequence:
    #
    # $ gem install git-hub
    # $ hub hub standalone > ~/bin/hub && chmod 755 ~/bin/hub
    # $ gem uninstall git-hub
    def hub(args)
      return help(args) unless args[1] == 'standalone'
      require 'hub/standalone'
      puts Hub::Standalone.build
      exit
    rescue LoadError
      abort "hub is running in standalone mode."
    end

    def alias(args)
      shells = {
        'sh'   => 'alias git=hub',
        'bash' => 'alias git=hub',
        'zsh'  => 'function git(){hub "$@"}',
        'csh'  => 'alias git hub',
        'fish' => 'alias git hub'
      }

      silent = args.delete('-s')

      if shell = args[1]
        if silent.nil?
          puts "Run this in your shell to start using `hub` as `git`:"
          print "  "
        end
      else
        puts "usage: hub alias [-s] SHELL", ""
        puts "You already have hub installed and available in your PATH,"
        puts "but to get the full experience you'll want to alias it to"
        puts "`git`.", ""
        puts "To see how to accomplish this for your shell, run the alias"
        puts "command again with the name of your shell.", ""
        puts "Known shells:"
        shells.map { |key, _| key }.sort.each do |key|
          puts "  " + key
        end
        puts "", "Options:"
        puts "  -s   Silent. Useful when using the output with eval, e.g."
        puts "       $ eval `hub alias -s bash`"

        exit
      end

      if shells[shell]
        puts shells[shell]
      else
        abort "fatal: never heard of `#{shell}'"
      end

      exit
    end

    # $ hub version
    # > git version
    # (print hub version)
    def version(args)
      args.after do
        puts "hub version %s" % Version
      end
    end
    alias_method "--version", :version

    # $ hub help
    # (print improved help text)
    def help(args)
      command = args.grep(/^[^-]/)[1]

      if command == 'hub'
        puts hub_manpage
        exit
      elsif command.nil?
        ENV['GIT_PAGER'] = '' if args.grep(/^-{1,2}p/).empty? # Use `cat`.
        puts improved_help_text
        exit
      end
    end
    alias_method "--help", :help

    # The text print when `hub help` is run, kept in its own method
    # for the convenience of the author.
    def improved_help_text
      <<-help
usage: git [--version] [--exec-path[=GIT_EXEC_PATH]] [--html-path]
    [-p|--paginate|--no-pager] [--bare] [--git-dir=GIT_DIR]
    [--work-tree=GIT_WORK_TREE] [--help] COMMAND [ARGS]

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

Advanced commands:
   reset      Reset your staging area or working directory to another point
   rebase     Re-apply a series of patches in one branch onto another
   bisect     Find by binary search the change that introduced a bug
   grep       Print files with lines matching a pattern in your codebase

See 'git help COMMAND' for more information on a specific command.
help
    end

  private
    #
    # Helper methods are private so they cannot be invoked
    # from the command line.
    #

    # Checks whether a command exists on this system in the $PATH.
    #
    # name - The String name of the command to check for.
    #
    # Returns a Boolean.
    def command?(name)
      `type -t #{name}`
      $?.success?
    end

    # Detects commands to launch the user's browser, checking $BROWSER
    # first then falling back to a few common launchers. Aborts with
    # an error if it can't find anything appropriate.
    #
    # Returns a launch command.
    def browser_launcher
      if ENV['BROWSER']
        ENV['BROWSER']
      elsif RUBY_PLATFORM.include?('darwin')
        "open"
      elsif command?("xdg-open")
        "xdg-open"
      elsif command?("cygstart")
        "cygstart"
      else
        abort "Please set $BROWSER to a web launcher to use this command."
      end
    end

    # Handles common functionality of browser commands like `browse`
    # and `compare`. Yields a block that returns params for `github_url`.
    def browse_command(args)
      url_only = args.delete('-u')
      secure = args.delete('-p')
      params = yield

      args.executable = url_only ? 'echo' : browser_launcher
      args.push github_url({:web => true, :private => secure}.update(params))
    end


    # Returns the terminal-formatted manpage, ready to be printed to
    # the screen.
    def hub_manpage
      return "** Can't find groff(1)" unless command?('groff')

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
      return unless $stdout.tty?

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
    end

    # Determines whether a user has a fork of the current repo on GitHub.
    def repo_exists?(user)
      require 'net/http'
      url = API_REPO % [user, repo_name]
      Net::HTTPSuccess === Net::HTTP.get_response(URI(url))
    end

    # Forks the current repo using the GitHub API.
    #
    # Returns nothing.
    def fork_repo
      url = API_FORK % [repo_owner, repo_name]
      Net::HTTP.post_form(URI(url), 'login' => github_user, 'token' => github_token)
    end
  end
end
