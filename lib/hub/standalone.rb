module Hub
  module Standalone
    extend self

    HUB_ROOT = File.expand_path('../../..', __FILE__)

    PREAMBLE = <<-preamble
#
# This file is generated code. DO NOT send patches for it.
#
# Original source files with comments are at:
# https://github.com/github/hub
#

preamble

    def save(filename, path = '.')
      target = File.join(File.expand_path(path), filename)
      File.open(target, 'w') do |f|
        build f
        f.chmod 0755
      end
    end

    def build io
      io.puts "#!#{ruby_shebang}"
      io << PREAMBLE

      each_source_file do |filename|
        File.open(filename, 'r') do |source|
          source.each_line do |line|
            next if line =~ /^\s*#/
            if line.include?(' VERSION =')
              line.sub!(/'(.+?)'/, "'#{detailed_version}'")
            end
            io << line
          end
        end
        io.puts ''
      end

      io.puts "Hub::Runner.execute(*ARGV)"
      io.puts "\n__END__"
      io << File.read(File.join(HUB_ROOT, 'man/hub.1'))
    end

    def each_source_file
      File.open(File.join(HUB_ROOT, 'lib/hub.rb'), 'r') do |main|
        main.each_line do |req|
          if req =~ /^require\s+["'](.+)["']/
            yield File.join(HUB_ROOT, 'lib', "#{$1}.rb")
          end
        end
      end
    end

    def detailed_version
      version = `git describe --tags HEAD 2>/dev/null`.chomp
      if version.empty?
        version = Hub::VERSION
        head_sha = `git rev-parse --short HEAD 2>/dev/null`.chomp
        version += "-g#{head_sha}" unless head_sha.empty?
        version
      else
        version.sub(/^v/, '')
      end
    end

    def ruby_executable
      if File.executable? '/usr/bin/ruby' then '/usr/bin/ruby'
      else
        require 'rbconfig'
        File.join RbConfig::CONFIG['bindir'], RbConfig::CONFIG['ruby_install_name']
      end
    end

    def ruby_shebang
      ruby = ruby_executable
      `#{ruby_executable} --disable-gems -e0 2>/dev/null`
      if $?.success?
        "#{ruby} --disable-gems"
      else
        ruby
      end
    end
  end
end
