require 'test/unit'

begin
  require 'redgreen'
rescue LoadError
end

$LOAD_PATH.unshift File.dirname(__FILE__) + '/../lib'
require 'hub'
require 'hub/standalone'

# We're looking for `open` in the tests.
ENV['BROWSER'] = 'open'

# Setup path with fake executables in case a test hits them
fakebin_dir = File.expand_path('../fakebin', __FILE__)
ENV['PATH'] = "#{fakebin_dir}:#{ENV['PATH']}"

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
  # STDOUT as a string. Pass it what you would normally pass `hub` on
  # the command line, e.g.
  #
  # shell: hub clone rtomayko/tilt
  #  test: hub("clone rtomayko/tilt")
  #
  # If a block is given it will be run in the child process before
  # execution begins. You can use this to monkeypatch or fudge the
  # environment before running hub.
  def hub(args, input = nil)
    parent_read, child_write = IO.pipe
    child_read, parent_write = IO.pipe if input

    fork do
      yield if block_given?
      $stdin.reopen(child_read) if input
      $stdout.reopen(child_write)
      $stderr.reopen(child_write)
      Hub(args).execute
    end
    
    if input
      parent_write.write input
      parent_write.close
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
    assert_equal expected, Hub(input).command, "$ git #{input}"
  end

  # Asserts that `hub` will show a specific alias command for a
  # specific shell.
  #
  # e.g.
  #  assert_alias_command "sh", "alias git=hub"
  #
  # Here we are saying that this:
  #   $ hub alias sh
  # Should display this:
  #   Run this in your shell to start using `hub` as `git`:
  #     alias git=hub
  def assert_alias_command(shell, command)
    expected = "Run this in your shell to start using `hub` as `git`:\n  %s\n"
    assert_equal(expected % command, hub("alias #{shell}"))
  end

  # Asserts that `haystack` includes `needle`.
  def assert_includes(needle, haystack)
    assert haystack.include?(needle),
      "expected #{needle.inspect} in #{haystack.inspect}"
  end

  # Asserts that `haystack` does not include `needle`.
  def assert_not_includes(needle, haystack)
    assert !haystack.include?(needle),
      "didn't expect #{needle.inspect} in #{haystack.inspect}"
  end
end
