require 'rbconfig'

module Hub
  module Context
    module System
      GENERAL_BROWSERS =
        %w(xdg-open cygstart x-www-browser firefox opera mozilla netscape)

      # Cross-platform web browser command; respects the value set in $BROWSER.
      #
      # Returns an array, e.g.: ['open']
      def browser_launcher
        browser =
          ENV['BROWSER'] ||
            case
            when osx?     then 'open'
            when windows? then  %w(cmd /c start)
            else
              GENERAL_BROWSERS.find { |command| which command }
            end

        unless browser
          abort 'Please set $BROWSER to a web launcher to use this command.'
        end

        Array(browser)
      end

      def osx?
        RbConfig::CONFIG['host_os'].to_s.include?('darwin')
      end

      def windows?
        RbConfig::CONFIG['host_os'] =~ /msdos|mswin|djgpp|mingw|windows/
      end

      def unix?
        RbConfig::CONFIG['host_os'] =~
          /(aix|darwin|linux|(net|free|open)bsd|cygwin|solaris|irix|hpux)/i
      end

      # Cross-platform way of finding an executable in the $PATH.
      #
      #   which('ruby') #=> /usr/bin/ruby
      def which(cmd)
        exts = ENV['PATHEXT'] ? ENV['PATHEXT'].split(';') : ['']

        ENV['PATH'].split(File::PATH_SEPARATOR).each do |path|
          exts.each do |ext|
            exe = "#{ path }/#{ cmd }#{ ext }"

            return exe if File.executable? exe
          end
        end
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
        width =
          if unix?
            %x(stty size 2>#{ NULL }).split[1].to_i.nonzero? ||
              %x(tput cols 2>#{ NULL }).to_i
          else
            0
          end

        width < 10 ? 78 : width
      end
    end
  end
end
