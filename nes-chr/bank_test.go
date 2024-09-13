package neschr

import (
	"os"
	"testing"
)

func TestCHRBankFromNESFile(t *testing.T) {
	nesFile, err := NewNESFile("../roms/test2.nes")
	if err != nil {
		t.Fatal(err)
	}

	chr := nesFile.CHR()

	// Get the first bank only
	bank, err := chr.Bank(0)
	if err != nil {
		t.Fatal(err)
	}

	// Saving its raw content as file
	err = Save(bank, "test2-bank0.bank")
	if err != nil {
		t.Fatal(err)
	}

	// Saving its content as image file (PNG)
	err = SaveImage(bank, "test2-bank0.png")
	if err != nil {
		t.Fatal(err)
	}

	// Opening another NES ROM
	nesFile2, err := NewNESFile("../roms/test.nes")
	if err != nil {
		t.Fatal(err)
	}

	chr2 := nesFile2.CHR()

	// Retrieve bank from its image file
	bank1 := NewCHREmptyBank()
	err = SetFromImageFile(bank1, "test2-bank0.png")
	if err != nil {
		t.Fatal(err)
	}

	// Update its 7th bank with another
	err = chr2.SetBank(1, bank1)
	if err != nil {
		t.Fatal(err)
	}

	// Saving its content as image file (PNG)
	err = SaveImage(chr2, "test-custom-bank1.png")
	if err != nil {
		t.Fatal(err)
	}

	// Update the NES ROM
	err = nesFile2.UpdateCHR(chr2)
	if err != nil {
		t.Fatal(err)
	}

	// Save the custom NES ROM
	err = nesFile2.Save("test2-custom.nes")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCHRBankCleanup(t *testing.T) {
	paths := []string{
		"test2-bank0.bank",
		"test2-bank0.png",
		"test-custom-bank1.png",
		"test2-custom.nes",
	}

	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			t.Fatal(err)
		}
	}
}
