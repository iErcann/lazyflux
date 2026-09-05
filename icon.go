package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// Draws the tray icon: a disc with a bite, colored by temperature.
func iconPNG(temp int) []byte {
	const size = 64
	const (
		cx, cy, r  = 32.0, 32.0, 24.0 // main disc
		bx, by, br = 45.0, 19.0, 16.0 // bite
	)
	fill := tint(temp)
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			px, py := float64(x)+0.5, float64(y)+0.5
			inDisc := (px-cx)*(px-cx)+(py-cy)*(py-cy) < r*r
			inBite := (px-bx)*(px-bx)+(py-by)*(py-by) < br*br
			if inDisc && !inBite {
				img.Set(x, y, fill)
			}
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

// Cool grey at 6500 K, warm orange at 2700 K.
func tint(temp int) color.NRGBA {
	cool := color.NRGBA{R: 150, G: 160, B: 175, A: 255} // 6500 K and above
	warm := color.NRGBA{R: 245, G: 140, B: 50, A: 255}  // 2700 K and below
	t := float64(min(max(temp, 2700), 6500)-2700) / (6500 - 2700)
	mix := func(a, b uint8) uint8 { return uint8(float64(b) + (float64(a)-float64(b))*t) }
	return color.NRGBA{R: mix(cool.R, warm.R), G: mix(cool.G, warm.G), B: mix(cool.B, warm.B), A: 255}
}
