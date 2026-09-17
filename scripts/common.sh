#!/usr/bin/env bash
# Shared setup helpers; no host supervisor.
CONTROL_DIR="${FACTORY_CONTROL_DIR:-$ROOT/.local/chapters}"
BIN_DIR="$CONTROL_DIR/bin"
info() { printf '%s\n' "$*" >&2; }
ok() { info "$@"; }
warn() { info "$@"; }
error() { info "$@"; }
die() { error "$*"; exit 1; }
beans_bin() { printf '%s' "$BIN_DIR/beans"; }
require_cmd() {
  [ -x "$1" ] || command -v "$1" >/dev/null 2>&1 || die "$2 unavailable. Run $3"
}
copy_file() {
  mkdir -p "$(dirname "$2")"
  cat "$1" > "$2" || return 1
  cmp -s "$1" "$2"
}
