hub(7) -- git + hub = github
============================

DESCRIPTION
-----------

`hub` is a command line utility which improves your `git` experience.

The goal is threefold:

* Augment existing `git` subcommands (such as `git clone`) with
  additional, often GitHub-aware functionality.
* Clarify many of git's famous error messages.
* Add new and useful subcommands.

`hub` can be used in place of `git` or you can alias the `git` command in
your shell to run `hub` - no existing functionality is removed. `hub`
simply adds and improves.


INSTALL
-------

hub can be installed using homebrew:

    $ brew install hub

It can also be installed with rubygems:

    $ gem install ron

Or installed from source:

    $ git clone git://github.com/defunkt/hub.git
    $ cd hub
    $ rake install

Once you've installed `hub`, you can invoke it directly from the
command line:

    $ hub --version

To get the full experience, alias your `git` command to run `hub` by
placing the following in your `.bash_profile` (or relevant startup
script):

    alias git=hub

Typing `hub install <shell>` will display install instructions for you
shell.


COMMANDS
--------

### git clone

    $ git clone schacon/ticgit
    $ git clone -p schacon/ticgit

### git remote add

    $ git remote add rtomayko
    $ git remote add -p pjhyett

### git init

   $ git init -g


PRIOR ART
---------

These projects also aim to either improve git or make interacting with
GitHub simpler:

* [eg](http://www.gnome.org/~newren/eg/)
* [github-gem](http://github.com/defunkt/github-gem)
* [gh](http://github.com/visionmedia/gh)


COPYRIGHT
---------

hub is Copyright (C) 2009 Chris Wanstrath
See the file COPYING for more information.
