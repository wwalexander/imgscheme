imgscheme
=========

Generates terminal color schemes from images

Building
--------

    go build

Usage
-----

Generate a color scheme from the image file at `PATH`:

    imgscheme [PATH]

PNG, JPEG, and GIF images are supported. Currently, imgscheme runs very slowly;
an interrupt signal can be sent at any time, at which point imgscheme will
print the colors generated so far. Colors are printed as RGB tuples in the
standard order. If no valid colors are found for a given place, imgscheme will
use black.
