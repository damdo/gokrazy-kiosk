FROM debian:11-slim

RUN apt-get update -y && \
    DEBIAN_FRONTEND='noninteractive' apt-get install --no-install-recommends -y \
    xorg xserver-xorg-input-evdev xserver-xorg-input-all chromium wmctrl xdotool

RUN printf "\
Section \"ServerFlags\"\n\
    Option \"AutoAddDevices\" \"False\"\n\
    Option \"StandbyTime\" \"0\"\n\
    Option \"SuspendTime\" \"0\"\n\
    Option \"OffTime\" \"0\"\n\
    Option \"BlankTime\" \"0\"\n\
EndSection\n\
\n\
Section \"ServerLayout\"\n\
    Identifier     \"Desktop\"\n\
    InputDevice    \"Mouse0\" \"CorePointer\"\n\
EndSection\n\
\n\
Section \"InputDevice\"\n\
    Identifier \"Mouse0\"\n\
    Option \"CoreKeyboard\"\n\
    Option \"Device\" \"/dev/input/event4\"\n\
    Driver \"evdev\"\n\
EndSection\n\
" > /etc/X11/xorg.conf.d/10-input.conf

ADD entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod +x /usr/bin/entrypoint.sh
ENTRYPOINT ["/usr/bin/entrypoint.sh"]
