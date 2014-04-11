require 'helper'
require 'webmock/minitest'
require 'rbconfig'
require 'yaml'
require 'forwardable'
require 'fileutils'
require 'tempfile'

WebMock::BodyPattern.class_eval do
  undef normalize_hash
  # override normalizing hash since it otherwise requires JSON
  def normalize_hash(hash) hash end

  # strip out the "charset" directive from Content-type value
  alias matches_with_dumb_content_type matches?
  def matches?(body, content_type = "")
    content_type = content_type.split(';').first if content_type.respond_to? :split
    matches_with_dumb_content_type(body, content_type)
  end
end

class HubTest < Minitest::Test
  extend Forwardable

  if defined? WebMock::API
    include WebMock::API
  else
    include WebMock
  end

  COMMANDS = []

  Hub::Context::System.class_eval do
    remove_method :which
    define_method :which do |name|
      COMMANDS.include?(name) ? "/usr/bin/#{name}" : nil
    end
  end

  attr_reader :git_reader
  include Hub::Context::GitReaderMethods
  def_delegators :git_reader, :stub_config_value, :stub_command_output

  def setup
    super
    COMMANDS.replace %w[open groff]
    Hub::Context::PWD.replace '/path/to/hub'
    Hub::SshConfig::CONFIG_FILES.replace []

    @prompt_stubs = prompt_stubs = []
    @password_prompt_stubs = password_prompt_stubs = []
    @repo_file_read = repo_file_read = {}

    Hub::GitHubAPI::Configuration.class_eval do
      undef prompt
      undef prompt_password

      define_method :prompt do |what|
        prompt_stubs.shift.call(what)
      end
      define_method :prompt_password do |host, user|
        password_prompt_stubs.shift.call(host, user)
      end
    end

    Hub::Context::LocalRepo.class_eval do
      undef file_read
      undef file_exist?

      define_method(:file_read) do |*args|
        name = File.join(*args)
        if value = repo_file_read[name]
          value.dup
        else
          raise Errno::ENOENT
        end
      end

      define_method(:file_exist?) do |*args|
        name = File.join(*args)
        !!repo_file_read[name]
      end
    end

    @git_reader = Hub::Context::GitReader.new 'git' do |cache, cmd|
      unless cmd.index('config --get alias.') == 0
        raise ArgumentError, "`git #{cmd}` not stubbed"
      end
    end

    Hub::Commands.instance_variable_set :@git_reader, @git_reader
    Hub::Commands.instance_variable_set :@local_repo, nil
    Hub::Commands.instance_variable_set :@api_client, nil

    FileUtils.rm_rf ENV['HUB_CONFIG']

    edit_hub_config do |data|
      data['github.com'] = [{'user' => 'tpw', 'oauth_token' => 'OTOKEN'}]
    end

    @git_reader.stub! \
      'remote -v' => "origin\tgit://github.com/defunkt/hub.git (fetch)\nmislav\tgit://github.com/mislav/hub.git (fetch)",
      'rev-parse --symbolic-full-name master@{upstream}' => 'refs/remotes/origin/master',
      'config --get --bool hub.http-clone' => 'false',
      'config --get hub.protocol' => nil,
      'config --get-all hub.host' => nil,
      'config --get push.default' => nil,
      'rev-parse -q --git-dir' => '.git'

    stub_branch('refs/heads/master')
    stub_remote_branch('origin/master')
  end

  def teardown
    super
    WebMock.reset!
  end

  def test_cherry_pick
    assert_forwarded "cherry-pick a319d88"
  end

  def test_cherry_pick_url
    url = 'http://github.com/mislav/hub/commit/a319d88'
    assert_commands "git fetch mislav", "git cherry-pick a319d88", "cherry-pick #{url}"
  end

  def test_cherry_pick_url_with_fragment
    url = 'http://github.com/mislav/hub/commit/abcdef0123456789#comments'
    assert_commands "git fetch mislav", "git cherry-pick abcdef0123456789", "cherry-pick #{url}"
  end

  def test_cherry_pick_url_with_remote_add
    url = 'https://github.com/xoebus/hub/commit/a319d88'
    assert_commands "git remote add -f xoebus git://github.com/xoebus/hub.git",
                    "git cherry-pick a319d88",
                    "cherry-pick #{url}"
  end

  def test_cherry_pick_origin_url
    url = 'https://github.com/defunkt/hub/commit/a319d88'
    assert_commands "git fetch origin", "git cherry-pick a319d88", "cherry-pick #{url}"
  end

  def test_cherry_pick_github_user_notation
    assert_commands "git fetch mislav", "git cherry-pick 368af20", "cherry-pick mislav@368af20"
  end

  def test_cherry_pick_github_user_repo_notation
    # not supported
    assert_forwarded "cherry-pick mislav/hubbub@a319d88"
  end

  def test_cherry_pick_github_notation_too_short
    assert_forwarded "cherry-pick mislav@a319"
  end

  def test_cherry_pick_github_notation_with_remote_add
    assert_commands "git remote add -f xoebus git://github.com/xoebus/hub.git",
                    "git cherry-pick a319d88",
                    "cherry-pick xoebus@a319d88"
  end

  def test_am_untouched
    assert_forwarded "am some.patch"
  end

  def test_am_pull_request
    stub_request(:get, "https://api.github.com/repos/defunkt/hub/pulls/55").
      with(:headers => {'Accept'=>'application/vnd.github.v3.patch', 'Authorization'=>'token OTOKEN'}).
      to_return(:status => 200)

    with_tmpdir('/tmp/') do
      assert_commands "git am --signoff /tmp/55.patch -p2",
                      "am --signoff https://github.com/defunkt/hub/pull/55#comment_123 -p2"

      cmd = Hub("am https://github.com/defunkt/hub/pull/55/files").command
      assert_includes '/tmp/55.patch', cmd
    end
  end

  def test_am_no_tmpdir
    stub_request(:get, "https://api.github.com/repos/defunkt/hub/pulls/55").
      to_return(:status => 200)

    with_tmpdir(nil) do
      cmd = Hub("am https://github.com/defunkt/hub/pull/55").command
      assert_includes '/tmp/55.patch', cmd
    end
  end

  def test_am_commit_url
    stub_request(:get, "https://api.github.com/repos/davidbalbert/hub/commits/fdb9921").
      with(:headers => {'Accept'=>'application/vnd.github.v3.patch', 'Authorization'=>'token OTOKEN'}).
      to_return(:status => 200)

    with_tmpdir('/tmp/') do
      url = 'https://github.com/davidbalbert/hub/commit/fdb9921'
      assert_commands "git am --signoff /tmp/fdb9921.patch -p2",
                      "am --signoff #{url} -p2"
    end
  end

  def test_am_gist
    stub_request(:get, "https://api.github.com/gists/8da7fb575debd88c54cf").
      with(:headers => {'Authorization'=>'token OTOKEN'}).
      to_return(:body => Hub::JSON.generate(:files => {
        'file.diff' => {
          :raw_url => "https://gist.github.com/raw/8da7fb575debd88c54cf/SHA/file.diff"
        }
      }))

    stub_request(:get, "https://gist.github.com/raw/8da7fb575debd88c54cf/SHA/file.diff").
      with(:headers => {'Accept'=>'text/plain'}).
      to_return(:status => 200)

    with_tmpdir('/tmp/') do
      url = 'https://gist.github.com/8da7fb575debd88c54cf'

      assert_commands "git am --signoff /tmp/gist-8da7fb575debd88c54cf.txt -p2",
                      "am --signoff #{url} -p2"
    end
  end

  def test_apply_untouched
    assert_forwarded "apply some.patch"
  end

  def test_apply_pull_request
    stub_request(:get, "https://api.github.com/repos/defunkt/hub/pulls/55").
      to_return(:status => 200)

    with_tmpdir('/tmp/') do
      assert_commands "git apply /tmp/55.patch -p2",
                      "apply https://github.com/defunkt/hub/pull/55 -p2"

      cmd = Hub("apply https://github.com/defunkt/hub/pull/55/files").command
      assert_includes '/tmp/55.patch', cmd
    end
  end

  def test_apply_commit_url
    stub_request(:get, "https://api.github.com/repos/davidbalbert/hub/commits/fdb9921").
      to_return(:status => 200)

    with_tmpdir('/tmp/') do
      url = 'https://github.com/davidbalbert/hub/commit/fdb9921'

      assert_commands "git apply /tmp/fdb9921.patch -p2",
                      "apply #{url} -p2"
    end
  end

  def test_apply_gist
    stub_request(:get, "https://api.github.com/gists/8da7fb575debd88c54cf").
      with(:headers => {'Authorization'=>'token OTOKEN'}).
      to_return(:body => Hub::JSON.generate(:files => {
        'file.diff' => {
          :raw_url => "https://gist.github.com/raw/8da7fb575debd88c54cf/SHA/file.diff"
        }
      }))

    stub_request(:get, "https://gist.github.com/raw/8da7fb575debd88c54cf/SHA/file.diff").
      to_return(:status => 200)

    with_tmpdir('/tmp/') do
      url = 'https://gist.github.com/mislav/8da7fb575debd88c54cf'

      assert_commands "git apply /tmp/gist-8da7fb575debd88c54cf.txt -p2",
                      "apply #{url} -p2"
    end
  end

  def test_init
    stub_no_remotes
    stub_no_git_repo
    assert_commands "git init", "git remote add origin git@github.com:tpw/hub.git", "init -g"
  end

  def test_init_enterprise
    stub_no_remotes
    stub_no_git_repo
    edit_hub_config do |data|
      data['git.my.org'] = [{'user'=>'myfiname'}]
    end

    with_host_env('git.my.org') do
      assert_commands "git init", "git remote add origin git@git.my.org:myfiname/hub.git", "init -g"
    end
  end

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

  def test_push_multiple_refs
    assert_commands "git push origin master new-feature",
                    "git push staging master new-feature",
                    "push origin,staging master new-feature"
  end

  def test_pullrequest_from_branch_tracking_local
    stub_config_value 'push.default', 'upstream'
    stub_branch('refs/heads/feature')
    stub_tracking('feature', 'refs/heads/master')

    stub_request(:post, "https://api.github.com/repos/defunkt/hub/pulls").
      with(:body => {'base' => "master", 'head' => "defunkt:feature", 'title' => "hereyougo" }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request -m hereyougo -f"
  end

  def test_pullrequest_enterprise_no_tracking
    stub_hub_host('git.my.org')
    stub_repo_url('git@git.my.org:defunkt/hub.git')
    stub_branch('refs/heads/feature')
    stub_remote_branch('origin/feature')
    stub_tracking_nothing('feature')
    stub_command_output "rev-list --cherry-pick --right-only --no-merges origin/feature...", nil
    edit_hub_config do |data|
      data['git.my.org'] = [{'user'=>'myfiname', 'oauth_token' => 'FITOKEN'}]
    end

    stub_request(:post, "https://git.my.org/api/v3/repos/defunkt/hub/pulls").
      with(:body => {'base' => "master", 'head' => "defunkt:feature", 'title' => "hereyougo" }).
      to_return(:body => mock_pullreq_response(1, 'api/v3/defunkt/hub', 'git.my.org'))

    expected = "https://git.my.org/api/v3/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request -m hereyougo -f"
  end

  def test_pullrequest_alias
    out = hub('e-note')
    assert_equal hub('pull-request'), out
  end

  def test_version
    out = hub('--version')
    assert_includes "git version 1.7.0.4", out
    assert_includes "hub version #{Hub::Version}", out
  end

  def test_exec_path
    out = hub('--exec-path')
    assert_equal "/usr/lib/git-core\n", out
  end

  def test_exec_path_arg
    out = hub('--exec-path=/home/wombat/share/my-l33t-git-core')
    assert_equal improved_help_text, out
  end

  def test_html_path
    out = hub('--html-path')
    assert_equal "/usr/share/doc/git-doc\n", out
  end

  def test_help
    assert_equal improved_help_text, hub("help")
  end

  def test_help_by_default
    assert_equal improved_help_text, hub("")
  end

  def test_help_with_pager
    assert_equal improved_help_text, hub("-p")
  end

  def test_help_hub
    help_manpage = strip_man_escapes hub("help hub")
    assert_includes "git + hub = github", help_manpage
    assert_includes "Hub will prompt for GitHub username & password", help_manpage.gsub(/ {2,}/, ' ')
  end

  def test_help_flag_on_command
    help_manpage = strip_man_escapes hub("browse --help")
    assert_includes "git + hub = github", help_manpage
    assert_includes "git browse", help_manpage
  end

  def test_help_custom_command
    help_manpage = strip_man_escapes hub("help fork")
    assert_includes "git fork [--no-remote]", help_manpage
  end

  def test_help_short_flag_on_command
    usage_help = hub("create -h")
    expected = "Usage: git create [NAME] [-p] [-d DESCRIPTION] [-h HOMEPAGE]\n"
    assert_equal expected, usage_help

    usage_help = hub("pull-request -h")
    expected = "Usage: git pull-request [-o|--browse] [-f] [-m MESSAGE|-F FILE|-i ISSUE|ISSUE-URL] [-b BASE] [-h HEAD]\n"
    assert_equal expected, usage_help
  end

  def test_help_hub_no_groff
    stub_available_commands()
    assert_equal "** Can't find groff(1)\n", hub("help hub")
  end

  def test_hub_standalone
    assert_includes 'This file is generated code', hub("hub standalone")
  end

  def test_hub_browse_no_repo
    stub_repo_url(nil)
    assert_equal "Usage: hub browse [<USER>/]<REPOSITORY>\n", hub("browse")
  end

  def test_hub_browse_ssh_alias
    with_ssh_config "Host gh\n User git\n HostName github.com" do
      stub_repo_url "gh:singingwolfboy/sekrit.git"
      assert_command "browse", "open https://github.com/singingwolfboy/sekrit"
    end
  end

  def test_hub_browse_ssh_github_alias
    with_ssh_config "Host github.com\n HostName ssh.github.com" do
      stub_repo_url "git@github.com:suan/git-sanity.git"
      assert_command "browse", "open https://github.com/suan/git-sanity"
    end
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

  def test_context_method_doesnt_hijack_git_command
    assert_forwarded 'remotes'
  end

  def test_not_choking_on_ruby_methods
    assert_forwarded 'id'
    assert_forwarded 'name'
  end

  def test_global_flags_preserved
    cmd = '--no-pager --bare -c core.awesome=true -c name=value --git-dir=/srv/www perform'
    assert_command cmd, 'git --bare -c core.awesome=true -c name=value --git-dir=/srv/www --no-pager perform'
    assert_equal %w[git --bare -c core.awesome=true -c name=value --git-dir=/srv/www], git_reader.executable
  end

  private

    def stub_repo_url(value, remote_name = 'origin')
      stub_command_output 'remote -v', "#{remote_name}\t#{value} (fetch)"
    end

    def stub_branch(value)
      @repo_file_read['HEAD'] = "ref: #{value}\n"
    end

    def stub_tracking(from, upstream, remote_branch = nil)
      stub_command_output "rev-parse --symbolic-full-name #{from}@{upstream}",
        remote_branch ? "refs/remotes/#{upstream}/#{remote_branch}" : upstream
    end

    def stub_tracking_nothing(from = 'master')
      stub_tracking(from, nil)
    end

    def stub_remote_branch(branch, sha = 'abc123')
      @repo_file_read["refs/remotes/#{branch}"] = sha
    end

    def stub_remotes_group(name, value)
      stub_config_value "remotes.#{name}", value
    end

    def stub_no_remotes
      stub_command_output 'remote -v', nil
    end

    def stub_no_git_repo
      stub_command_output 'rev-parse -q --git-dir', nil
    end

    def stub_alias(name, value)
      stub_config_value "alias.#{name}", value
    end

    def stub_existing_fork(user, repo = 'hub')
      stub_fork(user, repo, 200)
    end

    def stub_nonexisting_fork(user, repo = 'hub')
      stub_fork(user, repo, 404)
    end

    def stub_fork(user, repo, status)
      stub_request(:get, "https://api.github.com/repos/#{user}/#{repo}").
        to_return(:status => status)
    end

    def stub_available_commands(*names)
      COMMANDS.replace names
    end

    def stub_https_is_preferred
      stub_config_value 'hub.protocol', 'https'
    end

    def stub_hub_host(names)
      stub_config_value "hub.host", Array(names).join("\n"), '--get-all'
    end

    def with_browser_env(value)
      browser, ENV['BROWSER'] = ENV['BROWSER'], value
      yield
    ensure
      ENV['BROWSER'] = browser
    end

    def with_tmpdir(value)
      dir, ENV['TMPDIR'] = ENV['TMPDIR'], value
      yield
    ensure
      ENV['TMPDIR'] = dir
    end

    def with_host_env(value)
      host, ENV['GITHUB_HOST'] = ENV['GITHUB_HOST'], value
      yield
    ensure
      ENV['GITHUB_HOST'] = host
    end

    def assert_browser(browser)
      assert_command "browse", "#{browser} https://github.com/defunkt/hub"
    end

    def with_host_os(value)
      host_os = RbConfig::CONFIG['host_os']
      RbConfig::CONFIG['host_os'] = value
      begin
        yield
      ensure
        RbConfig::CONFIG['host_os'] = host_os
      end
    end

    def mock_pullreq_response(id, name_with_owner = 'defunkt/hub', host = 'github.com')
      Hub::JSON.generate :html_url => "https://#{host}/#{name_with_owner}/pull/#{id}"
    end

    def mock_pull_response(label, priv = false)
      Hub::JSON.generate :head => {
        :label => label,
        :repo => {:private => !!priv}
      }
    end

    def improved_help_text
      Hub::Commands.send :improved_help_text
    end

    def with_ssh_config content
      config_file = Tempfile.open 'ssh_config'
      config_file << content
      config_file.close

      begin
        Hub::SshConfig::CONFIG_FILES.replace [config_file.path]
        yield
      ensure
        config_file.unlink
      end
    end

    def strip_man_escapes(manpage)
      manpage.gsub(/_\010/, '').gsub(/\010./, '')
    end

end
