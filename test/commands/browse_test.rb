require 'test_helper'

class BrowseTest < Test::Unit::TestCase
  def test_hub_browse
    assert_command "browse mojombo/bert", "open https://github.com/mojombo/bert"
  end

  def test_hub_browse_commit
    assert_command "browse mojombo/bert commit/5d5582", "open https://github.com/mojombo/bert/commit/5d5582"
  end

  def test_hub_browse_tracking_nothing
    stub_tracking_nothing
    assert_command "browse mojombo/bert", "open https://github.com/mojombo/bert"
  end

  def test_hub_browse_url
    assert_command "browse -u mojombo/bert", "echo https://github.com/mojombo/bert"
  end

  def test_hub_browse_self
    assert_command "browse resque", "open https://github.com/tpw/resque"
  end

  def test_hub_browse_subpage
    assert_command "browse resque commits",
      "open https://github.com/tpw/resque/commits/master"
    assert_command "browse resque issues",
      "open https://github.com/tpw/resque/issues"
    assert_command "browse resque wiki",
      "open https://github.com/tpw/resque/wiki"
  end

  def test_hub_browse_on_branch
    stub_branch('refs/heads/feature')
    stub_tracking('feature', 'mislav', 'experimental')

    assert_command "browse resque", "open https://github.com/tpw/resque"
    assert_command "browse resque commits",
      "open https://github.com/tpw/resque/commits/master"

    assert_command "browse",
      "open https://github.com/mislav/hub/tree/experimental"
    assert_command "browse -- tree",
      "open https://github.com/mislav/hub/tree/experimental"
    assert_command "browse -- commits",
      "open https://github.com/mislav/hub/commits/experimental"
  end

  def test_hub_browse_current
    assert_command "browse", "open https://github.com/defunkt/hub"
    assert_command "browse --", "open https://github.com/defunkt/hub"
  end

  def test_hub_browse_commit_from_current
    assert_command "browse -- commit/6616e4", "open https://github.com/defunkt/hub/commit/6616e4"
  end

  def test_hub_browse_no_tracking
    stub_tracking_nothing
    assert_command "browse", "open https://github.com/defunkt/hub"
  end

  def test_hub_browse_no_tracking_on_branch
    stub_branch('refs/heads/feature')
    stub_tracking_nothing('feature')
    assert_command "browse", "open https://github.com/defunkt/hub"
  end

  def test_hub_browse_current_wiki
    stub_repo_url 'git://github.com/defunkt/hub.wiki.git'

    assert_command "browse", "open https://github.com/defunkt/hub/wiki"
    assert_command "browse -- wiki", "open https://github.com/defunkt/hub/wiki"
    assert_command "browse -- commits", "open https://github.com/defunkt/hub/wiki/_history"
    assert_command "browse -- pages", "open https://github.com/defunkt/hub/wiki/_pages"
  end

  def test_hub_browse_current_subpage
    assert_command "browse -- network",
      "open https://github.com/defunkt/hub/network"
    assert_command "browse -- anything/everything",
      "open https://github.com/defunkt/hub/anything/everything"
  end

  def test_hub_browse_deprecated_private
    with_browser_env('echo') do
      assert_includes "Warning: the `-p` flag has no effect anymore\n", hub("browse -p defunkt/hub")
    end
  end

  def test_hub_browse_no_repo
    stub_repo_url(nil)
    assert_equal "Usage: hub browse [<USER>/]<REPOSITORY>\n", hub("browse")
  end

  def test_custom_browser
    with_browser_env("custom") do
      assert_browser("custom")
    end
  end

  def test_linux_browser
    stub_available_commands "open", "xdg-open", "cygstart"
    with_browser_env(nil) do
      with_host_os("i686-linux") do
        assert_browser("xdg-open")
      end
    end
  end

  def test_cygwin_browser
    stub_available_commands "open", "cygstart"
    with_browser_env(nil) do
      with_host_os("i686-linux") do
        assert_browser("cygstart")
      end
    end
  end

  def test_no_browser
    stub_available_commands()
    expected = "Please set $BROWSER to a web launcher to use this command.\n"
    with_browser_env(nil) do
      with_host_os("i686-linux") do
        assert_equal expected, hub("browse")
      end
    end
  end

  def test_multiple_remote_urls
    stub_repo_url("git://example.com/other.git\ngit://github.com/my/repo.git")
    assert_command "browse", "open https://github.com/my/repo"
  end
end
