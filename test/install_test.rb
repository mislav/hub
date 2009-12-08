$LOAD_PATH.unshift File.dirname(__FILE__)
require 'helper'

class InstallTest < Test::Unit::TestCase
  def test_install
    instructions = hub("install")
    assert_includes "bash", instructions
    assert_includes "sh", instructions
    assert_includes "csh", instructions
    assert_includes "zsh", instructions
    assert_includes "fish", instructions
  end

  def test_install_silent
    assert_equal "alias git=hub\n", hub("install -s bash")
  end

  def test_install_bash
    assert_install_command "bash", "alias git=hub"
  end

  def test_install_sh
    assert_install_command "sh", "alias git=hub"
  end

  def test_install_zsh
    assert_install_command "zsh", "alias git=hub"
  end

  def test_install_csh
    assert_install_command "csh", "alias git hub"
  end

  def test_install_fish
    assert_install_command "fish", "alias git hub"
  end

  def test_install_blah
    assert_install_command "blah", "fatal: never heard of `blah'"
  end
end
