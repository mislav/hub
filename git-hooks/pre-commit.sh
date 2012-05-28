#!/bin/bash
## Prevents a commit that changes version but doesn't update HISTORY.md

set -e

new_version="$("$(dirname "$0")"/changed-version)" || exit 0

if git cat-file blob :HISTORY.md | head -1 | grep $new_version >/dev/null; then
  exit 0
else
  echo "aborted: History.md is out of date" >&2
  exit 1
fi
