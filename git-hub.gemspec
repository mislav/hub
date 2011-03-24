$LOAD_PATH.unshift File.expand_path('../lib', __FILE__)
require 'hub/version'

Gem::Specification.new do |s|
  s.name              = "git-hub"
  s.version           = Hub::VERSION
  s.date              = Time.now.strftime('%Y-%m-%d')
  s.summary           = "The `git-hub' gem has been renamed `hub'."
  s.homepage          = "http://github.com/defunkt/hub"
  s.email             = "chris@ozmm.org"
  s.authors           = [ "Chris Wanstrath" ]
  s.has_rdoc          = false
  s.add_dependency "hub"
end
