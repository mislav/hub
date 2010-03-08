$LOAD_PATH.unshift File.dirname(__FILE__)
require 'helper'

class AliasTest < Test::Unit::TestCase
  def test_alias
    instructions = hub("alias")
    assert_includes "bash", instructions
    assert_includes "sh", instructions
    assert_includes "csh", instructions
    assert_includes "zsh", instructions
    assert_includes "fish", instructions
  end

  def test_alias_silent
    assert_equal "alias git=hub\n", hub("alias -s bash")
  end

  def test_alias_bash
    assert_alias_command "bash", "alias git=hub"
  end

  def test_alias_sh
    assert_alias_command "sh", "alias git=hub"
  end

  def test_alias_zsh
    assert_alias_command "zsh", 'function git(){hub "$@"}'
  end

  def test_alias_csh
    assert_alias_command "csh", "alias git hub"
  end

  def test_alias_fish
    assert_alias_command "fish", "alias git hub"
  end

  def test_alias_blah
    assert_alias_command "blah", "fatal: never heard of `blah'"
  end
end
