## 1.10.2 (2012-07-24)

* fix pushing multiple refs to multiple remotes
* support ssh.github.com
* security: mode for ~/.config/hub is 0600
* fix integration with GitHub Enterprise API
* fix cloning repos that start with a period

## 1.10.1 (2012-05-28)

* don't choke on empty config file
* fix `browse` when not in git repo

## 1.10.0 (2012-05-08)

* improve improved help text
* fix GitHub username prompt in `create` command
* make `fetch` command work with private repos too
* add `merge` command to merge pull requests

## 1.9.0 (2012-05-04)

* internal refactoring and code reorganization
* switch to GitHub API v3 and authenticate via OAuth
* auth info is now stored in ~/.config/hub instead of ~/.gitconfig

## 1.8.4 (2012-03-20)

* add bash, zsh completion
* improve `hub alias` command
* change `git fork` so it fails when repo already exists under user
* teach custom commands to respect `-h` & `--help` flags
* `pull-request`: better error message for invalid remotes/URLs
* respect local SSH aliases for host names

## 1.8.3 (2012-03-02)

* fix `pull-request` from branch tracking another local branch
* fix `browse` command when not on any branch

## 1.8.2 (2012-02-07)

* if `pull-request` editor is vim, set appropriate filetype
* `pull-request` editor message defaults to single commit message
* fix cherry-picking from an existing remote
* fix `clone` from local repository
* `checkout` command forwards flags to internal checkout command,
  force-resets the existing local branch by default
* fix `am` command when given URLs that include the fragment

## 1.8.1 (2012-01-24)

* fix JSON parsing error while using GitHub API
* HTTP(S) proxy support

## 1.8.0 (2012-01-03)

* fix `pull-request` on GH Enterprise project branch without upstream
* ensure Content-Length for POST requests
* handle pull requests from private repos
* support branches with slashes in their name
* display server errors when creating pullrequest fails
* support GitHub Enterprise via multiple whitelisted host names
* GitHub remote urls don't have to necessarily end in ".git"
* fix `git init -g`
* authenticate all API requests, helps dealing with private repos
* ensure periods are allowed in repository names
* fix am/apply commands if TMPDIR environment variable isn't set

## 1.7.0 (2011-11-24)

* lock down standalone script to system ruby
* don't try to use command output pager on Windows
* opt in for HTTPS: `git config hub.protocol https`
* add `hub pull-request`
* improve detecting upstream configuration ("tracking" branches)
* introduce `hub --noop`
* `hub apply` now downloads GitHub patches same as `hub am`
* `hub create <name>` to explicitly name a repository
* switch API communication to HTTPS
* better handling of API HTTP exceptions
* replace two dots (`sha1..sha2`) with three for ranges in `compare`
* avoid ugly error & stack trace when git is not found on the system

## 1.6.1 (2011-05-13)

* `git push remote1,remote2` without branch name pushes the current branch
* fix `browse` command for current repo with no tracking setup
* preserve global flags to git such as `--bare` and `--git-dir=/some/path`
* true cross-platform command detection and browser launcher

## 1.6.0 (2011-03-24)

* `am` strips extra path from pull reqs URLs such as "pull/42/files"
* Fixed permissions on `hub(1)` when installing
* gem renamed from `git-hub` to `hub`!

## 1.5.1 (2011-03-18)

* support git aliases
* Bugfix: `browse/compare` for wiki repos
* gracefully handle HTTP errors in `create` and `fork`
* `hub am` supports Gist URLs
* Bugfix: `clone` command doesn't get confused by mixed arguments

## 1.5.0 (2010-12-28)

* compensate for GitHub switch to HTTPS
* `hub am`: cherry-pick pull request and commit URLs
* support multiple URLs for a single remote
* Bugfix: ensure that internal ruby methods can't pretend to be git commands
* Bugfix: don't show help when `--exec-path` or `--html-path` flags are used
* Support for `GITHUB_USER` and `GITHUB_TOKEN` env variables
* Eliminate some ruby warnings

## 1.4.1 (2010-08-08)
## 1.4.0 (2010-08-08)

* Added new `hub create` command
* Added support for `remote set-url`
* Bugfix: Don't try multiple git commands on a non-git dir when grabbing remote
* Bugfix: Adding remotes when no remotes exist

## 1.3.2 (2010-07-24)

* bugfix: cherry picking of commit URL
* bugfix: git init -g

## 1.3.1 (2010-04-29)
## 1.3.0 (2010-04-29)

* Tracking branches awareness
* `git browse` subpages (e.g. `git browse repo issues`)
* `git fetch <fork>` adds new remotes if missing
* `cherry-pick` supports GitHub commit URLs and "user@sha" notation

## 1.2.0 (2010-04-11)

* `hub compare` command - Thanks joshthecoder!

## 1.1.0 (2010-04-07)

* `hub fork` command - Thanks Mislav!

## 1.0.3 (2010-03-10)

* Bugfix: `hub remote` for repos with -, /, etc

## 1.0.2 (2010-03-07)

* Bugfix: `hub remote -f name` (for real this time)
* Bugfix: zsh quoting [thommay]

## 1.0.1 (2010-03-05)

* Bugfix: `hub remote -f name`

## 1.0.0 (2010-03-03)

* `hub browse` with no arguments uses the current repo.
* `hub submodule add user/repo directory
* `hub remote add rtomayko/tilt`
* `remote add -p origin rtomayko/tilt`

## 0.3.2 (2010-02-17)

* Fixed zshell git completion / aliasing - `hub alias zsh`.

## 0.3.1 (2010-02-13)

* Add `hub remote origin` shortcut. Assumes your GitHub login.

## 0.3.0 (2010-01-23)

* Add `hub browse` command for opening a repo in a browser.
* Add `hub standalone` for installation of standalone via RubyGems
* Bugfix: Don't run hub standalone in standalone mode
* Bugfix: `git clone` flags are now passed through.
* Bugfix: `git clone` with url and path works.
* Bugfix: basename call
* Bugfix: Check for local directories before cloning


## 0.2.0 (2009-12-24)

* Respected GIT_PAGER and core.pager
* Aliased `--help` to `help`
* Ruby 1.9 fixes
* Respect git behavior when pager is empty string
* `git push` multi-remote support
* `hub.http-clone` configuration setting
* Use the origin url to find the repo name

## 0.1.3 (2009-12-11)

* Homebrew!
* Fix inaccuracy in man page
* Better help arrangement
* Bugfix: Path problems in standalone.rb
* Bugfix: Standalone not loaded by default

## 0.1.2 (2009-12-10)

* Fixed README typos
* Better standalone install line
* Added man page
* Added `hub help hub`

## 0.1.1 (2009-12-08)

* Fixed gem problems

## 0.1.0 (2009-12-08)

* First release
