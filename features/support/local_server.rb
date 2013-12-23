# based on <github.com/jnicklas/capybara/blob/ab62b27/lib/capybara/server.rb>
require 'net/http'
require 'rack/handler/webrick'

module Hub
  class LocalServer
    class Identify < Struct.new(:app)
      def call(env)
        if env["PATH_INFO"] == "/__identify__"
          [200, {}, [app.object_id.to_s]]
        else
          app.call(env)
        end
      end
    end

    def self.ports
      @ports ||= {}
    end

    class JsonParamsParser < Struct.new(:app)
      def call(env)
        if env['rack.input'] and not input_parsed?(env) and type_match?(env)
          env['rack.request.form_input'] = env['rack.input']
          data = env['rack.input'].read
          env['rack.request.form_hash'] = data.empty?? {} : JSON.parse(data)
        end
        app.call(env)
      end

      def input_parsed? env
        env['rack.request.form_input'].eql? env['rack.input']
      end

      def type_match? env
        type = env['CONTENT_TYPE'] and
          type.split(/\s*[;,]\s*/, 2).first.downcase == 'application/json'
      end
    end

    def self.start_sinatra(&block)
      require 'json'
      require 'sinatra/base'
      klass = Class.new(Sinatra::Base)
      klass.use JsonParamsParser
      klass.set :environment, :test
      klass.disable :protection
      klass.class_eval(&block)
      klass.helpers do
        def json(value)
          content_type :json
          JSON.generate value
        end

        def assert(expected)
          expected.each do |key, value|
            if params[key] != value
              halt 422, json(
                :message => "expected %s to be %s; got %s" % [
                  key.inspect,
                  value.inspect,
                  params[key].inspect
                ]
              )
            end
          end
        end

        def assert_basic_auth(*expected)
          require 'rack/auth/basic'
          auth = Rack::Auth::Basic::Request.new(env)
          if auth.credentials != expected
            halt 401, json(
              :message => "expected %p; got %p" % [
                expected, auth.credentials
              ]
            )
          end
        end
      end

      new(klass.new).start
    end

    attr_reader :app, :host, :port
    attr_accessor :server

    def initialize(app, host = '127.0.0.1')
      @app = app
      @host = host
      @server = nil
      @server_thread = nil
    end

    def responsive?
      return false if @server_thread && @server_thread.join(0)

      res = Net::HTTP.start(host, port) { |http| http.get('/__identify__') }

      res.is_a?(Net::HTTPSuccess) and res.body == app.object_id.to_s
    rescue Errno::ECONNREFUSED, Errno::EBADF
      return false
    end

    def start
      @port = self.class.ports[app.object_id]

      if not @port or not responsive?
        @server_thread = start_handler(Identify.new(app)) do |server, host, port|
          self.server = server
          @port = self.class.ports[app.object_id] = port
        end

        Timeout.timeout(60) { @server_thread.join(0.01) until responsive? }
      end
    rescue TimeoutError
      raise "Rack application timed out during boot"
    else
      self
    end

    def start_handler(app)
      server = nil
      thread = Rack::Handler::WEBrick.run(app, server_options) { |s| server = s }
      addr = server.listeners[0].addr
      yield server, addr[3], addr[1]
      return thread
    end

    def server_options
      { :Port => 0,
        :BindAddress => '127.0.0.1',
        :ShutdownSocketWithoutClose => true,
        :ServerType => Thread,
        :AccessLog => [],
        :Logger => WEBrick::Log::new(nil, 0)
      }
    end

    def stop
      server.shutdown
      @server_thread.join
    end
  end
end

WEBrick::HTTPStatus::StatusMessage[422] = "Unprocessable Entity"
