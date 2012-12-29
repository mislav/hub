module Hub
  # The Args class exists to make it more convenient to work with
  # command line arguments intended for git from within the Hub
  # codebase.
  #
  # The ARGV array is converted into an Args instance by the Hub
  # instance when instantiated.
  class Args < Array
    attr_accessor :executable

    def initialize(*args)
      super
      @executable = ENV["GIT"] || "git"
      @skip = @noop = false
      @original_args = args.first
      @chain = [nil]
    end

    # Adds an `after` callback.
    # A callback can be a command or a proc.
    def after(cmd_or_args = nil, args = nil, &block)
      @chain.insert(-1, normalize_callback(cmd_or_args, args, block))
    end

    # Adds a `before` callback.
    # A callback can be a command or a proc.
    def before(cmd_or_args = nil, args = nil, &block)
      @chain.insert(@chain.index(nil), normalize_callback(cmd_or_args, args, block))
    end

    # Tells if there are multiple (chained) commands or not.
    def chained?
      @chain.size > 1
    end

    # Returns an array of all commands.
    def commands
      chain = @chain.dup
      chain[chain.index(nil)] = self.to_exec
      chain
    end

    # Skip running this command.
    def skip!
      @skip = true
    end

    # Boolean indicating whether this command will run.
    def skip?
      @skip
    end

    # Mark that this command shouldn't really run.
    def noop!
      @noop = true
    end

    def noop?
      @noop
    end

    # Array of `executable` followed by all args suitable as arguments
    # for `exec` or `system` calls.
    def to_exec(args = self)
      Array(executable) + args
    end

    def add_exec_flags(flags)
      self.executable = Array(executable).concat(flags)
    end

    # All the words (as opposed to flags) contained in this argument
    # list.
    #
    # args = Args.new([ 'remote', 'add', '-f', 'tekkub' ])
    # args.words == [ 'remote', 'add', 'tekkub' ]
    def words
      reject { |arg| arg.index('-') == 0 }
    end

    # All the flags (as opposed to words) contained in this argument
    # list.
    #
    # args = Args.new([ 'remote', 'add', '-f', 'tekkub' ])
    # args.flags == [ '-f' ]
    def flags
      self - words
    end

    # Tests if arguments were modified since instantiation
    def changed?
      chained? or self != @original_args
    end

    def has_flag?(*flags)
      pattern = flags.flatten.map { |f| Regexp.escape(f) }.join('|')
      !grep(/^#{pattern}(?:=|$)/).empty?
    end

    private

    def normalize_callback(cmd_or_args, args, block)
      if block
        block
      elsif args
        [cmd_or_args].concat args
      elsif Array === cmd_or_args
        self.to_exec cmd_or_args
      elsif cmd_or_args
        cmd_or_args
      else
        raise ArgumentError, "command or block required"
      end
    end
  end
end
