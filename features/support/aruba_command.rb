require 'open3'
require 'shellwords'

module Aruba
  remove_const :Command
  class Command
    attr_reader :commandline, :stdout, :stderr
    attr_reader :exit_timeout, :io_wait_timeout, :startup_wait_time, :environment, :stop_signal, :exit_status

    def initialize(cmd, mode:, exit_timeout:, io_wait_timeout:,
                   working_directory:, environment:, main_class:, stop_signal:,
                   startup_wait_time:, event_bus:)
      @commandline = cmd
      @working_directory = working_directory
      @event_bus = event_bus
      @exit_timeout = exit_timeout
      @io_wait_timeout = io_wait_timeout
      @startup_wait_time = startup_wait_time
      @environment = environment
      @stop_signal = stop_signal

      @stopped = false
      @exit_status = nil
      @stdout = nil
      @stderr = nil
    end

    def inspect
      %(#<Command "#{@commandline}" exited:#{@exit_status}>)
    end

    def output
      stdout + stderr
    end

    def start
      @event_bus.notify Events::CommandStarted.new(self)
      cmd = Shellwords.split @commandline
      @stdin_io, @stdout_io, @stderr_io, @thread = Open3.popen3(@environment, *cmd, chdir: @working_directory)
    end

    def write(input)
      @stdin_io.write input
      @stdin_io.flush
    end

    def close_io(io)
      case io
      when :stdin then @stdin_io.close
      else
        raise ArgumentError, io.to_s
      end
    end

    def stop
      return if @exit_status
      @event_bus.notify Events::CommandStopped.new(self)
      terminate
    end

    def terminate
      return if @exit_status

      close_io(:stdin)
      @stdout = @stdout_io.read
      @stderr = @stderr_io.read

      status = @thread.value
      @exit_status = status.exitstatus
      @thread = nil
    end

    def interactive?
      true
    end

    def timed_out?
      false
    end
  end
end
