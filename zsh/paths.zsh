# Tool install roots
export BUN_INSTALL="$HOME/.bun"
export VOLTA_HOME="$XDG_DATA_HOME/volta"

# PATH
# typeset -U keeps path entries unique across nested shells.
typeset -U path PATH
path=("$XDG_BIN_HOME" "$VOLTA_HOME/bin" "$BUN_INSTALL/bin" $path)
export PATH
