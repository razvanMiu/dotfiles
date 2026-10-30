#!/usr/bin/env bash
set -euo pipefail

# Copy curated Agent Skills from one or more Git repositories into this
# dotfiles repo. Pi and other agents load them via:
#
#   ~/.agents/skills -> ~/.config/agents/skills
#
# Sources are declared in ~/.config/agents/skill-sources.json.

XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"
XDG_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"

AGENTS_DIR="${AGENTS_DIR:-$XDG_CONFIG_HOME/agents}"
DEST_DIR="$AGENTS_DIR/skills"
SOURCES_FILE="${SKILL_SOURCES_FILE:-$AGENTS_DIR/skill-sources.json}"
SOURCE_ROOT="${SKILL_SOURCE_ROOT:-$XDG_DATA_HOME/agent-skill-sources}"
MANIFEST="$AGENTS_DIR/.managed-skills.manifest"

trap 'rm -f "$MANIFEST.tmp"' EXIT

validate_sources_json() {
  python - "$SOURCES_FILE" <<'PY'
import json
import re
import sys
from pathlib import PurePosixPath

path = sys.argv[1]
name_re = re.compile(r"^[a-zA-Z0-9._-]+$")
skill_name_re = re.compile(r"^[a-z0-9][a-z0-9-]{0,63}$")

with open(path, "r", encoding="utf-8") as f:
    data = json.load(f)

sources = data.get("sources")
if not isinstance(sources, list) or not sources:
    raise SystemExit("skill-sources.json must contain a non-empty 'sources' array")

seen_sources = set()
seen_destinations = set()

for source in sources:
    if not isinstance(source, dict):
        raise SystemExit("each source must be an object")

    name = source.get("name")
    repo = source.get("repo")
    ref = source.get("ref", "")
    skills = source.get("skills")

    if not isinstance(name, str) or not name_re.fullmatch(name):
        raise SystemExit(f"invalid source name: {name!r}")
    if name in seen_sources:
        raise SystemExit(f"duplicate source name: {name}")
    seen_sources.add(name)

    if not isinstance(repo, str) or not repo:
        raise SystemExit(f"source {name}: repo must be a non-empty string")
    if not isinstance(ref, str):
        raise SystemExit(f"source {name}: ref must be a string")
    if not isinstance(skills, list) or not skills:
        raise SystemExit(f"source {name}: skills must be a non-empty array")

    for skill in skills:
        if isinstance(skill, str):
            skill_path = skill
            dest_name = PurePosixPath(skill_path).name
        elif isinstance(skill, dict):
            skill_path = skill.get("path")
            dest_name = skill.get("name") or (PurePosixPath(skill_path).name if isinstance(skill_path, str) else None)
        else:
            raise SystemExit(f"source {name}: each skill must be a string or object")

        if not isinstance(skill_path, str) or not skill_path or skill_path.startswith("/") or ".." in PurePosixPath(skill_path).parts:
            raise SystemExit(f"source {name}: invalid skill path: {skill_path!r}")
        if not isinstance(dest_name, str) or not skill_name_re.fullmatch(dest_name):
            raise SystemExit(f"source {name}: invalid destination skill name: {dest_name!r}")
        if dest_name in seen_destinations:
            raise SystemExit(f"duplicate destination skill name: {dest_name}")
        seen_destinations.add(dest_name)
PY
}

emit_skill_rows() {
  python - "$SOURCES_FILE" <<'PY'
import json
import sys
from pathlib import PurePosixPath

with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)

for source in data["sources"]:
    source_name = source["name"]
    repo = source["repo"]
    ref = source.get("ref", "")
    for skill in source["skills"]:
        if isinstance(skill, str):
            skill_path = skill
            dest_name = PurePosixPath(skill_path).name
        else:
            skill_path = skill["path"]
            dest_name = skill.get("name") or PurePosixPath(skill_path).name
        print("\x1f".join([source_name, repo, ref, skill_path, dest_name]))
PY
}

