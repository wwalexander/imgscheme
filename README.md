imgscheme
=========

Generates terminal color schemes from images and existing color schemes

Building
--------

    go build

Usage
-----

    imgscheme [-base path] [path]

imgscheme generates a terminal color scheme from the image located at path. It
attempts to make the generated colors correspond to prominent colors in the
image, and also close to the corresponding colors in a base color scheme. A base
color scheme file can be specified using the -base flag; this file should
contain a newline-separated list of hex triplets. The generated colors will be
output in the same order as the base color scheme. If the -base flag is not set,
imgscheme will use the standard VGA color scheme; the first 8 triplets of the
output will be the 8 normal colors in order (black, red, green, yellow, blue,
magenta, cyan, and white) and the next 8 triplets will be the 8 bold colors in
the same order.

The output of imgscheme can be rendered on
[this page](https://wwalexander.github.io/imgscheme/).

Options
-------

`-base`: the base color scheme file to use
