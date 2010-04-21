$LOAD_PATH.unshift 'lib'
require 'hub/version'

Gem::Specification.new do |s|
  s.name              = "git-hub"
  s.version           = Hub::VERSION
  s.date              = Time.now.strftime('%Y-%m-%d')
  s.summary           = "hub introduces git to GitHub"
  s.homepage          = "http://github.com/defunkt/hub"
  s.email             = "chris@ozmm.org"
  s.authors           = [ "Chris Wanstrath" ]
  s.has_rdoc          = false

  s.files             = %w( README.md Rakefile LICENSE )
  s.files            += Dir.glob("lib/**/*")
  s.files            += Dir.glob("bin/**/*")
  s.files            += Dir.glob("man/**/*")
  s.files            += Dir.glob("test/**/*")

  s.executables       = %w( hub )
  s.description       = <<desc
  `hub` is a command line utility which adds GitHub knowledge to `git`.

  It can used on its own or as a `git` wrapper.

  Normal:

      $ hub clone rtomayko/tilt

      Expands to:
      $ git clone git://github.com/rtomayko/tilt.git

  Wrapping `git`:

      $ git clone rack/rack

      Expands to:
      $ git clone git://github.com/rack/rack.git
desc

  s.post_install_message = <<-message

------------------------------------------------------------

                  You there! Wait, I say!
                  =======================

       If you are a heavy user of `git` on the command
       line  you  may  want  to  install `hub` the old
       fashioned way.  Faster  startup  time,  you see.

       Check  out  the  installation  instructions  at
       http://github.com/defunkt/hub#readme  under the
       "Standalone" section.

       Cheers,
       defunkt

------------------------------------------------------------

message
end
