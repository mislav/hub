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
end
