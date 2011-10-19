module Hub
  # The Hub runner expects to be initialized with `ARGV` and primarily
  # exists to run a git command.
  #
  # The actual functionality, that is, the code it runs when it needs to
  # augment a git command, is kept in the `Hub::Commands` module.
  class Runner
    attr_reader :args
    
    def initialize(*args)

      # pre-process:
      #  go thru entire list of args searching for ["a","/","b"]
      #  and transforming them all by collapsing into ["a/b"]
      idx, result = 0, []
      while idx < args.length do
        current_item = args[idx]
        if idx <= args.length - 3
          next_item      = args[idx + 1]
          if next_item == "/"
            result.push args[idx..(idx+2)].join
            idx += 3
            next
          end
        end
        result.push current_item
        idx += 1
      end
      args = result

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
          cmd.map { |c| (c.index(' ') || c.empty?) ? "'#{c}'" : c }.join(' ')
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
      unless args.skip?
        if args.chained?
          execute_command_chain
        else
          exec(*args.to_exec)
        end
      end
    end

    # Runs multiple commands in succession; exits at first failure.
    def execute_command_chain
      commands = args.commands
      commands.each_with_index do |cmd, i|
        if cmd.respond_to?(:call) then cmd.call
        elsif i == commands.length - 1
          # last command in chain
          exec(*cmd)
        else
          exit($?.exitstatus) unless system(*cmd)
        end
      end
    end
  end
end
