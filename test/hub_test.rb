require 'test/unit'
load File.dirname(__FILE__) + '/../bin/hub'

class HubTest < Test::Unit::TestCase
  # Shortcut for creating a `Hub` instance. Pass it what you would
  # normally pass `hub` on the command line, e.g.
  #
  # shell: hub clone rtomayko/tilt
  #  test: Hub("clone rtomayko/tilt")
  def Hub(args)
    Hub.new(*args.split(' '))
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

  #
  # Assertions
  #

  def test_private_clone
    input   = "clone -p rtomayko/ron"
    command = "git clone git@github.com:rtomayko/ron.git"
    assert_command input, command
  end

  def test_public_clone
    input   = "clone rtomayko/ron"
    command = "git clone git://github.com/rtomayko/ron.git"
    assert_command input, command
  end

  def test_private_remote
    input   = "remote add -p rtomayko"
    command = "git remote add rtomayko git@github.com:rtomayko/hub.git"
    assert_command input, command
  end

  def test_public_remote
    input   = "remote add rtomayko"
    command = "git remote add rtomayko git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_init
    h = Hub("init -g")
    assert_equal "git init", h.command
    assert_equal "git remote add origin git@github.com:defunkt/hub.git", h.after
  end

  def test_version
    h = hub("--version")
    assert_equal "git version 1.6.4.2\nhub version 0.1.0\n", h
  end

  def test_help
    assert_equal Hub::Commands.improved_help_text, hub("help")
  end

  def test_help_by_default
    assert_equal Hub::Commands.improved_help_text, hub("")
  end
end
