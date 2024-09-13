package neschr

import (
	"os"
	"slices"
	"testing"
)

const (
	NESFilePath      = "../roms/test.nes"
	CHRROMPath       = "../roms/test.chr"
	NESTestImagePath = "../test.png"
)

func TestCHRFromNESFile(t *testing.T) {
	nesFile, err := NewNESFile(NESFilePath)
	if err != nil {
		t.Fatal(err)
	}

	chr := nesFile.CHR()

	err = SaveImage(chr, NESTestImagePath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCHRInjectImage(t *testing.T) {
	chr := NewEmptyCHR()

	err := SetFromImageFile(chr, NESTestImagePath)
	if err != nil {
		t.Fatal(err)
	}

	original, err := NewCHRFromFile(CHRROMPath)
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
