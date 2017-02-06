# Installation instructions

## bash + Homebrew

If you're using Homebrew, just run `brew install hub` and you should be all set with auto-completion.

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
