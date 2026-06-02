# XDG Base Directories
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_CACHE_HOME="$HOME/.cache"
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_STATE_HOME="$HOME/.local/state"
export XDG_BIN_HOME="$HOME/.local/bin"

# XDG Runtime
export XDG_RUNTIME_DIR="${XDG_RUNTIME_DIR:-$HOME/.local/run}"

# X11 Env
export XAUTHORITY="$XDG_RUNTIME_DIR/Xauthority"
export XINITRC="$XDG_CONFIG_HOME/X11/xinitrc"

# Oh My Zsh
export ZSH="$XDG_CONFIG_HOME/oh-my-zsh"
export ZSH_STATE="$XDG_STATE_HOME/oh-my-zsh"
export ZSH_COMPDUMP="$XDG_CACHE_HOME/oh-my-zsh/.zcompdump-$HOST"

# History
export HISTFILE="$ZSH_STATE/history"
export HISTSIZE=100000
export SAVEHIST=100000

# Editor
export EDITOR=nvim

# npm
export npm_config_cache="$XDG_CACHE_HOME/npm"
export npm_config_logs_dir="$XDG_STATE_HOME/npm/logs"

# bun
export BUN_INSTALL="$HOME/.bun"

# Pi.dev
export PI_CODING_AGENT_DIR="$XDG_CONFIG_HOME/pi/agent"
export PI_CODING_AGENT_SESSION_DIR="$XDG_STATE_HOME/pi/sessions"

# Paths
export FNM_PATH="$XDG_DATA_HOME/fnm"

export PATH="$XDG_BIN_HOME:$PATH"
export PATH="$FNM_PATH:$PATH"
export PATH="$BUN_INSTALL/bin:$PATH"

# Theme
ZSH_THEME="eastwood"

# Autosuggest settings
ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE="fg=60"
ZSH_AUTOSUGGEST_BUFFER_MAX_SIZE="20"
ZSH_AUTOSUGGEST_USE_ASYNC=1

