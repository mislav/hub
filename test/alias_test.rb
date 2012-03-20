require 'helper'

class AliasTest < Test::Unit::TestCase
  def test_alias_instructions
    expected = "# Wrap git automatically by adding the following to your profile:\n"
    expected << "\n"
    expected << 'eval "$(hub alias -s)"' << "\n"
    assert_equal expected, hub("alias sh")
  end

  def test_alias_instructions_bash
    with_shell('bash') do
      assert_includes '~/.bash_profile', hub("alias")
    end
  end

  def test_alias_instructions_zsh
    with_shell('zsh') do
      assert_includes '~/.zshrc', hub("alias")
    end
  end

  def test_alias_script_bash
    with_shell('bash') do
      assert_equal "alias git=hub\n", hub("alias -s")
    end
  end

  def test_alias_script_zsh
    with_shell('zsh') do
      script = hub("alias -s")
      assert_includes "alias git=hub\n", script
      assert_includes "compdef hub=git\n", script
    end
  end

  def test_unknown_shell
    with_shell(nil) do
      assert_equal "hub alias: unknown shell\n", hub("alias -s")
    end
  end

  def test_unsupported_shell
    with_shell('foosh') do
      expected = "hub alias: unsupported shell\n"
      expected << "supported shells: bash zsh sh ksh csh fish\n"
      assert_equal expected, hub("alias -s")
    end
  end

  private

  def with_shell(shell)
    old_shell, ENV['SHELL'] = ENV['SHELL'], shell
    yield
  ensure
    ENV['SHELL'] = old_shell
  end
end
