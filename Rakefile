require 'rake/testtask'

#
# Helpers
#

def command?(command)
  `type -t #{command}`
  $?.success?
end

task :load_hub do
  $LOAD_PATH.unshift 'lib'
  require 'hub'
end

task :check_dirty do
  if !`git status`.include?('nothing to commit')
    abort "dirty index - not publishing!"
  end
end


#
# Tests
#

task :default => :test

if command? :turn
  desc "Run tests"
  task :test do
    suffix = "-n #{ENV['TEST']}" if ENV['TEST']
    sh "turn test/*.rb #{suffix}"
  end
else
  Rake::TestTask.new do |t|
    t.libs << 'lib'
    t.ruby_opts << '-rubygems'
    t.pattern = 'test/**/*_test.rb'
    t.verbose = false
  end
end

if command? :kicker
  desc "Launch Kicker (like autotest)"
  task :kicker do
    puts "Kicking... (ctrl+c to cancel)"
    exec "kicker -e rake test lib"
  end
end


#
# Ron
#

if command? :ronn
  desc "Show the manual"
  task :man => "man:build" do
    exec "man man/hub.1"
  end

  desc "Build the manual"
  task "man:build" do
    sh "ronn -br5 --organization=DEFUNKT --manual='Git Manual' man/*.ronn"
  end
end


#
# Gems
#

desc "Build standalone script"
task :standalone => :load_hub do
  require 'hub/standalone'
  Hub::Standalone.save('hub')
end

begin
  require 'mg'
  MG.new('git-hub.gemspec')
rescue LoadError
  warn "mg not available."
  warn "Install it with: gem install mg"
end

desc "Install standalone script and man pages"
task :install => :standalone do
  prefix = ENV['PREFIX'] || ENV['prefix'] || '/usr/local'

  FileUtils.mkdir_p "#{prefix}/bin"
  FileUtils.cp "hub", "#{prefix}/bin"

  FileUtils.mkdir_p "#{prefix}/share/man/man1"
  FileUtils.cp "man/hub.1", "#{prefix}/share/man/man1"
end

desc "Push a new version."
task :publish => "gem:publish" do
  require 'hub/version'
  system "git tag v#{Hub::Version}"
  sh "git push origin v#{Hub::Version}"
  sh "git push origin master"
  sh "git clean -fd"
  exec "rake pages"
end

desc "Publish to GitHub Pages"
task :pages => [ "man:build", :check_dirty, :standalone ] do
  cp "man/hub.1.html", "html"
  sh "git checkout gh-pages"
  sh "mv hub standalone"
  sh "git add standalone*"
  sh "mv html hub.1.html"
  sh "git add hub.1.html"
  sh "git commit -m 'update standalone'"
  sh "git push origin gh-pages"
  sh "git checkout master"
  puts :done
end
