$LOAD_PATH.unshift File.dirname(__FILE__)
require 'helper'
require 'webmock/test_unit'

class HubTest < Test::Unit::TestCase
  include WebMock

  COMMANDS = []

  Hub::Commands.class_eval do
    remove_method :command?
    define_method :command? do |name|
      COMMANDS.include?(name)
    end
  end

  def setup
    COMMANDS.replace %w[open groff]

    @git = Hub::Context::GIT_CONFIG.replace(Hash.new { |h, k|
      raise ArgumentError, "`git #{k}` not stubbed"
    }).update(
      'remote' => "mislav\norigin",
      'symbolic-ref -q HEAD' => 'refs/heads/master',
      'config github.user'   => 'tpw',
      'config github.token'  => 'abc123',
      'config remote.origin.url'     => 'git://github.com/defunkt/hub.git',
      'config remote.mislav.url'     => 'git://github.com/mislav/hub.git',
      'config branch.master.remote'  => 'origin',
      'config branch.master.merge'   => 'refs/heads/master',
      'config branch.feature.remote' => 'mislav',
      'config branch.feature.merge'  => 'refs/heads/experimental',
      'config --bool hub.http-clone' => 'false'
    )
  end

  def test_private_clone
    input   = "clone -p rtomayko/ronn"
    command = "git clone git@github.com:rtomayko/ronn.git"
    assert_command input, command
  end

  def test_public_clone
    input   = "clone rtomayko/ronn"
    command = "git clone git://github.com/rtomayko/ronn.git"
    assert_command input, command
  end

  def test_your_private_clone
    input   = "clone -p resque"
    command = "git clone git@github.com:tpw/resque.git"
    assert_command input, command
  end

  def test_your_public_clone
    input   = "clone resque"
    command = "git clone git://github.com/tpw/resque.git"
    assert_command input, command
  end

  def test_clone_with_arguments_and_path
    input   = "clone --bare -o master -- resque"
    command = "git clone --bare -o master -- git://github.com/tpw/resque.git"
    assert_command input, command
  end

  def test_your_private_clone_fails_without_config
    out = hub("clone -p mustache") do
      stub_github_user(nil)
    end

    assert_equal "** No GitHub user set. See http://github.com/guides/local-github-config\n", out
  end

  def test_your_public_clone_fails_without_config
    out = hub("clone mustache") do
      stub_github_user(nil)
    end

    assert_equal "** No GitHub user set. See http://github.com/guides/local-github-config\n", out
  end

  def test_private_clone_left_alone
    input   = "clone git@github.com:rtomayko/ronn.git"
    command = "git clone git@github.com:rtomayko/ronn.git"
    assert_command input, command
  end

  def test_public_clone_left_alone
    input   = "clone git://github.com/rtomayko/ronn.git"
    command = "git clone git://github.com/rtomayko/ronn.git"
    assert_command input, command
  end

  def test_normal_public_clone_with_path
    input   = "clone git://github.com/rtomayko/ronn.git ronn-dev"
    command = "git clone git://github.com/rtomayko/ronn.git ronn-dev"
    assert_command input, command
  end

  def test_normal_clone_from_path
    input   = "clone ./test"
    command = "git clone ./test"
    assert_command input, command
  end

  def test_remote_origin
    input   = "remote add origin"
    command = "git remote add origin git://github.com/tpw/hub.git"
    assert_command input, command
  end

  def test_private_remote_origin
    input   = "remote add -p origin"
    command = "git remote add origin git@github.com:tpw/hub.git"
    assert_command input, command
  end

  def test_public_remote_origin_as_normal
    input   = "remote add origin http://github.com/defunkt/resque.git"
    command = "git remote add origin http://github.com/defunkt/resque.git"
    assert_command input, command
  end

  def test_remote_from_rel_path
    input = "remote add origin ./path"
    command = "git remote add origin ./path"
    assert_command input, command
  end

  def test_remote_from_abs_path
    input = "remote add origin /path"
    command = "git remote add origin /path"
    assert_command input, command
  end

  def test_private_remote_origin_as_normal
    input   = "remote add origin git@github.com:defunkt/resque.git"
    command = "git remote add origin git@github.com:defunkt/resque.git"
    assert_command input, command
  end

  def test_public_submodule
    input   = "submodule add wycats/bundler vendor/bundler"
    command = "git submodule add git://github.com/wycats/bundler.git vendor/bundler"
    assert_command input, command
  end

  def test_private_submodule
    input   = "submodule add -p grit vendor/grit"
    command = "git submodule add git@github.com:tpw/grit.git vendor/grit"
    assert_command input, command
  end

  def test_submodule_branch
    input   = "submodule add -b ryppl ryppl/pip vendor/pip"
    command = "git submodule add -b ryppl git://github.com/ryppl/pip.git vendor/pip"
    assert_command input, command
  end

  def test_submodule_with_args
    input   = "submodule -q add --bare -- grit grit"
    command = "git submodule -q add --bare -- git://github.com/tpw/grit.git grit"
    assert_command input, command
  end

  def test_private_remote
    input   = "remote add -p rtomayko"
    command = "git remote add rtomayko git@github.com:rtomayko/hub.git"
    assert_command input, command
  end

  def test_public_remote
    input   = "remote add rtomayko"
    command = "git remote add rtomayko git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_public_remote_f
    input   = "remote add -f rtomayko"
    command = "git remote add -f rtomayko git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_named_public_remote
    input   = "remote add origin rtomayko"
    command = "git remote add origin git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_named_public_remote_f
    input   = "remote add -f origin rtomayko"
    command = "git remote add -f origin git://github.com/rtomayko/hub.git"
    assert_command input, command
  end

  def test_private_remote_with_repo
    input   = "remote add -p jashkenas/coffee-script"
    command = "git remote add jashkenas git@github.com:jashkenas/coffee-script.git"
    assert_command input, command
  end

  def test_public_remote_with_repo
    input   = "remote add jashkenas/coffee-script"
    command = "git remote add jashkenas git://github.com/jashkenas/coffee-script.git"
    assert_command input, command
  end

  def test_public_remote_f_with_repo
    input   = "remote add -f jashkenas/coffee-script"
    command = "git remote add -f jashkenas git://github.com/jashkenas/coffee-script.git"
    assert_command input, command
  end

  def test_named_private_remote_with_repo
    input   = "remote add -p origin jashkenas/coffee-script"
    command = "git remote add origin git@github.com:jashkenas/coffee-script.git"
    assert_command input, command
  end

  def test_fetch_existing_remote
    assert_command "fetch mislav", "git fetch mislav"
  end

  def test_fetch_new_remote
    stub_remotes_group('xoebus', nil)
    stub_existing_fork('xoebus')

    h = Hub("fetch xoebus")
    assert_equal "git remote add xoebus git://github.com/xoebus/hub.git", h.command
    assert_equal "git fetch xoebus", h.after
  end

  def test_fetch_new_remote_with_options
    stub_remotes_group('xoebus', nil)
    stub_existing_fork('xoebus')

    h = Hub("fetch --depth=1 --prune xoebus")
    assert_equal "git remote add xoebus git://github.com/xoebus/hub.git", h.command
    assert_equal "git fetch --depth=1 --prune xoebus", h.after
  end

  def test_fetch_multiple_new_remotes
    stub_remotes_group('xoebus', nil)
    stub_remotes_group('rtomayko', nil)
    stub_existing_fork('xoebus')
    stub_existing_fork('rtomayko')

    h = Hub("fetch --multiple xoebus rtomayko")

    assert_equal "git remote add xoebus git://github.com/xoebus/hub.git", h.command
    expected = ["git remote add rtomayko git://github.com/rtomayko/hub.git"] <<
                "git fetch --multiple xoebus rtomayko"
    assert_equal expected.join('; '), h.after
  end

  def test_fetch_multiple_comma_separated_remotes
    stub_remotes_group('xoebus', nil)
    stub_remotes_group('rtomayko', nil)
    stub_existing_fork('xoebus')
    stub_existing_fork('rtomayko')

    h = Hub("fetch xoebus,rtomayko")

    assert_equal "git remote add xoebus git://github.com/xoebus/hub.git", h.command
    expected = ["git remote add rtomayko git://github.com/rtomayko/hub.git"] <<
                "git fetch --multiple xoebus rtomayko"
    assert_equal expected.join('; '), h.after
  end

  def test_fetch_multiple_new_remotes_with_filtering
    stub_remotes_group('xoebus', nil)
    stub_remotes_group('mygrp', 'one two')
    stub_remotes_group('typo', nil)
    stub_existing_fork('xoebus')
    stub_nonexisting_fork('typo')

    # mislav: existing remote; skipped
    # xoebus: new remote, fork exists; added
    # mygrp:  a remotes group; skipped
    # URL:    can't be a username; skipped
    # typo:   fork doesn't exist; skipped
    h = Hub("fetch --multiple mislav xoebus mygrp git://example.com typo")

    assert_equal "git remote add xoebus git://github.com/xoebus/hub.git", h.command
    expected = "git fetch --multiple mislav xoebus mygrp git://example.com typo"
    assert_equal expected, h.after
  end

  def test_cherry_pick
    h = Hub("cherry-pick a319d88")
    assert_equal "git cherry-pick a319d88", h.command
    assert !h.args.after?
  end

  def test_cherry_pick_url
    url = 'http://github.com/mislav/hub/commit/a319d88#comments'
    h = Hub("cherry-pick #{url}")
    assert_equal "git fetch mislav", h.command
    assert_equal "git cherry-pick a319d88", h.after
  end

  def test_cherry_pick_url_with_remote_add
    url = 'http://github.com/xoebus/hub/commit/a319d88'
    h = Hub("cherry-pick #{url}")
    assert_equal "git remote add -f xoebus git://github.com/xoebus/hub.git", h.command
    assert_equal "git cherry-pick a319d88", h.after
  end

  def test_cherry_pick_private_url_with_remote_add
    url = 'https://github.com/xoebus/hub/commit/a319d88'
    h = Hub("cherry-pick #{url}")
    assert_equal "git remote add -f xoebus git@github.com:xoebus/hub.git", h.command
    assert_equal "git cherry-pick a319d88", h.after
  end

  def test_cherry_pick_origin_url
    url = 'https://github.com/defunkt/hub/commit/a319d88'
    h = Hub("cherry-pick #{url}")
    assert_equal "git fetch origin", h.command
    assert_equal "git cherry-pick a319d88", h.after
  end

  def test_cherry_pick_github_user_notation
    h = Hub("cherry-pick mislav@a319d88")
    assert_equal "git fetch mislav", h.command
    assert_equal "git cherry-pick a319d88", h.after
  end

  def test_cherry_pick_github_user_repo_notation
    # not supported
    h = Hub("cherry-pick mislav/hubbub@a319d88")
    assert_equal "git cherry-pick mislav/hubbub@a319d88", h.command
    assert !h.args.after?
  end

  def test_cherry_pick_github_notation_too_short
    h = Hub("cherry-pick mislav@a319")
    assert_equal "git cherry-pick mislav@a319", h.command
    assert !h.args.after?
  end

  def test_cherry_pick_github_notation_with_remote_add
    h = Hub("cherry-pick xoebus@a319d88")
    assert_equal "git remote add -f xoebus git://github.com/xoebus/hub.git", h.command
    assert_equal "git cherry-pick a319d88", h.after
  end

  def test_init
    h = Hub("init -g")
    assert_equal "git init", h.command
    assert_equal "git remote add origin git@github.com:tpw/hub.git", h.after
  end

  def test_init_no_login
    out = hub("init -g") do
      stub_github_user(nil)
    end

    assert_equal "** No GitHub user set. See http://github.com/guides/local-github-config\n", out
  end

  def test_push_two
    h = Hub("push origin,staging cool-feature")
    assert_equal "git push origin cool-feature", h.command
    assert_equal "git push staging cool-feature", h.after
  end

  def test_push_more
    h = Hub("push origin,staging,qa cool-feature")
    assert_equal "git push origin cool-feature", h.command
    assert_equal "git push staging cool-feature; git push qa cool-feature", h.after
  end

  def test_fork
    stub_nonexisting_fork('tpw')
    stub_request(:post, "github.com/api/v2/yaml/repos/fork/defunkt/hub").with { |req|
      params = Hash[*req.body.split(/[&=]/)]
      params == { 'login'=>'tpw', 'token'=>'abc123' }
    }

    expected = "remote add -f tpw git@github.com:tpw/hub.git\n"
    expected << "new remote: tpw\n"
    assert_equal expected, hub("fork") { ENV['GIT'] = 'echo' }
  end

  def test_fork_no_remote
    stub_nonexisting_fork('tpw')
    stub_request(:post, "github.com/api/v2/yaml/repos/fork/defunkt/hub")

    assert_equal "", hub("fork --no-remote") { ENV['GIT'] = 'echo' }
  end

  def test_fork_already_exists
    stub_existing_fork('tpw')

    expected = "tpw/hub already exists on GitHub\n"
    expected << "remote add -f tpw git@github.com:tpw/hub.git\n"
    expected << "new remote: tpw\n"
    assert_equal expected, hub("fork") { ENV['GIT'] = 'echo' }
  end

  def test_version
    out = hub('--version')
    assert_includes "git version 1.7.0.4", out
    assert_includes "hub version #{Hub::Version}", out
  end

  def test_help
    assert_equal Hub::Commands.improved_help_text, hub("help")
  end

  def test_help_by_default
    assert_equal Hub::Commands.improved_help_text, hub("")
  end

  def test_help_with_pager
    assert_equal Hub::Commands.improved_help_text, hub("-p")
  end

  def test_help_hub
    help_manpage = hub("help hub")
    assert_includes "git + hub = github", help_manpage
    assert_includes "Chris Wanstrath :: chris@ozmm.org", help_manpage
    assert_includes <<-config, help_manpage
