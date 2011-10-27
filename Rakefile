require 'rake/testtask'

#
# Helpers
#

def command?(command)
  `which #{command} 2>/dev/null`
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

Rake::TestTask.new do |t|
  t.libs << 'test'
  t.ruby_opts << '-rubygems'
  t.pattern = 'test/**/*_test.rb'
  t.verbose = false
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
    File.popen("ronn --pipe --#{type} --organization=DEFUNKT --manual='Git Manual'", 'w+') { |io|
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
# Gems
#

desc "Build standalone script"
task :standalone => :load_hub do
  require 'hub/standalone'
  Hub::Standalone.save('hub')
end

desc "Install standalone script and man pages"
task :install => :standalone do
  prefix = ENV['PREFIX'] || ENV['prefix'] || '/usr/local'

  FileUtils.mkdir_p "#{prefix}/bin"
  FileUtils.cp "hub", "#{prefix}/bin", :preserve => true

  FileUtils.mkdir_p "#{prefix}/share/man/man1"
  FileUtils.cp "man/hub.1", "#{prefix}/share/man/man1"
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
