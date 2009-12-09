$LOAD_PATH.unshift File.dirname(__FILE__)
require 'helper'

class HubTest < Test::Unit::TestCase
  def setup
    Hub::Commands::USER.replace("tpw")
  end

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

  def test_your_private_clone
    input   = "clone -p resque"
    command = "git clone git@github.com:tpw/resque.git"
    assert_command input, command
  end

  def test_your_public_clone
    input   = "clone resque"
    command = "git clone git://github.com/tpw/resque.git"
    assert_command input, command
  end

  def test_your_private_clone_fails_without_config
    out = hub("clone -p mustache") do
      Hub::Commands::USER.replace("")
    end

    assert_equal "** No GitHub user set. See http://github.com/guides/local-github-config\n", out
  end

  def test_your_public_clone_fails_without_config
    out = hub("clone mustache") do
      Hub::Commands::USER.replace("")
    end

    assert_equal "** No GitHub user set. See http://github.com/guides/local-github-config\n", out
  end

  def test_private_clone_left_alone
    input   = "clone git@github.com:rtomayko/ron.git"
    command = "git clone git@github.com:rtomayko/ron.git"
    assert_command input, command
  end

  def test_public_clone_left_alone
    input   = "clone git://github.com/rtomayko/ron.git"
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
    assert_equal "git remote add origin git@github.com:tpw/hub.git", h.after
  end

  def test_init_no_login
    out = hub("init -g") do
      Hub::Commands::USER.replace("")
    end

    assert_equal "** No GitHub user set. See http://github.com/guides/local-github-config\n", out
  end

  def test_version
    out = hub('--version')
    assert_includes "git version 1.6", out
    assert_includes "hub version 0.1.0", out
  end

  def test_help
    assert_equal Hub::Commands.improved_help_text, hub("help")
  end

  def test_help_by_default
    assert_equal Hub::Commands.improved_help_text, hub("")
  end
end
