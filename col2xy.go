package col2xy

import (
	"image/color"
	"math"
)

const (
	maxColorF = float64(0xffff)
	maxByteF  = float64(255)
)

// ApplyGammaCorrection returns gamma corrected value of a given normalized
// (within range [0, 1]) color component.
func ApplyGammaCorrection(c float64) float64 {
	if c > 0.04045 {
		return math.Pow((c+0.055)/(1+0.055), 2.4)
	}

	return c / 12.92
}

// Normalized2XY returns x and y coordinates on CIE 1931 chromaticity diagram
// of given normalized (within range [0, 1]) color components.
// Theoretical source here:
// https://gist.github.com/popcorn245/30afa0f98eea1c2fd34d
// the same on StackOverflow:
// https://stackoverflow.com/questions/70612781/how-to-convert-ikea-light-bulb-color-xy-cie-1931-colorspace-to-rgb
// Other sources:
// https://en.wikipedia.org/wiki/CIE_1931_color_space#CIE_RGB_color_space
// https://stackoverflow.com/questions/54663997/convert-rgb-color-to-xy
// https://www.researchgate.net/publication/50863695_Comparison_between_Digital_Image_Processing_and_Spectrophotometric_Measurements_Methods
func Normalized2XY(r, g, b float64) (float64, float64) {
	cr := ApplyGammaCorrection(r)
	cg := ApplyGammaCorrection(g)
	cb := ApplyGammaCorrection(b)

	X := cr*0.6491852651246980 + cg*0.1034883891428110 + cb*0.1973263457324920
	Y := cr*0.2340599935483600 + cg*0.7433166037561910 + cb*0.0226234026954449
	Z := cr*0.0 + cg*0.0530940431254422 + cb*1.0369059568745600

	sum := X + Y + Z
	x := X / sum
	y := Y / sum

	return x, y
}

// NormalizeRGB returns normalized (within range [0, 1]) color components of
// given color components represented as bytes.
func NormalizeRGB(r, g, b byte) (float64, float64, float64) {
	return float64(r) / maxByteF,
		float64(g) / maxByteF,
		float64(b) / maxByteF
}

// RGB2XY returns x and y coordinates on CIE 1931 chromaticity diagram of
// given color components represented as bytes.
func RGB2XY(r, g, b byte) (float64, float64) {
	return Normalized2XY(NormalizeRGB(r, g, b))
}

// NormalizeColor returns normalized (within range [0, 1]) color components of
// a given color.Color.
func NormalizeColor(c color.Color) (float64, float64, float64) {
	r, g, b, _ := c.RGBA()
	return float64(r) / maxColorF,
		float64(g) / maxColorF,
		float64(b) / maxColorF
}

// Color2XY returns x and y coordinates on CIE 1931 chromaticity diagram of a
// given color.
func Color2XY(c color.Color) (float64, float64) {
	return Normalized2XY(NormalizeColor(c))
}
