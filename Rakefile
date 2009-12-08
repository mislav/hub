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
    gemspec.post_install_message = <<-message

------------------------------------------------------------

                  You there! Wait, I say!
                  =======================

       If you are a heavy user of `git` on the command
       line  you  may  want  to  install `hub` the old
       fashioned way!  Faster  startup  time,  you see.

       Check  out  the  installation  instructions  at
       http://github.com/defunkt/hub#readme or  simply
       use the `install` command:

       $ hub install

------------------------------------------------------------

message
  end
rescue LoadError
  puts "Jeweler not available."
  puts "Install it with: gem install jeweler"
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

desc "Publish to GitHub Pages"
task :pages => [ :check_dirty, :standalone ] do
  `git checkout gh-pages`
  `md5 -q standalone > standalone.md5`
  `git add standalone*`
  `git commit -m "update standalone"`
  `git push origin gh-pages`
  `git checkout master`
  puts :done
end

task :check_dirty do
  if !`git status`.include?('nothing to commit')
    abort "dirty index - not publishing!"
  end
end

module Standalone
  PREAMBLE = <<-premable
#!/usr/bin/env ruby
#
# This file, hub, is generated code.
# Please DO NOT EDIT or send patches for it.
#
# Please take a look at the source from
# http://github.com/defunkt/hub
# and submit patches against the individual files
# that build hub.
#

premable
  POSTAMBLE = "Hub::Runner.execute(*ARGV)"
end

desc "Build standalone script"
task :standalone => :test do
  File.open('standalone', 'w') do |f|
    f.puts Standalone::PREAMBLE
    Dir['lib/*/**'].each do |file|
      File.readlines(file).each do |line|
        next if line =~ /^\s*#/
        f.puts line
      end
    end
    f.puts Standalone::POSTAMBLE
  end
end
