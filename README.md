hub: git + hub = github
=======================

`hub` is a command line utility which adds GitHub knowledge to `git`.

It can be used on its own or as a `git` wrapper.

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

    curl http://defunkt.io/hub/standalone -sLo ~/bin/hub &&
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

    $ gem install hub

(It's not recommended for casual use because of the RubyGems startup
time. See [this gist][speed] for information.)

### Standalone via RubyGems

    $ gem install hub
    $ hub hub standalone > ~/bin/hub && chmod 755 ~/bin/hub

This installs a standalone version which doesn't require RubyGems to
run.

### Source

You can also install from source:

    $ git clone git://github.com/defunkt/hub.git
    $ cd hub
    $ rake install prefix=/usr/local

### Help! It's Slow!

Is your prompt slow? It may be hub.

1. Check that it's **not** installed using RubyGems.
2. Check that RUBYOPT isn't loading anything shady:

       $ echo $RUBYOPT

3. Check that your system Ruby is speedy:

       $ time /usr/bin/env ruby -e0

If #3 is slow, it may be your [GC settings][gc].


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

Assuming you've aliased `hub` to `git`, the following commands now have
superpowers:

### git clone

    $ git clone schacon/ticgit
    > git clone git://github.com/schacon/ticgit.git

    $ git clone -p schacon/ticgit
    > git clone git@github.com:schacon/ticgit.git

    $ git clone resque
    > git clone git@github.com/YOUR_USER/resque.git

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

### git am, git apply

    $ git am https://github.com/defunkt/hub/pull/55
    > curl https://github.com/defunkt/hub/pull/55.patch -o /tmp/55.patch
    > git am /tmp/55.patch

    $ git am --ignore-whitespace https://github.com/davidbalbert/hub/commit/fdb9921
    > curl https://github.com/davidbalbert/hub/commit/fdb9921.patch -o /tmp/fdb9921.patch
    > git am --ignore-whitespace /tmp/fdb9921.patch

    $ git apply https://gist.github.com/8da7fb575debd88c54cf
    > curl https://gist.github.com/8da7fb575debd88c54cf.txt -o /tmp/gist-8da7fb575debd88c54cf.txt
    > git apply /tmp/gist-8da7fb575debd88c54cf.txt

### git fork

    $ git fork
    [ repo forked on GitHub ]
    > git remote add -f YOUR_USER git@github.com:YOUR_USER/CURRENT_REPO.git

### git pull-request

    # while on a topic branch called "feature":
    $ git pull-request
    [ opens text editor to edit title & body for the request ]
    [ opened pull request on GitHub for "YOUR_USER:feature" ]

    # explicit title, pull base & head:
    $ git pull-request "I've implemented feature X" -b defunkt:master -h mislav:feature

    $ git pull-request -i 123
    [ attached pull request to issue #123 ]

    # while on a topic branch called "feature"
    # and while the base project's user is "defunkt"
    $ git pull-request -s
    [ opened pull request on GitHub for "defunkt:feature"]

### git checkout

    # $ git checkout https://github.com/defunkt/hub/pull/73
    # > git remote add -f -t feature git://github:com/mislav/hub.git
    # > git checkout -b mislav-feature mislav/feature

    # $ git checkout https://github.com/defunkt/hub/pull/73 custom-branch-name

### git create

    $ git create
    [ repo created on GitHub ]
    > git remote add origin git@github.com:YOUR_USER/CURRENT_REPO.git

    # with description:
    $ git create -d 'It shall be mine, all mine!'

    $ git create recipes
    [ repo created on GitHub ]
    > git remote add origin git@github.com:YOUR_USER/recipes.git

    $ git create sinatra/recipes
    [ repo created in GitHub organization ]
    > git remote add origin git@github.com:sinatra/recipes.git

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
    > open https://github.com/YOUR_USER/CURRENT_REPO

    $ git browse -- commit/SHA
    > open https://github.com/YOUR_USER/CURRENT_REPO/commit/SHA

    $ git browse -- issues
    > open https://github.com/YOUR_USER/CURRENT_REPO/issues

    $ git browse schacon/ticgit
    > open https://github.com/schacon/ticgit

    $ git browse schacon/ticgit commit/SHA
    > open https://github.com/schacon/ticgit/commit/SHA

    $ git browse resque
    > open https://github.com/YOUR_USER/resque

    $ git browse resque network
    > open https://github.com/YOUR_USER/resque/network

### git compare

    $ git compare refactor
    > open https://github.com/CURRENT_REPO/compare/refactor

    $ git compare 1.0..1.1
    > open https://github.com/CURRENT_REPO/compare/1.0...1.1

    $ git compare -u fix
    > (https://github.com/CURRENT_REPO/compare/fix)

    $ git compare other-user patch
    > open https://github.com/other-user/REPO/compare/patch

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
is stored locally in your Git config or environment variables.

To test it run this:

    $ git config --global github.user

If you see nothing, you need to set the config setting:

    $ git config --global github.user YOUR_USER

For commands that require write access to GitHub (such as `fork`), you'll want to
setup "github.token" as well. See [GitHub config guide][2] for more information.

If present, environment variables `GITHUB_USER` and `GITHUB_TOKEN` override the
values of "github.user" and "github.token".

Configuration
-------------

If you prefer using the HTTPS protocol for GitHub repositories instead of the git
protocol for read and ssh for write, you can set "hub.protocol" to "https".

For example:

    $ git clone defunkt/repl
    < git clone >
    
    $ git config --global hub.protocol https
    $ git clone defunkt/repl
    < https clone >

Prior Art
---------

These projects also aim to either improve git or make interacting with
GitHub simpler:

* [eg](http://www.gnome.org/~newren/eg/)
* [github-gem](https://github.com/defunkt/github-gem)


Contributing
------------

These instructions assume that you already have `hub` installed and that
you've set it up so it wraps `git` (see "Aliasing").

1. Clone hub:  
    `git clone defunkt/hub`
2. Verify that existing tests pass (see "Development dependencies"):  
    `rake test`
3. Create a topic branch:  
    `git checkout -b my_branch`
4. Make your changes – it helps a lot if you write tests first
5. Verify that tests still pass:  
    `rake test`
6. Fork hub on GitHub (adds a remote named "YOUR_USER"):  
    `git fork`
7. Push to your fork:  
    `git push -u YOUR_USER my_branch`
8. Open a pull request describing your changes:  
    `git pull-request`

### Development dependencies

You will need the following libraries for development:

* [ronn](https://github.com/rtomayko/ronn) (building man pages)
* [webmock](https://github.com/bblimke/webmock)

Meta
----

* Home: <https://github.com/defunkt/hub>
* Bugs: <https://github.com/defunkt/hub/issues>
* Gem: <https://rubygems.org/gems/hub>


Authors
-------

<https://github.com/defunkt/hub/contributors>

[speed]: http://gist.github.com/284823
[2]: http://help.github.com/set-your-user-name-email-and-github-token/
[gc]: https://twitter.com/brynary/status/49560668994674688
