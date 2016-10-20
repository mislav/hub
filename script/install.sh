#!/usr/bin/env bash
# Usage: [sudo] [prefix=/usr/local] ./install
set -e

case "$1" in
'-h' | '--help' )
  sed -ne '/^#/!q;s/.\{1,2\}//;1d;p' < "$0"
  exit 0
  ;;
esac

if [[ $BASH_SOURCE == */* ]]; then
  cd "${BASH_SOURCE%/*}"
fi

prefix="${PREFIX:-$prefix}"
prefix="${prefix:-/usr/local}"

for src in bin/hub share/man/*/*.1; do
  dest="${prefix}/${src}"
  mkdir -p "${dest%/*}"
  [[ $src == share/* ]] && mode="644" || mode=755
  install -m "$mode" "$src" "$dest"
done
