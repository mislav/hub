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

desc "Install `hub`"
task :setup => :standalone do
  path = ENV['BINPATH'] || %w( ~/bin /usr/local/bin /usr/bin ).detect do |dir|
    File.directory? File.expand_path(dir)
  end

  if path
    puts "Installing into #{path}"
    cp "standalone", hub = File.expand_path(File.join(path, 'hub'))
    chmod 0755, hub
    puts "Done. Type `hub version` to see if it worked!"
  else
    puts  "** Can't find a suitable installation location."
    abort "** Please set the BINPATH env variable and try again."
  end
end

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

desc "Build standalone script"
task :standalone => :test do
  $LOAD_PATH.unshift 'lib'
  require 'hub'
  Hub::Standalone.save('standalone')
end
