module Hub
  module Context
    class GitHubProject < Struct.new(:local_repo, :owner, :name, :host)
      def self.from_url(url, local_repo)
        if local_repo.known_hosts.include? url.host
          _, owner, name = url.path.split('/', 4)

          new(local_repo, owner, name.sub(/\.git$/, ''), url.host)
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
