require 'aruba/cucumber'
require 'fileutils'
require 'forwardable'
require 'tmpdir'

# Ruby 2.2.0 compat
Cucumber::Ast::Step.class_eval do
  undef_method :text_length
  def text_length(name=name())
    self.class::INDENT + self.class::INDENT +
      keyword.unpack('U*').length +
      name.unpack('U*').length
  end
end

system_git = `which git 2>/dev/null`.chomp
lib_dir = File.expand_path('../../../lib', __FILE__)
bin_dir = File.expand_path('../fakebin', __FILE__)
hub_dir = Dir.mktmpdir('hub_build')
raise 'hub build failed' unless system("./script/build -o #{hub_dir}/hub")

Before do
  # don't want hub to run in bundle
  unset_bundler_env_vars
  # have bin/hub load code from the current project
  set_env 'RUBYLIB', lib_dir
  # speed up load time by skipping RubyGems
  set_env 'RUBYOPT', '--disable-gems' if RUBY_VERSION > '1.9'
  # put fakebin on the PATH
  set_env 'PATH', "#{hub_dir}:#{bin_dir}:#{ENV['PATH']}"
  # clear out GIT if it happens to be set
  set_env 'GIT', nil
  # exclude this project's git directory from use in testing
  set_env 'GIT_CEILING_DIRECTORIES', File.dirname(lib_dir)
  # sabotage git commands that might try to access a remote host
  set_env 'GIT_PROXY_COMMAND', 'echo'
  # avoids reading from current user's "~/.gitconfig"
  set_env 'HOME', File.expand_path(File.join(current_dir, 'home'))
  # used in fakebin/git
  set_env 'HUB_SYSTEM_GIT', system_git
  # ensure that api.github.com is actually never hit in tests
  set_env 'HUB_TEST_HOST', 'http://127.0.0.1:0'
  # ensure we use fakebin `open` to test browsing
  set_env 'BROWSER', 'open'
  # sabotage opening a commit message editor interactively
  set_env 'GIT_EDITOR', 'false'
  # reset current localization settings
  set_env 'LANG', nil
  set_env 'LANGUAGE', nil
  set_env 'LC_ALL', 'en_US.UTF-8'
  # ignore current user's token
  set_env 'GITHUB_TOKEN', nil
  set_env 'GITHUB_USER', nil
  set_env 'GITHUB_PASSWORD', nil
  set_env 'GITHUB_HOST', nil

  author_name  = "Hub"
  author_email = "hub@test.local"
  set_env 'GIT_AUTHOR_NAME',     author_name
  set_env 'GIT_COMMITTER_NAME',  author_name
  set_env 'GIT_AUTHOR_EMAIL',    author_email
  set_env 'GIT_COMMITTER_EMAIL', author_email

  set_env 'HUB_VERSION', 'dev'
  set_env 'HUB_REPORT_CRASH', 'never'
  set_env 'HUB_PROTOCOL', nil

  FileUtils.mkdir_p ENV['HOME']

  # increase process exit timeout from the default of 3 seconds
  @aruba_timeout_seconds = 5

  if defined?(RUBY_ENGINE) and RUBY_ENGINE == 'jruby'
    @aruba_io_wait_seconds = 0.1
  else
    @aruba_io_wait_seconds = 0.02
  end
end

After do
  @server.stop if defined? @server and @server
  FileUtils.rm_f("#{bin_dir}/vim")
end

RSpec::Matchers.define :be_successful_command do
  match do |cmd|
    cmd.success?
  end

  failure_message do |cmd|
    %(command "#{cmd}" exited with status #{cmd.status}:) <<
      cmd.output.gsub(/^/, ' ' * 2)
  end
end

class SimpleCommand
  attr_reader :output
  extend Forwardable

  def_delegator :@status, :exitstatus, :status
  def_delegators :@status, :success?

  def initialize cmd
    @cmd = cmd
  end

  def to_s
    @cmd
  end

  def self.run cmd
    command = new(cmd)
    command.run
    command
  end

  def run
    @output = `#{@cmd} 2>&1`.chomp
    @status = $?
    $?.success?
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
    histfile = File.join(ENV['HOME'], '.history')
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
    config = File.join(ENV['HOME'], '.config/hub')
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
    File.open("#{bin_dir}/vim", 'w', 0755) { |exe|
      exe.puts "#!/bin/bash"
      exe.puts "set -e"
      exe.puts bash_code
    }
  end

  def run_silent cmd
    in_current_dir do
      command = SimpleCommand.run(cmd)
      expect(command).to be_successful_command
      command.output
    end
  end

  def empty_commit(message = nil)
    unless message
      @empty_commit_count = defined?(@empty_commit_count) ? @empty_commit_count + 1 : 1
      message = "empty #{@empty_commit_count}"
    end
    run_silent "git commit --quiet -m '#{message}' --allow-empty"
  end

  # Aruba unnecessarily creates new Announcer instance on each invocation
  def announcer
    @announcer ||= super
  end

  def shell_escape(message)
    message.to_s.gsub(/['"\\ $]/) { |m| "\\#{m}" }
  end
}
