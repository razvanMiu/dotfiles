#!/usr/bin/env bash
# Shared tmux plumbing for the fzf pickers in tmux/scripts.
#
# Sourced, never executed:
#   . "$SCRIPT_DIR/lib/pi-tmux.sh"
#
# Contract: sourcing only defines functions. It assigns no caller variables,
# runs no tmux command, and leaves no state a caller has to reset. Every
# function writes its result to stdout and is safe under `set -euo pipefail`,
# so callers capture output with command substitution:
#   cwd=$(pi_current_dir)
#
# Preview/list geometry and the ANSI palette for the inline python3 blocks live
# in lib/pi_preview.py; the bash-side geometry lives here.

# --- tmux context ------------------------------------------------------------

# Print the pane id the picker was launched from (empty outside tmux).
pi_target_pane() {
  tmux display-message -p "#{pane_id}" 2>/dev/null || true
}

# Print that pane's working directory, falling back to the caller's cwd.
pi_current_dir() {
  tmux display-message -p "#{pane_current_path}" 2>/dev/null || pwd
}

# Print the attached client width in cells, defaulting to 160 when detached.
pi_client_width() {
  tmux display-message -p "#{client_width}" 2>/dev/null || echo 160
}

# --- fzf popup geometry (bash side) -----------------------------------------

# Print the dim rule used inside an fzf --header, sized to the list column.
pi_list_divider() {
  local tmux_w="${1:-160}" popup_w inner_w list_w div_w
  popup_w=$((tmux_w * 90 / 100))
  inner_w=$((popup_w - 2))
  list_w=$((inner_w - (inner_w * 50 / 100)))
  div_w=$((list_w - 4))
  [ "$div_w" -lt 30 ] && div_w=30
  printf '─%.0s' $(seq 1 "$div_w")
}

# --- pane actions ------------------------------------------------------------

# Terminate the foreground command in a pane (empty id = current pane) and
# leave the shell prompt clean. Resolved next to this library, so neither the
# caller's cwd nor its own location matters.
pi_terminate_pane_process() {
  local pane_id="${1:-}" scripts_dir
  scripts_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
  "$scripts_dir/terminate-pane-process" "$pane_id"
}

# Replace whatever runs in the current pane with `command` (Enter appended).
# Falls back to the client's active pane when no pane id can be resolved.
pi_open_in_pane() {
  local command="$1" pane_id
  pane_id="$(pi_target_pane)"
  pi_terminate_pane_process "$pane_id"
  if [ -n "$pane_id" ]; then
    tmux send-keys -t "$pane_id" "$command" C-m
  else
    tmux send-keys "$command" C-m
  fi
}

# Open a new window in the current session, rooted at `dir`, named `name`.
# Optional trailing arguments are the command to run in it.
pi_open_in_new_window() {
  local dir="$1" name="$2"
  shift 2
  tmux new-window -n "$name" -c "$dir" "$@"
}

# --- session actions ---------------------------------------------------------

# Test whether a session with that exact name exists. Meant for `if`/`||`.
pi_session_exists() {
  tmux has-session -t "$1" 2>/dev/null
}

# Create a detached session rooted at `dir`; optional command to run in it.
pi_session_create_detached() {
  local name="$1" dir="$2"
  shift 2
  tmux new-session -d -s "$name" -c "$dir" "$@"
}

# Bring an existing session into view: switch inside tmux, attach outside.
pi_session_focus() {
  if [ -n "${TMUX:-}" ]; then
    tmux switch-client -t "$1"
  else
    tmux attach-session -t "$1"
  fi
}
