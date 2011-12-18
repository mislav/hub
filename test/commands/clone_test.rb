require 'test_helper'

class CloneTest < Test::Unit::TestCase
  def test_private_clone
    input   = "clone -p rtomayko/ronn"
    command = "git clone git@github.com:rtomayko/ronn.git"
    assert_command input, command
  end

  def test_private_clone_noop
    input   = "--noop clone -p rtomayko/ronn"
    command = "git clone git@github.com:rtomayko/ronn.git\n"
    assert_output command, hub(input)
  end

  def test_https_clone
    stub_https_is_preferred
    input   = "clone rtomayko/ronn"
    command = "git clone https://github.com/rtomayko/ronn.git"
    assert_command input, command
  end

  def test_public_clone
    input   = "clone rtomayko/ronn"
    command = "git clone git://github.com/rtomayko/ronn.git"
    assert_command input, command
  end

  def test_your_private_clone
    input   = "clone -p resque"
    command = "git clone git@github.com:tpw/resque.git"
    assert_command input, command
  end

  def test_your_clone_is_always_private
    input   = "clone resque"
    command = "git clone git@github.com:tpw/resque.git"
    assert_command input, command
  end

  def test_clone_repo_with_period
    input   = "clone hookio/hook.js"
    command = "git clone git://github.com/hookio/hook.js.git"
    assert_command input, command
  end

  def test_clone_with_arguments
    input   = "clone --bare -o master resque"
    command = "git clone --bare -o master git@github.com:tpw/resque.git"
    assert_command input, command
  end

  def test_clone_with_arguments_and_destination
    assert_forwarded "clone --template=one/two git://github.com/tpw/resque.git --origin master resquetastic"
  end

  def test_your_private_clone_fails_without_config
    out = hub("clone -p mustache") do
      stub_github_user(nil)
    end

    assert_equal "** No GitHub user set. See http://help.github.com/set-your-user-name-email-and-github-token/\n", out
  end

  def test_your_public_clone_fails_without_config
    out = hub("clone mustache") do
      stub_github_user(nil)
    end

    assert_equal "** No GitHub user set. See http://help.github.com/set-your-user-name-email-and-github-token/\n", out
  end

  def test_private_clone_left_alone
    assert_forwarded "clone git@github.com:rtomayko/ronn.git"
  end

  def test_public_clone_left_alone
    assert_forwarded "clone git://github.com/rtomayko/ronn.git"
  end

  def test_normal_public_clone_with_path
    assert_forwarded "clone git://github.com/rtomayko/ronn.git ronn-dev"
  end

  def test_normal_clone_from_path
    assert_forwarded "clone ./test"
  end

  def test_clone_with_host_alias
    assert_forwarded "clone server:git/repo.git"
  end

  def test_alias_expand
    stub_alias 'c', 'clone --bare'
    input   = "c rtomayko/ronn"
    command = "git clone --bare git://github.com/rtomayko/ronn.git"
    assert_command input, command
  end

  def test_alias_expand_advanced
    stub_alias 'c', 'clone --template="white space"'
    input   = "c rtomayko/ronn"
    command = "git clone '--template=white space' git://github.com/rtomayko/ronn.git"
    assert_command input, command
  end

  def test_alias_doesnt_expand_for_unknown_commands
    stub_alias 'c', 'compute --fast'
    assert_forwarded "c rtomayko/ronn"
  end
end
