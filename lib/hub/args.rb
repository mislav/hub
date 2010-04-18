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
      @after = nil
    end

    # With no arguments, returns the `after` callback.
    #
    # With a single argument, sets the `after` callback.
    # Can be set to a string or a proc.
    #
    # If proc:
    #   The proc is executed after the git command is executed. For
    #   example, the `hub version` command sets the following proc to
    #   print its information after running `git version`:
    #
    #     after { puts "hub version #{version_number}" }
    #
    # If string:
    #   The string is assumed to be a command and executed after the
    #   git command is executed:
    #
    #     after "echo 'hub version #{version_number}'"
    def after(command = nil, &block)
      @after ||= block ? block : command
    end

    # Boolean indicating whether an `after` callback has been set.
    def after?
      !!@after
    end

    # Array of `executable` followed by all args suitable as arguments
    # for `exec` or `system` calls.
    def to_exec
      [executable].concat self
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
  end
end
