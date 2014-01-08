require 'rake/testtask'

#
# Helpers
#

def command?(util)
  Rake::Task[:load_path].invoke
  context = Object.new
  require 'uri'
  require 'hub/context'
  context.extend Hub::Context
  context.send(:command?, util)
end

task :load_path do
  $LOAD_PATH.unshift File.expand_path('../lib', __FILE__)
end

task :check_dirty do
  unless system 'git', 'diff', '--quiet', 'HEAD'
    abort "Aborted: you have uncommitted changes"
  end
end


#
# Tests
#

task :default => [:test, :features]

Rake::TestTask.new do |t|
  t.libs << 'test'
  t.pattern = 'test/**/*_test.rb'
end

task :features do
  sh 'RUBYLIB=lib cucumber -f progress -t ~@wip features'
end

#
# Manual
#

if command? :ronn
  desc "Show man page"
  task :man => "man:build" do
    exec "man man/hub.1"
  end

  desc "Build man pages"
  task "man:build" => ["man/hub.1", "man/hub.1.html"]

  extract_examples = lambda { |readme_file|
    # split readme in sections
    examples = File.read(readme_file).split(/^-{4,}$/)[3].strip
    examples.sub!(/^.+?(###)/m, '\1')  # strip intro paragraph
    examples.sub!(/\n+.+\Z/, '')       # remove last line
    examples
  }

  # inject examples from README file to .ronn source
  source_with_examples = lambda { |source, readme|
    examples = extract_examples.call(readme)
    compiled = File.read(source)
    compiled.sub!('{{README}}', examples)
    compiled
  }

  # generate man page with ronn
  compile_ronn = lambda { |destination, type, contents|
    File.popen("ronn --pipe --#{type} --organization=GITHUB --manual='Hub Manual'", 'w+') { |io|
      io.write contents
      io.close_write
      File.open(destination, 'w') { |f| f << io.read }
    }
    abort "ronn --#{type} conversion failed" unless $?.success?
  }

  file "man/hub.1" => ["man/hub.1.ronn", "README.md"] do |task|
    contents = source_with_examples.call(*task.prerequisites)
    compile_ronn.call(task.name, 'roff', contents)
    compile_ronn.call("#{task.name}.html", 'html', contents)
  end

  file "man/hub.1.html" => ["man/hub.1.ronn", "README.md"] do |task|
    Rake::Task["man/hub.1"].invoke
  end
end


#
# Build
#

file "hub" => FileList.new("lib/hub.rb", "lib/hub/*.rb", "man/hub.1") do |task|
  Rake::Task[:load_path].invoke
  require 'hub/version'
  require 'hub/standalone'
  Hub::Standalone.save(task.name)
end

desc "Build standalone script"
task :standalone => "hub"

desc %{Install standalone script and man page.
On Unix-based OS, installs into PREFIX (default: `/usr/local`).
On Windows, installs into Ruby's main bin directory.}
task :install => "hub" do
  require 'rbconfig'
  if RbConfig::CONFIG['host_os'] =~ /mswin|mingw/
    bindir = RbConfig::CONFIG['bindir']
    File.open(File.join(bindir, 'hub.bat'), 'w') { |f| f.write('@"ruby.exe" "%~dpn0" %*') }
    FileUtils.cp 'hub', bindir
  else
    prefix = ENV['PREFIX'] || ENV['prefix'] || '/usr/local'
    prefix = File.join(ENV["DESTDIR"], prefix) if ENV["DESTDIR"]

    FileUtils.mkdir_p "#{prefix}/bin"
    FileUtils.cp "hub", "#{prefix}/bin", :preserve => true

    FileUtils.mkdir_p "#{prefix}/share/man/man1"
    FileUtils.cp "man/hub.1", "#{prefix}/share/man/man1"
  end
end

#
# Release
#

task :release => [:pages, :gem_release, :homebrew]

desc "Copy files to gh-pages branch, but don't publish"
task :gh_pages => [:check_dirty, "hub", "man/hub.1.html"] do
  cp "man/hub.1.html", "html"
  sh "git checkout gh-pages"
  # replace the specific shebang with a generic ruby one
  sh "echo '#!/usr/bin/env' ruby > standalone"
  sh "sed 1d hub >> standalone"
  mv "html", "hub.1.html"
  sh "git add standalone hub.1.html"
  sh "git commit -m 'update standalone'"
end

desc "Publish to GitHub Pages"
task :pages => :gh_pages do
  sh "git push origin gh-pages"
  sh "git checkout master"
  puts "Done."
end

task :gem_release do
  sh "gem release -t"
end

desc "Publish to Homebrew"
task :homebrew do
  require File.expand_path('../lib/hub/version', __FILE__)
  Dir.chdir `brew --prefix`.chomp do
    sh 'git checkout -q master'
    sh 'git pull -q origin master'

    formula_file = 'Library/Formula/hub.rb'
    sha = `curl -fsSL https://github.com/github/hub/archive/v#{Hub::VERSION}.tar.gz | shasum`.split(/\s+/).first
    abort unless $?.success? and sha.length == 40

    formula = File.read formula_file
    formula.sub!(/\bv\d+(\.\d+)*/, "v#{Hub::VERSION}")
    formula.sub!(/\b[0-9a-f]{40}\b/, sha)
    File.open(formula_file, 'w') {|f| f << formula }

    branch = "hub-v#{Hub::VERSION}"
    sh "git checkout -q -B #{branch}"
    sh "git commit -m 'hub v#{Hub::VERSION}' -- #{formula_file}"
    sh "git push -u mislav #{branch}"
    sh "hub pull-request -m 'upgrade hub to v#{Hub::VERSION}'"

    sh "git checkout -q master"
  end
end
