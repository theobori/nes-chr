package neschr

import (
	"fmt"
	"image"

	"github.com/theobori/nes-chr/internal/util"
)

func isCHRImageSizeValid(w int, h int) bool {
	return w == CHRBankImageWidth && h%CHRBankImageHeight == 0 &&
		h >= CHRBankImageHeight
}

type CHR struct {
	banks []CHRBank
}

func NewCHR(chunk []byte) *CHR {
	banks := []CHRBank{}

	if len(chunk) < CHRRomBankSize {
		return &CHR{
			banks: banks,
		}
	}

	for i := 0; i < len(chunk); i += CHRRomBankSize {
		bankChunk := chunk[i : i+CHRRomBankSize]
		banks = append(banks, *NewCHRBank(bankChunk))
	}

	return &CHR{
		banks: banks,
	}
}

func NewCHRFromFile(path string) (*CHR, error) {
	b, err := util.ReadChunk(path)
	if err != nil {
		return nil, err
	}

	return NewCHR(b), nil
}

func NewEmptyCHR() *CHR {
	return NewCHR([]byte{})
}

func (chr *CHR) BankAmount() int {
	return len(chr.banks)
}

func (chr *CHR) IsEmpty() bool {
	return chr.BankAmount() == 0
}

func (chr *CHR) Bank(i int) (*CHRBank, error) {
	if i < 0 || i >= len(chr.banks) {
		return nil, fmt.Errorf("bank index is out of range")
	}

	return &chr.banks[i], nil
}

func (chr *CHR) SetBank(i int, bank *CHRBank) error {
	srcBank, err := chr.Bank(i)
	if err != nil {
		return err
	}

	*srcBank = *bank

	return nil
}

func (chr *CHR) AddBank(bank *CHRBank) {
	chr.banks = append(chr.banks, *bank)
}

func (chr *CHR) Chunk() []byte {
	chunk := []byte{}

	for i := 0; i < chr.BankAmount(); i++ {
		bank, _ := chr.Bank(i)
		chunk = append(chunk, bank.Chunk()...)
	}

	return chunk
}

// Returns the image size in pixels
func (chr *CHR) ImageSize() (int, int) {
	if chr.BankAmount() == 0 {
		return 0, 0
	}

	return CHRBankImageWidth, CHRBankImageHeight * chr.BankAmount()
}

func (chr *CHR) Image() image.Image {
	w, h := chr.ImageSize()
	rect := image.Rect(0, 0, w, h)

	img := image.NewRGBA(rect)

	offset := 0

	for i := 0; i < chr.BankAmount(); i++ {
		bankImage := chr.banks[i].Image()

		for y := 0; y < CHRBankImageHeight; y++ {
			for x := 0; x < CHRBankImageWidth; x++ {
				c := bankImage.At(x, y)
				img.Set(x, offset+y, c)
			}
		}

		offset += CHRBankImageHeight
	}

	return img
}

func (chr *CHR) SetFromImage(img image.Image) error {
	bounds := img.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y

	isEmpty := chr.IsEmpty()

	if !isEmpty && w != CHRBankImageWidth && h != chr.BankAmount()*CHRBankImageHeight {
		return fmt.Errorf(
			"invalid image size, it must be %dx%d",
			CHRBankImageWidth,
			chr.BankAmount()*CHRBankImageHeight,
		)
	}

	if !isCHRImageSizeValid(w, h) {
		return fmt.Errorf("invalid image size, %dx%d", w, h)
	}

	for y := 0; y < h; y += CHRBankImageHeight {
		bankImage, err := util.ImageRectFromImage(
			img,
			0,
			y,
			CHRBankImageWidth,
			CHRBankImageHeight,
		)
		if err != nil {
			return err
		}

		if isEmpty {
			bank := NewCHREmptyBank()
			err := bank.SetFromImage(bankImage)
			if err != nil {
				return err
			}

			chr.AddBank(bank)
		} else {
			err := chr.banks[y/CHRBankImageHeight].SetFromImage(bankImage)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
