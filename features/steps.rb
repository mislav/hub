require 'fileutils'

Given /^HTTPS is preferred$/ do
  run_silent %(git config --global hub.protocol https)
end

Given /^there are no remotes$/ do
  run_silent('git remote').should be_empty
end

Given /^"([^"]*)" is a whitelisted Enterprise host$/ do |host|
  run_silent %(git config --global --add hub.host "#{host}")
end

Given /^the "([^"]*)" remote has url "([^"]*)"$/ do |remote_name, url|
  remotes = run_silent('git remote').split("\n")
  unless remotes.include? remote_name
    run_silent %(git remote add #{remote_name} "#{url}")
  else
    run_silent %(git remote set-url #{remote_name} "#{url}")
  end
end

Given /^I am "([^"]*)" on ([\w.-]+)(?: with OAuth token "([^"]*)")?$/ do |name, host, token|
  edit_hub_config do |cfg|
    entry = {'user' => name}
    entry['oauth_token'] = token if token
    cfg[host.downcase] = [entry]
  end
end

Given /^\$(\w+) is "([^"]*)"$/ do |name, value|
  set_env name, value
end

Given /^I am in "([^"]*)" git repo$/ do |dir_name|
  step %(a git repo in "#{dir_name}")
  step %(I cd to "#{dir_name}")
end

Given /^a git repo in "([^"]*)"$/ do |dir_name|
  step %(a directory named "#{dir_name}")
  dirs << dir_name
  step %(I successfully run `git init --quiet`)
  dirs.pop
end

Given /^there is a commit named "([^"]+)"$/ do |name|
  run_silent %(git commit --quiet --allow-empty --allow-empty-message -m '')
  run_silent %(git commit --quiet --allow-empty --allow-empty-message -m '')
  run_silent %(git tag #{name})
  run_silent %(git reset --quiet --hard HEAD^)
end

Given /^the current dir is not a repo$/ do
  in_current_dir do
    FileUtils.rm_rf '.git'
  end
end

Given /^the GitHub API server:$/ do |endpoints_str|
  @server = Hub::LocalServer.start_sinatra do
    eval endpoints_str, binding
  end
  # hit our Sinatra server instead of github.com
  set_env 'HUB_TEST_HOST', "127.0.0.1:#{@server.port}"
end

Then /^"([^"]*)" should be run$/ do |cmd|
  assert_command_run cmd
end

Then /^it should clone "([^"]*)"$/ do |repo|
  step %("git clone #{repo}" should be run)
end

Then /^"([^"]+)" should not be run$/ do |pattern|
  history.all? {|h| h.should_not include(pattern) }
end

Then /^there should be no output$/ do
  assert_exact_output('', all_output)
end

Then /^the git command should be unchanged$/ do
  @commands.should_not be_empty
  assert_command_run @commands.last.sub(/^hub\b/, 'git')
end

Then /^the url for "([^"]*)" should be "([^"]*)"$/ do |name, url|
  found = run_silent %(git config --get-all remote.#{name}.url)
  found.should eql(url)
end

Then /^the "([^"]*)" submodule url should be "([^"]*)"$/ do |name, url|
  found = run_silent %(git config --get-all submodule."#{name}".url)
  found.should eql(url)
end

Then /^there should be no "([^"]*)" remote$/ do |remote_name|
  remotes = run_silent('git remote').split("\n")
  remotes.should_not include(remote_name)
end
