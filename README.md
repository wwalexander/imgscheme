imgscheme
=========

Generates terminal color schemes from images

Building
--------

    go build

Usage
-----

Generate a color scheme from the image file at `PATH`:

    imgscheme [OPTIONS] [PATH]

PNG, JPEG, and GIF images are supported. Colors are printed as hex triplets in
the same order as the base color scheme. imgscheme might take some time to
complete, especially on large images.

Options
-------

`-base`: the base color scheme file to use, or a built-in color scheme:
    * `vga`
    * `xterm`

Color scheme files
------------------

Color scheme files are lists of RGB colors formatted as hex triplets (i.e.
strings which match the regular expression `^#[0-9a-fA-F]{6}$`. Standard color
schemes contain 16 colors, but the color scheme file can contain any number of
colors.

The output of imgscheme can be rendered on
[this page](https://wwalexander.github.io/imgscheme/).
