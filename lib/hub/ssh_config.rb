require_relative 'ssh_config/hostname_pattern'

module Hub
  # Reads ssh configuration files and records each setting under its host
  # pattern so it can be looked up by hostname.
  class SSHConfig
    CONFIG_FILES = %w(~/.ssh/config /etc/ssh_config /etc/ssh/ssh_config)

    def initialize(files = nil)
      @settings = Hash.new { |h, k| h[k] = {} }

      Array(files || CONFIG_FILES).each do |path|
        file = File.expand_path path

        parse_file file if File.readable?(file)
      end
    end

    # Public: Get a setting as it would apply to a specific hostname.
    #
    # Yields if not found.
    def get_value(hostname, key)
      key = key.to_s.downcase

      @settings.each do |pattern, settings|
        return settings[key] if pattern.match?(hostname) && settings[key]
      end

      yield
    end

    def parse_file(file)
      host_patterns = [HostPattern.new('*')]

      IO.foreach(file) do |line|
        case line
        when /^\s*(#|$)/ then
          next
        when /^\s*(\S+)\s*=/
          key, value = Regexp.last_match(1), $'
        else
          key, value = line.strip.split(/\s+/, 2)
        end

        next if value.nil?
        key.downcase!
        value = $1 if value =~ /^"(.*)"$/
        value.chomp!

        if 'host' == key
          host_patterns = value.split(/\s+/).map { |p| HostPattern.new p }
        else
          record_setting key, value, host_patterns
        end
      end
    end

    def record_setting(key, value, patterns)
      patterns.each do |pattern|
        @settings[pattern][key] ||= value
      end
    end
  end
end
