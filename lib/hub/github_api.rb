require 'uri'
require 'yaml'
require 'forwardable'
require 'fileutils'

module Hub
  # Client for the GitHub v3 API.
  #
  # First time around, user gets prompted for username/password in the shell.
  # Then this information is exchanged for an OAuth token which is saved in a file.
  #
  # Examples
  #
  #   @api_client ||= begin
  #     config_file = ENV['HUB_CONFIG'] || '~/.config/hub'
  #     file_store = GitHubAPI::FileStore.new File.expand_path(config_file)
  #     file_config = GitHubAPI::Configuration.new file_store
  #     GitHubAPI.new file_config, :app_url => 'http://defunkt.io/hub/'
  #   end
  class GitHubAPI
    attr_reader :config, :oauth_app_url

    # Public: Create a new API client instance
    #
    # Options:
    # - config: an object that implements:
    #   - username(host)
    #   - api_token(host, user)
    #   - password(host, user)
    #   - oauth_token(host, user)
    def initialize config, options
      @config = config
      @oauth_app_url = options.fetch(:app_url)
    end

    # Fake exception type for net/http exception handling.
    # Necessary because net/http may or may not be loaded at the time.
    module Exceptions
      def self.===(exception)
        exception.class.ancestors.map {|a| a.to_s }.include? 'Net::HTTPExceptions'
      end
    end

    def api_host host
      host = host.downcase
      'github.com' == host ? 'api.github.com' : host
    end

    # Public: Fetch data for a specific repo.
    def repo_info project
      get "https://%s/repos/%s/%s" %
        [api_host(project.host), project.owner, project.name]
    end

    # Public: Determine whether a specific repo exists.
    def repo_exists? project
      repo_info(project).success?
    end

    # Public: Fork the specified repo.
    def fork_repo project
      res = post "https://%s/repos/%s/%s/forks" %
        [api_host(project.host), project.owner, project.name]
      res.error! unless res.success?
    end

    # Public: Create a new project.
    def create_repo project, options = {}
      is_org = project.owner != config.username(api_host(project.host))
      params = { :name => project.name, :private => !!options[:private] }
      params[:description] = options[:description] if options[:description]
      params[:homepage]    = options[:homepage]    if options[:homepage]

      if is_org
        res = post "https://%s/orgs/%s/repos" % [api_host(project.host), project.owner], params
      else
        res = post "https://%s/user/repos" % api_host(project.host), params
      end
      res.error! unless res.success?
    end

    # Public: Fetch info about a pull request.
    def pullrequest_info project, pull_id
      res = get "https://%s/repos/%s/%s/pulls/%d" %
        [api_host(project.host), project.owner, project.name, pull_id]
      res.error! unless res.success?
      res.data
    end

    # Returns parsed data from the new pull request.
    def create_pullrequest options
      project = options.fetch(:project)
      params = {
        :base => options.fetch(:base),
        :head => options.fetch(:head)
      }

      if options[:issue]
        params[:issue] = options[:issue]
      else
        params[:title] = options[:title] if options[:title]
        params[:body]  = options[:body]  if options[:body]
      end

      res = post "https://%s/repos/%s/%s/pulls" %
        [api_host(project.host), project.owner, project.name], params

      res.error! unless res.success?
      res.data
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
        def error_message?() data? and data['errors'] || data['message'] end
        def error_message() error_sentences || data['message'] end
        def success?() Net::HTTPSuccess === self end
        def error_sentences
          data['errors'].map do |err|
            case err['code']
            when 'custom'        then err['message']
            when 'missing_field' then "field '%s' is missing" % err['field']
            end
          end.compact if data['errors']
        end
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
        url = URI.parse url unless url.respond_to? :host

        require 'net/https'
        req = Net::HTTP.const_get(type).new request_uri(url)
        # TODO: better naming?
        http = configure_connection(req, url) do |host_url|
          create_connection host_url
        end

        apply_authentication(req, url)
        yield req if block_given?
        res = http.start { http.request(req) }
        res.extend ResponseMethods
        res
      rescue SocketError => err
        raise Context::FatalError, "error with #{type.to_s.upcase} #{url} (#{err.message})"
      end

      def request_uri url
        str = url.request_uri
        str = '/api/v3' << str if url.host != 'api.github.com'
        str
      end

      def configure_connection req, url
        if ENV['HUB_TEST_HOST']
          req['Host'] = url.host
          url = url.dup
          url.scheme = 'http'
          url.host, test_port = ENV['HUB_TEST_HOST'].split(':')
          url.port = test_port.to_i if test_port
        end
        yield url
      end

      def apply_authentication req, url
        user = url.user || config.username(url.host)
        pass = config.password(url.host, user)
        req.basic_auth user, pass
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

    module OAuth
      def apply_authentication req, url
        if (req.path =~ /\/authorizations$/)
          super
        else
          user = url.user || config.username(url.host)
          token = config.oauth_token(url.host, user) {
            obtain_oauth_token url.host, user
          }
          req['Authorization'] = "token #{token}"
        end
      end

      def obtain_oauth_token host, user
        # first try to fetch existing authorization
        res = get "https://#{user}@#{host}/authorizations"
        res.error! unless res.success?

        if found = res.data.find {|auth| auth['app']['url'] == oauth_app_url }
          found['token']
        else
          # create a new authorization
          res = post "https://#{user}@#{host}/authorizations",
            :scopes => %w[repo], :note => 'hub', :note_url => oauth_app_url
          res.error! unless res.success?
          res.data['token']
        end
      end
    end

    include HttpMethods
    include OAuth

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
          # FIXME: more elegant handling of empty strings
          return nil if user.nil? or user.empty?
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
        existing_data = File.read(@filename)
        @data.update YAML.load(existing_data) unless existing_data.strip.empty?
      end

      def save
        FileUtils.mkdir_p File.dirname(@filename)
        File.open(@filename, 'w', 0600) {|f| f << YAML.dump(@data) }
      end
    end

    # Provides authentication info per GitHub host such as username, password,
    # and API/OAuth tokens.
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

      def password host, user
        host = normalize_host host
        @password_cache["#{user}@#{host}"] ||= prompt_password host, user
      end

      def oauth_token host, user, &block
        @data.fetch_value normalize_host(host), user, :oauth_token, &block
      end

      def prompt what
        print "#{what}: "
        $stdin.gets.chomp
      end

      # special prompt that has hidden input
      def prompt_password host, user
        print "#{host} password for #{user} (never stored): "
        if $stdin.tty?
          password = askpass
          puts ''
          password
        else
          # in testing
          $stdin.gets.chomp
        end
      end

      # FIXME: probably not cross-platform
      def askpass
        tty_state = `stty -g`
        system 'stty raw -echo -icanon isig' if $?.success?
        pass = ''
        while char = $stdin.getbyte and !(char == 13 or char == 10)
          if char == 127 or char == 8
            pass[-1,1] = '' unless pass.empty?
          else
            pass << char.chr
          end
        end
        pass
      ensure
        system "stty #{tty_state}" unless tty_state.empty?
      end

      def proxy_uri(with_ssl)
        env_name = "HTTP#{with_ssl ? 'S' : ''}_PROXY"
        if proxy = ENV[env_name] || ENV[env_name.downcase] and !proxy.empty?
          proxy = "http://#{proxy}" unless proxy.include? '://'
          URI.parse proxy
        end
      end
    end
  end
end
