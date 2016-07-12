desc "Show man page"
task :man => "man:build" do
  exec "man man/hub.1"
end

desc "Build man pages"
task "man:build" => ["man/hub.1", "man/hub.1.html"]

# split readme in sections
# and return the specified section
def split_readme(file, index)
  File.read(file).split(/^-{4,}$/)[index].strip
end

def extract_configs(readme_file)
  configs = split_readme(readme_file, 4)
  configs.gsub!(/\*\*(.+?)\*\*/, '<\1>')      # replace **xx** with <xx>
  configs.sub!(/\n+.+\Z/, '')                 # remove last line
  configs
end

def extract_examples(readme_file)
  examples = split_readme(readme_file, 3)
  examples.sub!(/^.+?(###)/m, '\1')  # strip intro paragraph
  examples.sub!(/\n+.+\Z/, '')       # remove last line
  examples
end

# inject configs and examples from README file to .ronn source
def compiled_source(source, readme)
  configs = extract_configs(readme)
  examples = extract_examples(readme)
  compiled = File.read(source)
  compiled.sub!('{{CONFIGS}}', configs)
  compiled.sub!('{{README}}', examples)
  compiled
end

# generate man page with ronn
def compile_ronn(destination, type, contents)
  File.popen("bin/ronn --pipe --#{type} --organization=GITHUB --manual='Hub Manual'", 'w+') { |io|
    io.write contents
    io.close_write
    File.open(destination, 'w') { |f| f << io.read }
  }
  abort "ronn --#{type} conversion failed" unless $?.success?
end

file "man/hub.1" => ["man/hub.1.ronn", "README.md"] do |task|
  contents = compiled_source(*task.prerequisites)
  compile_ronn(task.name, 'roff', contents)
  compile_ronn("#{task.name}.html", 'html', contents)
end

file "man/hub.1.html" => ["man/hub.1.ronn", "README.md"] do |task|
  Rake::Task["man/hub.1"].invoke
end
