# based on <github.com/jnicklas/capybara/blob/ab62b27/lib/capybara/server.rb>
require 'net/http'
require 'rack/handler/webrick'
require 'json'
require 'sinatra/base'

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
          type.split(/\s*[;,]\s*/, 2).first.downcase =~ /[\/+]json$/
      end
    end

    class App < Sinatra::Base
      def invoke
        res = super
        content_type :json unless response.content_type
        response.body = '{}' if blank_response?(response.body) ||
          (response.body.respond_to?(:[]) && blank_response?(response.body[0]))
        res
      end

      def blank_response?(obj)
        obj.nil? || (obj.respond_to?(:empty?) && obj.empty?)
      end
    end

    def self.start_sinatra(&block)
      klass = Class.new(App)
      klass.use JsonParamsParser
      klass.set :environment, :test
      klass.disable :protection
      klass.error(404, 401) { content_type :json; nil }
      klass.class_eval(&block)
      klass.helpers do
        def json(value)
          content_type :json
          JSON.generate value
        end

        def assert(expected)
          expected.each do |key, value|
            if :no == value
              halt 422, json(
                :message => "expected %s not to be passed; got %s" % [
                  key.inspect,
                  params[key].inspect
                ]
              ) if params.key?(key.to_s)
            elsif params[key] != value
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
            halt 401, json(:message => "Bad credentials")
          end
        end

        def generate_patch(subject)
        <<PATCH
From 7eb75a26ee8e402aad79fcf36a4c1461e3ec2592 Mon Sep 17 00:00:00 2001
From: Mislav <mislav.marohnic@gmail.com>
Date: Tue, 24 Jun 2014 11:07:05 -0700
Subject: [PATCH] #{subject}
---
diff --git a/README.md b/README.md
new file mode 100644
index 0000000..ce01362
--- /dev/null
+++ b/README.md
+hello
-- 
1.9.3
PATCH
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
