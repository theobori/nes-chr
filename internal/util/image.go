package util

import (
	"fmt"
	"image"
)

func ImageRectFromImage(img image.Image, x, y, w, h int) (image.Image, error) {
	bounds := img.Bounds()
	wMax := bounds.Max.X
	hMax := bounds.Max.Y

	if x+w > wMax || y+h > hMax {
		return nil, fmt.Errorf("out of bounds")
	}

	rect := image.Rect(0, 0, w, h)
	newImg := image.NewRGBA(rect)

	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			newImg.Set(xx, yy, img.At(xx+x, yy+y))
		}
	}

	return newImg, nil
}
