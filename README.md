hub: git + hub = github
=======================

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

hub requires you have `git` installed and in your `$PATH`. It also
requires Ruby 1.8.6+ or Ruby 1.9.1+. No other libraries necessary.


Install
-------

### Standalone

`hub` is most easily installed as a standalone script:

    curl -s http://defunkt.github.com/hub/standalone > ~/bin/hub &&
    chmod 755 ~/bin/hub

Assuming `~/bin/` is in your `$PATH`, you're ready to roll:

    $ hub version
    git version 1.7.0.4
    hub version 1.1.0

### Homebrew

    $ brew install hub
    $ which hub
    /usr/local/bin/hub
    $ hub version
    ...

### RubyGems

Though not recommended, `hub` can also be installed as a RubyGem:

    $ gem install git-hub

Yes, the gem name is "git-hub".

(It's not recommended for casual use because of the RubyGems startup
time. See [this gist][speed] for information.)

### Standalone via RubyGems

Yes, the gem name is still "git-hub":

    $ gem install git-hub
    $ hub hub standalone > ~/bin/hub && chmod 755 ~/bin/hub

This installs a standalone version which doesn't require RubyGems to
run.

### Source

You can also install from source:

    $ git clone git://github.com/defunkt/hub.git
    $ cd hub
    $ rake install prefix=/usr/local


Aliasing
--------

`hub` works best when it wraps `git`. This is not dangerous - your
normal git commands should all work. hub merely adds some sugar.

Typing `hub alias <shell>` will display alias instructions for
your shell. `hub alias` alone will show the known shells.

For example:

    $ hub alias bash
    Run this in your shell to start using `hub` as `git`:
      alias git=hub

You should place this command in your `.bash_profile` or other startup
script to ensure runs on login.

The alias command can also be eval'd directly using the `-s` flag:

    $ eval `hub alias -s bash`


Commands
--------

Assuming you've aliased `hub` to `git` the following commands now have
superpowers:

### git clone

    $ git clone schacon/ticgit
    > git clone git://github.com/schacon/ticgit.git

    $ git clone -p schacon/ticgit
    > git clone git@github.com:schacon/ticgit.git

    $ git clone resque
    > git clone git://github.com/YOUR_USER/resque.git

    $ git clone -p resque
    > git clone git@github.com:YOUR_USER/resque.git

### git remote add

    $ git remote add rtomayko
    > git remote add rtomayko git://github.com/rtomayko/CURRENT_REPO.git

    $ git remote add -p rtomayko
    > git remote add rtomayko git@github.com:rtomayko/CURRENT_REPO.git

    $ git remote add origin
    > git remote add origin git://github.com/YOUR_USER/CURRENT_REPO.git

### git fetch

    $ git fetch mislav
    > git remote add mislav git://github.com/mislav/REPO.git
    > git fetch mislav

    $ git fetch mislav,xoebus
    > git remote add mislav ...
    > git remote add xoebus ...
    > git fetch --multiple mislav xoebus

### git cherry-pick

    $ git cherry-pick http://github.com/mislav/REPO/commit/SHA
    > git remote add -f mislav git://github.com/mislav/REPO.git
    > git cherry-pick SHA

    $ git cherry-pick mislav@SHA
    > git remote add -f mislav git://github.com/mislav/CURRENT_REPO.git
    > git cherry-pick SHA

    $ git cherry-pick mislav@SHA
    > git fetch mislav
    > git cherry-pick SHA

### git fork

    $ git fork
    ... hardcore forking action ...
    > git remote add -f YOUR_USER git@github.com:YOUR_USER/CURRENT_REPO.git

Forks the original repo on GitHub and adds the new remote under your
username. It requires your GitHub token to be present; see "GitHub
login" below for details.

### git init

    $ git init -g
    > git init
    > git remote add origin git@github.com:YOUR_USER/REPO.git

