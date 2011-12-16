$LOAD_PATH.unshift File.expand_path("../../lib", __FILE__)

require 'test/unit'
require 'hub'
require 'hub/standalone'
require 'webmock/test_unit'
require 'rbconfig'
require 'yaml'

# We're checking for `open` in our tests
ENV['BROWSER'] = 'open'

# Setup path with fake executables in case a test hits them
fakebin_dir = File.expand_path('../fakebin', __FILE__)
ENV['PATH'] = "#{fakebin_dir}:#{ENV['PATH']}"

WebMock::BodyPattern.class_eval do
  undef normalize_hash
  # override normalizing hash since it otherwise requires JSON
  def normalize_hash(hash) hash end
end

class Test::Unit::TestCase

  # Stubbing commands and git configuration
  COMMANDS = []

  def setup
    Hub::Context.class_eval do
      remove_method :which
      define_method :which do |name|
        COMMANDS.include?(name) ? "/usr/bin/#{name}" : nil
      end
    end

    COMMANDS.replace %w[open groff]
    Hub::Context::DIRNAME.replace 'hub'
    Hub::Context::REMOTES.clear
    Hub::Context::Remote.clear!

    Hub::Context::GIT_CONFIG.executable = 'git'

    @git = Hub::Context::GIT_CONFIG.replace(Hash.new { |h, k|
      unless k.index('config alias.') == 0
        raise ArgumentError, "`git #{k}` not stubbed"
      end
    }).update(
      'remote' => "mislav\norigin",
      'symbolic-ref -q HEAD' => 'refs/heads/master',
      'config github.user'   => 'tpw',
      'config github.token'  => 'abc123',
      'config --get-all remote.origin.url' => 'git://github.com/defunkt/hub.git',
      'config --get-all remote.mislav.url' => 'git://github.com/mislav/hub.git',
      "name-rev master@{upstream} --name-only --refs='refs/remotes/*' --no-undefined" => 'remotes/origin/master',
      'config --bool hub.http-clone' => 'false',
      'config hub.protocol' => nil,
      'rev-parse --git-dir' => '.git'
    )
  end

  def stub_available_commands(*names)
    COMMANDS.replace names
  end

  def auth(user = @git['config github.user'], password = @git['config github.token'])
    "#{user}%2Ftoken:#{password}@"
  end

  def stub_repo_url(value)
    @git['config --get-all remote.origin.url'] = value
    Hub::Context::REMOTES.clear
  end

  def stub_branch(value)
    @git['symbolic-ref -q HEAD'] = value
  end

  def stub_tracking(from, remote_name, remote_branch)
    value = remote_branch ? "remotes/#{remote_name}/#{remote_branch}" : nil
    @git["name-rev #{from}@{upstream} --name-only --refs='refs/remotes/*' --no-undefined"] = value
  end

  def stub_tracking_nothing(from = 'master')
    stub_tracking(from, nil, nil)
  end

  def stub_alias(name, value)
    @git["config alias.#{name}"] = value
  end

  def stub_https_is_preferred
    @git['config hub.protocol'] = 'https'
  end

  def stub_github_user(name)
    @git['config github.user'] = name
  end

  def stub_github_token(token)
    @git['config github.token'] = token
  end

  def stub_no_git_repo
    @git.replace({})
  end

  def stub_fork(user, repo, status)
    stub_request(:get, "github.com/api/v2/yaml/repos/show/#{user}/#{repo}").
    to_return(:status => status)
  end

  def stub_existing_fork(user, repo = 'hub')
    stub_fork(user, repo, 200)
  end

  def stub_nonexisting_fork(user, repo = 'hub')
    stub_fork(user, repo, 404)
  end

  def stub_no_remotes
    @git['remote'] = ''
  end

  def stub_remotes_group(name, value)
    @git["config remotes.#{name}"] = value
  end

  # Mock responses

  def mock_pullreq_response(id, name_with_owner = 'defunkt/hub')
    YAML.dump('pull' => {
      'html_url' => "https://github.com/#{name_with_owner}/pull/#{id}"
    })
  end

  #
  # Scope calls with custom configurations
  #
  def with_tmpdir(value)
    dir, ENV['TMPDIR'] = ENV['TMPDIR'], value
    yield
  ensure
    ENV['TMPDIR'] = dir
  end

  def with_browser_env(value)
    browser, ENV['BROWSER'] = ENV['BROWSER'], value
    yield
  ensure
    ENV['BROWSER'] = browser
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

  # Shortcut for creating a `Hub` instance. Pass it what you would
  # normally pass `hub` on the command line, e.g.
  #
  # shell: hub clone rtomayko/tilt
  #  test: Hub("clone rtomayko/tilt")
  def Hub(args)
    Hub::Runner.new(*args.split(' '))
  end

  # Shortcut for running the `hub` command in a subprocess. Returns
  # STDOUT as a string. Pass it what you would normally pass `hub` on
  # the command line, e.g.
  #
  # shell: hub clone rtomayko/tilt
  #  test: hub("clone rtomayko/tilt")
  #
  # If a block is given it will be run in the child process before
  # execution begins. You can use this to monkeypatch or fudge the
  # environment before running hub.
  def hub(args, input = nil)
    parent_read, child_write = IO.pipe
    child_read, parent_write = IO.pipe if input

    fork do
      yield if block_given?
      $stdin.reopen(child_read) if input
      $stdout.reopen(child_write)
      $stderr.reopen(child_write)
      Hub(args).execute
    end
    
    if input
      parent_write.write input
      parent_write.close
    end
    child_write.close
    parent_read.read
  end

  # Asserts that `hub` will run a specific git command based on
  # certain input.
  #
  # e.g.
  #  assert_command "clone git/hub", "git clone git://github.com/git/hub.git"
  #
  # Here we are saying that this:
  #   $ hub clone git/hub
  # Should in turn execute this:
  #   $ git clone git://github.com/git/hub.git
  def assert_command(input, expected)
    assert_equal expected, Hub(input).command, "$ git #{input}"
  end

  def assert_commands(*expected)
    input = expected.pop
    assert_equal expected, Hub(input).commands
  end

  # Asserts that the command will be forwarded to git without changes
  def assert_forwarded(input)
    cmd = Hub(input)
    assert !cmd.args.changed?, "arguments were not supposed to change: #{cmd.args.inspect}"
  end

  # Asserts that `hub` will show a specific alias command for a
  # specific shell.
  #
  # e.g.
  #  assert_alias_command "sh", "alias git=hub"
  #
  # Here we are saying that this:
  #   $ hub alias sh
  # Should display this:
  #   Run this in your shell to start using `hub` as `git`:
  #     alias git=hub
  def assert_alias_command(shell, command)
    expected = "Run this in your shell to start using `hub` as `git`:\n  %s\n"
    assert_equal(expected % command, hub("alias #{shell}"))
  end

  # Asserts that `haystack` includes `needle`.
  def assert_includes(needle, haystack)
    assert haystack.include?(needle),
      "expected #{needle.inspect} in #{haystack.inspect}"
  end

  # Asserts that `haystack` does not include `needle`.
  def assert_not_includes(needle, haystack)
    assert !haystack.include?(needle),
      "didn't expect #{needle.inspect} in #{haystack.inspect}"
  end

  # Version of assert_equal tailored for big output
  def assert_output(expected, command)
    output = hub(command) { ENV['GIT'] = 'echo' }
    assert expected == output,
      "expected:\n#{expected}\ngot:\n#{output}"
  end

  # assert that a specific browser is called
  def assert_browser(browser)
    assert_command "browse", "#{browser} https://github.com/defunkt/hub"
  end
end
