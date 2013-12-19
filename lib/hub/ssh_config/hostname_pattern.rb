module Hub
  class SSHConfig
    class HostPattern
      def initialize(pattern)
        @pattern = pattern.to_s.downcase
      end

      def to_s
        @pattern
      end

      def ==(other)
        other.to_s == to_s
      end

      def matcher
        @matcher ||=
          if @pattern == '*'
            ->(*) { true }
          elsif @pattern !~ /[?*]/
            ->(hostname) { hostname.to_s.downcase == @pattern }
          else
            re = self.class.pattern_to_regexp @pattern
            ->(hostname) { re =~ hostname }
          end
      end

      def match?(hostname)
        matcher.call(hostname)
      end

      def self.pattern_to_regexp(pattern)
        escaped = Regexp.escape(pattern)
        escaped.gsub!('\*', '.*')
        escaped.gsub!('\?', '.')
        /^#{escaped}$/i
      end
    end
  end
end
