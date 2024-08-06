#! /usr/bin/python3

from PIL import Image
import sys
from datetime import datetime
import pytz
import subprocess

timezone = pytz.timezone('Europe/Rome')

# approach and code based on https://dev.to/kaiwalter/aw-snap-refresh-web-app-in-chromium-kiosk-mode-after-it-breaks-1p4a
def get_main_color(file):
    img = Image.open(file)
    colors = img.getcolors(256*1024) #put a higher value if there are many colors in your image
    max_occurence, most_present = 0, 0
    try:
        for c in colors:
            if c[0] > max_occurence:
                (max_occurence, most_present) = c
        print(datetime.now().strftime("%Y-%m-%d %H:%M:%S"), 'most present color: ', most_present, 'OK!')
        return most_present
    except TypeError:
        # Too many colors in the image
        print("error detecting main screen color")
    return (1, 1, 1)


# create a screenshot, and overwrite if it already exists.
subprocess.run(["scrot", "-o", "screen.png"])

# if main color of the screenshot is white (typical of the chromium Aw,snap crash page),
# issue a refresh of the chromium browser window:
if get_main_color('screen.png') == (255, 255, 255):
    # a chromium browser crash was detected
    print(datetime.now().strftime("%Y-%m-%d %H:%M:%S"), 'detected Aw,snap', 'attempting refresh...')
    # issue a chromium browser window refresh
    subprocess.run(["xdotool", "getwindowfocus", "key", "ctrl+F5"])
    # exit
    sys.exit(1)
