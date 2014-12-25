prevent_require = lambda do |name|
  $" << "#{name}.rb"
  require name # hax to avoid Ruby 2.0.0 bug
end

unless defined?(CGI)
  prevent_require.call 'cgi'

  module CGI
    ESCAPE_RE = /[^a-zA-Z0-9 .~_-]/

    def self.escape(s)
      s.to_s.gsub(ESCAPE_RE) {|match|
        '%' + match.unpack('H2' * match.bytesize).join('%').upcase
      }.tr(' ', '+')
    end

    def self.unescape(s)
      s.tr('+', ' ').gsub(/((?:%[0-9a-fA-F]{2})+)/) {
        [$1.delete('%')].pack('H*')
      }
    end
  end
end

unless defined?(URI)
  prevent_require.call 'uri'

  Kernel.module_eval do
    def URI(str)
      URI.parse(str)
    end
  end

  module URI
    Error = Class.new(StandardError)
    InvalidURIError = Class.new(Error)
    InvalidComponentError = Class.new(Error)

    def self.parse(str)
      m = str.to_s.match(%r{^ ([\w-]+): // (?:([^/@]+)@)? ([^/?#]+) }x)
      raise InvalidURIError unless m

      _, scheme, userinfo, host = m.to_a
      default_port = scheme == 'https' ? 443 : 80
      host, port = host.split(':', 2)
      port = port ? port.to_i : default_port

      path, fragment = m.post_match.split('#', 2)
      path, query = path.split('?', 2) if path
      path = path.to_s

      URI::HTTP.new(scheme, userinfo, host, port, nil, path, nil, query, fragment)
    end

    def self.encode_www_form(params)
      params.map { |k, v|
        if v.class == Array
          encode_www_form(v.map { |x| [k, x] })
        else
          ek = CGI.escape(k)
          v.nil? ? ek : "#{ek}=#{CGI.escape(v)}"
        end
      }.join("&")
    end

    def self.===(other)
      other.respond_to?(:host)
    end

    class HTTP
      attr_accessor :scheme, :user, :password, :host, :path, :query, :fragment
      attr_reader :port
      alias hostname host

      def initialize(scheme, userinfo, host, port, registry, path, opaque, query, fragment, *extra)
        self.scheme = scheme
        self.userinfo = userinfo
        self.host = host
        self.port = port
        self.path = path
        self.query = query
        self.fragment = fragment
      end

      def to_s
        url = "#{scheme}://"
        url << "#{userinfo}@" if user || password
        url << host << display_port
        url << path
        url << "?#{query}" if query
        url << "##{fragment}" if fragment
        url
      end

      def request_uri
        url = path
        url += "?#{query}" if query
        url
      end

      def port=(number)
        if number.is_a?(Fixnum) && number > 0
          @port = number
        else
          raise InvalidComponentError, "bad component(expected port component): %p" % number
        end
      end

      def userinfo=(info)
        self.user, self.password = info.to_s.split(':', 2)
        info
      end

      def userinfo
        if password then "#{user}:#{password}"
        elsif user then user
        end
      end

      def find_proxy
      end

      private

      def default_port
        self.scheme == 'https' ? 443 : 80
      end

      def display_port
        if port != default_port
          ":#{port}"
        else
          ""
        end
      end
    end
  end
end
