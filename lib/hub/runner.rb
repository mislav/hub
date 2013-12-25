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
      Commands.run(@args)
    end

    # Shortcut
    def self.execute(*args)
      new(*args).execute
    end

    # A string representation of the command that would run.
    def command
      if args.skip?
        ''
      else
        commands.join('; ')
      end
    end

    # An array of all commands as strings.
    def commands
      args.commands.map do |cmd|
        if cmd.respond_to?(:join)
          # a simplified `Shellwords.join` but it's OK since this is only used to inspect
          cmd.map { |arg| arg = arg.to_s; (arg.index(' ') || arg.empty?) ? "'#{arg}'" : arg }.join(' ')
        else
          cmd.to_s
        end
      end
    end

    # Runs the target git command with an optional callback. Replaces
    # the current process. 
    # 
    # If `args` is empty, this will skip calling the git command. This
    # allows commands to print an error message and cancel their own
    # execution if they don't make sense.
    def execute
      if args.noop?
        puts commands
      elsif not args.skip?
        execute_command_chain args.commands
      end
    end

    # Runs multiple commands in succession; exits at first failure.
    def execute_command_chain commands
      commands.each_with_index do |cmd, i|
        if cmd.respond_to?(:call) then cmd.call
        elsif i == commands.length - 1
          # last command in chain
          STDOUT.flush; STDERR.flush
          exec(*cmd)
        else
          exit($?.exitstatus) unless system(*cmd)
        end
      end
    end

    # Special-case `echo` for Windows
    def exec *args
      if args.first == 'echo' && Context::windows?
        puts args[1..-1].join(' ')
      else
        super
      end
    end
  end
end
