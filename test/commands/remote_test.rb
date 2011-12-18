require 'test_helper'

class RemoteTest < Test::Unit::TestCase
  def test_remote_origin
    input   = "remote add origin"
    command = "git remote add origin git://github.com/tpw/hub.git"
    assert_command input, command
  end

  def test_remote_add_with_name
    input   = "remote add another hookio/hub.js"
    command = "git remote add another git://github.com/hookio/hub.js.git"
    assert_command input, command
  end

  def test_private_remote_origin
    input   = "remote add -p origin"
    command = "git remote add origin git@github.com:tpw/hub.git"
    assert_command input, command
  end

  def test_remote_with_host_alias
    assert_forwarded "remote add origin server:/git/repo.git"
  end

  def test_public_remote_origin_as_normal
    input   = "remote add origin http://github.com/defunkt/resque.git"
    command = "git remote add origin http://github.com/defunkt/resque.git"
    assert_command input, command
  end

  def test_remote_from_rel_path
    assert_forwarded "remote add origin ./path"
  end

  def test_remote_from_abs_path
    assert_forwarded "remote add origin /path"
  end

  def test_private_remote_origin_as_normal
    assert_forwarded "remote add origin git@github.com:defunkt/resque.git"
  end

  def test_private_remote
    input   = "remote add -p rtomayko"
    command = "git remote add rtomayko git@github.com:rtomayko/hub.git"
    assert_command input, command
  end

  def test_https_protocol_remote
    stub_https_is_preferred
    input   = "remote add rtomayko"
    command = "git remote add rtomayko https://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_public_remote
    input   = "remote add rtomayko"
    command = "git remote add rtomayko git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_public_remote_f
    input   = "remote add -f rtomayko"
    command = "git remote add -f rtomayko git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_named_public_remote
    input   = "remote add origin rtomayko"
    command = "git remote add origin git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_named_public_remote_f
    input   = "remote add -f origin rtomayko"
    command = "git remote add -f origin git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_private_remote_with_repo
    input   = "remote add -p jashkenas/coffee-script"
    command = "git remote add jashkenas git@github.com:jashkenas/coffee-script.git"
    assert_command input, command
  end

  def test_public_remote_with_repo
    input   = "remote add jashkenas/coffee-script"
    command = "git remote add jashkenas git://github.com/jashkenas/coffee-script.git"
    assert_command input, command
  end

  def test_public_remote_f_with_repo
    input   = "remote add -f jashkenas/coffee-script"
    command = "git remote add -f jashkenas git://github.com/jashkenas/coffee-script.git"
    assert_command input, command
  end

  def test_named_private_remote_with_repo
    input   = "remote add -p origin jashkenas/coffee-script"
    command = "git remote add origin git@github.com:jashkenas/coffee-script.git"
    assert_command input, command
  end
end
