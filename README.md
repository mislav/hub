hub: git + hub = github
=======================

`hub` is a command line utility which injects `git` with GitHub
knowledge.

It can used on its own or can serve as a complete, backwards
compatible replacement for the `git` script.


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
    > git clone git://github.com/schacon/ticgit.git

    $ git clone -p schacon/ticgit
    > git clone git@github.com:schacon/ticgit.git

### git remote add

    $ git remote add rtomayko
    > git remote add rtomayko git://github.com/rtomayko/CURRENT_REPO.git

    $ git remote add -p pjhyett
    > git remote add rtomayko git@github.com:rtomayko/CURRENT_REPO.git

### git init

    $ git init -g
    > git init
    > git remote add origin git@github.com:USER/REPO.git

### git help

    $ git help
    > (improved git help)


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
