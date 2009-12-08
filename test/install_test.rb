$LOAD_PATH.unshift File.dirname(__FILE__)
require 'helper'

class InstallTest < Test::Unit::TestCase
  def test_install_bash
    assert_equal "alias git=hub\n", hub("install bash")
  end

  def test_install_sh
    assert_equal "alias git=hub\n", hub("install sh")
  end

  def test_install_zsh
    assert_equal "alias git=hub\n", hub("install zsh")
  end

  def test_install_csh
    assert_equal "alias git hub\n", hub("install csh")
  end

  def test_install_fish
    assert_equal "alias git hub\n", hub("install fish")
  end

  def test_install_blah
    assert_equal "fatal: never heard of `blah'\n", hub("install blah")
  end
end