### git push

    $ git push origin,staging,qa bert_timeout
    > git push origin bert_timeout
    > git push staging bert_timeout
    > git push qa bert_timeout

### git browse

    $ git browse
    > open http://github.com/CURRENT_REPO

    $ git browse -- issues
    > open http://github.com/CURRENT_REPO/issues

    $ git browse schacon/ticgit
    > open http://github.com/schacon/ticgit

    $ git browse -p schacon/ticgit
    > open https://github.com/schacon/ticgit

    $ git browse resque
    > open http://github.com/YOUR_USER/resque

    $ git browse resque network
    > open http://github.com/YOUR_USER/resque/network

    $ git browse -p resque
    > open https://github.com:YOUR_USER/resque

### git compare

    $ git compare refactor
    > open http://github.com/CURRENT_REPO/compare/refactor

    $ git compare 1.0...1.1
    > open http://github.com/CURRENT_REPO/compare/1.0...1.1

    $ git compare -u fix
    > (http://github.com/CURRENT_REPO/compare/fix)

    $ git compare other-user patch
    > open http://github.com/other-user/REPO/compare/patch

### git submodule

    $ hub submodule add wycats/bundler vendor/bundler
    > git submodule add git://github.com/wycats/bundler.git vendor/bundler

    $ hub submodule add -p wycats/bundler vendor/bundler
    > git submodule add git@github.com:wycats/bundler.git vendor/bundler

    $ hub submodule add -b ryppl ryppl/pip vendor/pip
    > git submodule add -b ryppl git://github.com/ryppl/pip.git vendor/pip


### git help

    $ git help
    > (improved git help)
    $ git help hub
    > (hub man page)


GitHub Login
------------

To get the most out of `hub`, you'll want to ensure your GitHub login
is stored locally in your Git config.

To test it run this:

    $ git config --global github.user

If you see nothing, you need to set the config setting:

    $ git config --global github.user YOUR_USER

For commands that require write access to GitHub (such as `fork`), you'll want to
setup "github.token" as well. See [local GitHub config guide][2] for more information.


Configuration
-------------

If you prefer `http://` clones to `git://` clones, you can set the
`hub.http-clone` option to true using `git-config`.

For example:

    $ git clone defunkt/repl
    < git clone >
    $ git config --global --bool hub.http-clone true
    $ git clone defunkt/repl
    < http clone >

Or you can enter this manually into your `~/.gitconfig` file:

    $ cat ~/.gitconfig
    [hub]
      http-clone = yes


Prior Art
---------

These projects also aim to either improve git or make interacting with
GitHub simpler:

* [eg](http://www.gnome.org/~newren/eg/)
* [github-gem](http://github.com/defunkt/github-gem)
* [gh](http://github.com/visionmedia/gh)


Contributing
------------

Once you've made your great commits:

1. [Fork][0] hub
2. Create a topic branch - `git checkout -b my_branch`
3. Push to your branch - `git push origin my_branch`
4. Create an [Issue][1] with a link to your branch
5. That's it!

### Development Gems
You will need the following gems (and their dependencies) to
contribute to `hub`:

* `rake` (`gem install rake`)
* `kicker` (`gem install kicker`)
* `turn` (`gem install turn`)
* `mg` (`gem install mg`)
* `ronn` (`gem install ronn`)
* `webhelper` (`gem install webhelper`)

Meta
----

* Code: `git clone git://github.com/defunkt/hub.git`
* Home: <http://github.com/defunkt/hub>
* Bugs: <http://github.com/defunkt/hub/issues>
* List: <http://groups.google.com/group/github>
* Test: <http://runcoderun.com/defunkt/hub>
* Gems: <http://gemcutter.org/gems/git-hub>


Author
------

Chris Wanstrath :: chris@ozmm.org :: @defunkt

[0]: http://help.github.com/forking/
[1]: http://github.com/defunkt/hub/issues
[speed]: http://gist.github.com/284823
[2]: http://github.com/guides/local-github-config
