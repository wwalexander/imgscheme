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

PNG, JPEG, and GIF images are supported. Colors are printed as hex triplets in
the standard order. imgscheme might take some time to complete, especially on
large images.
