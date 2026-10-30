# X11 startx
alias startx='startx "$XDG_CONFIG_HOME/X11/xinitrc"'

# DWM
alias cdwm='nvim "$XDG_CONFIG_HOME/dwm/config.h"'
alias mdwm='(cd "$XDG_CONFIG_HOME/dwm" && sudo make clean install)'

# Zed
alias zed-enable-biome='mkdir -p .zed && cp "$XDG_CONFIG_HOME/zed/templates/biome.json" .zed/settings.json'

# Tmux helper
# t          -> attach to most recent session or create one
# t <name>   -> attach to session <name> or create it
# tls        -> list all running sessions
# tk <name>  -> kill session <name>
t() {
  if [ $# -eq 0 ]; then
    tmux attach-session 2>/dev/null || tmux new-session
  else
    tmux attach-session -t "$1" 2>/dev/null || tmux new-session -s "$1"
  fi
}
alias tls='tmux list-sessions'
alias tk='tmux kill-session -t'
