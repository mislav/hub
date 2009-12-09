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
    hub("install standalone .")
    assert_equal Hub::Standalone.build + "\n", File.read('./hub')
  end

  def test_standalone_save_permission_denied
    out = hub("install standalone /tmp/_hub_private")
    assert_equal "** can't write to /tmp/_hub_private/hub\n", out
  end

  def test_standalone_save_doesnt_exist
    out = hub("install standalone /tmp/something/not/real")
    assert_equal "** can't write to /tmp/something/not/real/hub\n", out
  end

  def test_install
    out = hub("install")
    assert_includes "usage: hub", out
    assert_includes "check", out
    assert_includes "standalone", out
  end

  def test_install_check_up_to_date
    Hub::Commands.class_eval do
      alias_method :real_latest_md5, :latest_md5
      alias_method :latest_md5, :current_md5
    end

    assert_equal "* hub is up to date\n", hub("install check")
  end

  def test_install_check_not_up_to_date
    if Hub::Commands.instance_methods.include? 'real_latest_md5'
      Hub::Commands.class_eval do
        alias_method :latest_md5, :real_latest_md5
      end
    end

    assert_equal "* hub is not up to date\n", hub("install check")
  end
end
