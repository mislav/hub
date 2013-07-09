#!/usr/bin/env ruby
#
# This file is generated code. DO NOT send patches for it.
#
# Original source files with comments are at:
# https://github.com/defunkt/hub
#

module Hub
  Version = VERSION = '1.10.6'
end

module Hub
  class Args < Array
    attr_accessor :executable

    def initialize(*args)
      super
      @executable = ENV["GIT"] || "git"
      @skip = @noop = false
      @original_args = args.first
      @chain = [nil]
    end

    def after(cmd_or_args = nil, args = nil, &block)
      @chain.insert(-1, normalize_callback(cmd_or_args, args, block))
    end

    def before(cmd_or_args = nil, args = nil, &block)
      @chain.insert(@chain.index(nil), normalize_callback(cmd_or_args, args, block))
    end

    def chained?
      @chain.size > 1
    end

    def commands
      chain = @chain.dup
      chain[chain.index(nil)] = self.to_exec
      chain
    end

    def skip!
      @skip = true
    end

    def skip?
      @skip
    end

    def noop!
      @noop = true
    end

    def noop?
      @noop
    end

    def to_exec(args = self)
      Array(executable) + args
    end

    def add_exec_flags(flags)
      self.executable = Array(executable).concat(flags)
    end

    def words
      reject { |arg| arg.index('-') == 0 }
    end

    def flags
      self - words
    end

    def changed?
      chained? or self != @original_args
    end

    def has_flag?(*flags)
      pattern = flags.flatten.map { |f| Regexp.escape(f) }.join('|')
      !grep(/^#{pattern}(?:=|$)/).empty?
    end

    private

    def normalize_callback(cmd_or_args, args, block)
      if block
        block
      elsif args
        [cmd_or_args].concat args
      elsif Array === cmd_or_args
        self.to_exec cmd_or_args
      elsif cmd_or_args
        cmd_or_args
      else
        raise ArgumentError, "command or block required"
      end
    end
  end
end

module Hub
  class SshConfig
    CONFIG_FILES = %w(~/.ssh/config /etc/ssh_config /etc/ssh/ssh_config)

    def initialize files = nil
      @settings = Hash.new {|h,k| h[k] = {} }
      Array(files || CONFIG_FILES).each do |path|
        file = File.expand_path path
        parse_file file if File.exist? file
      end
    end

    def get_value hostname, key
      key = key.to_s.downcase
      @settings.each do |pattern, settings|
        if pattern.match? hostname and found = settings[key]
          return found
        end
      end
      yield
    end

    class HostPattern
      def initialize pattern
        @pattern = pattern.to_s.downcase
      end

      def to_s() @pattern end
      def ==(other) other.to_s == self.to_s end

      def matcher
        @matcher ||=
          if '*' == @pattern
            Proc.new { true }
          elsif @pattern !~ /[?*]/
            lambda { |hostname| hostname.to_s.downcase == @pattern }
          else
            re = self.class.pattern_to_regexp @pattern
            lambda { |hostname| re =~ hostname }
          end
      end

      def match? hostname
        matcher.call hostname
      end

      def self.pattern_to_regexp pattern
        escaped = Regexp.escape(pattern)
        escaped.gsub!('\*', '.*')
        escaped.gsub!('\?', '.')
        /^#{escaped}$/i
      end
    end

    def parse_file file
      host_patterns = [HostPattern.new('*')]

      IO.foreach(file) do |line|
        case line
        when /^\s*(#|$)/ then next
        when /^\s*(\S+)\s*=/
          key, value = $1, $'
        else
          key, value = line.strip.split(/\s+/, 2)
        end

        next if value.nil?
        key.downcase!
        value = $1 if value =~ /^"(.*)"$/
        value.chomp!

        if 'host' == key
          host_patterns = value.split(/\s+/).map {|p| HostPattern.new p }
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

require 'uri'
require 'yaml'
require 'forwardable'
require 'fileutils'

module Hub
  class GitHubAPI
    attr_reader :config, :oauth_app_url

    def initialize config, options
      @config = config
      @oauth_app_url = options.fetch(:app_url)
    end

    module Exceptions
      def self.===(exception)
        exception.class.ancestors.map {|a| a.to_s }.include? 'Net::HTTPExceptions'
      end
    end

    def api_host host
      host = host.downcase
      'github.com' == host ? 'api.github.com' : host
    end

    def repo_info project
      get "https://%s/repos/%s/%s" %
        [api_host(project.host), project.owner, project.name]
    end

    def repo_exists? project
      repo_info(project).success?
    end

    def fork_repo project
      res = post "https://%s/repos/%s/%s/forks" %
        [api_host(project.host), project.owner, project.name]
      res.error! unless res.success?
    end

    def create_repo project, options = {}
      is_org = project.owner.downcase != config.username(api_host(project.host)).downcase
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

    def pullrequest_info project, pull_id
      res = get "https://%s/repos/%s/%s/pulls/%d" %
        [api_host(project.host), project.owner, project.name, pull_id]
      res.error! unless res.success?
      res.data
    end

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

    module HttpMethods
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
        http = configure_connection(req, url) do |host_url|
          create_connection host_url
        end

        req['User-Agent'] = "Hub #{Hub::VERSION}"
        apply_authentication(req, url)
        yield req if block_given?

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
            proxy_args.concat proxy.userinfo.split(':', 2).map {|a| CGI.unescape a }
          end
        end

        http = Net::HTTP.new(url.host, url.port, *proxy_args)

        if http.use_ssl = use_ssl
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
          refresh = false
          user = url.user || config.username(url.host)
          token = config.oauth_token(url.host, user) {
            refresh = true
            obtain_oauth_token url.host, user
          }
          if refresh
            res = get "https://#{url.host}/user"
            res.error! unless res.success?
            config.update_username(url.host, user, res.data['login'])
          end
          req['Authorization'] = "token #{token}"
        end
      end

      def obtain_oauth_token host, user
        res = get "https://#{user}@#{host}/authorizations"
        res.error! unless res.success?

        if found = res.data.find {|auth| auth['app']['url'] == oauth_app_url }
          found['token']
        else
          res = post "https://#{user}@#{host}/authorizations",
            :scopes => %w[repo], :note => 'hub', :note_url => oauth_app_url
          res.error! unless res.success?
          res.data['token']
        end
      end
    end

    include HttpMethods
    include OAuth

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

    class Configuration
      def initialize store
        @data = store
        @password_cache = {}
      end

      def normalize_host host
        host = host.downcase
        'api.github.com' == host ? 'github.com' : host
      end

      def username host
        return ENV['GITHUB_USER'] unless ENV['GITHUB_USER'].to_s.empty?
        host = normalize_host host
        @data.fetch_user host do
          if block_given? then yield
          else prompt "#{host} username"
          end
        end
      end

      def update_username host, old_username, new_username
        entry = @data.entry_for_user(normalize_host(host), old_username)
        entry['user'] = new_username
        @data.save
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
        return ENV['GITHUB_PASSWORD'] unless ENV['GITHUB_PASSWORD'].to_s.empty?
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

      def prompt_password host, user
        print "#{host} password for #{user} (never stored): "
        if $stdin.tty?
          password = askpass
          puts ''
          password
        else
          $stdin.gets.chomp
        end
      end

      NULL = defined?(File::NULL) ? File::NULL :
               File.exist?('/dev/null') ? '/dev/null' : 'NUL'

      def askpass
        tty_state = `stty -g 2>#{NULL}`
        system 'stty raw -echo -icanon isig' if $?.success?
        pass = ''
        while char = getbyte($stdin) and !(char == 13 or char == 10)
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

require 'shellwords'
require 'forwardable'
require 'uri'

module Hub
  module Context
    extend Forwardable

    NULL = defined?(File::NULL) ? File::NULL : File.exist?('/dev/null') ? '/dev/null' : 'NUL'

    class GitReader
      attr_reader :executable

      def initialize(executable = nil, &read_proc)
        @executable = executable || 'git'
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

    module GitReaderMethods
      extend Forwardable

      def_delegator :git_reader, :read_config, :git_config
      def_delegator :git_reader, :read, :git_command

      def self.extended(base)
        base.extend Forwardable
        base.def_delegators :'self.class', :git_config, :git_command
      end
    end

    class Error < RuntimeError; end
    class FatalError < Error; end

    private

    def git_reader
      @git_reader ||= GitReader.new ENV['GIT']
    end

    include GitReaderMethods
    private :git_config, :git_command

    def local_repo(fatal = true)
      @local_repo ||= begin
        if is_repo?
          LocalRepo.new git_reader, current_dir
        elsif fatal
          raise FatalError, "Not a git repository"
        end
      end
    end

    repo_methods = [
      :current_branch,
      :current_project, :upstream_project,
      :repo_owner, :repo_host,
      :remotes, :remotes_group, :origin_remote
    ]
    def_delegator :local_repo, :name, :repo_name
    def_delegators :local_repo, *repo_methods
    private :repo_name, *repo_methods

    def master_branch
      if local_repo(false)
        local_repo.master_branch
      else
        Branch.new nil, 'refs/heads/master'
      end
    end

    class LocalRepo < Struct.new(:git_reader, :dir)
      include GitReaderMethods

      def name
        if project = main_project
          project.name
        else
          File.basename(dir)
        end
      end

      def repo_owner
        if project = main_project
          project.owner
        end
      end

      def repo_host
        project = main_project and project.host
      end

      def main_project
        remote = origin_remote and remote.project
      end

      def upstream_project
        if branch = current_branch and upstream = branch.upstream and upstream.remote?
          remote = remote_by_name upstream.remote_name
          remote.project
        end
      end

      def current_project
        upstream_project || main_project
      end

      def current_branch
        if branch = git_command('symbolic-ref -q HEAD')
          Branch.new self, branch
        end
      end

      def master_branch
        Branch.new self, 'refs/heads/master'
      end

      def remotes
        @remotes ||= begin
          list = git_command('remote').to_s.split("\n")
          main = list.delete('origin') and list.unshift(main)
          list.map { |name| Remote.new self, name }
        end
      end

      def remotes_group(name)
        git_config "remotes.#{name}"
      end

      def origin_remote
        remotes.first
      end

      def remote_by_name(remote_name)
        remotes.find {|r| r.name == remote_name }
      end

      def known_hosts
        hosts = git_config('hub.host', :all).to_s.split("\n")
        hosts << default_host
        hosts << "ssh.#{default_host}"
      end

      def self.default_host
        ENV['GITHUB_HOST'] || main_host
      end

      def self.main_host
        'github.com'
      end

      extend Forwardable
      def_delegators :'self.class', :default_host, :main_host

      def ssh_config
        @ssh_config ||= SshConfig.new
      end
    end

    class GithubProject < Struct.new(:local_repo, :owner, :name, :host)
      def self.from_url(url, local_repo)
        if local_repo.known_hosts.include? url.host
          _, owner, name = url.path.split('/', 4)
          GithubProject.new(local_repo, owner, name.sub(/\.git$/, ''), url.host)
        end
      end

      attr_accessor :repo_data

      def initialize(*args)
        super
        self.name = self.name.tr(' ', '-')
        self.host ||= (local_repo || LocalRepo).default_host
        self.host = host.sub(/^ssh\./i, '') if 'ssh.github.com' == host.downcase
      end

      def private?
        repo_data ? repo_data.fetch('private') :
          host != (local_repo || LocalRepo).main_host
      end

      def owned_by(new_owner)
        new_project = dup
        new_project.owner = new_owner
        new_project
      end

      def name_with_owner
        "#{owner}/#{name}"
      end

      def ==(other)
        name_with_owner == other.name_with_owner
      end

      def remote
        local_repo.remotes.find { |r| r.project == self }
      end

      def web_url(path = nil)
        project_name = name_with_owner
        if project_name.sub!(/\.wiki$/, '')
          unless '/wiki' == path
            path = if path =~ %r{^/commits/} then '/_history'
                   else path.to_s.sub(/\w+/, '_\0')
                   end
            path = '/wiki' + path
          end
        end
        "https://#{host}/" + project_name + path.to_s
      end

      def git_url(options = {})
        if options[:https] then "https://#{host}/"
        elsif options[:private] or private? then "git@#{host}:"
        else "git://#{host}/"
        end + name_with_owner + '.git'
      end
    end

    class GithubURL < URI::HTTPS
      extend Forwardable

      attr_reader :project
      def_delegator :project, :name, :project_name
      def_delegator :project, :owner, :project_owner

      def self.resolve(url, local_repo)
        u = URI(url)
        if %[http https].include? u.scheme and project = GithubProject.from_url(u, local_repo)
          self.new(u.scheme, u.userinfo, u.host, u.port, u.registry,
                   u.path, u.opaque, u.query, u.fragment, project)
        end
      rescue URI::InvalidURIError
        nil
      end

      def initialize(*args)
        @project = args.pop
        super(*args)
      end

      def project_path
        path.split('/', 4)[3]
      end
    end

    class Branch < Struct.new(:local_repo, :name)
      alias to_s name

      def short_name
        name.sub(%r{^refs/(remotes/)?.+?/}, '')
      end

      def master?
        short_name == 'master'
      end

      def upstream
        if branch = local_repo.git_command("rev-parse --symbolic-full-name #{short_name}@{upstream}")
          Branch.new local_repo, branch
        end
      end

      def remote?
        name.index('refs/remotes/') == 0
      end

      def remote_name
        name =~ %r{^refs/remotes/([^/]+)} and $1 or
          raise Error, "can't get remote name from #{name.inspect}"
      end
    end

    class Remote < Struct.new(:local_repo, :name)
      alias to_s name

      def ==(other)
        other.respond_to?(:to_str) ? name == other.to_str : super
      end

      def project
        urls.each_value { |url|
          if valid = GithubProject.from_url(url, local_repo)
            return valid
          end
        }
        nil
      end

      def urls
        return @urls if defined? @urls
        @urls = {}
        local_repo.git_command('remote -v').to_s.split("\n").map do |line|
          next if line !~ /^(.+?)\t(.+) \((.+)\)$/
          remote, uri, type = $1, $2, $3
          next if remote != self.name
          if uri =~ %r{^[\w-]+://} or uri =~ %r{^([^/]+?):}
            uri = "ssh://#{$1}/#{$'}" if $1
            begin
              @urls[type] = uri_parse(uri)
            rescue URI::InvalidURIError
            end
          end
        end
        @urls
      end

      def uri_parse uri
        uri = URI.parse uri
        uri.host = local_repo.ssh_config.get_value(uri.host, 'hostname') { uri.host }
        uri.user = local_repo.ssh_config.get_value(uri.host, 'user') { uri.user }
        uri
      end
    end


    def github_project(name, owner = nil)
      if owner and owner.index('/')
        owner, name = owner.split('/', 2)
      elsif name and name.index('/')
        owner, name = name.split('/', 2)
      else
        name ||= repo_name
        owner ||= github_user
      end

      if local_repo(false) and main_project = local_repo.main_project
        project = main_project.dup
        project.owner = owner
        project.name = name
        project
      else
        GithubProject.new(local_repo(false), owner, name)
      end
    end

    def git_url(owner = nil, name = nil, options = {})
      project = github_project(name, owner)
      project.git_url({:https => https_protocol?}.update(options))
    end

    def resolve_github_url(url)
      GithubURL.resolve(url, local_repo) if url =~ /^https?:/
    end

    def http_clone?
      git_config('--bool hub.http-clone') == 'true'
    end

    def https_protocol?
      git_config('hub.protocol') == 'https' or http_clone?
    end

    def git_alias_for(name)
      git_config "alias.#{name}"
    end

    def rev_list(a, b)
      git_command("rev-list --cherry-pick --right-only --no-merges #{a}...#{b}")
    end

    PWD = Dir.pwd

    def current_dir
      PWD
    end

    def git_dir
      git_command 'rev-parse -q --git-dir'
    end

    def is_repo?
      !!git_dir
    end

    def git_editor
      editor = git_command 'var GIT_EDITOR'
      editor = ENV[$1] if editor =~ /^\$(\w+)$/
      editor = File.expand_path editor if (editor =~ /^[~.]/ or editor.index('/')) and editor !~ /["']/
      if File.exist? editor then [editor]
      else editor.shellsplit
      end
    end

    module System
      def browser_launcher
        browser = ENV['BROWSER'] || (
          osx? ? 'open' : windows? ? %w[cmd /c start] :
          %w[xdg-open cygstart x-www-browser firefox opera mozilla netscape].find { |comm| which comm }
        )

        abort "Please set $BROWSER to a web launcher to use this command." unless browser
        Array(browser)
      end

      def osx?
        require 'rbconfig'
        RbConfig::CONFIG['host_os'].to_s.include?('darwin')
      end

      def windows?
        require 'rbconfig'
        RbConfig::CONFIG['host_os'] =~ /msdos|mswin|djgpp|mingw|windows/
      end

      def which(cmd)
        exts = ENV['PATHEXT'] ? ENV['PATHEXT'].split(';') : ['']
        ENV['PATH'].split(File::PATH_SEPARATOR).each do |path|
          exts.each { |ext|
            exe = "#{path}/#{cmd}#{ext}"
            return exe if File.executable? exe
          }
        end
        return nil
      end

      def command?(name)
        !which(name).nil?
      end

      def tmp_dir
        ENV['TMPDIR'] || ENV['TEMP'] || '/tmp'
      end
    end

    include System
    extend System
  end
end

require 'strscan'
require 'forwardable'

class Hub::JSON
  def self.parse(data) new(data).parse end

  WSP = /\s+/
  OBJ = /[{\[]/;    HEN = /\}/;  AEN = /\]/
  COL = /\s*:\s*/;  KEY = /\s*,\s*/
  NUM = /-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?/
  BOL = /true|false/;  NUL = /null/

  extend Forwardable

  attr_reader :scanner
  alias_method :s, :scanner
  def_delegators :scanner, :scan, :matched
  private :s, :scan, :matched

  def initialize data
    @scanner = StringScanner.new data.to_s
  end

  def parse
    space
    object
  end

  private

  def space() scan WSP end

  def endkey() scan(KEY) or space end

  def object
    matched == '{' ? hash : array if scan(OBJ)
  end

  def value
    object or string or
      scan(NUL) ? nil :
      scan(BOL) ? matched.size == 4:
      scan(NUM) ? eval(matched) :
      error
  end

  def hash
    obj = {}
    space
    repeat_until(HEN) { k = string; scan(COL); obj[k] = value; endkey }
    obj
  end

  def array
    ary = []
    space
    repeat_until(AEN) { ary << value; endkey }
    ary
  end

  SPEC = {'b' => "\b", 'f' => "\f", 'n' => "\n", 'r' => "\r", 't' => "\t"}
  UNI = 'u'; CODE = /[a-fA-F0-9]{4}/
  STR = /"/; STE = '"'
  ESC = '\\'

  def string
    if scan(STR)
      str, esc = '', false
      while c = s.getch
        if esc
          str << (c == UNI ? (s.scan(CODE) || error).to_i(16).chr : SPEC[c] || c)
          esc = false
        else
          case c
          when ESC then esc = true
          when STE then break
          else str << c
          end
        end
      end
      str
    end
  end

  def error
    raise "parse error at: #{scan(/.{1,10}/m).inspect}"
  end

  def repeat_until reg
    until scan(reg)
      pos = s.pos
      yield
      error unless s.pos > pos
    end
  end

  module Generator
    def generate(obj)
      raise ArgumentError unless obj.is_a? Array or obj.is_a? Hash
      generate_type(obj)
    end
    alias dump generate

    private

    def generate_type(obj)
      type = obj.is_a?(Numeric) ? :Numeric : obj.class.name
      begin send(:"generate_#{type}", obj)
      rescue NoMethodError; raise ArgumentError, "can't serialize #{type}"
      end
    end

    ESC_MAP = Hash.new {|h,k| k }.update \
      "\r" => 'r',
      "\n" => 'n',
      "\f" => 'f',
      "\t" => 't',
      "\b" => 'b'

    def generate_String(str)
      escaped = str.gsub(/[\r\n\f\t\b"\\]/) { "\\#{ESC_MAP[$&]}"}
      %("#{escaped}")
    end

    def generate_simple(obj) obj.inspect end
    alias generate_Numeric generate_simple
    alias generate_TrueClass generate_simple
    alias generate_FalseClass generate_simple

    def generate_Symbol(sym) generate_String(sym.to_s) end

    def generate_NilClass(*) 'null' end

    def generate_Array(ary) '[%s]' % ary.map {|o| generate_type(o) }.join(', ') end

    def generate_Hash(hash)
      '{%s}' % hash.map { |key, value|
        "#{generate_String(key.to_s)}: #{generate_type(value)}"
      }.join(', ')
    end
  end

  extend Generator
end

module Hub
  module Commands
    instance_methods.each { |m| undef_method(m) unless m =~ /(^__|send|to\?$)/ }
    extend self

    extend Context

    NAME_RE = /[\w.][\w.-]*/
    OWNER_RE = /[a-zA-Z0-9][a-zA-Z0-9-]*/
    NAME_WITH_OWNER_RE = /^(?:#{NAME_RE}|#{OWNER_RE}\/#{NAME_RE})$/

    CUSTOM_COMMANDS = %w[alias create browse compare fork pull-request ci-status]

    def run(args)
      slurp_global_flags(args)

      args.unshift 'help' if args.empty?

      cmd = args[0]
      if expanded_args = expand_alias(cmd)
        cmd = expanded_args[0]
        expanded_args.concat args[1..-1]
      end

      respect_help_flags(expanded_args || args) if custom_command? cmd

      cmd = cmd.gsub(/(\w)-/, '\1_')
      if method_defined?(cmd) and cmd != 'run'
        args.replace expanded_args if expanded_args
        send(cmd, args)
      end
    rescue Errno::ENOENT
      if $!.message.include? "No such file or directory - git"
        abort "Error: `git` command not found"
      else
        raise
      end
    rescue Context::FatalError => err
      abort "fatal: #{err.message}"
    end


    def ci_status(args)
      args.shift
      ref = args.words.first || 'HEAD'

      unless head_project = local_repo.current_project
        abort "Aborted: the origin remote doesn't point to a GitHub repository."
      end

      unless sha = local_repo.git_command("rev-parse -q #{ref}")
        abort "Aborted: no revision could be determined from '#{ref}'"
      end

      statuses = api_client.statuses(head_project, sha)
      status = statuses.first
      ref_state = status ? status['state'] : 'no status'

      exit_code = case ref_state
        when 'success'          then 0
        when 'failure', 'error' then 1
        when 'pending'          then 2
        else 3
        end

      $stdout.puts ref_state
      exit exit_code
    end

    def pull_request(args)
      args.shift
      options = { }
      force = explicit_owner = false
      base_project = local_repo.main_project
      head_project = local_repo.current_project

      unless current_branch
        abort "Aborted: not currently on any branch."
      end

      unless base_project
        abort "Aborted: the origin remote doesn't point to a GitHub repository."
      end

      from_github_ref = lambda do |ref, context_project|
        if ref.index(':')
          owner, ref = ref.split(':', 2)
          project = github_project(context_project.name, owner)
        end
        [project || context_project, ref]
      end

      while arg = args.shift
        case arg
        when '-f'
          force = true
        when '-F', '--file'
          file = args.shift
          text = file == '-' ? $stdin.read : File.read(file)
          options[:title], options[:body] = read_msg(text)
        when '-m', '--message'
          text = args.shift
          options[:title], options[:body] = read_msg(text)
        when '-b'
          base_project, options[:base] = from_github_ref.call(args.shift, base_project)
        when '-h'
          head = args.shift
          explicit_owner = !!head.index(':')
          head_project, options[:head] = from_github_ref.call(head, head_project)
        when '-i'
          options[:issue] = args.shift
        else
          if url = resolve_github_url(arg) and url.project_path =~ /^issues\/(\d+)/
            options[:issue] = $1
            base_project = url.project
          elsif !options[:title]
            options[:title] = arg
            warn "hub: Specifying pull request title without a flag is deprecated."
            warn "Please use one of `-m' or `-F' options."
          else
            abort "invalid argument: #{arg}"
          end
        end
      end

      options[:project] = base_project
      options[:base] ||= master_branch.short_name

      if tracked_branch = options[:head].nil? && current_branch.upstream
        if !tracked_branch.remote?
          tracked_branch = nil
        elsif base_project == head_project and tracked_branch.short_name == options[:base]
          $stderr.puts "Aborted: head branch is the same as base (#{options[:base].inspect})"
          warn "(use `-h <branch>` to specify an explicit pull request head)"
          abort
        end
      end
      options[:head] ||= (tracked_branch || current_branch).short_name

      user = github_user(head_project.host)
      if head_project.owner != user and !tracked_branch and !explicit_owner
        head_project = head_project.owned_by(user)
      end

      remote_branch = "#{head_project.remote}/#{options[:head]}"
      options[:head] = "#{head_project.owner}:#{options[:head]}"

      if !force and tracked_branch and local_commits = rev_list(remote_branch, nil)
        $stderr.puts "Aborted: #{local_commits.split("\n").size} commits are not yet pushed to #{remote_branch}"
        warn "(use `-f` to force submit a pull request anyway)"
        abort
      end

      if args.noop?
        puts "Would request a pull to #{base_project.owner}:#{options[:base]} from #{options[:head]}"
        exit
      end

      unless options[:title] or options[:issue]
        base_branch = "#{base_project.remote}/#{options[:base]}"
        commits = rev_list(base_branch, remote_branch).to_s.split("\n")

        case commits.size
        when 0
          default_message = commit_summary = nil
        when 1
          format = '%w(78,0,0)%s%n%+b'
          default_message = git_command "show -s --format='#{format}' #{commits.first}"
          commit_summary = nil
        else
          format = '%h (%aN, %ar)%n%w(78,3,3)%s%n%+b'
          default_message = nil
          commit_summary = git_command "log --no-color --format='%s' --cherry %s...%s" %
            [format, base_branch, remote_branch]
        end

        options[:title], options[:body] = pullrequest_editmsg(commit_summary) { |msg, initial_message|
          initial_message ||= default_message
          msg.puts initial_message if initial_message
          msg.puts ""
          msg.puts "# Requesting a pull to #{base_project.owner}:#{options[:base]} from #{options[:head]}"
          msg.puts "#"
          msg.puts "# Write a message for this pull request. The first block"
          msg.puts "# of text is the title and the rest is description."
        }
      end

      pull = api_client.create_pullrequest(options)

      args.executable = 'echo'
      args.replace [pull['html_url']]
    rescue GitHubAPI::Exceptions
      response = $!.response
      display_api_exception("creating pull request", response)
      if 404 == response.status
        base_url = base_project.web_url.split('://', 2).last
        warn "Are you sure that #{base_url} exists?"
      end
      exit 1
    else
      delete_editmsg
    end

    def clone(args)
      ssh = args.delete('-p')
      has_values = /^(--(upload-pack|template|depth|origin|branch|reference|name)|-[ubo])$/

      idx = 1
      while idx < args.length
        arg = args[idx]
        if arg.index('-') == 0
          idx += 1 if arg =~ has_values
        else
          if arg =~ NAME_WITH_OWNER_RE and !File.directory?(arg)
            name, owner = arg, nil
            owner, name = name.split('/', 2) if name.index('/')
            project = github_project(name, owner || github_user)
            ssh ||= args[0] != 'submodule' && project.owner == github_user(project.host) { }
            args[idx] = project.git_url(:private => ssh, :https => https_protocol?)
          end
          break
        end
        idx += 1
      end
    end

    def submodule(args)
      return unless index = args.index('add')
      args.delete_at index

      clone(args)
      args.insert index, 'add'
    end

    def remote(args)
      if %w[add set-url].include?(args[1])
        name = args.last
        if name =~ /^(#{OWNER_RE})$/ || name =~ /^(#{OWNER_RE})\/(#{NAME_RE})$/
          user, repo = $1, $2 || repo_name
        end
      end
      return unless user # do not touch arguments

      ssh = args.delete('-p')

      if args.words[2] == 'origin' && args.words[3].nil?
        user, repo = github_user, repo_name
      elsif args.words[-2] == args.words[1]
        idx = args.index( args.words[-1] )
        args[idx] = user
      else
        args.pop
      end

      args << git_url(user, repo, :private => ssh)
    end

    def fetch(args)
      if args.include?('--multiple')
        names = args.words[1..-1]
      elsif remote_name = args.words[1]
        if remote_name =~ /^\w+(,\w+)+$/
          index = args.index(remote_name)
          args.delete(remote_name)
          names = remote_name.split(',')
          args.insert(index, *names)
          args.insert(index, '--multiple')
        else
          names = [remote_name]
        end
      else
        names = []
      end

      projects = names.map { |name|
        unless name !~ /^#{OWNER_RE}$/ or remotes.include?(name) or remotes_group(name)
          project = github_project(nil, name)
          repo_info = api_client.repo_info(project)
          if repo_info.success?
            project.repo_data = repo_info.data
            project
          end
        end
      }.compact

      if projects.any?
        projects.each do |project|
          args.before ['remote', 'add', project.owner, project.git_url(:https => https_protocol?)]
        end
      end
    end

    def checkout(args)
      _, url_arg, new_branch_name = args.words
      if url = resolve_github_url(url_arg) and url.project_path =~ /^pull\/(\d+)/
        pull_id = $1
        pull_data = api_client.pullrequest_info(url.project, pull_id)

        args.delete new_branch_name
        user, branch = pull_data['head']['label'].split(':', 2)
        abort "Error: #{user}'s fork is not available anymore" unless pull_data['head']['repo']
        new_branch_name ||= "#{user}-#{branch}"

        if remotes.include? user
          args.before ['remote', 'set-branches', '--add', user, branch]
          args.before ['fetch', user, "+refs/heads/#{branch}:refs/remotes/#{user}/#{branch}"]
        else
          url = github_project(url.project_name, user).git_url(:private => pull_data['head']['repo']['private'],
                                                               :https => https_protocol?)
          args.before ['remote', 'add', '-f', '-t', branch, user, url]
        end
        idx = args.index url_arg
        args.delete_at idx
        args.insert idx, '--track', '-B', new_branch_name, "#{user}/#{branch}"
      end
    end

    def merge(args)
      _, url_arg = args.words
      if url = resolve_github_url(url_arg) and url.project_path =~ /^pull\/(\d+)/
        pull_id = $1
        pull_data = api_client.pullrequest_info(url.project, pull_id)

        user, branch = pull_data['head']['label'].split(':', 2)
        abort "Error: #{user}'s fork is not available anymore" unless pull_data['head']['repo']

        url = github_project(url.project_name, user).git_url(:private => pull_data['head']['repo']['private'],
                                                             :https => https_protocol?)

        merge_head = "#{user}/#{branch}"
        args.before ['fetch', url, "+refs/heads/#{branch}:refs/remotes/#{merge_head}"]

        idx = args.index url_arg
        args.delete_at idx
        args.insert idx, merge_head, '--no-ff', '-m',
                    "Merge pull request ##{pull_id} from #{merge_head}\n\n#{pull_data['title']}"
      end
    end

    def cherry_pick(args)
      unless args.include?('-m') or args.include?('--mainline')
        ref = args.words.last
        if url = resolve_github_url(ref) and url.project_path =~ /^commit\/([a-f0-9]{7,40})/
          sha = $1
          project = url.project
        elsif ref =~ /^(#{OWNER_RE})@([a-f0-9]{7,40})$/
          owner, sha = $1, $2
          project = local_repo.main_project.owned_by(owner)
        end

        if project
          args[args.index(ref)] = sha

          if remote = project.remote and remotes.include? remote
            args.before ['fetch', remote.to_s]
          else
            args.before ['remote', 'add', '-f', project.owner, project.git_url(:https => https_protocol?)]
          end
        end
      end
    end

    def am(args)
      if url = args.find { |a| a =~ %r{^https?://(gist\.)?github\.com/} }
        idx = args.index(url)
        gist = $1 == 'gist.'
        url = url.sub(/#.+/, '')
        url = url.sub(%r{(/pull/\d+)/\w*$}, '\1') unless gist
        ext = gist ? '.txt' : '.patch'
        url += ext unless File.extname(url) == ext
        patch_file = File.join(tmp_dir, "#{gist ? 'gist-' : ''}#{File.basename(url)}")
        args.before 'curl', ['-#LA', "hub #{Hub::Version}", url, '-o', patch_file]
        args[idx] = patch_file
      end
    end

    alias_method :apply, :am

    def init(args)
      if args.delete('-g')
        project = github_project(File.basename(current_dir))
        url = project.git_url(:private => true, :https => https_protocol?)
        args.after ['remote', 'add', 'origin', url]
      end
    end

    def fork(args)
      unless project = local_repo.main_project
        abort "Error: repository under 'origin' remote is not a GitHub project"
      end
      forked_project = project.owned_by(github_user(project.host))

      existing_repo = api_client.repo_info(forked_project)
      if existing_repo.success?
        parent_data = existing_repo.data['parent']
        parent_url  = parent_data && resolve_github_url(parent_data['html_url'])
        if !parent_url or parent_url.project != project
          abort "Error creating fork: %s already exists on %s" %
            [ forked_project.name_with_owner, forked_project.host ]
        end
      else
        api_client.fork_repo(project) unless args.noop?
      end

      if args.include?('--no-remote')
        exit
      else
        url = forked_project.git_url(:private => true, :https => https_protocol?)
        args.replace %W"remote add -f #{forked_project.owner} #{url}"
        args.after 'echo', ['new remote:', forked_project.owner]
      end
    rescue GitHubAPI::Exceptions
      display_api_exception("creating fork", $!.response)
      exit 1
    end

    def create(args)
      if !is_repo?
        abort "'create' must be run from inside a git repository"
      else
        owner = github_user
        args.shift
        options = {}
        options[:private] = true if args.delete('-p')
        new_repo_name = nil

        until args.empty?
          case arg = args.shift
          when '-d'
            options[:description] = args.shift
          when '-h'
            options[:homepage] = args.shift
          else
            if arg =~ /^[^-]/ and new_repo_name.nil?
              new_repo_name = arg
              owner, new_repo_name = new_repo_name.split('/', 2) if new_repo_name.index('/')
            else
              abort "invalid argument: #{arg}"
            end
          end
        end
        new_repo_name ||= repo_name
        new_project = github_project(new_repo_name, owner)

        if api_client.repo_exists?(new_project)
          warn "#{new_project.name_with_owner} already exists on #{new_project.host}"
          action = "set remote origin"
        else
          action = "created repository"
          unless args.noop?
            repo_data = api_client.create_repo(new_project, options)
            new_project = github_project(repo_data['full_name'])
          end
        end

        url = new_project.git_url(:private => true, :https => https_protocol?)

        if remotes.first != 'origin'
          args.replace %W"remote add -f origin #{url}"
        else
          args.replace %W"remote -v"
        end

        args.after 'echo', ["#{action}:", new_project.name_with_owner]
      end
    rescue GitHubAPI::Exceptions
      display_api_exception("creating repository", $!.response)
      exit 1
    end

    def push(args)
      return if args[1].nil? || !args[1].index(',')

      refs    = args.words[2..-1]
      remotes = args[1].split(',')
      args[1] = remotes.shift

      if refs.empty?
        refs = [current_branch.short_name]
        args.concat refs
      end

      remotes.each do |name|
        args.after ['push', name, *refs]
      end
    end

    def browse(args)
      args.shift
      browse_command(args) do
        dest = args.shift
        dest = nil if dest == '--'

        if dest
          project = github_project dest
          branch = master_branch
        else
          project = current_project
          branch = current_branch && current_branch.upstream || master_branch
        end

        abort "Usage: hub browse [<USER>/]<REPOSITORY>" unless project

        require 'cgi'
        path = case subpage = args.shift
        when 'commits'
          "/commits/#{branch_in_url(branch)}"
        when 'tree', NilClass
          "/tree/#{branch_in_url(branch)}" if branch and !branch.master?
        else
          "/#{subpage}"
        end

        project.web_url(path)
      end
    end

    def compare(args)
      args.shift
      browse_command(args) do
        if args.empty?
          branch = current_branch.upstream
          if branch and not branch.master?
            range = branch.short_name
            project = current_project
          else
            abort "Usage: hub compare [USER] [<START>...]<END>"
          end
        else
          sha_or_tag = /((?:#{OWNER_RE}:)?\w[\w.-]+\w)/
          range = args.pop.sub(/^#{sha_or_tag}\.\.#{sha_or_tag}$/, '\1...\2')
          project = if owner = args.pop then github_project(nil, owner)
                    else current_project
                    end
        end

        project.web_url "/compare/#{range}"
      end
    end

    def hub(args)
      return help(args) unless args[1] == 'standalone'
      require 'hub/standalone'
      Hub::Standalone.build $stdout
      exit
    rescue LoadError
      abort "hub is already running in standalone mode."
    rescue Errno::EPIPE
      exit # ignore broken pipe
    end

    def alias(args)
      shells = %w[bash zsh sh ksh csh fish]

      script = !!args.delete('-s')
      shell = args[1] || ENV['SHELL']
      abort "hub alias: unknown shell" if shell.nil? or shell.empty?
      shell = File.basename shell

      unless shells.include? shell
        $stderr.puts "hub alias: unsupported shell"
        warn "supported shells: #{shells.join(' ')}"
        abort
      end

      if script
        puts "alias git=hub"
      else
        profile = case shell
          when 'bash' then '~/.bash_profile'
          when 'zsh'  then '~/.zshrc'
          when 'ksh'  then '~/.profile'
          else
            'your profile'
          end

        puts "# Wrap git automatically by adding the following to #{profile}:"
        puts
        puts 'eval "$(hub alias -s)"'
      end

      exit
    end

    def version(args)
      args.after 'echo', ['hub version', Version]
    end
    alias_method "--version", :version

    def help(args)
      command = args.words[1]

      if command == 'hub' || custom_command?(command)
        puts hub_manpage
        exit
      elsif command.nil?
        if args.has_flag?('-a', '--all')
          args.after 'echo', ["\nhub custom commands\n"]
          args.after 'echo', CUSTOM_COMMANDS.map {|cmd| "  #{cmd}" }
        else
          ENV['GIT_PAGER'] = '' unless args.has_flag?('-p', '--paginate') # Use `cat`.
          puts improved_help_text
          exit
        end
      end
    end
    alias_method "--help", :help

  private

    def branch_in_url(branch)
      require 'cgi'
      CGI.escape(branch.short_name).gsub("%2F", "/")
    end

    def api_client
      @api_client ||= begin
        config_file = ENV['HUB_CONFIG'] || '~/.config/hub'
        file_store = GitHubAPI::FileStore.new File.expand_path(config_file)
        file_config = GitHubAPI::Configuration.new file_store
        GitHubAPI.new file_config, :app_url => 'http://github.github.com/hub/'
      end
    end

    def github_user host = nil, &block
      host ||= (local_repo(false) || Context::LocalRepo).default_host
      api_client.config.username(host, &block)
    end

    def custom_command? cmd
      CUSTOM_COMMANDS.include? cmd
    end

    def respect_help_flags args
      return if args.size > 2
      case args[1]
      when '-h'
        pattern = /(git|hub) #{Regexp.escape args[0].gsub('-', '\-')}/
        hub_raw_manpage.each_line { |line|
          if line =~ pattern
            $stderr.print "Usage: "
            $stderr.puts line.gsub(/\\f./, '').gsub('\-', '-')
            abort
          end
        }
        abort "Error: couldn't find usage help for #{args[0]}"
      when '--help'
        puts hub_manpage
        exit
      end
    end

    def improved_help_text
      <<-help
usage: git [--version] [--exec-path[=<path>]] [--html-path] [--man-path] [--info-path]
           [-p|--paginate|--no-pager] [--no-replace-objects] [--bare]
           [--git-dir=<path>] [--work-tree=<path>] [--namespace=<name>]
           [-c name=value] [--help]
           <command> [<args>]

Basic Commands:
   init       Create an empty git repository or reinitialize an existing one
   add        Add new or modified files to the staging area
   rm         Remove files from the working directory and staging area
   mv         Move or rename a file, a directory, or a symlink
   status     Show the status of the working directory and staging area
   commit     Record changes to the repository

History Commands:
   log        Show the commit history log
   diff       Show changes between commits, commit and working tree, etc
   show       Show information about commits, tags or files

Branching Commands:
   branch     List, create, or delete branches
   checkout   Switch the active branch to another branch
   merge      Join two or more development histories (branches) together
   tag        Create, list, delete, sign or verify a tag object

Remote Commands:
   clone      Clone a remote repository into a new directory
   fetch      Download data, tags and branches from a remote repository
   pull       Fetch from and merge with another repository or a local branch
   push       Upload data, tags and branches to a remote repository
   remote     View and manage a set of remote repositories

Advanced Commands:
   reset      Reset your staging area or working directory to another point
   rebase     Re-apply a series of patches in one branch onto another
   bisect     Find by binary search the change that introduced a bug
   grep       Print files with lines matching a pattern in your codebase

GitHub Commands:
   pull-request   Open a pull request on GitHub
   fork           Make a fork of a remote repository on GitHub and add as remote
   create         Create this repository on GitHub and add GitHub as origin
   browse         Open a GitHub page in the default browser
   compare        Open a compare page on GitHub
   ci-status      Show the CI status of a commit

See 'git help <command>' for more information on a specific command.
help
    end

    def slurp_global_flags(args)
      flags = %w[ --noop -c -p --paginate --no-pager --no-replace-objects --bare --version --help ]
      flags2 = %w[ --exec-path= --git-dir= --work-tree= ]

      globals = []
      locals = []

      while args[0] && (flags.include?(args[0]) || flags2.any? {|f| args[0].index(f) == 0 })
        flag = args.shift
        case flag
        when '--noop'
          args.noop!
        when '--version', '--help'
          args.unshift flag.sub('--', '')
        when '-c'
          config_pair = args.shift
          key, value = config_pair.split('=', 2)
          git_reader.stub_config_value(key, value)

          globals << flag << config_pair
        when '-p', '--paginate', '--no-pager'
          locals << flag
        else
          globals << flag
        end
      end

      git_reader.add_exec_flags(globals)
      args.add_exec_flags(globals)
      args.add_exec_flags(locals)
    end

    def browse_command(args)
      url_only = args.delete('-u')
      warn "Warning: the `-p` flag has no effect anymore" if args.delete('-p')
      url = yield

      args.executable = url_only ? 'echo' : browser_launcher
      args.push url
    end

    def hub_manpage
      abort "** Can't find groff(1)" unless command?('groff')

      require 'open3'
      out = nil
      Open3.popen3(groff_command) do |stdin, stdout, _|
        stdin.puts hub_raw_manpage
        stdin.close
        out = stdout.read.strip
      end
      out
    end

    def groff_command
      "groff -Wall -mtty-char -mandoc -Tascii"
    end

    def hub_raw_manpage
      if File.exists? file = File.dirname(__FILE__) + '/../../man/hub.1'
        File.read(file)
      else
        DATA.read
      end
    end

    def puts(*args)
      page_stdout
      super
    end

    def page_stdout
      return if not $stdout.tty? or windows?

      read, write = IO.pipe

      if Kernel.fork
        $stdin.reopen(read)
        read.close
        write.close

        ENV['LESS'] = 'FSRX'

        Kernel.select [STDIN]

        pager = ENV['GIT_PAGER'] ||
          `git config --get-all core.pager`.split("\n").first ||
          ENV['PAGER'] ||
          'less -isr'

        pager = 'cat' if pager.empty?

        exec pager rescue exec "/bin/sh", "-c", pager
      else
        $stdout.reopen(write)
        $stderr.reopen(write) if $stderr.tty?
        read.close
        write.close
      end
    rescue NotImplementedError
    end

    def pullrequest_editmsg(changes)
      message_file = pullrequest_editmsg_file

      if valid_editmsg_file?(message_file)
        title, body = read_editmsg(message_file)
        previous_message = [title, body].compact.join("\n\n") if title
      end

      File.open(message_file, 'w') { |msg|
        yield msg, previous_message
        if changes
          msg.puts "#\n# Changes:\n#"
          msg.puts changes.gsub(/^/, '# ').gsub(/ +$/, '')
        end
      }

      edit_cmd = Array(git_editor).dup
      edit_cmd << '-c' << 'set ft=gitcommit tw=0 wrap lbr' if edit_cmd[0] =~ /^[mg]?vim$/
      edit_cmd << message_file
      system(*edit_cmd)

      unless $?.success?
        delete_editmsg(message_file)
        abort "error using text editor for pull request message"
      end

      title, body = read_editmsg(message_file)
      abort "Aborting due to empty pull request title" unless title
      [title, body]
    end

    def valid_editmsg_file?(message_file)
      File.exists?(message_file) &&
        File.mtime(message_file) > File.mtime(__FILE__)
    end

    def read_msg(message)
      message.split("\n\n", 2).each {|s| s.strip! }.reject {|s| s.empty? }
    end

    def pullrequest_editmsg_file
      File.join(git_dir, 'PULLREQ_EDITMSG')
    end

    def read_editmsg(file)
      title, body = '', ''
      File.open(file, 'r') { |msg|
        msg.each_line do |line|
          next if line.index('#') == 0
          ((body.empty? and line =~ /\S/) ? title : body) << line
        end
      }
      title.tr!("\n", ' ')
      title.strip!
      body.strip!

      [title =~ /\S/ ? title : nil, body =~ /\S/ ? body : nil]
    end

    def delete_editmsg(file = pullrequest_editmsg_file)
      File.delete(file) if File.exist?(file)
    end

    def expand_alias(cmd)
      if expanded = git_alias_for(cmd)
        if expanded.index('!') != 0
          require 'shellwords' unless defined?(::Shellwords)
          Shellwords.shellwords(expanded)
        end
      end
    end

    def display_api_exception(action, response)
      $stderr.puts "Error #{action}: #{response.message.strip} (HTTP #{response.status})"
      if 422 == response.status and response.error_message?
        msg = response.error_message
        msg = msg.join("\n") if msg.respond_to? :join
        warn msg
      end
    end

  end
end

module Hub
  class Runner
    attr_reader :args
    
    def initialize(*args)
      @args = Args.new(args)
      Commands.run(@args)
    end

    def self.execute(*args)
      new(*args).execute
    end

    def command
      if args.skip?
        ''
      else
        commands.join('; ')
      end
    end

    def commands
      args.commands.map do |cmd|
        if cmd.respond_to?(:join)
          cmd.map { |arg| arg = arg.to_s; (arg.index(' ') || arg.empty?) ? "'#{arg}'" : arg }.join(' ')
        else
          cmd.to_s
        end
      end
    end

    def execute
      if args.noop?
        puts commands
      elsif not args.skip?
        execute_command_chain args.commands
      end
    end

    def execute_command_chain commands
      commands.each_with_index do |cmd, i|
        if cmd.respond_to?(:call) then cmd.call
        elsif i == commands.length - 1
          exec(*cmd)
        else
          exit($?.exitstatus) unless system(*cmd)
        end
      end
    end

    def exec *args
      if args.first == 'echo' && Context::windows?
        puts args[1..-1].join(' ')
      else
        super
      end
    end
  end
end

Hub::Runner.execute(*ARGV)

__END__
.\" generated with Ronn/v0.7.3
.\" http://github.com/rtomayko/ronn/tree/0.7.3
.
.TH "HUB" "1" "July 2013" "GITHUB" "Git Manual"
.
.SH "NAME"
\fBhub\fR \- git + hub = github
.
.SH "SYNOPSIS"
\fBhub\fR [\fB\-\-noop\fR] \fICOMMAND\fR \fIOPTIONS\fR
.
.br
\fBhub alias\fR [\fB\-s\fR] [\fISHELL\fR]
.
.SS "Expanded git commands:"
\fBgit init \-g\fR \fIOPTIONS\fR
.
.br
\fBgit clone\fR [\fB\-p\fR] \fIOPTIONS\fR [\fIUSER\fR/]\fIREPOSITORY\fR \fIDIRECTORY\fR
.
.br
\fBgit remote add\fR [\fB\-p\fR] \fIOPTIONS\fR \fIUSER\fR[/\fIREPOSITORY\fR]
.
.br
\fBgit remote set\-url\fR [\fB\-p\fR] \fIOPTIONS\fR \fIREMOTE\-NAME\fR \fIUSER\fR[/\fIREPOSITORY\fR]
.
.br
\fBgit fetch\fR \fIUSER\-1\fR,[\fIUSER\-2\fR,\.\.\.]
.
.br
\fBgit checkout\fR \fIPULLREQ\-URL\fR [\fIBRANCH\fR]
.
.br
\fBgit merge\fR \fIPULLREQ\-URL\fR
.
.br
\fBgit cherry\-pick\fR \fIGITHUB\-REF\fR
.
.br
\fBgit am\fR \fIGITHUB\-URL\fR
.
.br
\fBgit apply\fR \fIGITHUB\-URL\fR
.
.br
\fBgit push\fR \fIREMOTE\-1\fR,\fIREMOTE\-2\fR,\.\.\.,\fIREMOTE\-N\fR [\fIREF\fR]
.
.br
\fBgit submodule add\fR [\fB\-p\fR] \fIOPTIONS\fR [\fIUSER\fR/]\fIREPOSITORY\fR \fIDIRECTORY\fR
.
.SS "Custom git commands:"
\fBgit create\fR [\fINAME\fR] [\fB\-p\fR] [\fB\-d\fR \fIDESCRIPTION\fR] [\fB\-h\fR \fIHOMEPAGE\fR]
.
.br
\fBgit browse\fR [\fB\-u\fR] [[\fIUSER\fR\fB/\fR]\fIREPOSITORY\fR] [SUBPAGE]
.
.br
\fBgit compare\fR [\fB\-u\fR] [\fIUSER\fR] [\fISTART\fR\.\.\.]\fIEND\fR
.
.br
\fBgit fork\fR [\fB\-\-no\-remote\fR]
.
.br
\fBgit pull\-request\fR [\fB\-f\fR] [\fB\-m\fR \fIMESSAGE\fR|\fB\-F\fR \fIFILE\fR|\fB\-i\fR \fIISSUE\fR|\fIISSUE\-URL\fR] [\fB\-b\fR \fIBASE\fR] [\fB\-h\fR \fIHEAD\fR]
.
.br
\fBgit ci\-status\fR [\fICOMMIT\fR]
.
.SH "DESCRIPTION"
hub enhances various git commands to ease most common workflows with GitHub\.
.
.TP
\fBhub \-\-noop\fR \fICOMMAND\fR
Shows which command(s) would be run as a result of the current command\. Doesn\'t perform anything\.
.
.TP
\fBhub alias\fR [\fB\-s\fR] [\fISHELL\fR]
Shows shell instructions for wrapping git\. If given, \fISHELL\fR specifies the type of shell; otherwise defaults to the value of SHELL environment variable\. With \fB\-s\fR, outputs shell script suitable for \fBeval\fR\.
.
.TP
\fBgit init\fR \fB\-g\fR \fIOPTIONS\fR
Create a git repository as with git\-init(1) and add remote \fBorigin\fR at "git@github\.com:\fIUSER\fR/\fIREPOSITORY\fR\.git"; \fIUSER\fR is your GitHub username and \fIREPOSITORY\fR is the current working directory\'s basename\.
.
.TP
\fBgit clone\fR [\fB\-p\fR] \fIOPTIONS\fR [\fIUSER\fR\fB/\fR]\fIREPOSITORY\fR \fIDIRECTORY\fR
Clone repository "git://github\.com/\fIUSER\fR/\fIREPOSITORY\fR\.git" into \fIDIRECTORY\fR as with git\-clone(1)\. When \fIUSER\fR/ is omitted, assumes your GitHub login\. With \fB\-p\fR, clone private repositories over SSH\. For repositories under your GitHub login, \fB\-p\fR is implicit\.
.
.TP
\fBgit remote add\fR [\fB\-p\fR] \fIOPTIONS\fR \fIUSER\fR[\fB/\fR\fIREPOSITORY\fR]
Add remote "git://github\.com/\fIUSER\fR/\fIREPOSITORY\fR\.git" as with git\-remote(1)\. When /\fIREPOSITORY\fR is omitted, the basename of the current working directory is used\. With \fB\-p\fR, use private remote "git@github\.com:\fIUSER\fR/\fIREPOSITORY\fR\.git"\. If \fIUSER\fR is "origin" then uses your GitHub login\.
.
.TP
\fBgit remote set\-url\fR [\fB\-p\fR] \fIOPTIONS\fR \fIREMOTE\-NAME\fR \fIUSER\fR[/\fIREPOSITORY\fR]
Sets the url of remote \fIREMOTE\-NAME\fR using the same rules as \fBgit remote add\fR\.
.
.TP
\fBgit fetch\fR \fIUSER\-1\fR,[\fIUSER\-2\fR,\.\.\.]
Adds missing remote(s) with \fBgit remote add\fR prior to fetching\. New remotes are only added if they correspond to valid forks on GitHub\.
.
.TP
\fBgit checkout\fR \fIPULLREQ\-URL\fR [\fIBRANCH\fR]
Checks out the head of the pull request as a local branch, to allow for reviewing, rebasing and otherwise cleaning up the commits in the pull request before merging\. The name of the local branch can explicitly be set with \fIBRANCH\fR\.
.
.TP
\fBgit merge\fR \fIPULLREQ\-URL\fR
Merge the pull request with a commit message that includes the pull request ID and title, similar to the GitHub Merge Button\.
.
.TP
\fBgit cherry\-pick\fR \fIGITHUB\-REF\fR
Cherry\-pick a commit from a fork using either full URL to the commit or GitHub\-flavored Markdown notation, which is \fBuser@sha\fR\. If the remote doesn\'t yet exist, it will be added\. A \fBgit fetch <user>\fR is issued prior to the cherry\-pick attempt\.
.
.TP
\fBgit [am|apply]\fR \fIGITHUB\-URL\fR
Downloads the patch file for the pull request or commit at the URL and applies that patch from disk with \fBgit am\fR or \fBgit apply\fR\. Similar to \fBcherry\-pick\fR, but doesn\'t add new remotes\. \fBgit am\fR creates commits while preserving authorship info while \fBapply\fR only applies the patch to the working copy\.
.
.TP
\fBgit push\fR \fIREMOTE\-1\fR,\fIREMOTE\-2\fR,\.\.\.,\fIREMOTE\-N\fR [\fIREF\fR]
Push \fIREF\fR to each of \fIREMOTE\-1\fR through \fIREMOTE\-N\fR by executing multiple \fBgit push\fR commands\.
.
.TP
\fBgit submodule add\fR [\fB\-p\fR] \fIOPTIONS\fR [\fIUSER\fR/]\fIREPOSITORY\fR \fIDIRECTORY\fR
Submodule repository "git://github\.com/\fIUSER\fR/\fIREPOSITORY\fR\.git" into \fIDIRECTORY\fR as with git\-submodule(1)\. When \fIUSER\fR/ is omitted, assumes your GitHub login\. With \fB\-p\fR, use private remote "git@github\.com:\fIUSER\fR/\fIREPOSITORY\fR\.git"\.
.
.TP
\fBgit help\fR
Display enhanced git\-help(1)\.
.
.P
hub also adds some custom commands that are otherwise not present in git:
.
.TP
\fBgit create\fR [\fINAME\fR] [\fB\-p\fR] [\fB\-d\fR \fIDESCRIPTION\fR] [\fB\-h\fR \fIHOMEPAGE\fR]
Create a new public GitHub repository from the current git repository and add remote \fBorigin\fR at "git@github\.com:\fIUSER\fR/\fIREPOSITORY\fR\.git"; \fIUSER\fR is your GitHub username and \fIREPOSITORY\fR is the current working directory name\. To explicitly name the new repository, pass in \fINAME\fR, optionally in \fIORGANIZATION\fR/\fINAME\fR form to create under an organization you\'re a member of\. With \fB\-p\fR, create a private repository, and with \fB\-d\fR and \fB\-h\fR set the repository\'s description and homepage URL, respectively\.
.
.TP
\fBgit browse\fR [\fB\-u\fR] [[\fIUSER\fR\fB/\fR]\fIREPOSITORY\fR] [SUBPAGE]
Open repository\'s GitHub page in the system\'s default web browser using \fBopen(1)\fR or the \fBBROWSER\fR env variable\. If the repository isn\'t specified, \fBbrowse\fR opens the page of the repository found in the current directory\. If SUBPAGE is specified, the browser will open on the specified subpage: one of "wiki", "commits", "issues" or other (the default is "tree")\. With \fB\-u\fR, outputs the URL rather than opening the browser\.
.
.TP
\fBgit compare\fR [\fB\-u\fR] [\fIUSER\fR] [\fISTART\fR\.\.\.]\fIEND\fR
Open a GitHub compare view page in the system\'s default web browser\. \fISTART\fR to \fIEND\fR are branch names, tag names, or commit SHA1s specifying the range of history to compare\. If a range with two dots (\fBa\.\.b\fR) is given, it will be transformed into one with three dots\. If \fISTART\fR is omitted, GitHub will compare against the base branch (the default is "master")\. With \fB\-u\fR, outputs the URL rather than opening the browser\.
.
.TP
\fBgit fork\fR [\fB\-\-no\-remote\fR]
Forks the original project (referenced by "origin" remote) on GitHub and adds a new remote for it under your username\.
.
.TP
\fBgit pull\-request\fR [\fB\-f\fR] [\fB\-m\fR \fIMESSAGE\fR|\fB\-F\fR \fIFILE\fR|\fB\-i\fR \fIISSUE\fR|\fIISSUE\-URL\fR] [\fB\-b\fR \fIBASE\fR] [\fB\-h\fR \fIHEAD\fR]
Opens a pull request on GitHub for the project that the "origin" remote points to\. The default head of the pull request is the current branch\. Both base and head of the pull request can be explicitly given in one of the following formats: "branch", "owner:branch", "owner/repo:branch"\. This command will abort operation if it detects that the current topic branch has local commits that are not yet pushed to its upstream branch on the remote\. To skip this check, use \fB\-f\fR\.
.
.IP
Without \fIMESSAGE\fR or \fIFILE\fR, a text editor will open in which title and body of the pull request can be entered in the same manner as git commit message\. Pull request message can also be passed via stdin with \fB\-F \-\fR\.
.
.IP
If instead of normal \fITITLE\fR an issue number is given with \fB\-i\fR, the pull request will be attached to an existing GitHub issue\. Alternatively, instead of title you can paste a full URL to an issue on GitHub\.
.
.TP
\fBgit ci\-status\fR [\fICOMMIT\fR]
Looks up the SHA for \fICOMMIT\fR in GitHub Status API and displays the latest status\. Exits with one of:
.
.br
success (0), error (1), failure (1), pending (2), no status (3)
.
.SH "CONFIGURATION"
Hub will prompt for GitHub username & password the first time it needs to access the API and exchange it for an OAuth token, which it saves in "~/\.config/hub"\.
.
.P
To avoid being prompted, use \fIGITHUB_USER\fR and \fIGITHUB_PASSWORD\fR environment variables\.
.
.P
If you prefer the HTTPS protocol for GitHub repositories, you can set "hub\.protocol" to "https"\. This will affect \fBclone\fR, \fBfork\fR, \fBremote add\fR and other operations that expand references to GitHub repositories as full URLs that otherwise use git and ssh protocols\.
.
.IP "" 4
.
.nf

$ git config \-\-global hub\.protocol https
.
.fi
.
.IP "" 0
.
.SS "GitHub Enterprise"
By default, hub will only work with repositories that have remotes which point to github\.com\. GitHub Enterprise hosts need to be whitelisted to configure hub to treat such remotes same as github\.com:
.
.IP "" 4
.
.nf

$ git config \-\-global \-\-add hub\.host my\.git\.org
.
.fi
.
.IP "" 0
.
.P
The default host for commands like \fBinit\fR and \fBclone\fR is still github\.com, but this can be affected with the \fIGITHUB_HOST\fR environment variable:
.
.IP "" 4
.
.nf

$ GITHUB_HOST=my\.git\.org git clone myproject
.
.fi
.
.IP "" 0
.
.SH "EXAMPLES"
.
.SS "git clone"
.
.nf

$ git clone schacon/ticgit
> git clone git://github\.com/schacon/ticgit\.git

$ git clone \-p schacon/ticgit
> git clone git@github\.com:schacon/ticgit\.git

$ git clone resque
> git clone git@github\.com/YOUR_USER/resque\.git
.
.fi
.
.SS "git remote add"
.
.nf

$ git remote add rtomayko
> git remote add rtomayko git://github\.com/rtomayko/CURRENT_REPO\.git

$ git remote add \-p rtomayko
> git remote add rtomayko git@github\.com:rtomayko/CURRENT_REPO\.git

$ git remote add origin
> git remote add origin git://github\.com/YOUR_USER/CURRENT_REPO\.git
.
.fi
.
.SS "git fetch"
.
.nf

$ git fetch mislav
> git remote add mislav git://github\.com/mislav/REPO\.git
> git fetch mislav

$ git fetch mislav,xoebus
> git remote add mislav \.\.\.
> git remote add xoebus \.\.\.
> git fetch \-\-multiple mislav xoebus
.
.fi
.
.SS "git cherry\-pick"
.
.nf

$ git cherry\-pick http://github\.com/mislav/REPO/commit/SHA
> git remote add \-f mislav git://github\.com/mislav/REPO\.git
> git cherry\-pick SHA

$ git cherry\-pick mislav@SHA
> git remote add \-f mislav git://github\.com/mislav/CURRENT_REPO\.git
> git cherry\-pick SHA

$ git cherry\-pick mislav@SHA
> git fetch mislav
> git cherry\-pick SHA
.
.fi
.
.SS "git am, git apply"
.
.nf

$ git am https://github\.com/defunkt/hub/pull/55
> curl https://github\.com/defunkt/hub/pull/55\.patch \-o /tmp/55\.patch
> git am /tmp/55\.patch

$ git am \-\-ignore\-whitespace https://github\.com/davidbalbert/hub/commit/fdb9921
> curl https://github\.com/davidbalbert/hub/commit/fdb9921\.patch \-o /tmp/fdb9921\.patch
> git am \-\-ignore\-whitespace /tmp/fdb9921\.patch

$ git apply https://gist\.github\.com/8da7fb575debd88c54cf
> curl https://gist\.github\.com/8da7fb575debd88c54cf\.txt \-o /tmp/gist\-8da7fb575debd88c54cf\.txt
> git apply /tmp/gist\-8da7fb575debd88c54cf\.txt
.
.fi
.
.SS "git fork"
.
.nf

$ git fork
[ repo forked on GitHub ]
> git remote add \-f YOUR_USER git@github\.com:YOUR_USER/CURRENT_REPO\.git
.
.fi
.
.SS "git pull\-request"
.
.nf

# while on a topic branch called "feature":
$ git pull\-request
[ opens text editor to edit title & body for the request ]
[ opened pull request on GitHub for "YOUR_USER:feature" ]

# explicit title, pull base & head:
$ git pull\-request \-m "Implemented feature X" \-b defunkt:master \-h mislav:feature

$ git pull\-request \-i 123
[ attached pull request to issue #123 ]
.
.fi
.
.SS "git checkout"
.
.nf

$ git checkout https://github\.com/defunkt/hub/pull/73
> git remote add \-f \-t feature git://github:com/mislav/hub\.git
> git checkout \-\-track \-B mislav\-feature mislav/feature

$ git checkout https://github\.com/defunkt/hub/pull/73 custom\-branch\-name
.
.fi
.
.SS "git merge"
.
.nf

$ git merge https://github\.com/defunkt/hub/pull/73
> git fetch git://github\.com/mislav/hub\.git +refs/heads/feature:refs/remotes/mislav/feature
> git merge mislav/feature \-\-no\-ff \-m \'Merge pull request #73 from mislav/feature\.\.\.\'
.
.fi
.
.SS "git create"
.
.nf

$ git create
[ repo created on GitHub ]
> git remote add origin git@github\.com:YOUR_USER/CURRENT_REPO\.git

# with description:
$ git create \-d \'It shall be mine, all mine!\'

$ git create recipes
[ repo created on GitHub ]
> git remote add origin git@github\.com:YOUR_USER/recipes\.git

$ git create sinatra/recipes
[ repo created in GitHub organization ]
> git remote add origin git@github\.com:sinatra/recipes\.git
.
.fi
.
.SS "git init"
.
.nf

$ git init \-g
> git init
> git remote add origin git@github\.com:YOUR_USER/REPO\.git
.
.fi
.
.SS "git push"
.
.nf

$ git push origin,staging,qa bert_timeout
> git push origin bert_timeout
> git push staging bert_timeout
> git push qa bert_timeout
.
.fi
.
.SS "git browse"
.
.nf

$ git browse
> open https://github\.com/YOUR_USER/CURRENT_REPO

$ git browse \-\- commit/SHA
> open https://github\.com/YOUR_USER/CURRENT_REPO/commit/SHA

$ git browse \-\- issues
> open https://github\.com/YOUR_USER/CURRENT_REPO/issues

$ git browse schacon/ticgit
> open https://github\.com/schacon/ticgit

$ git browse schacon/ticgit commit/SHA
> open https://github\.com/schacon/ticgit/commit/SHA

$ git browse resque
> open https://github\.com/YOUR_USER/resque

$ git browse resque network
> open https://github\.com/YOUR_USER/resque/network
.
.fi
.
.SS "git compare"
.
.nf

$ git compare refactor
> open https://github\.com/CURRENT_REPO/compare/refactor

$ git compare 1\.0\.\.1\.1
> open https://github\.com/CURRENT_REPO/compare/1\.0\.\.\.1\.1

$ git compare \-u fix
> (https://github\.com/CURRENT_REPO/compare/fix)

$ git compare other\-user patch
> open https://github\.com/other\-user/REPO/compare/patch
.
.fi
.
.SS "git submodule"
.
.nf

$ hub submodule add wycats/bundler vendor/bundler
> git submodule add git://github\.com/wycats/bundler\.git vendor/bundler

$ hub submodule add \-p wycats/bundler vendor/bundler
> git submodule add git@github\.com:wycats/bundler\.git vendor/bundler

$ hub submodule add \-b ryppl \-\-name pip ryppl/pip vendor/pip
> git submodule add \-b ryppl \-\-name pip git://github\.com/ryppl/pip\.git vendor/pip
.
.fi
.
.SS "git ci\-status"
.
.nf

$ hub ci\-status [commit]
> (prints CI state of commit and exits with appropriate code)
> One of: success (0), error (1), failure (1), pending (2), no status (3)
.
.fi
.
.SS "git help"
.
.nf

$ git help
> (improved git help)
$ git help hub
> (hub man page)
.
.fi
.
.SH "BUGS"
\fIhttps://github\.com/github/hub/issues\fR
.
.SH "AUTHORS"
\fIhttps://github\.com/github/hub/contributors\fR
.
.SH "SEE ALSO"
git(1), git\-clone(1), git\-remote(1), git\-init(1), \fIhttp://github\.com\fR, \fIhttps://github\.com/github/hub\fR
