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
        @ssh_config ||= SSHConfig.new
      end
    end
  end
end
