package neschr

import (
	"testing"
)

const (
	NESFilePath = "../roms/test.nes"
)

func TestNESFile(t *testing.T) {
	nesFile, err := NewNESFile(NESFilePath)
	if err != nil {
		t.Fatal(err)
	}

	chr := nesFile.CHR()

	err = chr.SetColorPalette(1)
	if err != nil {
		t.Fatal(err)
	}

	err = chr.SaveImage("./test.png")
	if err != nil {
		t.Fatal(err)
	}
}
