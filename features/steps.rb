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
  run_silent %(git remote add #{remote_name} "#{url}")
end

Given /^I am "([^"]*)" on ([\w.-]+)$/ do |name, host|
  edit_hub_config do |cfg|
    cfg[host.downcase] = [{'user' => name}]
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

Then /^"([^"]*)" should be run$/ do |cmd|
  assert_command_run cmd
end

Then /^it should clone "([^"]*)"$/ do |repo|
  step %("git clone #{repo}" should be run)
end

Then /^nothing should be run$/ do
  history.should be_empty
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
