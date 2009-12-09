$LOAD_PATH.unshift File.dirname(__FILE__)
require 'helper'
require 'fileutils'

class InstallTest < Test::Unit::TestCase
  include FileUtils

  def setup
    rm "hub" if File.exists? 'hub'
    rm_rf "/tmp/_hub_private" if File.exists? '/tmp/_hub_private'
    mkdir "/tmp/_hub_private"
    chmod 0400, "/tmp/_hub_private"
  end

  def teardown
    rm "hub" if File.exists? 'hub'
    rm_rf "/tmp/_hub_private" if File.exists? "/tmp/_hub_private"
  end

  def test_standalone
    standalone = Hub::Standalone.build
    assert_includes "This file, hub, is generated code", standalone
    assert_includes "Runner", standalone
    assert_includes "Args", standalone
    assert_includes "Commands", standalone
    assert_includes ".execute(*ARGV)", standalone
    assert_not_includes "module Standalone", standalone
  end

  def test_standalone_save
    Hub::Standalone.save("hub")
    assert_equal Hub::Standalone.build + "\n", File.read('./hub')
  end

  def test_standalone_save_permission_denied
    assert_raises Errno::EACCES do
      Hub::Standalone.save("hub", "/tmp/_hub_private")
    end
  end

  def test_standalone_save_doesnt_exist
    assert_raises Errno::ENOENT do
      Hub::Standalone.save("hub", "/tmp/something/not/real")
    end
  end

  def test_install
    out = hub("install")
    assert_includes "usage: hub", out
    assert_includes "check", out
    assert_includes "update", out
  end

  def test_install_check_up_to_date
    fake_up_to_date do
      assert_equal "* hub is up to date\n", hub("install check")
    end
  end

  def test_install_check_not_up_to_date
    assert_equal "* hub is not up to date\n", hub("install check")
  end

  def test_update_nothing_to_do
    fake_up_to_date do
      out = hub("install update")
      assert_equal "* hub is up to date\n", out
    end
  end

  def test_update_standalone
    out = hub("install update")
    assert_equal "0.1.1", out
  end

  def test_update_rubygems
    out = hub("install update")
    assert_equal "rubygems", out
  end
end
