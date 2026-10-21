#!/usr/bin/env bash
# Assemble the public repository from this one: the code, the patches, the
# scripts, the example configurations and the tests, plus the documents and
# license files of public/. Nothing else goes: no other document, no git
# history, no build output.
#
# Usage: public/export-public.sh [--dirty] <destination>
#   destination  a directory that does not exist yet or is empty
#   --dirty      export even if the exported paths have uncommitted changes
#                (trial run); without it the export matches a commit
#
# Only files tracked by git are taken from the code paths; public/ is copied
# as it is, without this script. The destination gets no .git: run git init
# there. At the end, the export is checked for references that only make
# sense in this repository.
set -euo pipefail

usage() { echo "usage: $0 [--dirty] <destination>" >&2; exit 2; }
dirty=no
dest=""
for a in "$@"; do
    case "$a" in
    --dirty) dirty=yes ;;
    -*) usage ;;
    *) [ -z "$dest" ] || usage; dest="$a" ;;
    esac
done
[ -n "$dest" ] || usage

repo="$(cd "$(dirname "$0")/.." && pwd)"
paths=(
    shared
    reference/kdf
    free5gc/aanf free5gc/config free5gc/patches free5gc/scripts
    af/go.mod af/go.sum af/af af/akmatls af/config af/scripts af/uastar af/ue af/wolftls
)

if [ -e "$dest" ] && [ -n "$(ls -A "$dest")" ]; then
    echo "$dest is not empty" >&2
    exit 1
fi
if [ "$dirty" = no ] && [ -n "$(git -C "$repo" status --porcelain -- "${paths[@]}" public)" ]; then
    echo "uncommitted changes in the exported paths; commit them or use --dirty:" >&2
    git -C "$repo" status --short -- "${paths[@]}" public >&2
    exit 1
fi

mkdir -p "$dest"
dest="$(cd "$dest" && pwd)"
case "$dest/" in
"$repo"/*) echo "$dest is inside the repository" >&2; exit 1 ;;
esac

git -C "$repo" ls-files -z -- "${paths[@]}" | (cd "$repo" && xargs -0 cp --parents -t "$dest")
tar -C "$repo/public" --exclude=./export-public.sh -cf - . | tar -C "$dest" -xf -

# References to documents and records that stay here.
leaks="$(grep -rnIE 'MVP|DV-[0-9]|ADR-[0-9]|docs/(plano|conformidade|decisions|traceability|baseline|literatura|af-side|free5gc-side|shared|tutorial)|00-baseline|/home/[a-z]' "$dest" || true)"
if [ -n "$leaks" ]; then
    echo "$leaks" >&2
    echo "the export has references that only make sense in this repository" >&2
    exit 1
fi

echo "exported $(find "$dest" -type f | wc -l) files to $dest"
echo "next: cd $dest && git init"
