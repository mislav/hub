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

### `rake install` from source

This is the preferred installation method when no package manager that
supports hub is available:

~~~ sh
# Download or clone the project from GitHub:
$ git clone git://github.com/github/hub.git
$ cd hub
$ rake install
~~~

On a Unix-based OS, this installs under `PREFIX`, which is `/usr/local` by default.

Now you should be ready to roll:

~~~ sh
$ hub version
git version 1.7.6
hub version 1.8.3
~~~

#### Windows "Git Bash" (msysGit) note

Avoid aliasing hub as `git` due to the fact that msysGit automatically
configures your prompt to include git information, and you want to avoid slowing
that down. See [Is your shell prompt slow?](#is-your-shell-prompt-slow)

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

### Help! It's slow!

#### Is `hub` noticeably slower than plain git?

That is inconvenient, especially if you want to alias hub as `git`. Few things
you can try:

* Find out which ruby is used for the hub executable:

    ``` sh
    head -1 `which hub`
    ```

* That ruby should be speedy. Time it with:

    ``` sh
    time /usr/bin/ruby -e0
    #=> it should be below 0.01 s total
    ```

* Check that Ruby isn't loading something shady:

    ``` sh
    echo $RUBYOPT
    ```

* Check your [GC settings][gc]

General recommendation: you should change hub's shebang line to run with system
ruby (usually `/usr/bin/ruby`) instead of currently active ruby (`/usr/bin/env
ruby`). Also, Ruby 1.8 is speedier than 1.9.

#### Is your shell prompt slow?

Does your prompt show git information? Hub may be slowing down your prompt.

This can happen if you've aliased hub as `git`. This is fine when you use `git`
manually, but may be unacceptable for your prompt, which doesn't need hub
features anyway!

The solution is to identify which shell functions are calling `git`, and replace
each occurrence of that with `command git`. This is a shell feature that enables
you to call a command directly and skip aliases and functions wrapping it.


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

* [hub bash completion](https://github.com/github/hub/blob/master/etc/hub.bash_completion.sh)
* [hub zsh completion](https://github.com/github/hub/blob/master/etc/hub.zsh_completion)


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
    [ downloads patch via API ]
    > git am /tmp/55.patch

    $ git am --ignore-whitespace https://github.com/davidbalbert/hub/commit/fdb9921
    [ downloads patch via API ]
    > git am --ignore-whitespace /tmp/fdb9921.patch

    $ git apply https://gist.github.com/8da7fb575debd88c54cf
    [ downloads patch via API ]
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
    $ git pull-request -m "Implemented feature X" -b defunkt:master -h mislav:feature

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

    $ git submodule add wycats/bundler vendor/bundler
    > git submodule add git://github.com/wycats/bundler.git vendor/bundler

    $ git submodule add -p wycats/bundler vendor/bundler
    > git submodule add git@github.com:wycats/bundler.git vendor/bundler

    $ git submodule add -b ryppl --name pip ryppl/pip vendor/pip
    > git submodule add -b ryppl --name pip git://github.com/ryppl/pip.git vendor/pip

### git ci-status

    $ git ci-status [commit]
    > (prints CI state of commit and exits with appropriate code)
    > One of: success (0), error (1), failure (1), pending (2), no status (3)


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


Meta
----

* Home: <https://github.com/github/hub>
* Bugs: <https://github.com/github/hub/issues>
* Gem: <https://rubygems.org/gems/hub>
* Authors: <https://github.com/github/hub/contributors>

### Prior art

These projects also aim to either improve git or make interacting with
GitHub simpler:

* [eg](http://www.gnome.org/~newren/eg/)
* [github-gem](https://github.com/defunkt/github-gem)


[speed]: http://gist.github.com/284823
[gc]: https://twitter.com/brynary/status/49560668994674688
