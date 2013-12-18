module Hub
  module Context
    class GitReader
      attr_reader :executable

      def initialize(executable = nil, &read_proc)
        @executable = executable || 'git'
        # caches output when shelling out to git
        read_proc ||= lambda { |cache, cmd|
          result = %x{#{command_to_string(cmd)} 2>#{NULL}}.chomp
          cache[cmd] = $?.success? && !result.empty? ? result : nil
        }
        @cache = Hash.new(&read_proc)
      end

      def add_exec_flags(flags)
        @executable = Array(executable).concat(flags)
      end

      def read_config(cmd, all = false)
        config_cmd = ['config', (all ? '--get-all' : '--get'), *cmd]
        config_cmd = config_cmd.join(' ') unless cmd.respond_to? :join
        read config_cmd
      end

      def read(cmd)
        @cache[cmd]
      end

      def stub_config_value(key, value, get = '--get')
        stub_command_output "config #{get} #{key}", value
      end

      def stub_command_output(cmd, value)
        @cache[cmd] = value.nil? ? nil : value.to_s
      end

      def stub!(values)
        @cache.update values
      end

      private

      def to_exec(args)
        args = Shellwords.shellwords(args) if args.respond_to? :to_str
        Array(executable) + Array(args)
      end

      def command_to_string(cmd)
        full_cmd = to_exec(cmd)
        full_cmd.respond_to?(:shelljoin) ? full_cmd.shelljoin : full_cmd.join(' ')
      end
    end
  end
end
