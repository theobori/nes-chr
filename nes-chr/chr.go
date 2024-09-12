package neschr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/theobori/nes-chr/internal/util"
)

const (
	CHRImagePixelWidth  = 128
	CHRImagePixelHeight = 128
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

func isColorEqual(src, dest color.Color) (bool, error) {
	err := fmt.Errorf("unable to assert to RGBA")

	srcRGBA, ok := src.(color.RGBA)
	if !ok {
		return false, err
	}

	destRGBA, ok := dest.(color.RGBA)
	if !ok {
		return false, err
	}

	return srcRGBA == destRGBA, nil
}

func (cp *Colorpalette) Index(dest color.Color) (int, error) {
	for i, src := range cp {
		ok, err := isColorEqual(src, dest)
		if err != nil {
			return -1, err
		}

		if ok {
			return i, nil
		}
	}

	return -1, fmt.Errorf("this color does not exists")
}

func isImageSizeValid(w int, h int) bool {
	return w == CHRImagePixelWidth && h%CHRImagePixelHeight == 0 &&
		h >= CHRImagePixelHeight
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

func NewNESCHRFromFile(path string) (*NESCHR, error) {
	b, err := util.ReadChunk(path)
	if err != nil {
		return nil, err
	}

	return NewNESCHR(b), nil
}

func NewEmptyNESCHR() *NESCHR {
	return NewNESCHR([]byte{})
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
	palette, err := ColorPaletteFromIndex(index)
	if err != nil {
		return err
	}

	chr.SetCustomColorPalette(*palette)

	return nil
}

// Returns the image size in pixels
func (chr *NESCHR) ImageSize() (int, int) {
	if chr.bankAmount == 0 {
		return 0, 0
	}

	return CHRImagePixelWidth, CHRImagePixelHeight * chr.bankAmount
}

func (chr *NESCHR) extractBank(img *image.RGBA, bank []byte, offsetY int) {
	memY, memX := 0, 0

	for b := 0; b < CHRRomBankSize; b += 16 {
		for y := 0; y < 8; y++ {
			if memX >= CHRImagePixelWidth {
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
	w, h := chr.ImageSize()
	rect := image.Rect(0, 0, w, h)

	img := image.NewRGBA(rect)

	for i := 0; i < chr.bankAmount; i++ {
		bank, err := chr.Bank(i)
		if err != nil {
			return nil, err
		}

		chr.extractBank(img, bank, CHRImagePixelWidth*i)
	}

	return img, nil
}

func (chr *NESCHR) copyFromChunk(chunk []byte) *NESCHR {
	dest := NewNESCHR(chunk)
	dest.SetCustomColorPalette(chr.colorPalette)

	*chr = *dest

	return chr
}

func (chr *NESCHR) injectPixels(img image.Image, w, h int) error {
	// Prepare the chunk to store the CHR data
	chunk := make([]byte, (h/w)*CHRRomBankSize)
	chunkIndex := 0

	for tileY := 0; tileY < h; tileY += 8 {
		for tileX := 0; tileX < w; tileX += 8 {
			for y := 0; y < 8; y++ {
				var lo, hi byte = 0, 0

				for bit := 0; bit < 8; bit++ {
					color := img.At(tileX+bit, tileY+y)
					colorValue, err := chr.colorPalette.Index(color)
					if err != nil {
						return err
					}

					lo |= (byte(colorValue) & 1) << (7 - bit)
					hi |= (byte(colorValue) >> 1 & 1) << (7 - bit)
				}

				chunk[chunkIndex] = lo
				chunk[chunkIndex+8] = hi
				chunkIndex++
			}

			chunkIndex += 8
		}
	}

	chr.copyFromChunk(chunk)

	return nil
}

func (chr *NESCHR) InjectImage(img image.Image) error {
	// Check if the image has the same width and height
	// as the extracted one with `ExtractImage`.
	bounds := img.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y

	romW, romH := chr.ImageSize()

	// If there is an existing ROM, it compares the ROM and the image size
	if len(chr.chunk) != 0 && (w != romW || h != romH) {
		return fmt.Errorf("invalid image size, it must be %dx%d", romW, romH)
	}

	// Check image width and height
	if !isImageSizeValid(w, h) {
		return fmt.Errorf("invalid image size")
	}

	// Encode pixels into the CHR memory
	return chr.injectPixels(img, w, h)
}

func (chr *NESCHR) InjectImageFromFile(path string) error {
	chunk, err := util.ReadChunk(path)
	if err != nil {
		return err
	}

	r := bytes.NewReader(chunk)

	img, err := png.Decode(r)
	if err != nil {
		return err
	}

	return chr.InjectImage(img)
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
