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

    # Templates and useful information.
    USER       = `git config --global github.user`.chomp
    ORIGIN     = `git config remote.origin.url`.chomp
    HTTP_CLONE = `git config --global hub.http-clone`.chomp == 'yes'
    PUBLIC     = (HTTP_CLONE ? 'http' : 'git') + '://github.com/%s/%s.git'
    PRIVATE    = 'git@github.com:%s/%s.git'
    LGHCONF    = "http://github.com/guides/local-github-config"

    # Set the repo name based on the current origin or, as a fallback,
    # the cwd.
    if ORIGIN =~ %r{\bgithub\.com[:/](.+)/(.+).git$}
      OWNER, REPO = $1, $2
    else
      REPO = File.basename(Dir.pwd)
      OWNER = ''
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
        elsif arg.scan('/').size == 1 && !arg.include?(':')
          # $ hub clone rtomayko/tilt
          url = ssh ? PRIVATE : PUBLIC
          args[args.index(arg)] = url % arg.split('/')
          break
        elsif arg !~ /:|\//
          # $ hub clone tilt
          url = ssh ? PRIVATE : PUBLIC
          args[args.index(arg)] = url % [ github_user, arg ]
          break
        end
      end
    end

    # $ hub submodule add wycats/bundler vendor/bundler
    # > git submodule add git://github.com/wycats/bundler.git vendor/bundler
    #
    # $ hub submodule add -p wycats/bundler vendor/bundler
    # > git submodule add git@github.com:wycats/bundler.git vendor/bundler
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
      return unless args[1] == 'add'

      ssh = args.delete('-p')
      url = ssh ? PRIVATE : PUBLIC

      if args.last =~ /\b(\w+)\/(\w+)/
        # user/repo
        user, repo = $1, $2

        if args.words[-2] == args.words[1]
          # rtomayko/tilt => rtomayko
          args[-1] = user
        else
          # They're specifying the remote name manually (e.g.
          # git remote add blah rtomayko/tilt), so just drop the last
          # argument.
          args.replace args[0...-1]
        end

        args << url % [ user, repo ]
      elsif args.last !~ /:|\//
        if args[2] == 'origin' && args[3].nil?
          # Origin special case.
          user = github_user
        else
          # Assume no : or / means GitHub user.
          user = args.last
        end

        if args[-2] != args[1]
          # They're specifying the remote name manually (e.g.
          # git remote add blah rtomayko), so just drop the last
          # argument.
          args.replace args[0...-1]
        end

        args << url % [ user, REPO ]
      end
    end

    # $ hub init -g
    # > git init
    # > git remote add origin git@github.com:USER/REPO.git
    def init(args)
      if args.delete('-g')
        # Can't do anything if we don't have a USER set.

        url = PRIVATE % [ github_user, REPO ]
        args.after "git remote add origin #{url}"
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
    # $ hub browse pjhyett/github-services
    # > open http://github.com/pjhyett/github-services
    #
    # $ hub browse -p pjhyett/github-fi
    # > open https://github.com/pjhyett/github-fi
    #
    # $ hub browse github-services
    # > open http://github.com/YOUR_LOGIN/github-services
    #
    # $ hub browse -p github-fi
    # > open https://github.com/YOUR_LOGIN/github-fi
    def browse(args)
      args.shift
      protocol = args.delete('-p') ? 'https' : 'http'
      dest = args.pop

      if dest
        if dest.include? '/'
          # $ hub browse pjhyett/github-services
          user, repo = dest.split('/')
        else
          # $ hub browse github-services
          user, repo = github_user, dest
        end
      elsif !OWNER.empty?
        # $ hub browse
        user, repo = OWNER, REPO
      else
        warn "Usage: hub browse [<USER>/]<REPOSITORY>"
        exit(1)
      end

      args.executable = ENV['BROWSER'] || 'open'
      args.push "#{protocol}://github.com/#{user}/#{repo}"
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
        'zsh'  => 'function git(){hub $@}',
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

    # Either returns the GitHub user as set by git-config(1) or aborts
    # with an error message.
    def github_user
      if USER.empty?
        abort "** No GitHub user set. See #{LGHCONF}"
      else
        USER
      end
    end

    # Returns the terminal-formatted manpage, ready to be printed to
    # the screen.
    def hub_manpage
      return "** Can't find groff(1)" unless groff?

      require 'open3'
      out = nil
      Open3.popen3(groff_command) do |stdin, stdout, _|
        stdin.puts hub_raw_manpage
        stdin.close
        out = stdout.read.strip
      end
      out
    end

    # Returns true if groff is installed and in our path, false if
    # not.
    def groff?
      system("which groff")
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
  end
end
