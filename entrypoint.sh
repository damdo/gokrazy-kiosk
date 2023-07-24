#!/bin/bash

if [[ $# -eq 0 ]] ; then
    /usr/bin/X :0 -nolisten tcp vt1
else
    # maximize the main window
    echo 'nohup bash -c "sleep 10 ; xdotool getwindowfocus windowsize $(xdotool getdisplaygeometry)" &' > /usr/bin/entrypoint-command.sh

    echo "$@" >> /usr/bin/entrypoint-command.sh
    chmod +x /usr/bin/entrypoint-command.sh
    /usr/bin/xinit /usr/bin/entrypoint-command.sh -- :0 -nolisten tcp vt1
fi
