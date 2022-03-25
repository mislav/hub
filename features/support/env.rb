require 'aruba/cucumber'
require 'fileutils'
require 'forwardable'
require 'tmpdir'
require 'open3'

system_git = `which git 2>/dev/null`.chomp
bin_dir = File.expand_path('../fakebin', __FILE__)

tmpdir = Dir.mktmpdir('hub_test')
tmp_bin_dir = "#{tmpdir}/bin"
Aruba.configure do |aruba|
  aruba.send(:find_option, :root_directory).value = tmpdir
end

hub_dir = Dir.mktmpdir('hub_build')
raise 'hub build failed' unless system("./script/build -o #{hub_dir}/hub")

Before do
  author_name  = "Hub"
  author_email = "hub@test.local"

  aruba.environment.update(
    # speed up load time by skipping RubyGems
    'RUBYOPT' => '--disable-gems',
    # put fakebin on the PATH
    'PATH' => "#{hub_dir}:#{tmp_bin_dir}:#{bin_dir}:#{ENV['PATH']}",
    # clear out GIT if it happens to be set
    'GIT' => nil,
    # exclude this project's git directory from use in testing
    'GIT_CEILING_DIRECTORIES' => File.expand_path('../../..', __FILE__),
    # sabotage git commands that might try to access a remote host
    'GIT_PROXY_COMMAND' => 'echo',
    # avoids reading from current user's "~/.gitconfig"
    'HOME' => expand_path('home'),
    'TMPDIR' => tmpdir,
    # https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html#variables
    'XDG_CONFIG_HOME' => nil,
    'XDG_CONFIG_DIRS' => nil,
    # used in fakebin/git
    'HUB_SYSTEM_GIT' => system_git,
    # ensure that api.github.com is actually never hit in tests
    'HUB_TEST_HOST' => 'http://127.0.0.1:0',
    # ensure we use fakebin `open` to test browsing
    'BROWSER' => 'open',
    # sabotage opening a commit message editor interactively
    'GIT_EDITOR' => 'false',
    # reset current localization settings
    'LANG' => nil,
    'LANGUAGE' => nil,
    'LC_ALL' => 'C.UTF-8',
    # ignore current user's token
    'GITHUB_TOKEN' => nil,
    'GITHUB_USER' => nil,
    'GITHUB_PASSWORD' => nil,
    'GITHUB_HOST' => nil,
    'GITHUB_REPOSITORY' => nil,

    'GIT_AUTHOR_NAME' =>     author_name,
    'GIT_COMMITTER_NAME' =>  author_name,
    'GIT_AUTHOR_EMAIL' =>    author_email,
    'GIT_COMMITTER_EMAIL' => author_email,

    'HUB_VERSION' => 'dev',
    'HUB_REPORT_CRASH' => 'never',
    'HUB_PROTOCOL' => nil,
  )

  FileUtils.mkdir_p(expand_path('~'))
end

After do
  @server.stop if defined? @server and @server
  FileUtils.rm_f("#{tmp_bin_dir}/vim")
end

After('@cache_clear') do
  FileUtils.rm_rf("#{tmpdir}/hub/api")
end

RSpec::Matchers.define :be_successfully_executed do
  match do |cmd|
    expect(cmd).to have_exit_status(0)
  end

  failure_message do |cmd|
    msg = %(command `#{cmd.commandline}` exited with status #{cmd.exit_status})
    stderr = cmd.stderr
    msg << ":\n" << stderr.gsub(/^/, '  ') unless stderr.empty?
    msg
  end
end

World Module.new {
  # If there are multiple inputs, e.g., type in username and then type in password etc.,
  # the Go program will freeze on the second input. Giving it a small time interval
  # temporarily solves the problem.
  # See https://github.com/cucumber/aruba/blob/7afbc5c0cbae9c9a946d70c4c2735ccb86e00f08/lib/aruba/api.rb#L379-L382
  def type(*args)
    super.tap { sleep 0.1 }
  end

  def history
    histfile = expand_path('~/.history')
    if File.exist? histfile
      File.readlines histfile
    else
      []
    end
  end

  def assert_command_run cmd
    cmd += "\n" unless cmd[-1..-1] == "\n"
    expect(history).to include(cmd)
  end

  def edit_hub_config
    config = expand_path('~/.config/hub')
    FileUtils.mkdir_p File.dirname(config)
    if File.exist? config
      data = YAML.load File.read(config)
    else
      data = {}
    end
    yield data
    File.open(config, 'w') { |cfg| cfg << YAML.dump(data) }
  end

  define_method(:text_editor_script) do |bash_code|
    FileUtils.mkdir_p(tmp_bin_dir)
    File.open("#{tmp_bin_dir}/vim", 'w', 0755) { |exe|
      exe.puts "#!/bin/bash"
      exe.puts "set -e"
      exe.puts bash_code
    }
  end

  def empty_commit(message = nil)
    unless message
      @empty_commit_count = defined?(@empty_commit_count) ? @empty_commit_count + 1 : 1
      message = "empty #{@empty_commit_count}"
    end
    run_ignored_command "git commit --quiet -m '#{message}' --allow-empty"
  end

  def shell_escape(message)
    message.to_s.gsub(/['"\\ $]/) { |m| "\\#{m}" }
  end

  # runs a command entirely outside of Aruba's command system and returns its stdout
  def run_ignored_command(cmd_string)
    stdout, stderr, status = Open3.capture3(aruba.environment, cmd_string, chdir: expand_path('.'))
    expect(status).to be_success
    stdout
  end
}
