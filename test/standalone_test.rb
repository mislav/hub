require 'helper'
require 'hub/standalone'
require 'fileutils'
require 'stringio'

class StandaloneTest < Minitest::Test
  include FileUtils

  def setup
    rm "hub" if File.exist? 'hub'
    rm_rf "/tmp/_hub_private" if File.exist? '/tmp/_hub_private'
    mkdir "/tmp/_hub_private"
    chmod 0400, "/tmp/_hub_private"
  end

  def teardown
    rm "hub" if File.exist? 'hub'
    rm_rf "/tmp/_hub_private" if File.exist? "/tmp/_hub_private"
  end

  def test_standalone
    io = StringIO.new
    Hub::Standalone.build io
    standalone = io.string

    assert_includes "This file is generated code", standalone
    assert_includes "Runner", standalone
    assert_includes "Args", standalone
    assert_includes "Commands", standalone
    assert_includes ".execute(*ARGV)", standalone
    assert_not_includes "module Standalone", standalone

    standalone =~ /__END__\s*(.+)/m
    assert_equal File.read('man/hub.1'), $1
  end

  def test_standalone_save
    Hub::Standalone.save("hub")
    output = `RUBYOPT= RUBYLIB= ./hub version 2>&1`
    assert_equal <<-OUT, output
git version 1.7.0.4
hub version #{Hub::VERSION}
    OUT
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
end
