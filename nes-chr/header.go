package neschr

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var (
	magicBytes = [4]byte{'N', 'E', 'S', '\x41'}
)

type NESHeader struct {
	Magic      [4]byte
	PRGROMSize byte
	CHRROMSize byte
	Flags6     byte
	Flags7     byte
	PRGRAMSize byte
	Flags9     byte
	Flags10    byte
	_          [5]byte // Reserved
}

func (h *NESHeader) IsValid() bool {
	return h.Magic == magicBytes
}

func (h *NESHeader) HasTrainer() bool {
	return h.Flags6&0b0000_0100 != 0
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
