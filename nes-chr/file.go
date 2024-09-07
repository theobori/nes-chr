package neschr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/theobori/nes-chr/internal/util"
)

const (
	HeaderSize     = 16
	TrainerSize    = 512
	PRGRomUnitSize = 0x4000
	CHRRomUnitSize = 0x2000
	CHRRomBankSize = 0x1000
)

type BlockMetadata struct {
	Size     int
	Position int
}

func (bm *BlockMetadata) EndPosition() int {
	return bm.Position + bm.Size
}

func NewNESHeader(chunk []byte) (*NESHeader, error) {
	var header NESHeader

	if len(chunk) < 16 {
		return nil, fmt.Errorf("chunk size must be at least 16 bytes")
	}

	r := bytes.NewReader(chunk)

	err := binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}

type NESFile struct {
	header      *NESHeader
	chunk       []byte
	CHRMetadata *BlockMetadata
}

func NewNESFileFromBytes(chunk []byte) (*NESFile, error) {
	header, err := NewNESHeader(chunk)
	if err != nil {
		return nil, err
	}

	nesFile := &NESFile{
		header:      header,
		chunk:       chunk,
		CHRMetadata: nil,
	}

	err = nesFile.Parse()
	if err != nil {
		return nil, err
	}

	return nesFile, nil
}

func NewNESFile(path string) (*NESFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	chunk, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return NewNESFileFromBytes(chunk)
}

func (nf *NESFile) parseCHR() error {
	prgSize := int(nf.header.PRGROMSize) * PRGRomUnitSize
	chrSize := int(nf.header.CHRROMSize) * CHRRomUnitSize

	if chrSize == 0 {
		return fmt.Errorf("missing CHR ROM")
	}

	pos := HeaderSize

	if nf.header.HasTrainer() {
		pos += TrainerSize
	}

	// Just after the program ROM
	pos += prgSize

	if len(nf.chunk) < pos+chrSize {
		return fmt.Errorf("nes file rom is smaller than %d", pos+chrSize)
	}

	nf.CHRMetadata = &BlockMetadata{
		Size:     chrSize,
		Position: pos,
	}

	return nil
}

func (nf *NESFile) Parse() error {
	return nf.parseCHR()
}

func (nf *NESFile) CHR() *NESCHR {
	m := nf.CHRMetadata
	chrChunk := nf.chunk[m.Position:m.EndPosition()]

	return NewNESCHR(chrChunk)
}

func (nf *NESFile) UpdateCHR(chr *NESCHR) error {
	m := nf.CHRMetadata

	chrChunk := chr.Chunk()

	if len(chrChunk) != m.Size {
		return fmt.Errorf("inconsistent CHR, %d bytes and %d bytes", len(chrChunk), m.Size)
	}

	for i := 0; i < m.Size; i++ {
		nf.chunk[m.Position+i] = chrChunk[i]
	}

	return nil
}

func (nf *NESFile) Save(path string) error {
	return util.WriteChunk(path, nf.chunk)
}