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

### 2.x

Starting from 2.2.0, `hub` is [powered by Go](https://github.com/github/hub/issues/475).
It should maintain a high backward compatibility with 1.x. Note that 2.x is in preview.

#### Homebrew

`hub` 2.x can be installed through Homebrew. Since it's in preview, you need to compile
from HEAD:

~~~ sh
$ brew install --HEAD hub
~~~

#### Standalone

`hub` 2.x can be easily installed as an executable. Download the latest
[compiled binaries](https://github.com/github/hub/releases) and put it anywhere
in your executable path.

#### Source

To install `hub` 2.x from source, you need to have a [Go development environment](http://golang.org/doc/install),
version 1.4 or better:

~~~ sh
$ git clone https://github.com/github/hub.git
$ cd hub
$ ./script/build
$ cp hub YOUR_BIN_PATH
~~~

Or, if you've done Go development before and your $GOPATH/bin
directory is already in your PATH:

~~~ sh
$ go get github.com/github/hub
~~~

### 1.x

#### Homebrew

Installing on OS X is easiest with Homebrew:

~~~ sh
$ brew install hub
~~~

Now you should be ready to roll:

~~~ sh
$ hub version
git version 1.7.6
hub version 1.8.3
~~~


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
    > git remote add -f -t feature mislav git://github.com/mislav/hub.git
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

    $ git browse -- issues/10
    > open https://github.com/YOUR_USER/CURRENT_REPO/issues/10

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

To avoid being prompted, use **GITHUB_USER** and **GITHUB_PASSWORD** environment
variables.

### HTTPS instead of git protocol

If you prefer the HTTPS protocol for GitHub repositories, you can set
"hub.protocol" to "https". This will affect `clone`, `fork`, `remote add`
and other operations that expand references to GitHub repositories as full
URLs that otherwise use git and ssh protocols.

"hub.protocol" only applies when the "OWNER/REPO" shorthand is used instead of
a full git URL.

    # default behavior
    $ git clone defunkt/repl
    < git clone >

    # opt into HTTPS:
    $ git config --global hub.protocol https
    $ git clone defunkt/repl
    < https clone >

### GitHub Enterprise

By default, hub will only work with repositories that have remotes which
point to github.com. GitHub Enterprise hosts need to be whitelisted to
configure hub to treat such remotes same as github.com:

    $ git config --global --add hub.host my.git.org

The default host for commands like `init` and `clone` is still
github.com, but this can be affected with the `GITHUB_HOST` environment
variable:

    $ GITHUB_HOST=my.git.org git clone myproject

Meta
----

* Home: <https://github.com/github/hub>
* Bugs: <https://github.com/github/hub/issues>
* Authors: <https://github.com/github/hub/contributors>

### Prior art

These projects also aim to either improve git or make interacting with
GitHub simpler:

* [hub in Ruby](https://github.com/github/hub/tree/1.12-stable) (previous implementation)
* [eg](http://www.gnome.org/~newren/eg/)
* [github-gem](https://github.com/defunkt/github-gem)
