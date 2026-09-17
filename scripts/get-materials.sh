#!/usr/bin/env bash
# Download pinned app fixtures and the native adapter; never run app code on the host.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
. "$ROOT/scripts/common.sh"
. "$ROOT/scripts/versions.env"
case "$(uname -s)/$(uname -m)" in
  Darwin/arm64) platform=darwin_arm64 ;;
  Darwin/x86_64) platform=darwin_amd64 ;;
  Linux/aarch64|Linux/arm64) platform=linux_arm64 ;;
  Linux/x86_64) platform=linux_amd64 ;;
  *) die 'This workshop downloader supports macOS and Linux; this host is not covered.' ;;
esac
asset="beans-mcp_${BEANS_MCP_VERSION}_${platform}"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
gh release download materials-v0.1.0 --repo shelajev/wad-sbx-workshop \
  --pattern incident-triage-board.bundle --pattern "$asset" \
  --pattern SHA256SUMS --pattern LICENSE --pattern THIRD-PARTY-NOTICES.md \
  --pattern release-manifest.json --dir "$work"
for name in incident-triage-board.bundle "$asset" LICENSE THIRD-PARTY-NOTICES.md release-manifest.json; do
  expected="$(awk -v n="$name" '$2 == n {print $1}' "$work/SHA256SUMS")"
  [ -n "$expected" ] || die "No published checksum for $name"
  actual="$(shasum -a 256 "$work/$name" | awk '{print $1}')"
  [ "$actual" = "$expected" ] || die "Checksum mismatch: $name"
done
mkdir -p "$ROOT/.local" "$ROOT/dist/beans-mcp/$BEANS_MCP_VERSION"
for name in "$asset" SHA256SUMS LICENSE THIRD-PARTY-NOTICES.md release-manifest.json; do
  copy_file "$work/$name" "$ROOT/dist/beans-mcp/$BEANS_MCP_VERSION/$name"
done
if [ -e "$ROOT/.local/app" ]; then
  info 'Keeping your existing .local/app; no application files changed.'
else
  git clone "$work/incident-triage-board.bundle" "$ROOT/.local/app"
  git -C "$ROOT/.local/app" remote remove origin
fi
for ref in app-00-starter app-01-warmup-solution app-02-feature-solution; do
  git -C "$ROOT/.local/app" rev-parse --verify "$ref^{commit}" >/dev/null
done
"$ROOT/scripts/prepare-app.sh"
ok 'Materials ready. Start at chapters/00-setup/README.md.'
