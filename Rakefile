require 'rake/testtask'

task :default => :test

Rake::TestTask.new do |t|
  t.libs << 'lib'
  t.pattern = 'test/**/*_test.rb'
  t.verbose = false
end

desc "Launch Kicker (like autotest)"
task :kicker do
  puts "Kicking... (ctrl+c to cancel)"
  exec "kicker -e rake test bin"
end

desc "Build a gem"
task :gem => [ :gemspec, :build ]

begin
  require 'jeweler'
  $LOAD_PATH.unshift 'lib'
  require 'hub'
  Jeweler::Tasks.new do |gemspec|
    gemspec.name = "hub/version"
    gemspec.summary = gemspec.description = "hub introduces git to GitHub"
    gemspec.homepage = "http://github.com/defunkt/hub"
    gemspec.version = Hub::Version
    gemspec.authors = ["Chris Wanstrath"]
    gemspec.email = "chris@ozmm.org"
  end
rescue LoadError
  puts "Jeweler not available."
  puts "Install it with: gem install jeweler"
end

begin
  require 'sdoc_helpers'
rescue LoadError
  puts "sdoc support not enabled. Please gem install sdoc-helpers."
end

desc "Push a new version to Gemcutter"
task :publish => [ :test, :gemspec, :build ] do
  system "git tag v#{Hub::Version}"
  system "git push origin v#{Hub::Version}"
  system "git push origin master"
  system "gem push pkg/hub-#{Hub::Version}.gem"
  system "git clean -fd"
  exec "rake pages"
end

