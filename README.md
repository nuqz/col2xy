# col2xy

This is a simple package for converting [colors](https://pkg.go.dev/image/color#Color) from the Go standard library and colors specified by RGB color components as bytes (within range [0-255]) into `x` and `y` coordinates on the CIE 1931 chromaticity diagram.

## Visual demo

The image shown here shows the transformation only approximately, since the diagram was cut manually in a graphical editor.

<img src="https://github.com/nuqz/col2xy/blob/master/demo.gif" width="325">

## Usage

See [tests](https://github.com/nuqz/col2xy/blob/master/col2xy_test.go#L27) to get an idea how to use this package.
