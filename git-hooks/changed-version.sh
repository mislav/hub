#!/bin/bash
## If version file is staged to be committed, display the new version number.

version_file="${1:-lib/hub/version.rb}"

if git rev-parse --verify HEAD >/dev/null 2>&1; then
  if git diff-index --quiet --cached HEAD -- "$version_file"; then
    exit 1
  else
    git cat-file blob ":$version_file" | grep VERSION | head -1 | cut -d\' -f2
  fi
else
  exit 1
fi
