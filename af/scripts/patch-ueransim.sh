#!/usr/bin/env bash
# Apply or revert the project patches on UERANSIM, and rebuild nr-ue. Runs on
# the host where UERANSIM is built (the free5GC host, decision D5).
#
# Usage: patch-ueransim.sh apply|revert|status [patch-dir]
#   patch-dir     default: af/ue/ueransim/patches in the repository
#   UERANSIM_DIR  UERANSIM directory (default: ~/UERANSIM)
#
# apply   checks every patch, applies them in order with git apply, keeps a
#         copy of the original binary as build/nr-ue.orig (once) and rebuilds
#         nr-ue only, incrementally, in cmake-build-release.
# revert  reverses the patches in the opposite order and restores
#         build/nr-ue from build/nr-ue.orig.
# status  tells whether the patches are applied.
set -euo pipefail

action="${1:?usage: patch-ueransim.sh apply|revert|status [patch-dir]}"
here="$(cd "$(dirname "$0")/.." && pwd)"
patch_dir="${2:-$here/ue/ueransim/patches}"
dir="${UERANSIM_DIR:-$HOME/UERANSIM}"
bin="$dir/build/nr-ue"

git -C "$dir" rev-parse --show-toplevel >/dev/null 2>&1 || { echo "$dir is not a git checkout" >&2; exit 1; }
gitapply() { git -C "$dir" apply "$@"; }
mapfile -t patches < <(find "$patch_dir" -maxdepth 1 -name '*.patch' | sort)
[ "${#patches[@]}" -gt 0 ] || { echo "no patch in $patch_dir" >&2; exit 1; }
# Only the files touched by the patches are checked for local changes, so
# that an edited gNB configuration does not block the UE patches.
mapfile -t files < <(for p in "${patches[@]}"; do gitapply --numstat "$p"; done | cut -f3 | sort -u)
gitstatus() { git -C "$dir" status --porcelain "$@" -- "${files[@]}"; }

applied() {
    local p
    for p in "${patches[@]}"; do
        gitapply --reverse --check "$p" 2>/dev/null || return 1
    done
}

build() {
    if [ -d "$dir/cmake-build-release" ]; then
        cmake --build "$dir/cmake-build-release" --target nr-ue
    else
        make -C "$dir"
    fi
    # Rename into place: works while an nr-ue is running.
    cp -p "$dir/cmake-build-release/nr-ue" "$bin.new"
    mv "$bin.new" "$bin"
}

case "$action" in
status)
    if applied; then echo "UERANSIM: patches applied"; else echo "UERANSIM: patches not applied"; fi
    gitstatus
    ;;
apply)
    if applied; then
        echo "UERANSIM: patches already applied"
    else
        if [ -n "$(gitstatus --untracked-files=no)" ]; then
            echo "$dir has local changes; not applying" >&2
            gitstatus --untracked-files=no >&2
            exit 1
        fi
        for p in "${patches[@]}"; do gitapply --check "$p"; done
        for p in "${patches[@]}"; do gitapply "$p"; echo "applied $(basename "$p")"; done
    fi
    if [ -f "$bin" ] && [ ! -f "$bin.orig" ]; then
        cp -p "$bin" "$bin.orig"
        echo "original binary kept as build/nr-ue.orig"
    fi
    build
    echo "nr-ue rebuilt: $bin"
    ;;
revert)
    if applied; then
        for ((i=${#patches[@]}-1; i>=0; i--)); do
            gitapply --reverse "${patches[$i]}"
            echo "reverted $(basename "${patches[$i]}")"
        done
    else
        echo "UERANSIM: patches not applied; nothing to reverse in the source"
    fi
    if [ -f "$bin.orig" ]; then
        cp -p "$bin.orig" "$bin.new"
        mv "$bin.new" "$bin"
        echo "build/nr-ue restored from build/nr-ue.orig"
    fi
    ;;
*)
    echo "unknown action $action" >&2
    exit 2
    ;;
esac
