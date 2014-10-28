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
