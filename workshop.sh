#!/bin/bash

xdotool key Ctrl+plus
xdotool key Ctrl+plus
glow --pager --width 100 ./workshop-slides/part-1.md
glow --pager --width 100 ./workshop-slides/part-2.md
xdotool key Ctrl+minus
xdotool key Ctrl+minus
glow --pager --width 150 ./workshop-slides/part-3.md
xdotool key Ctrl+plus
xdotool key Ctrl+plus
glow --pager --width 120 ./workshop-slides/part-4.md
xdotool key Ctrl+minus
xdotool key Ctrl+minus
