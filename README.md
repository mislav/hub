hub: git + hub = github
=======================

`hub` is a command line utility which improves your `git` experience.

The goal is threefold:

* Augment existing `git` subcommands (such as `git clone`) with
  additional, often GitHub-aware functionality.
* Clarify many of git's famous error messages.
* Add new and useful subcommands.

`hub` can be used in place of `git` or you can alias the `git` command in
your shell to run `hub` - no existing functionality is removed. `hub`
simply adds and improves.


Install
-------

hub can be installed using rubygems:

    $ gem install hub -s http://gemcutter.org/

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


Commands
--------

### git clone

    $ git clone schacon/ticgit
    $ git clone -p schacon/ticgit

### git remote add

    $ git remote add rtomayko
    $ git remote add -p pjhyett

### git init

    $ git init -g


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


Meta
----

* Code: `git clone git://github.com/defunkt/hub.git`
* Home: <http://github.com/defunkt/hub>
* Docs: <http://defunkt.github.com/hub/>
* Bugs: <http://github.com/defunkt/hub/issues>
* List: <http://groups.google.com/group/github>
* Gems: <http://gemcutter.org/gems/hub>


Author
------

Chris Wanstrath :: chris@ozmm.org :: @defunkt

[0]: http://help.github.com/forking/
[1]: http://github.com/defunkt/hub/issues
