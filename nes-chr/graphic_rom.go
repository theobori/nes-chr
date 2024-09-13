package neschr

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"github.com/theobori/nes-chr/internal/util"
)

type GraphicROM interface {
	Chunk() []byte
	ImageSize() (int, int)
	Image() image.Image
	SetFromImage(img image.Image) error
}

func SetFromImageFile(gROM GraphicROM, path string) error {
	chunk, err := util.ReadChunk(path)
	if err != nil {
		return err
	}

	r := bytes.NewReader(chunk)

	img, err := png.Decode(r)
	if err != nil {
		return err
	}

	return gROM.SetFromImage(img)
}

func Save(gROM GraphicROM, path string) error {
	return util.WriteChunk(path, gROM.Chunk())
}

func SaveImage(gROM GraphicROM, path string) error {
	img := gROM.Image()

	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
