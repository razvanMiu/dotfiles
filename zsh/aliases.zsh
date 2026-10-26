# X11 startx
alias startx='startx "$XDG_CONFIG_HOME/X11/xinitrc"'

# DWM
alias cdwm='nvim "$XDG_CONFIG_HOME/dwm/config.h"'
alias mdwm='(cd "$XDG_CONFIG_HOME/dwm" && sudo make clean install)'

# Zed
alias zed-enable-biome='mkdir -p .zed && cp "$XDG_CONFIG_HOME/zed/templates/biome.json" .zed/settings.json'
