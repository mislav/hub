module Hub
  # The Hub runner expects to be initialized with `ARGV` and primarily
  # exists to run a git command.
  #
  # The actual functionality, that is, the code it runs when it needs to
  # augment a git command, is kept in the `Hub::Commands` module.
  class Runner
    attr_reader :args
    
    def initialize(*args)
      @args = Args.new(args)

      # Hack to emulate git-style
      @args.unshift 'help' if @args.grep(/^[^-]|version/).empty?

      # git commands can have dashes
      cmd = @args[0].sub(/(\w)-/, '\1_')
      Commands.send(cmd, @args) if Commands.respond_to?(cmd)
    end

    # Shortcut
    def self.execute(*args)
      new(*args).execute
    end

    # Returns the current after callback, which (if set) is run after
    # the target git command.
    #
    # See the `Hub::Args` class for more information on the `after`
    # callback.
    def after
      args.after.to_s
    end

    # A string representation of the git command we would run if
    # #execute were called.
    def command
      args.to_exec.join(' ')
    end

    # Runs the target git command with an optional callback. Replaces
    # the current process.
    def execute
      if args.after?
        execute_with_after_callback
      else
        exec(*args.to_exec)
      end
    end

    # Runs the target git command then executes the `after` callback.
    #
    # See the `Hub::Args` class for more information on the `after`
    # callback.
    def execute_with_after_callback
      after = args.after
      if system(*args.to_exec)
        after.respond_to?(:call) ? after.call : exec(after)
        exit
      else
        exit 1
      end
    end
  end
end
