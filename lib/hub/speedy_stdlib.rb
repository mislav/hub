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
    InvalidURIError = Class.new(StandardError)

    def self.parse(str)
      URI::HTTP.new(str)
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
      attr_writer :port
      alias hostname host

      def initialize(str)
        m = str.to_s.match(%r{^ ([\w-]+): // (?:([^/@]+)@)? ([^/?#]+) }x)
        raise InvalidURIError unless m
        _, self.scheme, self.userinfo, host = m.to_a
        self.host, self.port = host.split(':', 2)
        path, self.fragment = m.post_match.split('#', 2)
        self.path, self.query = path.to_s.split('?', 2)
      end

      def to_s
        url = "#{scheme}://"
        url << "#{userinfo}@" if user || password
        url << host
        url << ":#{@port}" if @port
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

      def port
        (@port || (scheme == 'https' ? 443 : 80)).to_i
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
    end
  end
end
