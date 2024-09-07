package neschr

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/theobori/nes-chr/internal/util"
)

const (
	CHRImageWidth  = 128
	CHRImageHeight = 128
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

type NESCHR struct {
	chunk        []byte
	colorPalette Colorpalette
	bankAmount   int
}

func NewNESCHR(chunk []byte) *NESCHR {
	return &NESCHR{
		chunk:        chunk,
		colorPalette: ColorPaletteFromScheme(DefaultColorScheme),
		bankAmount:   len(chunk) / CHRRomBankSize,
	}
}

func (chr *NESCHR) Chunk() []byte {
	return chr.chunk
}

func (chr *NESCHR) Bank(n int) ([]byte, error) {
	if n < 0 || n >= chr.bankAmount {
		return nil, fmt.Errorf("%d is out of the CHR ROM", n)
	}

	start := CHRRomBankSize * n
	bank := chr.chunk[start : start+CHRRomBankSize]

	return bank, nil
}

func (chr *NESCHR) SetCustomColorPalette(palette Colorpalette) {
	chr.colorPalette = palette
}

func (chr *NESCHR) SetColorPalette(index int) error {
	if index < 0 || index >= len(ColorSchemes) {
		return fmt.Errorf("palette are between %d and %d", 0, len(ColorSchemes))
	}

	palette := ColorPaletteFromScheme(ColorSchemes[index])

	chr.SetCustomColorPalette(palette)

	return nil
}

func (chr *NESCHR) extractBank(img *image.RGBA, bank []byte, offsetY int) {
	memY, memX := 0, 0

	for b := 0; b < CHRRomBankSize; b += 16 {
		for y := 0; y < 8; y++ {
			if memX >= CHRImageWidth {
				memY += 8
				memX = 0
			}

			lo := bank[b+y]
			hi := bank[b+y+8]

			for bit := 0; bit < 8; bit++ {
				left := lo >> (7 - bit) & 1
				right := hi >> (7 - bit) & 1
				color := chr.colorPalette[right<<1|left]

				img.Set(bit+memX, y+memY+offsetY, color)
			}
		}

		memX += 8
	}
}

func (chr *NESCHR) ExtractImage() (*image.RGBA, error) {
	rect := image.Rect(0, 0, CHRImageWidth, CHRImageHeight*chr.bankAmount)
	img := image.NewRGBA(rect)

	for i := 0; i < chr.bankAmount; i++ {
		bank, err := chr.Bank(i)
		if err != nil {
			return nil, err
		}

		chr.extractBank(img, bank, CHRImageWidth*i)
	}

	return img, nil
}

func (chr *NESCHR) InjectImage() {

}

func (chr *NESCHR) Save(path string) error {
	return util.WriteChunk(path, chr.chunk)
}

func (chr *NESCHR) SaveImage(path string) error {
	img, err := chr.ExtractImage()
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return nil
	}

	return nil
}
