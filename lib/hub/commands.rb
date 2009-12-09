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
    PRIVATE = 'git@github.com:%s/%s.git'
    PUBLIC  = 'git://github.com/%s/%s.git'
    USER    = `git config --global github.user`.chomp
    REPO    = `basename $(pwd)`.chomp

    # $ hub clone rtomayko/tilt
    # > git clone git://github.com/rtomayko/tilt.
    #
    # $ hub clone -p kneath/hemingway
    # > git clone git@github.com:kneath/hemingway.git
    def clone(args)
      ssh = args.delete('-p')
      args.each_with_index do |arg, i|
        if arg.scan('/').size == 1 && !arg.include?(':')
          url = ssh ? PRIVATE : PUBLIC
          args[i] = url % arg.split('/')
        end
      end
    end

    # $ hub remote add pjhyett
    # > git remote add pjhyett git://github.com/pjhyett/THIS_REPO.git
    #
    # $ hub remote add -p mojombo
    # > git remote add mojombo git@github.com:mojombo/THIS_REPO.git
    def remote(args)
      return unless args[1] == 'add'

      # Assume GitHub usernames don't ever contain : or /, while URLs
      # do.
      if args[-1] !~ /:\//
        ssh  = args.delete('-p')
        user = args.last
        url  = ssh ? PRIVATE : PUBLIC
        args << url % [ user, REPO ]
      end
    end

    # $ hub init -g
    # > git init
    # > git remote add origin git@github.com:USER/REPO.git
    def init(args)
      if args.delete('-g')
        url = PRIVATE % [ USER, REPO ]
        args.after "git remote add origin #{url}"
      end
    end

    def alias(args)
      shells = {
        'sh'   => 'alias git=hub',
        'bash' => 'alias git=hub',
        'zsh'  => 'alias git=hub',
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
      return if args.size > 1
      puts improved_help_text
      exit
    end

    # The text print when `hub help` is run, kept in its own method
    # for the convenience of the author.
    def improved_help_text
      <<-help
usage: git [--version] [--exec-path[=GIT_EXEC_PATH]] [--html-path]
    [-p|--paginate|--no-pager] [--bare] [--git-dir=GIT_DIR]
    [--work-tree=GIT_WORK_TREE] [--help] COMMAND [ARGS]

Creating a git repository:
   clone      Clone a repository into a new directory
   init       Create an empty git repository or reinitialize an existing one

Working with content:
   add        Add file contents to the index
   branch     List, create, or delete branches
   checkout   Checkout a branch or paths to the working tree
   commit     Record changes to the repository
   diff       Show changes between commits, commit and working tree, etc
   log        Show commit logs
   merge      Join two or more development histories together
   mv         Move or rename a file, a directory, or a symlink
   rm         Remove files from the working tree and from the index
   status     Show the working tree status
   show       Show various types of objects
   tag        Create, list, delete or verify a tag object signed with GPG

Over the network:
   fetch      Download objects and refs from another repository
   pull       Fetch from and merge with another repository or a local branch
   push       Update remote refs along with associated objects
   remote     Manage a set of tracked repositories

Advanced commands:
   bisect     Find by binary search the change that introduced a bug
   grep       Print lines matching a pattern
   reset      Reset current HEAD to the specified state
   rebase     Forward-port local commits to the updated upstream head

See 'git help COMMAND' for more information on a specific command.
help
    end

    # $ hub install
    # $ hub install check
    # $ hub install standalone
    # $ hub install standalone ~/bin
    def install(args)
      command, subcommand, target = args

      if subcommand.to_s == 'standalone'
        Standalone.save('hub', target.empty? ? '.' : target)
      elsif subcommand.to_s == 'check'
        begin
          raise 'blah'
          if up_to_date?
            puts "*".green + " hub is up to date"
          else
            puts "*".red + " hub is " + "not".bold.underline + " up to date"
          end
        rescue Object => e
          puts "*".bold.yellow + " error checking status: #{e.class}"
        end
      else
        puts <<-output
usage: hub install COMMAND [ARGS]

Commands:
  standalone    Installs the standalone `hub` script locally. If
                a path is provided, attempts to install it there.
                If not path is provided asks you to choose from
                possible install locations.

  check         Checks if the current installation is up to date
                by phoning home.
output
      end

      exit
    end

  private
    #
    # Helper methods are private so they cannot be invoked
    # from the command line
    #

    # Is the current running version of hub up to date?
    def up_to_date?
      latest_md5 == current_md5
    end

    # Grab the latest standalone's md5 from the web
    def latest_md5
      require 'open-uri'
      md5_url = "http://defunkt.github.com/hub/standalone.md5"
      md5 = open(md5_url).read.chomp
    end

    # Compute the current standalone's md5
    def current_md5
      require 'digest/md5'
      hub = defined?(Standalone) ? Standalone.build : File.read(__FILE__)
      Digest::MD5.hexdigest(hub)
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

        pager = ENV['PAGER'] || 'less'
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
