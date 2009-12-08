require 'test/unit'
load File.dirname(__FILE__) + '/../bin/hub'

class HubTest < Test::Unit::TestCase

  #
  # Test helpers
  #

  def Hub(args)
    Hub.new(*args.split(' '))
  end

  def hub(args)
    parent_read, child_write = IO.pipe

    fork do
      $stdout.reopen(child_write)
      Hub(args).execute
    end

    child_write.close
    parent_read.read
  end

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
    input   = "remote add -g -p rtomayko"
    command = "git remote add rtomayko git@github.com:rtomayko/hub.git"
    assert_command input, command
  end

  def test_public_remote
    input   = "remote add -g rtomayko"
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
end
