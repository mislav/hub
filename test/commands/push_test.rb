require 'test_helper'

class PushTest < Test::Unit::TestCase
  def test_push_untouched
    assert_forwarded "push"
  end

  def test_push_two
    assert_commands "git push origin cool-feature", "git push staging cool-feature",
                    "push origin,staging cool-feature"
  end

  def test_push_current_branch
    stub_branch('refs/heads/cool-feature')
    assert_commands "git push origin cool-feature", "git push staging cool-feature",
                    "push origin,staging"
  end

  def test_push_more
    assert_commands "git push origin cool-feature",
                    "git push staging cool-feature",
                    "git push qa cool-feature",
                    "push origin,staging,qa cool-feature"
  end
end
