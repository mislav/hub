require 'aruba/cucumber'
require 'fileutils'
require 'hub/context'

unless system_git = Hub::Context.which('git')
  abort "Error: `git` not found in PATH"
end

lib_dir = File.expand_path('../../../lib', __FILE__)
bin_dir = File.expand_path('../fakebin', __FILE__)

Before do
  # don't want hub to run in bundle
  unset_bundler_env_vars
  # have bin/hub load code from the current project
  set_env 'RUBYLIB', lib_dir
  # put fakebin on the PATH
  set_env 'PATH', "#{bin_dir}:#{ENV['PATH']}"
  # exclude this project's git directory from use in testing
  set_env 'GIT_CEILING_DIRECTORIES', File.dirname(lib_dir)
  # sabotage git commands that might try to access a remote host
  set_env 'GIT_PROXY_COMMAND', 'echo'
  # avoids reading from current user's "~/.gitconfig"
  set_env 'HOME', File.expand_path(File.join(current_dir, 'home'))
  # used in fakebin/git
  set_env 'HUB_SYSTEM_GIT', system_git

  FileUtils.mkdir_p ENV['HOME']
end

Before '~@noexec' do
  set_env 'GIT', nil
end

Before '@noexec' do
  set_env 'GIT', 'echo'
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

  def run_silent cmd
    in_current_dir do
      output = `#{cmd} 2>&1`.chomp
      $?.should be_success
      output
    end
  end
}
