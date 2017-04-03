imgscheme
=========

Generates terminal color schemes from images and existing color schemes

Building
--------

    go build

Usage
-----

    imgscheme [-[n] path] [path]

imgscheme generates a terminal color scheme from the image at the named path. It
attempts to make the generated colors correspond to prominent colors in the
image, and also close to the corresponding colors in a base color scheme. The
colors in the base scheme can be set using numeric flags corresponding to the 16
ANSI colors (0-15). The generated colors will be output in order as hex
triplets.

The output of imgscheme can be rendered on
[this page](https://wwalexander.github.io/imgscheme/).
