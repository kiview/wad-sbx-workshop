#!/usr/bin/env bash
# Download pinned app fixtures and the native adapter; never run app code on the host.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
. "$ROOT/scripts/common.sh"
. "$ROOT/scripts/versions.env"
asset="$(beans_mcp_asset)"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
release_url="https://github.com/shelajev/wad-sbx-workshop/releases/download/materials-v0.1.0"
for name in incident-triage-board.bundle "$asset" SHA256SUMS LICENSE THIRD-PARTY-NOTICES.md release-manifest.json; do
  curl --fail --location --silent --show-error --retry 3 "$release_url/$name" --output "$work/$name"
done
for name in incident-triage-board.bundle "$asset" LICENSE THIRD-PARTY-NOTICES.md release-manifest.json; do
  expected="$(awk -v n="$name" '$2 == n {print $1}' "$work/SHA256SUMS")"
  [ -n "$expected" ] || die "No published checksum for $name"
  actual="$(sha256_file "$work/$name")"
  [ "$actual" = "$expected" ] || die "Checksum mismatch: $name"
done
mkdir -p "$ROOT/.local" "$ROOT/dist/beans-mcp/$BEANS_MCP_VERSION"
for name in "$asset" SHA256SUMS LICENSE THIRD-PARTY-NOTICES.md release-manifest.json; do
  copy_file "$work/$name" "$ROOT/dist/beans-mcp/$BEANS_MCP_VERSION/$name"
done
if [ -e "$ROOT/.local/app" ]; then
  info 'Keeping your existing .local/app; no application files changed.'
else
  git clone --config core.autocrlf=false --config core.eol=lf "$work/incident-triage-board.bundle" "$ROOT/.local/app"
  git -C "$ROOT/.local/app" remote remove origin
fi
for ref in app-00-starter app-01-warmup-solution app-02-feature-solution; do
  git -C "$ROOT/.local/app" rev-parse --verify "$ref^{commit}" >/dev/null
done
"$ROOT/scripts/prepare-app.sh"
ok 'Materials ready. Start at chapters/00-setup/README.md.'
