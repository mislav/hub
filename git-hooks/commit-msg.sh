#!/bin/bash
## Checks that commit that changes version number includes it in subject line.

set -e

new_version="$("$(dirname "$0")"/changed-version)" || exit 0

head -1 "$1" | grep "$new_version" >/dev/null || {
  echo "aborted: version $new_version not present in subject line"
  exit 1
}
