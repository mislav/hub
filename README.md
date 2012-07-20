git + hub = github
==================

hub is a command line tool that wraps `git` in order to extend it with extra
features and commands that make working with GitHub easier.

~~~ sh
$ hub clone rtomayko/tilt

# expands to:
$ git clone git://github.com/rtomayko/tilt.git
~~~

hub is best aliased as `git`, so you can type `$ git <command>` in the shell and
get all the usual `hub` features. See "Aliasing" below.


Installation
------------

Dependencies:

* **git 1.7.3** or newer
* **Ruby 1.8.6** or newer

### Homebrew

Installing on OS X is easiest with Homebrew:

~~~ sh
$ brew install hub
~~~

### Standalone

`hub` is easily installed as a standalone script:

~~~ sh
$ curl http://defunkt.io/hub/standalone -sLo ~/bin/hub &&
  chmod +x ~/bin/hub
~~~

Assuming "~/bin/" is in your `$PATH`, you're ready to roll:

~~~ sh
$ hub version
git version 1.7.6
hub version 1.8.3
~~~

### RubyGems

Though not recommended, hub can also be installed as a RubyGem:

~~~ sh
$ gem install hub
~~~

(It's not recommended for casual use because of the RubyGems startup
time. See [this gist][speed] for information.)

#### Standalone via RubyGems

~~~ sh
$ gem install hub
$ hub hub standalone > ~/bin/hub && chmod +x ~/bin/hub
~~~

This installs a standalone version which doesn't require RubyGems to
run, so it's faster.

### Source

You can also install from source:

~~~ sh
$ git clone git://github.com/defunkt/hub.git
$ cd hub
$ rake install prefix=/usr/local
~~~

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

Using hub feels best when it's aliased as `git`. This is not dangerous; your
_normal git commands will all work_. hub merely adds some sugar.

`hub alias` displays instructions for the current shell. With the `-s` flag, it
outputs a script suitable for `eval`.

You should place this command in your `.bash_profile` or other startup script:

~~~ sh
eval "$(hub alias -s)"
~~~

### Shell tab-completion

hub repository contains tab-completion scripts for bash and zsh. These scripts
complement existing completion scripts that ship with git.

* [hub bash completion](https://github.com/defunkt/hub/blob/master/etc/hub.bash_completion.sh)
* [hub zsh completion](https://github.com/defunkt/hub/blob/master/etc/hub.zsh_completion)


Commands
--------

Assuming you've aliased hub as `git`, the following commands now have
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

### git checkout

    $ git checkout https://github.com/defunkt/hub/pull/73
    > git remote add -f -t feature git://github:com/mislav/hub.git
    > git checkout --track -B mislav-feature mislav/feature

    $ git checkout https://github.com/defunkt/hub/pull/73 custom-branch-name

### git merge

    $ git merge https://github.com/defunkt/hub/pull/73
    > git fetch git://github.com/mislav/hub.git +refs/heads/feature:refs/remotes/mislav/feature
    > git merge mislav/feature --no-ff -m 'Merge pull request #73 from mislav/feature...'

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


Configuration
-------------

### GitHub OAuth authentication

Hub will prompt for GitHub username & password the first time it needs to access
the API and exchange it for an OAuth token, which it saves in "~/.config/hub".

### HTTPS instead of git protocol

If you prefer using the HTTPS protocol for GitHub repositories instead of the git
protocol for read and ssh for write, you can set "hub.protocol" to "https".

~~~ sh
# default behavior
$ git clone defunkt/repl
< git clone >

# opt into HTTPS:
$ git config --global hub.protocol https
$ git clone defunkt/repl
< https clone >
~~~


Contributing
------------

These instructions assume that _you already have hub installed_ and aliased as
`git` (see "Aliasing").

1. Clone hub:  
    `git clone defunkt/hub && cd hub`
1. Ensure Bundler is installed:  
    `which bundle || gem install bundler`
1. Install development dependencies:  
    `bundle install`
2. Verify that existing tests pass:  
    `bundle exec rake`
3. Create a topic branch:  
    `git checkout -b feature`
4. **Make your changes.** (It helps a lot if you write tests first.)
5. Verify that tests still pass:  
    `bundle exec rake`
6. Fork hub on GitHub (adds a remote named "YOUR_USER"):  
    `git fork`
7. Push to your fork:  
    `git push -u YOUR_USER feature`
8. Open a pull request describing your changes:  
    `git pull-request`


Meta
----

* Home: <https://github.com/defunkt/hub>
* Bugs: <https://github.com/defunkt/hub/issues>
* Gem: <https://rubygems.org/gems/hub>
* Authors: <https://github.com/defunkt/hub/contributors>

### Prior art

These projects also aim to either improve git or make interacting with
GitHub simpler:

* [eg](http://www.gnome.org/~newren/eg/)
* [github-gem](https://github.com/defunkt/github-gem)


[speed]: http://gist.github.com/284823
[gc]: https://twitter.com/brynary/status/49560668994674688
