#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
. "$ROOT/scripts/common.sh"
BEANS_MCP_VERSION=0.1.0

assert_equal() {
  [ "$1" = "$2" ] || die "$3: expected '$2', got '$1'"
}
uname() {
  case "$1" in
    -s) printf '%s' "$test_os" ;;
    -m) printf '%s' "$test_arch" ;;
    *) die "unexpected uname argument: $1" ;;
  esac
}

while read -r test_os test_arch asset binary; do
  assert_equal "$(beans_mcp_asset)" "beans-mcp_0.1.0_$asset" "$test_os/$test_arch adapter"
  assert_equal "$(beans_bin)" "$BIN_DIR/$binary" "$test_os/$test_arch Beans"
done <<'PLATFORMS'
Darwin arm64 darwin_arm64 beans
Darwin x86_64 darwin_amd64 beans
Linux aarch64 linux_arm64 beans
Linux arm64 linux_arm64 beans
Linux x86_64 linux_amd64 beans
MINGW64_NT-10.0-26200 x86_64 windows_amd64.exe beans.exe
MSYS_NT-10.0 amd64 windows_amd64.exe beans.exe
CYGWIN_NT-10.0 x86_64 windows_amd64.exe beans.exe
PLATFORMS

while read -r test_os test_arch; do
  if (beans_mcp_asset >/dev/null 2>&1); then
    die "unexpected adapter support: $test_os/$test_arch"
  fi
done <<'UNSUPPORTED'
MINGW64_NT-10.0-26200 arm64
FreeBSD x86_64
Linux riscv64
UNSUPPORTED

empty_sha=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
assert_equal "$(sha256_file /dev/null)" "$empty_sha" 'SHA-256'
if command -v sha256sum >/dev/null 2>&1; then
  hash_command=("$(command -v sha256sum)")
else
  hash_command=("$(command -v shasum)" -a 256)
fi
command() {
  case "$*" in
    '-v sha256sum') return 1 ;;
    '-v shasum') return 0 ;;
    *) builtin command "$@" ;;
  esac
}
shasum() {
  [ "$#" -eq 3 ] && [ "$1" = -a ] && [ "$2" = 256 ] || die 'expected shasum SHA-256 options'
  "${hash_command[@]}" "$3"
}
assert_equal "$(sha256_file /dev/null)" "$empty_sha" 'shasum fallback'
printf 'Platform and checksum checks passed.\n'
