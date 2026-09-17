# Source this file from Bash or Zsh: source ./scripts/workshop-env.sh
# --completed selects the supplied end-of-chapter-1 app without replacing your work.
if [ -n "${BASH_VERSION:-}" ]; then
  _workshop_source="${BASH_SOURCE[0]}"
elif [ -n "${ZSH_VERSION:-}" ]; then
  _workshop_source="${(%):-%x}"
else
  printf 'Use Bash or Zsh to source this file.\n' >&2
  return 1
fi
_workshop_setup() {
  local root choice ref target
  root="$(cd "$(dirname "$_workshop_source")/.." && pwd)" || return 1
  [ $# -eq 0 ] || { [ $# -eq 1 ] && [ "$1" = --completed ]; } || {
    printf 'Usage: source scripts/workshop-env.sh [--completed]\n' >&2; return 1;
  }
  git -C "$root/.local/app" rev-parse --git-dir >/dev/null 2>&1 || {
    printf 'First run: ./scripts/get-materials.sh\n' >&2; return 1;
  }
  choice=warmup-app
  if [ -f "$root/.local/warmup-selection" ]; then
    choice="$(cat "$root/.local/warmup-selection")"
  fi
  [ "${1:-}" != --completed ] || choice=warmup-completed
  case "$choice" in
    warmup-app) ref=app-00-starter ;;
    warmup-completed) ref=app-01-warmup-solution ;;
    *) printf 'Unexpected workshop selection: %s\n' "$choice" >&2; return 1 ;;
  esac
  target="$root/.local/$choice"
  if [ ! -e "$target" ]; then
    git clone --quiet --no-hardlinks "$root/.local/app" "$target" || return 1
    git -C "$target" switch --quiet -c workshop "$ref" || return 1
    git -C "$target" remote remove origin || return 1
  fi
  git -C "$target" rev-parse --git-dir >/dev/null 2>&1 || return 1
  if [ ! -f "$target/WORKSHOP-TASK.md" ]; then
    cat "$root/backlog/seed/wad-101--warm-up-active-filter-count.md" > "$target/WORKSHOP-TASK.md" || return 1
  fi
  if ! grep -qxF /WORKSHOP-TASK.md "$target/.git/info/exclude"; then
    printf '\n/WORKSHOP-TASK.md\n' >> "$target/.git/info/exclude" || return 1
  fi
  mkdir -p "$root/.local/chapters" || return 1
  printf '%s\n' "$choice" > "$root/.local/warmup-selection" || return 1
  printf 'This file stays outside the mounted application.\n' > "$root/.local/host-only.txt" || return 1
  export WORKSHOP="$root" APP_REPO="$root/.local/app"
  export CONTROL="$root/.local/chapters" FACTORY_CONTROL_DIR="$root/.local/chapters"
  export WARMUP="$target"
  printf 'Workshop ready. App working directory: %s\n' "$WARMUP"
}
if _workshop_setup "$@"; then
  unset -f _workshop_setup
  unset _workshop_source
else
  unset -f _workshop_setup
  unset _workshop_source
  return 1
fi
