package neschr

import (
	"fmt"
	"image/color"
)

type ColorScheme [4][3]byte

var (
	DefaultColorScheme = ColorScheme{
		{0, 0, 0},
		{126, 126, 126},
		{189, 189, 189},
		{255, 255, 255},
	}

	ColorSchemes = []ColorScheme{
		DefaultColorScheme,
		{
			{34, 139, 34},
			{139, 69, 19},
			{210, 180, 140},
			{160, 82, 45},
		},
		{
			{0, 105, 148},
			{135, 206, 235},
			{70, 130, 180},
			{240, 248, 255},
		},
		{
			{255, 69, 0},
			{255, 140, 0},
			{255, 215, 0},
			{139, 0, 0},
		},
		{
			{169, 169, 169},
			{192, 192, 192},
			{105, 105, 105},
			{220, 220, 220},
		},
	}
)

type Colorpalette [4]color.Color

func ColorPaletteFromScheme(scheme ColorScheme) Colorpalette {
	var palette Colorpalette

	for i, c := range scheme {
		palette[i] = color.RGBA{c[0], c[1], c[2], 255}
	}

	return palette
}

func ColorPaletteFromIndex(index int) (*Colorpalette, error) {
	if index < 0 || index >= len(ColorSchemes) {
		return nil, fmt.Errorf("palette are between %d and %d", 0, len(ColorSchemes))
	}

	scheme := ColorSchemes[index]
	palette := ColorPaletteFromScheme(scheme)

	return &palette, nil
}

func (cp *Colorpalette) Index(dest color.Color) int {
	for i, src := range cp {
		ok := isColorEqual(src, dest)
		if ok {
			return i
		}
	}

	return -1
}

// Removed the assertion `src.(color.RGBA)` because it was faily
// for some reasons
func isColorEqual(src, dest color.Color) bool {
	r1, g1, b1, a1 := src.RGBA()
	r2, g2, b2, a2 := dest.RGBA()

	return r1 == r2 &&
		g1 == g2 &&
		b1 == b2 &&
		a1 == a2
}
