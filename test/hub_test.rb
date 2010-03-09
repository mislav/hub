$LOAD_PATH.unshift File.dirname(__FILE__)
require 'helper'

class HubTest < Test::Unit::TestCase
  def setup
    Hub::Commands::REPO.replace("hub")
    Hub::Commands::USER.replace("tpw")
    Hub::Commands::OWNER.replace("defunkt")
  end

  def test_private_clone
    input   = "clone -p rtomayko/ron"
    command = "git clone git@github.com:rtomayko/ron.git"
    assert_command input, command
  end

  def test_public_clone
    input   = "clone rtomayko/ron"
    command = "git clone git://github.com/rtomayko/ron.git"
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
      Hub::Commands::USER.replace("")
    end

    assert_equal "** No GitHub user set. See http://github.com/guides/local-github-config\n", out
  end

  def test_your_public_clone_fails_without_config
    out = hub("clone mustache") do
      Hub::Commands::USER.replace("")
    end

    assert_equal "** No GitHub user set. See http://github.com/guides/local-github-config\n", out
  end

  def test_private_clone_left_alone
    input   = "clone git@github.com:rtomayko/ron.git"
    command = "git clone git@github.com:rtomayko/ron.git"
    assert_command input, command
  end

  def test_public_clone_left_alone
    input   = "clone git://github.com/rtomayko/ron.git"
    command = "git clone git://github.com/rtomayko/ron.git"
    assert_command input, command
  end

  def test_normal_public_clone_with_path
    input   = "clone git://github.com/rtomayko/ron.git ron-dev"
    command = "git clone git://github.com/rtomayko/ron.git ron-dev"
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

  def test_remote_origin_as_normal
    input   = "remote add origin git@github.com:defunkt/resque.git"
    command = "git remote add origin git@github.com:defunkt/resque.git"
    assert_command input, command
  end

  def test_public_submodule
    input   = "submodule add wycats/bundler vendor/bundler"
    command = "git submodule add git://github.com/wycats.bundler.git"
  end

  def test_private_submodule
    input   = "submodule add -p grit vendor/grit"
    command = "git submodule add git@github.com:tpw/grit.git"
  end

  def test_submodule_with_args
    input   = "submodule -q add --bare -- grit grit"
    command = "git submodule -q add --bare -- git://github.com/tpw/grit.git grit"
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
    input   = "remote add -p rtomayko/tilt"
    command = "git remote add rtomayko git@github.com:rtomayko/tilt.git"
    assert_command input, command
  end

  def test_public_remote_with_repo
    input   = "remote add rtomayko/tilt"
    command = "git remote add rtomayko git://github.com/rtomayko/tilt.git"
    assert_command input, command
  end

  def test_public_remote_f_with_repo
    input   = "remote add -f rtomayko/tilt"
    command = "git remote add -f rtomayko git://github.com/rtomayko/tilt.git"
    assert_command input, command
  end

  def test_named_private_remote_with_repo
    input   = "remote add -p origin rtomayko/tilt"
    command = "git remote add origin git@github.com:rtomayko/tilt.git"
    assert_command input, command
  end

  def test_init
    h = Hub("init -g")
    assert_equal "git init", h.command
    assert_equal "git remote add origin git@github.com:tpw/hub.git", h.after
  end

  def test_init_no_login
    out = hub("init -g") do
      Hub::Commands::USER.replace("")
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

  def test_version
    out = hub('--version')
    assert_includes "git version", out
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
    help_manpage = hub("help hub") do
      Hub::Commands.class_eval do
        remove_method :groff?
        def groff?; false end
      end
    end
    assert_equal "** Can't find groff(1)\n", help_manpage
  end

  def test_hub_standalone
    help_standalone = hub("hub standalone")
    assert_equal Hub::Standalone.build, help_standalone
  end

  def test_hub_open
    assert_command "browse mojombo/bert", "open http://github.com/mojombo/bert"
  end

  def test_hub_open_private
    assert_command "browse -p bmizerany/sinatra", "open https://github.com/bmizerany/sinatra"
  end

  def test_hub_open_self
    assert_command "browse resque", "open http://github.com/tpw/resque"
  end

  def test_hub_open_self_private
    assert_command "browse -p github", "open https://github.com/tpw/github"
  end

  def test_hub_open_current
    assert_command "browse", "open http://github.com/defunkt/hub"
  end

  def test_hub_open_current_private
    assert_command "browse -p", "open https://github.com/defunkt/hub"
  end

  def test_hub_open_no_repo
    Hub::Commands::OWNER.replace("")
    input = "browse"
    assert_equal "Usage: hub browse [<USER>/]<REPOSITORY>\n", hub(input)
  end
end
