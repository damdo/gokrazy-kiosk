#!/bin/bash

if [[ $# -eq 0 ]] ; then
    /usr/bin/X :0 -nolisten tcp vt1
else
    # keep maximizing the main window
    echo 'nohup bash -c "sleep 10; while true; do xdotool getwindowfocus windowsize $(xdotool getdisplaygeometry) || true; sleep 20; done" &' > /usr/bin/entrypoint-command.sh
    # start chromium refresher python script
    chmod +x /usr/bin/chromium-refresher.py
    echo 'nohup bash -c "sleep 60; while true; do /usr/bin/chromium-refresher.py || true; sleep 60; done" &' >> /usr/bin/entrypoint-command.sh
    # add the main command
    echo "$@" >> /usr/bin/entrypoint-command.sh
    # make the script executable
    chmod +x /usr/bin/entrypoint-command.sh
    # start the script via xinit
    /usr/bin/xinit /usr/bin/entrypoint-command.sh -- :0 -nolisten tcp vt1
fi
