# Installation instructions

## zsh

Create a new folder for completions:

```sh
mkdir -p ~/.zsh/completions
```

Download the file hub.zsh_completion and rename it to `_hub`:

```sh
curl https://github.com/github/hub/blob/master/etc/hub.zsh_completion >  ~/.zsh/completions/_hub
```

Then add the following lines to your .zshrc file:

```sh
fpath=(~/.zsh/completions $fpath) 
autoload -U compinit && compinit
```