package neschr

import (
	"fmt"
	"image"

	"github.com/theobori/nes-chr/internal/util"
)

const (
	CHRBankImageWidth  = 128
	CHRBankImageHeight = 128
)

type CHRBank struct {
	chunk        []byte
	colorPalette Colorpalette
}

func NewCHRBank(chunk []byte) *CHRBank {
	return &CHRBank{
		chunk:        chunk,
		colorPalette: ColorPaletteFromScheme(DefaultColorScheme),
	}
}

func NewCHRBankFromFile(path string) (*CHRBank, error) {
	b, err := util.ReadChunk(path)
	if err != nil {
		return nil, err
	}

	return NewCHRBank(b), nil
}

func NewCHREmptyBank() *CHRBank {
	return NewCHRBank([]byte{})
}

func (b *CHRBank) SetCustomColorPalette(palette Colorpalette) {
	b.colorPalette = palette
}

func (b *CHRBank) SetColorPalette(index int) error {
	palette, err := ColorPaletteFromIndex(index)
	if err != nil {
		return err
	}

	b.SetCustomColorPalette(*palette)

	return nil
}

func (b *CHRBank) Chunk() []byte {
	return b.chunk
}

func (b *CHRBank) IsEmpty() bool {
	return len(b.chunk) == 0
}

func (b *CHRBank) ImageSize() (int, int) {
	if b.IsEmpty() {
		return 0, 0
	}

	return CHRBankImageWidth, CHRBankImageHeight
}

func (b *CHRBank) Image() image.Image {
	rect := image.Rect(0, 0, CHRBankImageWidth, CHRBankImageHeight)
	img := image.NewRGBA(rect)

	memY, memX := 0, 0
	for i := 0; i < CHRRomBankSize; i += 16 {
		for y := 0; y < 8; y++ {
			if memX >= CHRBankImageWidth {
				memY += 8
				memX = 0
			}

			lo := b.chunk[i+y]
			hi := b.chunk[i+y+8]

			for bit := 0; bit < 8; bit++ {
				left := lo >> (7 - bit) & 1
				right := hi >> (7 - bit) & 1
				color := b.colorPalette[right<<1|left]

				img.Set(bit+memX, y+memY, color)
			}
		}

		memX += 8
	}

	return img
}

func (b *CHRBank) SetFromImage(img image.Image) error {
	bounds := img.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y

	if w != CHRBankImageWidth || h != CHRBankImageHeight {
		return fmt.Errorf(
			"invalid image size, it must be %dx%d",
			CHRBankImageWidth,
			CHRBankImageHeight,
		)
	}

	chunk := make([]byte, (h/w)*CHRRomBankSize)
	chunkIndex := 0

	for tileY := 0; tileY < h; tileY += 8 {
		for tileX := 0; tileX < w; tileX += 8 {
			for y := 0; y < 8; y++ {
				var lo, hi byte = 0, 0

				for bit := 0; bit < 8; bit++ {
					color := img.At(tileX+bit, tileY+y)
					colorValue := b.colorPalette.Index(color)

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

	b.chunk = make([]byte, len(chunk))
	copy(b.chunk, chunk)

	return nil
}
