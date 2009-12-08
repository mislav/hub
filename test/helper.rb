require 'test/unit'
$LOAD_PATH.unshift File.dirname(__FILE__) + '/../lib'
require 'hub'

class Test::Unit::TestCase
  # Shortcut for creating a `Hub` instance. Pass it what you would
  # normally pass `hub` on the command line, e.g.
  #
  # shell: hub clone rtomayko/tilt
  #  test: Hub("clone rtomayko/tilt")
  def Hub(args)
    Hub::Runner.new(*args.split(' '))
  end

  # Shortcut for running the `hub` command in a subprocess. Returns
  # STDOUT as a string. Pass it what you would normall pass `hub` on
  # the command line, e.g.
  #
  # shell: hub clone rtomayko/tilt
  #  test: hub("clone rtomayko/tilt")
  def hub(args)
    parent_read, child_write = IO.pipe

    fork do
      $stdout.reopen(child_write)
      $stderr.reopen(child_write)
      Hub(args).execute
    end

    child_write.close
    parent_read.read
  end

  # Asserts that `hub` will run a specific git command based on
  # certain input.
  #
  # e.g.
  #  assert_command "clone git/hub", "git clone git://github.com/git/hub.git"
  #
  # Here we are saying that this:
  #   $ hub clone git/hub
  # Should in turn execute this:
  #   $ git clone git://github.com/git/hub.git
  def assert_command(input, expected)
    assert_equal expected, Hub(input).command
  end

  # Asserts that `hub` will show a specific installation command for a
  # specific shell.
  #
  # e.g.
  #  assert_install_command "sh", "alias git=hub"
  #
  # Here we are saying that this:
  #   $ hub install sh
  # Should display this:
  #   Run this in your shell to start using `hub` as `git`:
  #     alias git=hub
  def assert_install_command(shell, command)
    expected = "Run this in your shell to start using `hub` as `git`:\n  %s\n"
    assert_equal(expected % command, hub("install #{shell}"))
  end

  # Asserts that `haystack` includes `needle`.
  def assert_includes(needle, haystack)
    assert haystack.include?(needle)
  end
end
