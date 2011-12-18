require 'test_helper'

class CompareTest < Test::Unit::TestCase
  def test_hub_compare
    assert_command "compare refactor",
      "open https://github.com/defunkt/hub/compare/refactor"
  end

  def test_hub_compare_nothing
    expected = "Usage: hub compare [USER] [<START>...]<END>\n"
    assert_equal expected, hub("compare")
  end

  def test_hub_compare_tracking_nothing
    stub_tracking_nothing
    expected = "Usage: hub compare [USER] [<START>...]<END>\n"
    assert_equal expected, hub("compare")
  end

  def test_hub_compare_tracking_branch
    stub_branch('refs/heads/feature')
    stub_tracking('feature', 'mislav', 'experimental')

    assert_command "compare",
      "open https://github.com/mislav/hub/compare/experimental"
  end

  def test_hub_compare_range
    assert_command "compare 1.0...fix",
      "open https://github.com/defunkt/hub/compare/1.0...fix"
  end

  def test_hub_compare_range_fixes_two_dots_for_tags
    assert_command "compare 1.0..fix",
      "open https://github.com/defunkt/hub/compare/1.0...fix"
  end

  def test_hub_compare_range_fixes_two_dots_for_shas
    assert_command "compare 1234abc..3456cde",
      "open https://github.com/defunkt/hub/compare/1234abc...3456cde"
  end

  def test_hub_compare_range_ignores_two_dots_for_complex_ranges
    assert_command "compare @{a..b}..@{c..d}",
      "open https://github.com/defunkt/hub/compare/@{a..b}..@{c..d}"
  end

  def test_hub_compare_on_wiki
    stub_repo_url 'git://github.com/defunkt/hub.wiki.git'
    assert_command "compare 1.0...fix",
      "open https://github.com/defunkt/hub/wiki/_compare/1.0...fix"
  end

  def test_hub_compare_fork
    assert_command "compare myfork feature",
      "open https://github.com/myfork/hub/compare/feature"
  end

  def test_hub_compare_url
    assert_command "compare -u 1.0...1.1",
      "echo https://github.com/defunkt/hub/compare/1.0...1.1"
  end
end
