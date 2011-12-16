require 'test_helper'

class InitTest < Test::Unit::TestCase
  def test_init
    assert_commands "git init", "git remote add origin git@github.com:tpw/hub.git", "init -g"
  end

  def test_init_no_login
    out = hub("init -g") do
      stub_github_user(nil)
    end

    assert_equal "** No GitHub user set. See http://help.github.com/git-email-settings/\n", out
  end
end
