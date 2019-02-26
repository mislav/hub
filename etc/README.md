# Installation instructions

## Homebrew

If you're using Homebrew, just run `brew install hub` and you should be all set with auto-completion (for both Bash and Zsh).

### Zsh

hub's Zsh completions are [designed to work with Zsh's official Git completions](https://github.com/github/hub/pull/295#issuecomment-14324233), *not* Git's official Zsh completions.

If you installed Zsh and Git via Homebrew, you will need to delete Git's official Zsh completions to allow Zsh's official Git completions to work, which in turn will allow hub's Zsh completions to work.

You can do this deletion in `.zshrc`, which ensures the file remains deleted even after upgrades of Git.

``` sh
# Delete Git's official completions to allow Zsh's official Git completions to work.
# This is also necessary for hub's Zsh completions to work:
# https://github.com/github/hub/issues/1956.
function () {
  GIT_ZSH_COMPLETIONS_FILE_PATH="$(brew --prefix)/share/zsh/site-functions/_git"
  if [ -f $GIT_ZSH_COMPLETIONS_FILE_PATH ]
  then
    rm $GIT_ZSH_COMPLETIONS_FILE_PATH
  fi
}
```

Note that Zsh's official Git completions are much more powerful than Git's official Zsh completions, so they are generally preferred: https://github.com/Homebrew/homebrew-core/issues/33275#issuecomment-432528793.

## bash

Open your `.bashrc` file if you're on Linux, or your `.bash_profile` if you're on macOS and add:

```sh
if [ -f /path/to/hub.bash_completion ]; then
  . /path/to/hub.bash_completion
fi
```

## zsh

Copy the file `etc/hub.zsh_completion` from the location where you downloaded
`hub` to the folder `~/.zsh/completions/` and rename it to `_hub`:

```sh
mkdir -p ~/.zsh/completions
cp etc/hub.zsh_completion ~/.zsh/completions/_hub
```

Then add the following lines to your `.zshrc` file:

```sh
fpath=(~/.zsh/completions $fpath) 
autoload -U compinit && compinit
```

## fish

Copy the file `etc/hub.fish_completion` from the location where you downloaded
`hub` to the folder `~/.config/fish/completions/` and rename it to `hub.fish`:

```sh
mkdir -p ~/.config/fish/completions
cp etc/hub.fish_completion ~/.config/fish/completions/hub.fish
```
