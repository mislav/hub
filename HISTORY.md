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
