hub is a command line tool that wraps `git` in order to extend it with extra
features and commands that make working with GitHub easier.

For an official, potentially more user-friendly command-line interface to GitHub,
see [cli.github.com](https://cli.github.com) and
[this comparison](https://github.com/cli/cli/blob/trunk/docs/gh-vs-hub.md).

This repository and its issue tracker is **not for reporting problems with
GitHub.com** web interface. If you have a problem with GitHub itself, please
[contact Support](https://github.com/contact).

Usage
-----

``` sh
$ hub clone rtomayko/tilt
#=> git clone git://github.com/rtomayko/tilt.git

# if you prefer HTTPS to git/SSH protocols:
$ git config --global hub.protocol https
$ hub clone rtomayko/tilt
#=> git clone https://github.com/rtomayko/tilt.git
```

See [usage examples](https://hub.github.com/#developer) or the [full reference
documentation](https://hub.github.com/hub.1.html) to see all available commands
and flags.

hub can also be used to make shell scripts that [directly interact with the
GitHub API](https://hub.github.com/#scripting).

hub can be safely [aliased](#aliasing) as `git`, so you can type `$ git
<command>` in the shell and have it expanded with `hub` features.

Installation
------------

The `hub` executable has no dependencies, but since it was designed to wrap
`git`, it's recommended to have at least **git 1.7.3** or newer.

platform | manager | command to run
---------|---------|---------------
macOS, Linux | [Homebrew](https://docs.brew.sh/Installation) | `brew install hub`
Windows | [Scoop](http://scoop.sh/) | `scoop install hub`
Windows | [Chocolatey](https://chocolatey.org/) | `choco install hub`
Fedora Linux | [DNF](https://fedoraproject.org/wiki/DNF) | `sudo dnf install hub`
Arch Linux | [pacman](https://wiki.archlinux.org/index.php/pacman) | `sudo pacman -S hub`
FreeBSD | [pkg(8)](http://man.freebsd.org/pkg/8) | `pkg install hub`
Debian | [apt(8)](https://manpages.debian.org/buster/apt/apt.8.en.html) | `sudo apt install hub`
Ubuntu | [Snap](https://snapcraft.io) | [We do not recommend installing the snap anymore.](https://github.com/github/hub/issues?q=is%3Aissue+snap)
openSUSE | [Zypper](https://en.opensuse.org/SDB:Zypper_manual) | `sudo zypper install hub`
Void Linux | [xbps](https://github.com/void-linux/xbps) | `sudo xbps-install -S hub`
Gentoo | [Portage](https://wiki.gentoo.org/wiki/Portage) | `sudo emerge dev-vcs/hub`
_any_ | [conda](https://docs.conda.io/en/latest/) | `conda install -c conda-forge hub`


Packages other than Homebrew are community-maintained (thank you!) and they
are not guaranteed to match the [latest hub release][latest]. Check `hub
version` after installing a community package.

#### Standalone

`hub` can be easily installed as an executable. Download the [latest
binary][latest] for your system and put it anywhere in your executable path.

#### GitHub Actions

hub is ready to be used in your [GitHub Actions][] workflows:
```yaml
steps:
- uses: actions/checkout@v2

- name: List open pull requests
  run: hub pr list
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Note that the default `secrets.GITHUB_TOKEN` will only work for API operations
scoped to the repository that runs this workflow. If you need to interact with other
repositories, [generate a Personal Access Token][pat] with at least the `repo` scope
and add it to your [repository secrets][].


[github actions]: https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions
[pat]: https://github.com/settings/tokens
[repository secrets]: https://docs.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets

#### Source

Prerequisites for building from source are:

* `make`
* [Go 1.11+](https://golang.org/doc/install)

Clone this repository and run `make install`:

```sh
git clone \
  --config transfer.fsckobjects=false \
  --config receive.fsckobjects=false \
  --config fetch.fsckobjects=false \
  https://github.com/github/hub.git

cd hub
make install prefix=/usr/local
```

Aliasing
--------

Some hub features feel best when it's aliased as `git`. This is not dangerous; your
_normal git commands will all work_. hub merely adds some sugar.

`hub alias` displays instructions for the current shell. With the `-s` flag, it
outputs a script suitable for `eval`.

You should place this command in your `.bash_profile` or other startup script:

``` sh
eval "$(hub alias -s)"
```

#### PowerShell

If you're using PowerShell, you can set an alias for `hub` by placing the
following in your PowerShell profile (usually
`~/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1`):

``` sh
Set-Alias git hub
```

A simple way to do this is to run the following from the PowerShell prompt:

``` sh
Add-Content $PROFILE "`nSet-Alias git hub"
```

Note: You'll need to restart your PowerShell console in order for the changes to be picked up.

If your PowerShell profile doesn't exist, you can create it by running the following:

``` sh
New-Item -Type file -Force $PROFILE
```

### Shell tab-completion

hub repository contains [tab-completion scripts](./etc) for bash, zsh and fish.
These scripts complement existing completion scripts that ship with git.

Meta
----

* Bugs: <https://github.com/github/hub/issues>
* Authors: <https://github.com/github/hub/contributors>
* Our [Code of Conduct](https://github.com/github/hub/blob/master/CODE_OF_CONDUCT.md)


[latest]: https://github.com/github/hub/releases/latest
