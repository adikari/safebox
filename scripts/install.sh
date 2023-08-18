#!/bin/sh
#
# Licensed under the MIT license
# <LICENSE-MIT or https://opensource.org/licenses/MIT>, at your
# option. This file may not be copied, modified, or distributed
# except according to those terms.

# curl -sSL https://raw.githubusercontent.com/monebag/safebox/main/scripts/install.sh | sh

set -u

get_latest_release() {
  curl --silent "https://api.github.com/repos/$1/releases/latest" |
    grep '"tag_name":' |
    sed -E 's/.*"([^"]+)".*/\1/'
}

BINARY_DOWNLOAD_PREFIX="https://github.com/monebag/safebox/releases/download"
PACKAGE_VERSION=$(get_latest_release monebag/safebox)

download_binary_and_run_installer() {
    downloader --check
    need_cmd mktemp
    need_cmd chmod
    need_cmd mkdir
    need_cmd rm
    need_cmd rmdir
    need_cmd which
    need_cmd dirname
    need_cmd awk
    need_cmd cut
    need_cmd tar

    get_architecture || return 1
    local _arch="$RETVAL"
    assert_nz "$_arch" "arch"

    local _ext=""
    case "$_arch" in
        *windows*)
            _ext=".exe"
            ;;
    esac

		local _current_dir=$(pwd)
    
    local _tar="safebox_""${PACKAGE_VERSION:1}""_${_arch}.tar.gz"
    local _url="$BINARY_DOWNLOAD_PREFIX/$PACKAGE_VERSION/${_tar}"
    local _dir="$(mktemp -d 2>/dev/null || ensure mktemp -d -t test)"
		local _bin_dir="/usr/local/bin"

		cd $_dir

    if [ ! -f "$_dir/$_tar" ];then 
      say "downloading safebox from $_url" 1>&2
      ensure mkdir -p "$_dir"
      downloader "$_url" "$_dir/$_tar"
      if [ $? != 0 ]; then
        say "failed to download $_url"
        exit 1
      fi
    fi

    tar -xf $_tar

		sudo mv "$_dir/dist/safebox$_ext" "$_bin_dir"
    local _retval=$?

    chmod +x "${_bin_dir}/safebox"

    ignore rm -rf "$_dir"

		cd $_current_dir

    echo "safebox was installed successfully to $_bin_dir/safebox"
    return "$_retval"
}

get_architecture() {
    local _ostype="$(uname -s)"
    local _cputype="$(uname -m)"

    [[ $_cputype = "aarch64" ]] && _cputype="arm64"

    RETVAL="$_ostype"_"$_cputype"
}

say() {
    local green=`tput setaf 2 2>/dev/null || echo ''`
    local reset=`tput sgr0 2>/dev/null || echo ''`
    echo "$1"
}

err() {
    local red=`tput setaf 1 2>/dev/null || echo ''`
    local reset=`tput sgr0 2>/dev/null || echo ''`
    say "${red}ERROR${reset}: $1" >&2
    exit 1
}

need_cmd() {
    if ! check_cmd "$1"
    then err "need '$1' (command not found)"
    fi
}

check_cmd() {
    command -v "$1" > /dev/null 2>&1
    return $?
}

need_ok() {
    if [ $? != 0 ]; then err "$1"; fi
}

assert_nz() {
    if [ -z "$1" ]; then err "assert_nz $2"; fi
}

# Run a command that should never fail. If the command fails execution
# will immediately terminate with an error showing the failing
# command.
ensure() {
    "$@"
    need_ok "command failed: $*"
}

# This is just for indicating that commands' results are being
# intentionally ignored. Usually, because it's being executed
# as part of error handling.
ignore() {
    "$@"
}

# This wraps curl or wget. Try curl first, if not installed,
# use wget instead.
downloader() {
    if check_cmd curl
    then _dld=curl
    elif check_cmd wget
    then _dld=wget
    else _dld='curl or wget' # to be used in error message of need_cmd
    fi

    if [ "$1" = --check ]
    then need_cmd "$_dld"
    elif [ "$_dld" = curl ]
    then curl -sSfL "$1" -o "$2"
    elif [ "$_dld" = wget ]
    then wget "$1" -O "$2"
    else err "Unknown downloader"   # should not reach here
    fi
}

download_binary_and_run_installer "$@" || exit 1
