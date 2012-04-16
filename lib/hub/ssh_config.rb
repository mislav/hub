module Hub
  # Reads ssh configuration files and records each setting under its host
  # pattern so it can be looked up by hostname.
  class SshConfig
    CONFIG_FILES = %w(~/.ssh/config /etc/ssh_config /etc/ssh/ssh_config)

    def initialize files = nil
      @settings = Hash.new {|h,k| h[k] = {} }
      Array(files || CONFIG_FILES).each do |path|
        file = File.expand_path path
        parse_file file if File.exist? file
      end
    end

    # Public: Get a setting as it would apply to a specific hostname.
    #
    # Yields if not found.
    def get_value hostname, key
      key = key.to_s.downcase
      @settings.each do |pattern, settings|
        if pattern.match? hostname and found = settings[key]
          return found
        end
      end
      yield
    end

    class HostPattern
      def initialize pattern
        @pattern = pattern.to_s.downcase
      end

      def to_s() @pattern end
      def ==(other) other.to_s == self.to_s end

      def matcher
        @matcher ||=
          if '*' == @pattern
            Proc.new { true }
          elsif @pattern !~ /[?*]/
            lambda { |hostname| hostname.to_s.downcase == @pattern }
          else
            re = self.class.pattern_to_regexp @pattern
            lambda { |hostname| re =~ hostname }
          end
      end

      def match? hostname
        matcher.call hostname
      end

      def self.pattern_to_regexp pattern
        escaped = Regexp.escape(pattern)
        escaped.gsub!('\*', '.*')
        escaped.gsub!('\?', '.')
        /^#{escaped}$/i
      end
    end

    def parse_file file
      host_patterns = [HostPattern.new('*')]

      IO.foreach(file) do |line|
        case line
        when /^\s*(#|$)/ then next
        when /^\s*(\S+)\s*=/
          key, value = $1, $'
        else
          key, value = line.strip.split(/\s+/, 2)
        end

        next if value.nil?
        key.downcase!
        value = $1 if value =~ /^"(.*)"$/
        value.chomp!

        if 'host' == key
          host_patterns = value.split(/\s+/).map {|p| HostPattern.new p }
        else
          record_setting key, value, host_patterns
        end
      end
    end

    def record_setting key, value, patterns
      patterns.each do |pattern|
        @settings[pattern][key] ||= value
      end
    end
  end
end
