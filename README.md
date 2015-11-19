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
the standard order. imgscheme might fail to find a suitable color for a place,
in which case it will print `<nil>`.
