# Tool install roots
export BUN_INSTALL="$HOME/.bun"
export FNM_PATH="$XDG_DATA_HOME/fnm"

# PATH
# typeset -U keeps path entries unique across nested shells.
typeset -U path PATH
path=("$XDG_BIN_HOME" "$FNM_PATH" "$BUN_INSTALL/bin" $path)
export PATH
