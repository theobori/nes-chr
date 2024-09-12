package neschr

import (
	"os"
	"slices"
	"testing"
)

const (
	NESFilePath      = "../roms/test.nes"
	NESCHRROMPath    = "../roms/test.chr"
	NESTestImagePath = "../test.png"
)

func TestCHRFromFile(t *testing.T) {
	nesFile, err := NewNESFile(NESFilePath)
	if err != nil {
		t.Fatal(err)
	}

	chr := nesFile.CHR()

	err = chr.SetColorPalette(1)
	if err != nil {
		t.Fatal(err)
	}

	err = chr.SaveImage(NESTestImagePath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCHRInjectImage(t *testing.T) {
	chr := NewEmptyNESCHR()

	err := chr.SetColorPalette(1)
	if err != nil {
		t.Fatal(err)
	}

	err = chr.InjectImageFromFile(NESTestImagePath)
	if err != nil {
		t.Fatal(err)
	}

	original, err := NewNESCHRFromFile(NESCHRROMPath)
	if err != nil {
		t.Fatal(err)
	}

	err = chr.SaveImage("test.png")
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(chr.Chunk(), original.Chunk()) {
		t.Fatal()
	}
}

func TestCHRCleanup(t *testing.T) {
	err := os.Remove(NESTestImagePath)
	if err != nil {
		t.Fatal(err)
	}
}
