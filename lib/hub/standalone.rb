module Hub
  module Standalone
    extend self

    RUBY_BIN = if File.executable? '/usr/bin/ruby' then '/usr/bin/ruby'
               else
                 require 'rbconfig'
                 File.join RbConfig::CONFIG['bindir'], RbConfig::CONFIG['ruby_install_name']
               end

    PREAMBLE = <<-preamble
#!#{RUBY_BIN}
#
# This file, hub, is generated code.
# Please DO NOT EDIT or send patches for it.
#
# Please take a look at the source from
# https://github.com/defunkt/hub
# and submit patches against the individual files
# that build hub.
#

preamble

    POSTAMBLE = "Hub::Runner.execute(*ARGV)\n"
    __DIR__   = File.dirname(__FILE__)
    MANPAGE   = "__END__\n#{File.read(__DIR__ + '/../../man/hub.1')}"

    def save(filename, path = '.')
      target = File.join(File.expand_path(path), filename)
      File.open(target, 'w') do |f|
        f.puts build
        f.chmod 0755
      end
    end

    def build
      root = File.dirname(__FILE__)

      standalone = ''
      standalone << PREAMBLE

      files = Dir["#{root}/*.rb"].sort - [__FILE__]
      # ensure context.rb appears before others
      ctx = files.find {|f| f['context.rb'] } and files.unshift(files.delete(ctx))

      files.each do |file|
        File.readlines(file).each do |line|
          standalone << line if line !~ /^\s*#/
        end
      end

      standalone << POSTAMBLE
      standalone << MANPAGE
      standalone
    end
  end
end