Use git-config(1) to display the currently configured GitHub username:
config
  end

  def test_help_hub_no_groff
    stub_available_commands()
    assert_equal "** Can't find groff(1)\n", hub("help hub")
  end

  def test_hub_standalone
    help_standalone = hub("hub standalone")
    assert_equal Hub::Standalone.build, help_standalone
  end

  def test_hub_compare
    assert_command "compare refactor",
      "open http://github.com/defunkt/hub/compare/refactor"
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

    assert_command "compare",
      "open http://github.com/mislav/hub/compare/experimental"
  end

  def test_hub_compare_range
    assert_command "compare 1.0...fix",
      "open http://github.com/defunkt/hub/compare/1.0...fix"
  end

  def test_hub_compare_fork
    assert_command "compare myfork feature",
      "open http://github.com/myfork/hub/compare/feature"
  end

  def test_hub_compare_private
    assert_command "compare -p myfork topsecret",
      "open https://github.com/myfork/hub/compare/topsecret"
  end

  def test_hub_compare_url
    assert_command "compare -u 1.0...1.1",
      "echo http://github.com/defunkt/hub/compare/1.0...1.1"
  end

  def test_hub_browse
    assert_command "browse mojombo/bert", "open http://github.com/mojombo/bert"
  end

  def test_hub_browse_tracking_nothing
    stub_tracking_nothing
    assert_command "browse mojombo/bert", "open http://github.com/mojombo/bert"
  end

  def test_hub_browse_url
    assert_command "browse -u mojombo/bert", "echo http://github.com/mojombo/bert"
  end

  def test_hub_browse_private
    assert_command "browse -p bmizerany/sinatra",
      "open https://github.com/bmizerany/sinatra"
  end

  def test_hub_browse_self
    assert_command "browse resque", "open http://github.com/tpw/resque"
  end

  def test_hub_browse_subpage
    assert_command "browse resque commits",
      "open http://github.com/tpw/resque/commits/master"
    assert_command "browse resque issues",
      "open http://github.com/tpw/resque/issues"
    assert_command "browse resque wiki",
      "open http://wiki.github.com/tpw/resque/"
  end

  def test_hub_browse_on_branch
    stub_branch('refs/heads/feature')

    assert_command "browse resque", "open http://github.com/tpw/resque"
    assert_command "browse resque commits",
      "open http://github.com/tpw/resque/commits/master"

    assert_command "browse",
      "open http://github.com/mislav/hub/tree/experimental"
    assert_command "browse -- tree",
      "open http://github.com/mislav/hub/tree/experimental"
    assert_command "browse -- commits",
      "open http://github.com/mislav/hub/commits/experimental"
  end

  def test_hub_browse_self_private
    assert_command "browse -p github", "open https://github.com/tpw/github"
  end

  def test_hub_browse_current
    assert_command "browse", "open http://github.com/defunkt/hub"
    assert_command "browse --", "open http://github.com/defunkt/hub"
  end

  def test_hub_browse_current_subpage
    assert_command "browse -- network",
      "open http://github.com/defunkt/hub/network"
    assert_command "browse -- anything/everything",
      "open http://github.com/defunkt/hub/anything/everything"
  end

  def test_hub_browse_current_private
    assert_command "browse -p", "open https://github.com/defunkt/hub"
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
      with_ruby_platform("i686-linux") do
        assert_browser("xdg-open")
      end
    end
  end

  def test_cygwin_browser
    stub_available_commands "open", "cygstart"
    with_browser_env(nil) do
      with_ruby_platform("i686-linux") do
        assert_browser("cygstart")
      end
    end
  end

  def test_no_browser
    stub_available_commands()
    expected = "Please set $BROWSER to a web launcher to use this command.\n"
    with_browser_env(nil) do
      with_ruby_platform("i686-linux") do
        assert_equal expected, hub("browse")
      end
    end
  end

  protected

    def stub_github_user(name)
      @git['config github.user'] = name
    end

    def stub_repo_url(value)
      @git['config remote.origin.url'] = value
      Hub::Context::REMOTES.clear
    end

    def stub_branch(value)
      @git['symbolic-ref -q HEAD'] = value
    end

    def stub_tracking_nothing
      @git['config branch.master.remote'] = nil
      @git['config branch.master.merge'] = nil
    end

    def stub_remotes_group(name, value)
      @git["config remotes.#{name}"] = value
    end

    def stub_existing_fork(user)
      stub_fork(user, 200)
    end

    def stub_nonexisting_fork(user)
      stub_fork(user, 404)
    end

    def stub_fork(user, status)
      stub_request(:get, "github.com/api/v2/yaml/repos/show/#{user}/hub").
        to_return(:status => status)
    end

    def stub_available_commands(*names)
      COMMANDS.replace names
    end

    def with_browser_env(value)
      browser, ENV['BROWSER'] = ENV['BROWSER'], value
      yield
    ensure
      ENV['BROWSER'] = browser
    end

    def assert_browser(browser)
      assert_command "browse", "#{browser} http://github.com/defunkt/hub"
    end

    def with_ruby_platform(value)
      platform = RUBY_PLATFORM
      Object.send(:remove_const, :RUBY_PLATFORM)
      Object.const_set(:RUBY_PLATFORM, value)
      yield
    ensure
      Object.send(:remove_const, :RUBY_PLATFORM)
      Object.const_set(:RUBY_PLATFORM, platform)
    end

end
