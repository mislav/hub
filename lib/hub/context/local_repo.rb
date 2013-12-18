module Hub
  module Context
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
        if remote = origin_remote
          default_branch = git_command("rev-parse --symbolic-full-name #{remote}")
        end
        Branch.new(self, default_branch || 'refs/heads/master')
      end

      def remotes
        @remotes ||= begin
          # TODO: is there a plumbing command to get a list of remotes?
          list = git_command('remote').to_s.split("\n")
          # force "origin" to be first in the list
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
        # support ssh.github.com
        # https://help.github.com/articles/using-ssh-over-the-https-port
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
  end
end
