# fnm shell integration
if [[ -d "$FNM_PATH" ]] && command -v fnm >/dev/null 2>&1; then
  eval "$(fnm env --shell zsh)"
fi

# bun completions
[[ -s "$BUN_INSTALL/_bun" ]] && source "$BUN_INSTALL/_bun"
