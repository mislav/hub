#!/bin/sh

_arch() {
    case "$(uname -m)" in
                 x86_64) _arch__type="amd64" ;;
    i386/i486/i586/i686) _arch__type="386"   ;;
                   arm*) _arch__type="arm"   ;;
    esac

    printf "%s\\n" "${_arch__type}"
}

_platform() {
    case "$(uname)" in
        Linux*)   _platform__type="linux"   ;;
        Darwin*)  _platform__type="darwin"  ;;
        FreeBSD*) _platform__type="freebsd" ;;
        CYGWIN*|MINGW*|MSYS*) _platform__type="windows" ;;
    esac

    printf "%s\\n" "${_platform__type}"
}

_curl() {
    [ -z "${1}" ] && return 0
    if command -v "curl" >/dev/null 2>&1; then
        curl -L -s "${1}"
    elif command -v "wget" >/dev/null 2>&1; then
        wget --q -O- "${1}"
    fi
}

_uniq() {
    awk '!seen[$0]++'
}

_err() {
    printf "%s\\n" "${1}" >&2
}

_die() {
    _err "${1}"
    exit 2
}

_check_deps() {
    if ! command -v "curl" >/dev/null 2>&1 || ! command -v "curl" >/dev/null 2>&1; then
        _die "install 'curl' or 'wget' to continue, exiting ..."
    fi
}

_extract() {
    for _extract__file; do
        if [ -f "${_extract__file}" ] ; then
            case "${_extract__file}" in
                *.tar.gz|*.tgz) zcat  < "${_extract__file}" | tar xf -   ;;
                *.zip|*.xpi|*.war|*.jar|*.ear) unzip "${_extract__file}" ;;
            esac
        else
            _die "${progname}: '${_extract__file}' does not exist or is not readable" >&2
        fi
    done
}

_filter_prefix() {
    sed -e 's:.tgz::g' -e 's:.zip::g'
}

_check_deps

progname="$(basename "${0}")"
arch="$(_arch)"
platform="$(_platform)"
all_releases="$(_curl "https://api.github.com/repos/github/hub/releases" | \
                      awk '/browser_download_url/ {print $2}')"
[ -z "${all_releases}" ] && _die "https://api.github.com/repos/github/hub/releases api timeout"

stable_releases="$(printf "%s\\n" "${all_releases}" | awk '!/-pre[0-9]/ && !/-rc[0-9]/ && !/-preview[0-9]/')"
pre_releases="$(printf "%s\\n" "${all_releases}"    | awk ' /-pre[0-9]/ ||  /-rc[0-9]/ ||  /-preview[0-9]/')"

stable_releases_versions="$(printf "%s\\n" "${stable_releases}" | awk -F/ '{print $8}' | _uniq)"
pre_releases_versions="$(printf "%s\\n" "${pre_releases}"       | awk -F/ '{print $8}' | _uniq)"

stable_releases_latest_version="$(printf "%s\\n" "${stable_releases_versions}" | awk 'NR==1')"
pre_releases_latest_version="$(printf "%s\\n" "${pre_releases_versions}"       | awk 'NR==1')"

binary_url="$(printf "%s\\n" "${stable_releases}" | \
    grep "${stable_releases_latest_version}"      | \
    grep "${arch}"                                | \
    grep "${platform}"                            | head -1 | sed 's:"::g')"
binary_pkg="$(pwd)/$(basename "${binary_url}")"

[ -s "$(basename "${binary_pkg}")" ] || _curl "${binary_url}" > "${binary_pkg}"
[ -s "$(basename "${binary_pkg}")" ] && printf "%s\\n" "$(basename "${binary_pkg}")"

if _extract "${binary_pkg}"; then
    printf "%s\\n" "$(basename "${binary_pkg}")" | _filter_prefix | sed 's:$:/:'
fi
