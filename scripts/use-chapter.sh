#!/usr/bin/env bash
# Put a reference configuration in factory/. App checkpoints are an explicit extra choice.
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
chapter="${1:?Usage: scripts/use-chapter.sh 03-pi [app-checkpoint]}"
case "$chapter" in 02-launcher|02.5-acr|03-pi|04-team|05-mcp|06-human|07-factory) ;; *) echo 'Choose a chapter directory from 02-launcher through 07-factory'; exit 2;; esac
mkdir -p "$root/.local"
if [ -d "$root/factory" ]; then
  backup="$(mktemp -d "$root/.local/factory-backup.XXXXXX")"
  cp -R "$root/factory/." "$backup/"
  echo "Previous configuration saved in $backup"
fi
mkdir -p "$root/factory"
for file in sbxenv.yaml chapter.env PROMPT.md; do cp "$root/chapters/$chapter/$file" "$root/factory/$file"; done
# Adjust reference paths from chapters/NN/ to factory/.
sed -i.bak -e 's|../kits/|../chapters/kits/|g' -e 's|../../scripts/|../scripts/|g' -e 's|../../sample-app|../sample-app|g' "$root/factory/sbxenv.yaml"
rm "$root/factory/sbxenv.yaml.bak"
if [ -f "$root/chapters/$chapter/team.tsv" ]; then cp "$root/chapters/$chapter/team.tsv" "$root/factory/team.tsv"; else rm -f "$root/factory/team.tsv"; fi
[ $# -lt 2 ] || "$root/scripts/prepare-app.sh" "$2"
echo 'Ready: factory/. Your sample-app/ is unchanged unless you selected a checkpoint.'
