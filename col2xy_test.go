package col2xy

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"testing"

	"image/color"
	"image/draw"
	"image/gif"
	_ "image/jpeg"

	"github.com/stretchr/testify/assert"
)

type rgb2xyTestCase struct {
	r byte
	g byte
	b byte
	x float64
	y float64
}

func TestRGB2XY(t *testing.T) {
	testCases := map[string]rgb2xyTestCase{
		"Red": {
			0xff, 0x00, 0x00,
			0.7350000000000004, 0.26499999999999957,
		},
		"Green": {
			0x00, 0xff, 0x00,
			0.11499999999999991, 0.8260000000000001,
		},
		"Blue": {
			0x00, 0x00, 0xff,
			0.157, 0.017999999999999964,
		},
		"White": {
			0xff, 0xff, 0xff,
			0.3125000000000004, 0.3289473684210514,
		},
	}

	for txt, tc := range testCases {
		t.Run(txt, func(t *testing.T) {
			x, y := RGB2XY(tc.r, tc.g, tc.b)
			assert.Equal(t, tc.x, x)
			assert.Equal(t, tc.y, y)
		})
	}
}

type color2XYTestCase struct {
	c color.Color
	x float64
	y float64
}

func TestColor2XY(t *testing.T) {
	testCases := map[string]color2XYTestCase{
		"Red": {
			color.RGBA64{0xffff, 0x0000, 0x0000, 0xffff},
			0.7350000000000004, 0.26499999999999957,
		},
		"Green": {
			color.RGBA64{0x0000, 0xffff, 0x0000, 0xffff},
			0.11499999999999991, 0.8260000000000001,
		},
		"Blue": {
			color.RGBA64{0x0000, 0x0000, 0xffff, 0xffff},
			0.157, 0.017999999999999964,
		},
		"White": {
			color.RGBA64{0xffff, 0xffff, 0xffff, 0xffff},
			0.3125000000000004, 0.3289473684210514,
		},
	}

	for color, tc := range testCases {
		t.Run(color, func(t *testing.T) {
			x, y := Color2XY(tc.c)
			assert.Equal(t, tc.x, x)
			assert.Equal(t, tc.y, y)

		})
	}
}

func ExampleRGB2XY() {
	decodeImgFile := func(filename string) (image.Image, int, int) {
		bs, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		img, _, err := image.Decode(bytes.NewReader(bs))
		if err != nil {
			log.Fatal(err)
		}

		b := img.Bounds()

		return img, b.Dx(), b.Dy()
	}

	chromaImg, chromaW, chromaH := decodeImgFile("cie_1931_chromaticity_diagram.jpg")
	chromaPixPerX, chromaPixPerY := float64(chromaW)/0.8, float64(chromaH)/0.9

	trackImg, trackW, trackH := decodeImgFile("gradient_track.jpg")

	outputW, outputH := chromaW, trackH+chromaH+5 // add some whitespace
	outputRect := image.Rect(0, 0, outputW, outputH)
	outputTrackRect := image.Rect(0, 0, trackW, trackH)
	outputChromaRect := image.Rect(0, trackH+5, outputW, outputH)
	outputImages := []*image.Paletted{}
	outputDelays := []int{}

	palette := []color.Color{
		image.White,
		image.Black,
	}
	getColorAtXPercent := func(p int) color.Color {
		x := trackW * p / 100
		if x != 0 {
			x--
		}

		return trackImg.At(x, trackH/2-1)
	}
	// Assume we have a slider input with min at 0 and max at 100
	for p := 0; p <= 100; p++ {
		palette = append(palette, getColorAtXPercent(p))
	}

	crossAt := func(x, y float64, img *image.Paletted) {
		cx, cy := int(x*chromaPixPerX), outputH-int(y*chromaPixPerY)
		img.Set(cx, cy, image.Black)

		for j := 0; j < 10; j++ {
			img.Set(cx+j, cy, image.Black)
			img.Set(cx-j, cy, image.Black)
			img.Set(cx, cy+j, image.Black)
			img.Set(cx, cy-j, image.Black)
		}
	}

	lineAt := func(p int, img *image.Paletted) {
		x := trackW * p / 100
		if x != 0 {
			x--
		}

		for j := 0; j < trackH; j++ {
			img.Set(x, j, image.Black)
		}
	}

	addFrame := func(p int) {
		outputImg := image.NewPaletted(outputRect, palette)

		draw.FloydSteinberg.Draw(outputImg, outputTrackRect, trackImg, image.Point{})
		draw.FloydSteinberg.Draw(outputImg, outputChromaRect, chromaImg, image.Point{})
		lineAt(p, outputImg)

		c := getColorAtXPercent(p)
		x, y := Color2XY(c)
		crossAt(x, y, outputImg)

		outputImages = append(outputImages, outputImg)
		outputDelays = append(outputDelays, 10)
	}
	for p := 0; p <= 100; p++ {
		addFrame(p)
	}
	for p := 100; p >= 0; p-- {
		addFrame(p)
	}

	outF, err := os.OpenFile("demo.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer outF.Close()

	if err := gif.EncodeAll(outF, &gif.GIF{
		Image: outputImages,
		Delay: outputDelays,
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
	// Output: DONE
}
