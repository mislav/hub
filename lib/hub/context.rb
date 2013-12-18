require 'shellwords'
require 'forwardable'
require 'uri'

require_relative 'context/git_reader'
require_relative 'context/git_reader_methods'
require_relative 'context/system'
require_relative 'context/local_repo'

module Hub
  # Methods for inspecting the environment, such as reading git config,
  # repository info, and other.
  module Context
    extend Forwardable

    NULL = defined?(File::NULL) ? File::NULL : File.exist?('/dev/null') ? '/dev/null' : 'NUL'

    Error = Class.new(RuntimeError)
    FatalError = Class.new(Error)

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
        # FIXME: duplicates functionality of LocalRepo#master_branch
        Branch.new nil, 'refs/heads/master'
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

      # segment of path after the project owner and name
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
        master_name = if local_repo then local_repo.master_branch.short_name
        else 'master'
        end
        short_name == master_name
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

    ## helper methods for local repo, GH projects

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

    # legacy setting
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
      # possible: ~/bin/vi, $SOME_ENVIRONMENT_VARIABLE, "C:\Program Files\Vim\gvim.exe" --nofork
      editor = git_command 'var GIT_EDITOR'
      editor.gsub!(/\$(\w+|\{\w+\})/) { ENV[$1.tr('{}', '')] }
      editor = ENV[$1] if editor =~ /^\$(\w+)$/
      editor = File.expand_path editor if (editor =~ /^[~.]/ or editor.index('/')) and editor !~ /["']/
      # avoid shellsplitting "C:\Program Files"
      if File.exist? editor then [editor]
      else editor.shellsplit
      end
    end

    include System
    extend System
  end
end
