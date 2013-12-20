require 'aruba/cucumber'
require 'fileutils'
require 'forwardable'

system_git = `which git 2>/dev/null`.chomp
lib_dir = File.expand_path('../../../lib', __FILE__)
bin_dir = File.expand_path('../fakebin', __FILE__)

Before do
  # don't want hub to run in bundle
  unset_bundler_env_vars
  # have bin/hub load code from the current project
  set_env 'RUBYLIB', lib_dir
  # speed up load time by skipping RubyGems
  set_env 'RUBYOPT', '--disable-gems' if RUBY_VERSION > '1.9'
  # put fakebin on the PATH
  set_env 'PATH', "#{bin_dir}:#{ENV['PATH']}"
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
  set_env 'HUB_TEST_HOST', '127.0.0.1:0'
  # ensure we use fakebin `open` to test browsing
  set_env 'BROWSER', 'open'
  # sabotage opening a commit message editor interactively
  set_env 'GIT_EDITOR', 'false'

  author_name  = "Hub"
  author_email = "hub@test.local"
  set_env 'GIT_AUTHOR_NAME',     author_name
  set_env 'GIT_COMMITTER_NAME',  author_name
  set_env 'GIT_AUTHOR_EMAIL',    author_email
  set_env 'GIT_COMMITTER_EMAIL', author_email

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

  failure_message_for_should do |cmd|
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
    history.should include(cmd)
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
      command.should be_successful_command
      command.output
    end
  end

  def empty_commit
    run_silent "git commit --quiet -m empty --allow-empty"
  end

  # Aruba unnecessarily creates new Announcer instance on each invocation
  def announcer
    @announcer ||= super
  end
}
