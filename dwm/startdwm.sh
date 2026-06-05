#!/bin/sh

state_dir="${XDG_STATE_HOME:-$HOME/.local/state}/dwm"
mkdir -p "$state_dir"

xrandr --output DP-2 --mode 1920x1080 --rate 144
setxkbmap us -option ctrl:nocaps,ctrl:swap_lalt_lctl

if command -v xwallpaper >/dev/null 2>&1; then
	xwallpaper --zoom "${XDG_CONFIG_HOME:-$HOME/.config}/wallpapers/current.png"
fi

if command -v picom >/dev/null 2>&1 && ! pgrep -x picom >/dev/null 2>&1; then
	picom --config "${XDG_CONFIG_HOME:-$HOME/.config}/picom/picom.conf" \
		> "$state_dir/picom.log" 2>&1 &
fi

status_cmd="${HOME}/.local/bin/dwm-status"
if [ ! -x "$status_cmd" ]; then
	status_cmd="$(command -v dwm-status 2>/dev/null || true)"
fi
if [ -n "$status_cmd" ] && ! { [ -r "$state_dir/status.pid" ] && kill -0 "$(cat "$state_dir/status.pid")" 2>/dev/null; }; then
	"$status_cmd" > "$state_dir/status.log" 2>&1 &
	printf '%s\n' "$!" > "$state_dir/status.pid"
fi

# relaunch DWM if the binary changes, otherwise bail
csum=""
new_csum=$(sha1sum "$(command -v dwm)")
while true
do
    if [ "$csum" != "$new_csum" ]
    then
        csum=$new_csum
	dwm 2> "$state_dir/stderror.log"
    else
        exit 0
    fi
    new_csum=$(sha1sum "$(command -v dwm)")
    sleep 0.5
done
