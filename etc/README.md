# Installation instructions

## bash

Open your `.bashrc` file and add:

> If you want to set up Git to automatically have Bash shell completion for all users, copy the `hub.bash_completion` script to the `/opt/local/etc/bash_completion.d` directory on Mac systems or to the `/etc/bash_completion.d/` directory on Linux systems. ([Source](https://git-scm.com/book/en/v1/Git-Basics-Tips-and-Tricks#Auto-Completion))

* [Link to git auto-completion bash file](https://github.com/git/git/blob/master/contrib/completion/git-completion.bash)

```sh
# Make sure you've aliased hub to git
eval "$(hub alias -s)"

# And make sure that the git auto-completion is being loaded
if [ -f /path/to/git-completion.bash ]; then
    . /path/to/git-completion.bash
fi

# Load hub autocompletion
if [ -f /path/to/hub.bash_completion ]; then
    . /path/to/hub.bash_completion
fi
```

## zsh

Create a new folder for completions:

```sh
mkdir -p ~/.zsh/completions
```

Copy the file `/etc/hub.zsh_completion` from the location where you downloaded `hub` to the folder `~/.zsh/completions/` and rename it to `_hub`:

```sh
cp /path/to/etc/hub.zsh_completion ~/.zsh/completions/ \
    mv ~/.zsh/completions/hub.zsh_completion ~/.zsh/completions/_hub
```

Then add the following lines to your `.zshrc` file:

```sh
fpath=(~/.zsh/completions $fpath) 
autoload -U compinit && compinit
```
