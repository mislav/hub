$LOAD_PATH.unshift File.dirname(__FILE__)
require 'helper'

class StandaloneTest < Test::Unit::TestCase
  def test_standalone
    standalone = Hub::Standalone.build
    assert_includes "This file, hub, is generated code", standalone
    assert_includes "Runner", standalone
    assert_includes "Args", standalone
    assert_includes "Commands", standalone
    assert_includes ".execute(*ARGV)", standalone
    assert_not_includes "module Standalone", standalone
  end
end