copy_skill() {
  local source_name="$1"
  local source_dir="$2"
  local rel_path="$3"
  local dest_name="$4"
  local src="$source_dir/$rel_path"
  local dest="$DEST_DIR/$dest_name"

  if [[ ! -f "$src/SKILL.md" ]]; then
    echo "Skipping $source_name:$rel_path: missing SKILL.md" >&2
    return
  fi

  if [[ -e "$dest" ]]; then
    echo "Error: destination already exists before copy: $dest" >&2
    echo "If this is a local skill, choose a different destination name in skill-sources.json." >&2
    exit 1
  fi

  echo "Copying $source_name:$rel_path -> $dest_name"
  cp -a "$src" "$dest"
  printf '%s\n' "$dest_name" >> "$MANIFEST.tmp"
}

update_source_repo() {
  local name="$1"
  local repo_url="$2"
  local ref="$3"
  local source_dir="$SOURCE_ROOT/$name"

  mkdir -p "$SOURCE_ROOT"

  if [[ -d "$source_dir/.git" ]]; then
    echo "Updating $name in $source_dir" >&2
    git -C "$source_dir" fetch --prune origin >&2
  else
    echo "Cloning $repo_url into $source_dir" >&2
    git clone "$repo_url" "$source_dir" >&2
  fi

  if [[ -n "$ref" ]]; then
    git -C "$source_dir" checkout "$ref" >&2
    git -C "$source_dir" pull --ff-only origin "$ref" >&2 2>/dev/null || true
  else
    default_branch="$(git -C "$source_dir" symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null || true)"
    if [[ -n "$default_branch" ]]; then
      git -C "$source_dir" checkout "${default_branch#origin/}" >&2 2>/dev/null || true
    fi
    git -C "$source_dir" pull --ff-only >&2
  fi

  printf '%s' "$source_dir"
}

mkdir -p "$DEST_DIR"

if [[ ! -f "$SOURCES_FILE" ]]; then
  echo "Missing skill sources file: $SOURCES_FILE" >&2
  exit 1
fi

validate_sources_json

# Pre-fetch and update all source repositories before modifying any destination skills.
# This prevents wiping existing managed skills if network or git fetch fails.
rows="$(emit_skill_rows)"
declare -a updated_keys=()

while IFS=$'\x1f' read -r source_name repo_url ref _rel_path _dest_name; do
  source_key="$source_name"$'\x1f'"$repo_url"$'\x1f'"$ref"
  already_updated=0
  for k in "${updated_keys[@]:-}"; do
    if [[ "$k" == "$source_key" ]]; then
      already_updated=1
      break
    fi
  done
  if [[ "$already_updated" -eq 0 ]]; then
    update_source_repo "$source_name" "$repo_url" "$ref" >/dev/null
    updated_keys+=("$source_key")
  fi
done <<< "$rows"

# Remove only skills previously managed by this script. Leave custom local skills
# in ~/.config/agents/skills untouched.
if [[ -f "$MANIFEST" ]]; then
  while IFS= read -r skill_name; do
    [[ -n "$skill_name" ]] || continue
    if [[ "$skill_name" =~ ^[a-z0-9][a-z0-9-]{0,63}$ ]]; then
      rm -rf -- "$DEST_DIR/$skill_name"
    fi
  done < "$MANIFEST"
fi

: > "$MANIFEST.tmp"

while IFS=$'\x1f' read -r source_name _repo_url _ref rel_path dest_name; do
  source_dir="$SOURCE_ROOT/$source_name"
  copy_skill "$source_name" "$source_dir" "$rel_path" "$dest_name"
done <<< "$rows"

sort -u "$MANIFEST.tmp" > "$MANIFEST"

# Shared agent convention used by Pi, Claude Code, and other harnesses.
if [[ -L "$HOME/.agents" || ! -e "$HOME/.agents" ]]; then
  ln -sfn "$AGENTS_DIR" "$HOME/.agents"
  echo "Ensured ~/.agents -> $AGENTS_DIR"
else
  echo "Warning: ~/.agents exists and is not a symlink; leaving it unchanged." >&2
fi

echo "Done. Skills are available under $DEST_DIR"
