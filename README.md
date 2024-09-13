# NES CHR ROM manipulation library

[![build](https://github.com/theobori/nes-chr/actions/workflows/build.yml/badge.svg)](https://github.com/theobori/nes-chr/actions/workflows/build.yml) [![lint](https://github.com/theobori/nes-chr/actions/workflows/lint.yml/badge.svg)](https://github.com/theobori/nes-chr/actions/workflows/lint.yml)

This is a highly flexible toy KISS Go module designed to help you extract the CHR ROM and its banks from a NES game or inject an images into the game.

## 📖 Build and run

For the build, you only need the following requirements:

- [Go](https://golang.org/doc/install) >= 1.22.0

## 🤝 Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).

## 🧪 Tests

The test ROMS are legal and can be found on the following sources.
- [ppu_read_buffer.zip](http://bisqwit.iki.fi/src/nes_tests/ppu_read_buffer.zip)
- [holydiverbatman-bin-0.01.7z](https://pineight.com/nes/holydiverbatman-bin-0.01.7z)

To run the tests you can run the following command.

```bash
make test
```

## 📎 Some examples

Here is a basic example of how you could use the module.

```go
package main

import (
	"log"

	"github.com/theobori/nes-chr/nes-chr"
)

func main() {
	// Open the NES game file
	nesFile, err := neschr.NewNESFile("./smb3.nes")
	if err != nil {
		log.Fatal(err)
	}

	// Get the CHR
	chr := nesFile.CHR()

	// Save the CHR ROM as a file
	err = neschr.Save(chr, "smb3.chr")
	if err != nil {
		log.Fatal(err)
	}

	// Save the CHR as a PNG file
	err = neschr.SaveImage(chr, "smb3.png")
	if err != nil {
		log.Fatal(err)
	}

	// Custom color scheme
	scheme := neschr.ColorScheme{
		{255, 0, 0},
		{0, 255, 0},
		{0, 0, 255},
		{0, 0, 0},
	}

	// Copy the CHR ROM
	chrRGB := neschr.NewCHR(chr.Chunk())

	// Create a color palette from the custom color scheme
	customPalette := neschr.ColorPaletteFromScheme(scheme)

	// Set this palette to banks
	// Note that if you want to create a CHR from this one
	// You will need to have the same exact palette values and order
	// You will have to create am empty CHR with empty banks and set their palette
	// before loading the PNG file
	for i := 0; i < chrRGB.BankAmount(); i++ {
		if i%2 == 0 {
			bank, _ := chrRGB.Bank(i)
			bank.SetCustomColorPalette(customPalette)
		}
	}

	// Save the RGB CHR as a PNG file
	err = neschr.SaveImage(chrRGB, "smb3-rgb-colors.png")
	if err != nil {
		log.Fatal(err)
	}

	// Get the 30th bank of the original CHR ROM
	bank, _ := chr.Bank(30)
	err = neschr.Save(bank, "smb3-bank-30.bank")
	if err != nil {
		log.Fatal(err)
	}

	// Saving it as a PNG
	err = neschr.SaveImage(bank, "smb3-bank-30.png")
	if err != nil {
		log.Fatal(err)
	}

	// Updating the bank 30th from the original CHR ROM
	err = neschr.SetFromImageFile(bank, "smb3-bank-30-custom.png")
	if err != nil {
		log.Fatal(err)
	}

	// Update the NES ROM
	err = nesFile.UpdateCHR(chr)
	if err != nil {
		log.Fatal(err)
	}

	// Save the custom NES ROM
	err = nesFile.Save("smb3-custom.nes")
	if err != nil {
		log.Fatal(err)
	}
}
```

## 🎉 Tasks

- [ ] Support other image formats than PNG
