require 'uri'
require 'yaml'
require 'forwardable'

module Hub
  # A client for the GitHub v2 or v3 APIs, depending on the host.
  #
  # When API v3 is used, user gets prompted for username/password in the shell,
  # then this information is exchanged for an OAuth token which is saved in a file.
  # When API v2 is used, username/token combination is read from git config.
  class GitHubAPI
    attr_reader :config

    # Public: Create a new API client instance
    #
    # Options:
    # - config: an object that implements:
    #   - username(host)
    #   - api_token(host, user)
    def initialize config
      @config = config
    end

    # Fake exception type for net/http exception handling.
    # Necessary because net/http may or may not be loaded at the time.
    module Exceptions
      def self.===(exception)
        exception.class.ancestors.map {|a| a.to_s }.include? 'Net::HTTPExceptions'
      end
    end

    def api_v2? host
      return true # TODO
      host = host.host if host.respond_to? :host
      host !~ /(^|\.)github\.com$/i
    end

    # Public: Determine whether a specific repo already exists.
    def repo_exists? project
      if api_v2? project
        res = get "https://%s/api/v2/yaml/repos/show/%s/%s" %
          [project.host, project.owner, project.name]
        res.success?
      else
        raise NotImplementedError
      end
    end

    # Public: Fork the specified repo.
    def fork_repo project
      if api_v2? project
        res = post "https://%s/api/v2/yaml/repos/fork/%s/%s" %
          [project.host, project.owner, project.name]
        res.error! unless res.success?
      else
        raise NotImplementedError
      end
    end

    # Public: Create a new project.
    def create_repo project, options = {}
      if api_v2? project
        is_org = project.owner != config.username(project.host)
        params = {'name' => is_org ? project.name_with_owner : project.name}
        params['public'] = '0' if options[:private]
        params['description'] = options[:description] if options[:description]
        params['homepage'] = options[:homepage] if options[:homepage]

        res = post_form "https://%s/api/v2/json/repos/create" % project.host, params
        res.error! unless res.success?
      else
        raise NotImplementedError
      end
    end

    # Public: Fetch info about a pull request.
    def pullrequest_info project, pull_id
      if api_v2? project
        res = get "https://%s/api/v2/json/pulls/%s/%s/%d" %
          [project.host, project.owner, project.name, pull_id]
        res.error! unless res.success?
        res.data['pull']
      else
        raise NotImplementedError
      end
    end

    # Returns parsed data from the new pull request.
    def create_pullrequest options
      project = options.fetch(:project)

      if api_v2? project
        params = {
          'pull[base]' => options.fetch(:base),
          'pull[head]' => options.fetch(:head)
        }
        params['pull[issue]'] = options[:issue] if options[:issue]
        params['pull[title]'] = options[:title] if options[:title]
        params['pull[body]'] = options[:body] if options[:body]

        res = post_form "https://%s/api/v2/json/pulls/%s/%s" %
          [project.host, project.owner, project.name], params
        res.error! unless res.success?
        res.data['pull']
      else
        raise NotImplementedError
      end
    end

    # Methods for performing HTTP requests
    #
    # Requires access to a `config` object that implements `proxy_uri(with_ssl)`
    module HttpMethods
      # Decorator for Net::HTTPResponse
      module ResponseMethods
        def status() code.to_i end
        def data?() content_type =~ /\bjson\b/ end
        def data() @data ||= JSON.parse(body) end
        def error_message?() data? and data['error'] end
        def error_message() data['error'] end
        def success?() Net::HTTPSuccess === self end
      end

      def get url, &block
        perform_request url, :Get, &block
      end

      def post url, params = nil
        perform_request url, :Post do |req|
          if params
            req.body = JSON.dump params
            req['Content-Type'] = 'application/json'
          end
          yield req if block_given?
          req['Content-Length'] = req.body ? req.body.length : 0
        end
      end

      def post_form url, params
        post(url) {|req| req.set_form_data params }
      end

      def perform_request url, type
        url = URI.parse url unless url.respond_to? :hostname

        require 'net/https'
        req = Net::HTTP.const_get(type).new(url.request_uri)
        http = create_connection(url)

        apply_authentication(req, url)
        yield req if block_given?
        res = http.start { http.request(req) }
        res.extend ResponseMethods
        res
      rescue SocketError => err
        raise Context::FatalError, "error with #{type.to_s.upcase} #{url} (#{err.message})"
      end

      def apply_authentication req, url
        user = config.username(url.host)
        token = config.api_token(url.host, user)
        req.basic_auth "#{user}/token", token
      end

      def create_connection url
        use_ssl = 'https' == url.scheme

        proxy_args = []
        if proxy = config.proxy_uri(use_ssl)
          proxy_args << proxy.host << proxy.port
          if proxy.userinfo
            require 'cgi'
            # proxy user + password
            proxy_args.concat proxy.userinfo.split(':', 2).map {|a| CGI.unescape a }
          end
        end

        http = Net::HTTP.new(url.host, url.port, *proxy_args)

        if http.use_ssl = use_ssl
          # FIXME: enable SSL peer verification!
          http.verify_mode = OpenSSL::SSL::VERIFY_NONE
        end
        return http
      end
    end

    include HttpMethods

    # Filesystem store suitable for Configuration
    class FileStore
      extend Forwardable
      def_delegator :@data, :[], :get
      def_delegator :@data, :[]=, :set

      def initialize filename
        @filename = filename
        @data = Hash.new {|d, host| d[host] = [] }
        load if File.exist? filename
      end

      def fetch_user host
        unless entry = get(host).first
          user = yield
          entry = entry_for_user(host, user)
        end
        entry['user']
      end

      def fetch_value host, user, key
        entry = entry_for_user host, user
        entry[key.to_s] || begin
          value = yield
          if value and !value.empty?
            entry[key.to_s] = value
            save
            value
          else
            raise "no value"
          end
        end
      end

      def entry_for_user host, username
        entries = get(host)
        entries.find {|e| e['user'] == username } or
          (entries << {'user' => username}).last
      end

      def load
        @data.update YAML.load(File.read(@filename))
      end

      def save
        File.open(@filename, 'w') {|f| f << YAML.dump(@data) }
      end
    end

    # Provides authentication info per GitHub host such as username and token.
    class Configuration
      def initialize store
        @data = store
        # passwords are cached in memory instead of persistent store
        @password_cache = {}
      end

      def normalize_host host
        host = host.downcase
        'api.github.com' == host ? 'github.com' : host
      end

      def username host
        host = normalize_host host
        @data.fetch_user host do
          if block_given? then yield
          else prompt "#{host} username"
          end
        end
      end

      def api_token host, user
        host = normalize_host host
        @data.fetch_value host, user, :api_token do
          if block_given? then yield
          else prompt "#{host} API token for #{user}"
          end
        end
      end

      def prompt what
        print "#{what}: "
        $stdin.gets.chomp
      end

      def proxy_uri(with_ssl)
        env_name = "HTTP#{with_ssl ? 'S' : ''}_PROXY"
        if proxy = ENV[env_name] || ENV[env_name.downcase]
          proxy = "http://#{proxy}" unless proxy.include? '://'
          URI.parse proxy
        end
      end
    end

    class LegacyConfiguration
      # Options:
      # git - a Hub::Context::GitReader
      def initialize git
        @git = git
      end

      def config_prefix host
        case host.downcase
        when 'api.github.com', 'github.com' then 'github'
        else %(github."#{host.downcase}")
        end
      end

      # Read username for a specific GitHub host from git-config.
      #
      # Yields if not found.
      def username host
        if user = ENV['GITHUB_USER'] and !user.empty?
          user
        else
          config_key = "#{config_prefix host}.user"
          @git.read_config config_key or yield
        end
      end

      # Read the API v2 token for a specific GitHub host from git-config.
      #
      # Yields if not found.
      def api_token host, user
        if token = ENV['GITHUB_TOKEN'] and !token.empty?
          token
        else
          config_key = "#{config_prefix host}.token"
          @git.read_config config_key or yield
        end
      end
    end

    # Wraps multiple Configuration-like objects, tries each one in order and
    # returns first value found.
    class CascadingConfiguration
      def initialize configs
        @configs = Array(configs)
      end

      def method_missing method, *args, &block
        try_config 0, method, args, &block
      end

      def respond_to? method, with_private = false
        !with_private && @configs.any? {|c| c.respond_to? method } or super
      end

      def try_config idx, method, args, &block
        if config = @configs[idx]
          if config.respond_to? method
            if @configs.length >= idx + 2 or block_given?
              config.send(method, *args) {
                try_config idx + 1, method, args, &block
              }
            else
              config.send(method, *args)
            end
          else
            try_config idx + 1, method, args, &block
          end
        else
          yield
        end
      end
    end
  end
end
