# encoding: utf-8
require File.expand_path('../lib/hub/version', __FILE__)

Gem::Specification.new do |s|
  s.name              = "hub"
  s.version           = Hub::VERSION
  s.summary           = "Command-line wrapper for git and GitHub"
  s.homepage          = "http://hub.github.com/"
  s.email             = "mislav.marohnic@gmail.com"
  s.authors           = [ "Chris Wanstrath", "Mislav Marohnić" ]
  s.license           = "MIT"

  s.files             = %w( README.md Rakefile LICENSE )
  s.files            += Dir.glob("lib/**/*")
  s.files            += Dir.glob("bin/**/*")
  s.files            += Dir.glob("man/**/*")

  # include only files in version control
  git_dir = File.expand_path('../.git', __FILE__)
  if File.directory?(git_dir)
    s.files &= `git --git-dir='#{git_dir}' ls-files -z`.split("\0")
  end

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
       https://github.com/github/hub#readme  under the
       "Standalone" section.

       Cheers,
       defunkt

------------------------------------------------------------

message
end
