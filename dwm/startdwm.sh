#!/bin/sh

xrandr --output DP-2 --mode 1920x1080 --rate 144
setxkbmap us -option ctrl:nocaps,ctrl:swap_lalt_lctl

# relaunch DWM if the binary changes, otherwise bail
csum=""
new_csum=$(sha1sum $(which dwm))
while true
do
    if [ "$csum" != "$new_csum" ]
    then
        csum=$new_csum
	dwm 2> ~/.local/state/dwm/stderror.log
    else
        exit 0
    fi
    new_csum=$(sha1sum $(which dwm))
    sleep 0.5
done
