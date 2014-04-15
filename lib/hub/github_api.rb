require 'forwardable'

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
  #     GitHubAPI.new file_config, :app_url => 'http://hub.github.com/'
  #   end
  class GitHubAPI
    attr_reader :config, :oauth_app_url

    # Public: Create a new API client instance
    #
    # Options:
    # - config: an object that implements:
    #   - username(host)
    #   - password(host, user)
    #   - oauth_token(host, user)
    def initialize config, options
      @config = config
      @oauth_app_url = options.fetch(:app_url)
      @verbose = options.fetch(:verbose, false)
    end

    def verbose?() @verbose end

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

    def username_via_auth_dance host
      host = api_host(host)
      config.username(host) do
        if block_given?
          yield
        else
          res = get("https://%s/user" % host)
          res.error! unless res.success?
          config.value_to_persist(res.data['login'])
        end
      end
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
      is_org = project.owner.downcase != username_via_auth_dance(project.host).downcase
      params = { :name => project.name, :private => !!options[:private] }
      params[:description] = options[:description] if options[:description]
      params[:homepage]    = options[:homepage]    if options[:homepage]

      if is_org
        res = post "https://%s/orgs/%s/repos" % [api_host(project.host), project.owner], params
      else
        res = post "https://%s/user/repos" % api_host(project.host), params
      end
      res.error! unless res.success?
      res.data
    end

    # Public: Fetch info about a pull request.
    def pullrequest_info project, pull_id
      res = get "https://%s/repos/%s/%s/pulls/%d" %
        [api_host(project.host), project.owner, project.name, pull_id]
      res.error! unless res.success?
      res.data
    end

    # Public: Fetch a pull request's patch
    def pullrequest_patch project, pull_id
      res = get "https://%s/repos/%s/%s/pulls/%d" %
        [api_host(project.host), project.owner, project.name, pull_id] do |req|
          req["Accept"] = "application/vnd.github.v3.patch"
        end
      res.error! unless res.success?
      res.body
    end

    # Public: Fetch the patch from a commit
    def commit_patch project, sha
      res = get "https://%s/repos/%s/%s/commits/%s" %
        [api_host(project.host), project.owner, project.name, sha] do |req|
          req["Accept"] = "application/vnd.github.v3.patch"
        end
      res.error! unless res.success?
      res.body
    end

    # Public: Fetch the first raw blob from a gist
    def gist_raw gist_id
      res = get("https://%s/gists/%s" % [api_host('github.com'), gist_id])
      res.error! unless res.success?
      raw_url = res.data['files'].values.first['raw_url']
      res = get(raw_url) do |req|
        req['Accept'] = 'text/plain'
      end
      res.error! unless res.success?
      res.body
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

    def statuses project, sha
      res = get "https://%s/repos/%s/%s/statuses/%s" %
        [api_host(project.host), project.owner, project.name, sha]

      res.error! unless res.success?
      res.data
    end

    # Methods for performing HTTP requests
    #
    # Requires access to a `config` object that implements:
    # - proxy_uri(with_ssl)
    # - username(host)
    # - password(host, user)
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
            when 'missing_field'
              %(Missing field: "%s") % err['field']
            when 'already_exists'
              %(Duplicate value for "%s") % err['field']
            when 'invalid'
              %(Invalid value for "%s": "%s") % [ err['field'], err['value'] ]
            when 'unauthorized'
              %(Not allowed to change field "%s") % err['field']
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
            req['Content-Type'] = 'application/json;charset=utf-8'
          end
          yield req if block_given?
          req['Content-Length'] = byte_size req.body
        end
      end

      def byte_size str
        if    str.respond_to? :bytesize then str.bytesize
        elsif str.respond_to? :length   then str.length
        else  0
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

        req['User-Agent'] = "Hub #{Hub::VERSION}"
        apply_authentication(req, url)
        yield req if block_given?
        finalize_request(req, url)

        begin
          res = http.start { http.request(req) }
          res.extend ResponseMethods
          return res
        rescue SocketError => err
          raise Context::FatalError, "error with #{type.to_s.upcase} #{url} (#{err.message})"
        end
      end

      def request_uri url
        str = url.request_uri
        str = '/api/v3' << str if url.host != 'api.github.com' && url.host != 'gist.github.com'
        str
      end

      def configure_connection req, url
        url.scheme = config.protocol(url.host)
        if ENV['HUB_TEST_HOST']
          req['Host'] = url.host
          req['X-Original-Scheme'] = url.scheme
          url = url.dup
          url.scheme = 'http'
          url.host, test_port = ENV['HUB_TEST_HOST'].split(':')
          url.port = test_port.to_i if test_port
        end
        yield url
      end

      def apply_authentication req, url
        user = url.user ? CGI.unescape(url.user) : config.username(url.host)
        pass = config.password(url.host, user)
        req.basic_auth user, pass
      end

      def finalize_request(req, url)
        if !req['Accept'] || req['Accept'] == '*/*'
          req['Accept'] = 'application/vnd.github.v3+json'
        end
      end

      def create_connection url
        use_ssl = 'https' == url.scheme

        proxy_args = []
        if proxy = config.proxy_uri(use_ssl)
          proxy_args << proxy.host << proxy.port
          if proxy.userinfo
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
        if req.path =~ %r{^(/api/v3)?/authorizations$}
          super
        else
          user = url.user ? CGI.unescape(url.user) : config.username(url.host)
          token = config.oauth_token(url.host, user) {
            obtain_oauth_token url.host, user
          }
          req['Authorization'] = "token #{token}"
        end
      end

      def obtain_oauth_token host, user, two_factor_code = nil
        auth_url = URI.parse("https://%s@%s/authorizations" % [CGI.escape(user), host])

        # dummy request to trigger a 2FA SMS since a HTTP GET won't do it
        post(auth_url) if !two_factor_code

        # first try to fetch existing authorization
        res = get(auth_url) do |req|
          req['X-GitHub-OTP'] = two_factor_code if two_factor_code
        end
        unless res.success?
          if !two_factor_code && res['X-GitHub-OTP'].to_s.include?('required')
            two_factor_code = config.prompt_auth_code
            return obtain_oauth_token(host, user, two_factor_code)
          else
            res.error!
          end
        end

        if found = res.data.find {|auth| auth['note'] == 'hub' || auth['note_url'] == oauth_app_url }
          found['token']
        else
          # create a new authorization
          res = post auth_url,
            :scopes => %w[repo], :note => 'hub', :note_url => oauth_app_url do |req|
              req['X-GitHub-OTP'] = two_factor_code if two_factor_code
            end
          res.error! unless res.success?
          res.data['token']
        end
      end
    end

    module GistAuth
      def apply_authentication(req, url)
        super unless url.host == 'gist.github.com'
      end
    end

    module Verbose
      def finalize_request(req, url)
        super
        dump_request_info(req, url) if verbose?
      end

      def perform_request(*)
        res = super
        dump_response_info(res) if verbose?
        res
      end

      def verbose_puts(msg)
        msg = "\e[36m%s\e[m" % msg if $stderr.tty?
        $stderr.puts msg
      end

      def dump_request_info(req, url)
        verbose_puts "> %s %s://%s%s" % [
          req.method.to_s.upcase,
          url.scheme,
          url.host,
          req.path,
        ]
        dump_headers(req, '> ')
        dump_body(req)
      end

      def dump_response_info(res)
        verbose_puts "< HTTP %s" % res.status
        dump_headers(res, '< ')
        dump_body(res)
      end

      def dump_body(obj)
        verbose_puts obj.body if obj.body
      end

      DUMP_HEADERS = %w[ Authorization X-GitHub-OTP Location ]

      def dump_headers(obj, indent)
        DUMP_HEADERS.each do |header|
          if value = obj[header]
            verbose_puts '%s%s: %s' % [
              indent,
              header,
              value.sub(/^(basic|token) (.+)/i, '\1 [REDACTED]'),
            ]
          end
        end
      end
    end

    include HttpMethods
    include OAuth
    include GistAuth
    include Verbose

    # Filesystem store suitable for Configuration
    class FileStore
      extend Forwardable
      def_delegator :@data, :[], :get
      def_delegator :@data, :[]=, :set

      def initialize filename
        @filename = filename
        @data = Hash.new {|d, host| d[host] = [] }
        @persist_next_change = false
        load if File.exist? filename
      end

      def fetch_value host, user, key
        entries = get(host)
        entries << {} if entries.empty?
        entry = entries.first
        entry.fetch(key.to_s) {
          value = yield
          raise "no value for key :#{key}" if value.nil? || value.empty?
          entry[key.to_s] = value
          save_if_needed
          value
        }
      end

      def persist_next_change!
        @persist_next_change = true
      end

      def save_if_needed
        @persist_next_change && save
        @persist_next_change = false
      end

      def load
        existing_data = File.read(@filename)
        @data.update yaml_load(existing_data) unless existing_data.strip.empty?
      end

      def save
        mkdir_p File.dirname(@filename)
        File.open(@filename, 'w', 0600) {|f| f << yaml_dump(@data) }
      end

      def mkdir_p(dir)
        dir.split('/').inject do |parent, name|
          d = File.join(parent, name)
          Dir.mkdir(d) unless File.exist?(d)
          d
        end
      end

      def yaml_load(string)
        hash = {}
        host = nil
        string.split("\n").each do |line|
          case line
          when /^---\s*$/, /^\s*(?:#|$)/
            # ignore
          when /^(.+):\s*$/
            host = hash[$1] = []
          when /^([- ]) (.+?): (.+)/
            key, value = $2, $3
            host << {} if $1 == '-'
            host.last[key] = value.gsub(/^'|'$/, '')
          else
            raise "unsupported YAML line: #{line}"
          end
        end
        hash
      end

      def yaml_dump(data)
        yaml = ['---']
        data.each do |host, values|
          yaml << "#{host}:"
          values.each do |hash|
            dash = '-'
            hash.each do |key, value|
              yaml << "#{dash} #{key}: #{value}"
              dash = ' '
            end
          end
        end
        yaml.join("\n")
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
        return ENV['GITHUB_USER'] unless ENV['GITHUB_USER'].to_s.empty?
        host = normalize_host host
        @data.fetch_value(host, nil, :user) do
          if block_given? then yield
          else prompt "#{host} username"
          end
        end
      end

      def password host, user
        return ENV['GITHUB_PASSWORD'] unless ENV['GITHUB_PASSWORD'].to_s.empty?
        host = normalize_host host
        @password_cache["#{user}@#{host}"] ||= prompt_password host, user
      end

      def oauth_token host, user
        host = normalize_host(host)
        @data.fetch_value(host, user, :oauth_token) do
          value_to_persist(yield)
        end
      end

      def protocol host
        host = normalize_host host
        @data.fetch_value(host, nil, :protocol) { 'https' }
      end

      def value_to_persist(value = nil)
        @data.persist_next_change!
        value
      end

      def prompt what
        print "#{what}: "
        $stdin.gets.chomp
      rescue Interrupt
        puts
        abort
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
      rescue Interrupt
        puts
        abort
      end

      def prompt_auth_code
        print "two-factor authentication code: "
        $stdin.gets.chomp
      rescue Interrupt
        puts
        abort
      end

      NULL = defined?(File::NULL) ? File::NULL :
               File.exist?('/dev/null') ? '/dev/null' : 'NUL'

      def askpass
        noecho $stdin do |input|
          input.gets.chomp
        end
      end

      def noecho io
        require 'io/console'
        io.noecho { yield io }
      rescue LoadError
        fallback_noecho io
      end

      def fallback_noecho io
        tty_state = `stty -g 2>#{NULL}`
        system 'stty raw -echo -icanon isig' if $?.success?
        pass = ''
        while char = getbyte(io) and !(char == 13 or char == 10)
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

      def getbyte(io)
        if io.respond_to?(:getbyte)
          io.getbyte
        else
          # In Ruby <= 1.8.6, getc behaved the same
          io.getc
        end
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
