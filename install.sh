#!/usr/bin/env bash

PROJECT_NAME="hub"

: ${USE_SUDO:="true"}
: ${HUB_INSTALL_DIR:="/usr/local"}

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
  esac
}


# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  local supported="darwin-386\ndarwin-amd64\nfreebsd-386\nfreebas-amd64\nlinux-386\nlinux-amd64\nlinux-arm\nlinux-arm64\nlinux-ppc64le\nwindows-386\nwindows-amd64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuilt binary for ${OS}-${ARCH}."
    echo "To build from source, go to https://github.com/github/hub"
    exit 1
  fi

  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# checkDesiredVersion checks if the desired version is available.
latest_release() {
	TAG=$(curl --silent "https://api.github.com/repos/github/hub/releases/latest" | grep "tag_name" | sed -E 's/.*"([^"]+)".*/\1/')
	VERSION=$(echo $TAG | sed 's/v//')
}

# checkHubInstalledVersion checks which version of hub is installed and
# if it needs to be changed.
checkHubInstalledVersion() {
  if [[ -f "${HUB_INSTALL_DIR}/${PROJECT_NAME}" ]]; then
    local version="v"$("${HUB_INSTALL_DIR}/${PROJECT_NAME}" --version | grep 'hub version' | cut -d' ' -f3)
    if [[ "$version" == "$TAG" ]]; then
      echo "HUB ${version} is already ${DESIRED_VERSION:-latest}"
      return 0
    else
      echo "HUB ${TAG} is available. Changing from version ${version}."
      return 1
    fi
  else
    return 1
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  HUB_DIST="hub-$OS-$ARCH-$TAG.tgz"
	DOWNLOAD_URL=$(curl -s https://api.github.com/repos/github/hub/releases/tags/$TAG | grep -E "browser_download_url\": \".+$OS-$ARCH.+\.tgz\"" | sed -E 's|.+(https://[^"]+).+|\1|')

  HUB_TMP_ROOT="$(mktemp -dt hub-installer-XXXXXX)"
  HUB_TMP_FILE="$HUB_TMP_ROOT/$HUB_DIST"
  
  echo "Downloading $DOWNLOAD_URL"
  if type "curl" > /dev/null; then
    curl -SsL "$DOWNLOAD_URL" -o "$HUB_TMP_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$HUB_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  HUB_TMP="$HUB_TMP_ROOT/$PROJECT_NAME"

  mkdir -p "$HUB_TMP"
  tar xf "$HUB_TMP_FILE" -C "$HUB_TMP"
	HUB_DIST="hub-$OS-$ARCH-$VERSION"
  HUB_INSTALL_FILE="$HUB_TMP/$HUB_DIST/install"

  echo "Preparing to install $PROJECT_NAME into ${HUB_INSTALL_DIR}"
	prefix=$HUB_INSTALL_DIR $HUB_INSTALL_FILE
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    if [[ -n "$INPUT_ARGUMENTS" ]]; then
      echo "Failed to install $PROJECT_NAME with the arguments provided: $INPUT_ARGUMENTS"
      help
    else
      echo "Failed to install $PROJECT_NAME"
    fi
  fi
  cleanup
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  HUB="$(which $PROJECT_NAME)"
  if [ "$?" = "1" ]; then
    echo "$PROJECT_NAME not found. Is $HUB_INSTALL_DIR on your "'$PATH?'
    exit 1
  fi
  set -e
  echo "Run '$PROJECT_NAME init' to configure $PROJECT_NAME."
}

# help provides possible cli installation arguments
help () {
  echo "Accepted cli arguments are:"
  echo -e "\t[--help|-h ] ->> prints this help"
  echo -e "\t[--version|-v <desired_version>] . When not defined it defaults to latest"
  echo -e "\te.g. --version v2.4.0  or -v latest"
  echo -e "\t[--no-sudo]  ->> install without sudo"
}

cleanup() {
  if [[ -d "${HUB_TMP_ROOT:-}" ]]; then
    rm -rf "$HUB_TMP_ROOT"
  fi
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e

# Parsing input arguments (if any)
export INPUT_ARGUMENTS="${@}"
set -u
while [[ $# -gt 0 ]]; do
  case $1 in
    '--version'|-v)
       shift
       if [[ $# -ne 0 ]]; then
           export DESIRED_VERSION="${1}"
       else
           echo -e "Please provide the desired version. e.g. --version v2.4.0 or -v latest"
           exit 0
       fi
       ;;
    '--no-sudo')
       USE_SUDO="false"
       ;;
    '--help'|-h)
       help
       exit 0
       ;;
    *) exit 1
       ;;
  esac
  shift
done
set +u

initArch
initOS
verifySupported
latest_release
if ! checkHubInstalledVersion; then
  downloadFile
  installFile
fi
testVersion
cleanup