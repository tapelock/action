#!/usr/bin/env bash
# Runs `tapelock check` against a cassette file, or every *.jsonl cassette
# in a directory, and fails if any of them fails. This is the whole job of
# the action's v1: install tapelock, run tapelock check, let its exit code
# become the workflow's pass/fail, see CONTRIBUTING.md.
set -euo pipefail

cassette_path="${1:?usage: check.sh <cassette-path> [extra tapelock check args...]}"
shift

if [ -d "$cassette_path" ]; then
  shopt -s nullglob
  files=("$cassette_path"/*.jsonl)
  shopt -u nullglob
  if [ ${#files[@]} -eq 0 ]; then
    echo "tapelock-action: no .jsonl cassette files found in $cassette_path" >&2
    exit 1
  fi
else
  files=("$cassette_path")
fi

status=0
for f in "${files[@]}"; do
  echo "::group::tapelock check $f"
  if ! tapelock check --cassette "$f" "$@"; then
    status=1
  fi
  echo "::endgroup::"
done

exit "$status"
