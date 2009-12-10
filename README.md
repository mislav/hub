hub: git + hub = github
=======================

`hub` is a command line utility which adds GitHub knowledge to `git`.

It can used on its own or as a `git` wrapper.

Normal:

    $ hub clone rtomayko/tilt
    Initialized empty Git repository in /Users/chris/sandbox/tilt/.git/
    remote: Counting objects: 307, done.
    remote: Compressing objects: 100% (219/219), done.
    remote: Total 307 (delta 175), reused 85 (delta 45)
    Receiving objects: 100% (307/307), 48.91 KiB, done.
    Resolving deltas: 100% (175/175), done.

Wrapping `git`:

    $ git clone rack/rack
    Initialized empty Git repository in /Users/chris/sandbox/rack/.git/
    remote: Counting objects: 4005, done.
    remote: Compressing objects: 100% (1738/1738), done.
    remote: Total 4005 (delta 2505), reused 3620 (delta 2208)
    Receiving objects: 100% (4005/4005), 785.82 KiB | 129 KiB/s, done.
    Resolving deltas: 100% (2505/2505), done.

hub requires you have `git` installed and in your `$PATH`. It also
requires Ruby 1.8.6+ or Ruby 1.9.1+. No other libraries necessary.


Install
-------

### Standalone

`hub` is most easily installed as a standalone script:

    curl -s http://defunkt.github.com/hub/standalone > ~/bin/hub && chmod 755 !#:4

Assuming `~/bin/` is in your `$PATH`, you're ready to roll:

    $ hub version
    git version 1.6.4.2
    hub version 0.1.0

### RubyGems

Though not recommended, `hub` can also be installed as a RubyGem:

    $ gem install git-hub -s http://gemcutter.org/

(Yes, the gem name is `git-hub`.)

### Source

You can also install from source:

    $ git clone git://github.com/defunkt/hub.git
    $ cd hub
    $ rake standalone
    $ cp hub /usr/local/bin/


Aliasing
--------

hub works best when it wraps `git`. This is not dangerous - your
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

### git init

    $ git init -g
    > git init
    > git remote add origin git@github.com:YOUR_USER/REPO.git

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

See <http://github.com/guides/local-github-config> for more information.


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
* Bugs: <http://github.com/defunkt/hub/issues>
* List: <http://groups.google.com/group/github>
* Test: <http://runcoderun.com/defunkt/hub>
* Gems: <http://gemcutter.org/gems/git-hub>


Author
------

Chris Wanstrath :: chris@ozmm.org :: @defunkt

[0]: http://help.github.com/forking/
[1]: http://github.com/defunkt/hub/issues
