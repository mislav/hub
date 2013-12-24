module Hub
  # Reads ssh_config(5) files and records "Host" to "HostName" mappings to
  # provide resolving of ssh aliases.
  class SshConfig
    CONFIG_FILES = %w(~/.ssh/config /etc/ssh_config /etc/ssh/ssh_config)
    CONFIG_RE = /^\s*(Host|HostName)\s+/

    def initialize files = nil
      @settings = Hash.new {|h,k| h[k] = {} }
      Array(files || CONFIG_FILES).each do |path|
        file = File.expand_path path
        parse_file file if File.readable?(file)
      end
    end

    # Public: Get a setting as it would apply to a specific hostname.
    #
    # Yields if not found.
    def get_value(host, key, &fallback)
      host = host.to_s.downcase
      @settings[host].fetch(key.to_sym, &fallback)
    end

    def parse_file file
      host_patterns = %w[*]

      IO.foreach(file) do |line|
        next unless line =~ CONFIG_RE
        key, value = $1.to_sym, $'.chomp
        value.gsub!(/^"|"$/, '')

        if :Host == key
          host_patterns = value.downcase.split(/\s+/)
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
