package util

import "os"

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
