# XDG base directories
export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"
export XDG_CACHE_HOME="${XDG_CACHE_HOME:-$HOME/.cache}"
export XDG_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"
export XDG_STATE_HOME="${XDG_STATE_HOME:-$HOME/.local/state}"
export XDG_BIN_HOME="${XDG_BIN_HOME:-$HOME/.local/bin}"
export XDG_RUNTIME_DIR="${XDG_RUNTIME_DIR:-$HOME/.local/run}"
export XDG_DOCUMENTS_DIR="${XDG_DOCUMENTS_DIR:-$HOME/Documents}"

# X11
# Preserve display-manager supplied XAUTHORITY (ly uses $XDG_RUNTIME_DIR/lyxauth).
export XAUTHORITY="${XAUTHORITY:-$XDG_RUNTIME_DIR/Xauthority}"
export XINITRC="$XDG_CONFIG_HOME/X11/xinitrc"

# History
export ZSH_STATE="$XDG_STATE_HOME/oh-my-zsh"
export HISTFILE="$ZSH_STATE/history"
export HISTSIZE=100000
export SAVEHIST=100000

# Editor
export EDITOR=nvim

# npm
export npm_config_cache="$XDG_CACHE_HOME/npm"
export npm_config_logs_dir="$XDG_STATE_HOME/npm/logs"

# Claude Code
export CLAUDE_CONFIG_DIR="$XDG_CONFIG_HOME/claude"

# Pi coding agent
export PI_CODING_AGENT_DIR="$XDG_CONFIG_HOME/pi/agent"
export PI_CODING_AGENT_SESSION_DIR="$XDG_STATE_HOME/pi/sessions"
