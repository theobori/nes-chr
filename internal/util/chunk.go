package util

import (
	"io"
	"os"
)

// Write a chunk (bytes) into a file, existing or not
func WriteChunk(path string, chunk []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return nil
	}

	defer f.Close()

	_, err = f.Write(chunk)
	if err != nil {
		return err
	}

	return nil
}

// Read a chunk (bytes) from a file
func ReadChunk(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	chunk, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}
