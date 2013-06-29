#!/bin/bash

# git-sh (https://github.com/rtomayko/git-sh) support for hub.

eval "$(hub alias -s)"
source /usr/local/etc/hub.bash_completion.sh

gitalias pull-request="git pull-request"
gitalias fork="git fork"
gitalias create="git create"
gitalias browse="git browse"
gitalias compare="git compare"
gitalias ci-status="hub ci-status"
