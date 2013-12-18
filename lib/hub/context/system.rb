module Hub
  module Context
    module System
      # Cross-platform web browser command; respects the value set in $BROWSER.
      #
      # Returns an array, e.g.: ['open']
      def browser_launcher
        browser = ENV['BROWSER'] || (
          osx? ? 'open' : windows? ? %w[cmd /c start] :
          %w[xdg-open cygstart x-www-browser firefox opera mozilla netscape].find { |comm| which comm }
        )

        abort "Please set $BROWSER to a web launcher to use this command." unless browser
        Array(browser)
      end

      def osx?
        require 'rbconfig'
        RbConfig::CONFIG['host_os'].to_s.include?('darwin')
      end

      def windows?
        require 'rbconfig'
        RbConfig::CONFIG['host_os'] =~ /msdos|mswin|djgpp|mingw|windows/
      end

      def unix?
        require 'rbconfig'
        RbConfig::CONFIG['host_os'] =~ /(aix|darwin|linux|(net|free|open)bsd|cygwin|solaris|irix|hpux)/i
      end

      # Cross-platform way of finding an executable in the $PATH.
      #
      #   which('ruby') #=> /usr/bin/ruby
      def which(cmd)
        exts = ENV['PATHEXT'] ? ENV['PATHEXT'].split(';') : ['']
        ENV['PATH'].split(File::PATH_SEPARATOR).each do |path|
          exts.each { |ext|
            exe = "#{path}/#{cmd}#{ext}"
            return exe if File.executable? exe
          }
        end
        return nil
      end

      # Checks whether a command exists on this system in the $PATH.
      #
      # name - The String name of the command to check for.
      #
      # Returns a Boolean.
      def command?(name)
        !which(name).nil?
      end

      def tmp_dir
        ENV['TMPDIR'] || ENV['TEMP'] || '/tmp'
      end

      def terminal_width
        if unix?
          width = %x{stty size 2>#{NULL}}.split[1].to_i
          width = %x{tput cols 2>#{NULL}}.to_i if width.zero?
        else
          width = 0
        end
        width < 10 ? 78 : width
      end
    end
  end
end
