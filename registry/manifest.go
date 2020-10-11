package registry

import (
	"io"
	"os"
)

type Manifest struct {
	image        string
	fileLocation string
	digest       string
}

func (m *Manifest) WriteManifest(r io.Reader) error {
	f, err := os.Create(m.fileLocation)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, r)

	if err != nil {
		return err
	}

	return nil
}

func (m *Manifest) SendData(w io.Writer) error {
	f, err := os.Open(m.fileLocation)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(w, f)

	return err
}

func (m *Manifest) Delete() error {
	err := os.Remove(m.fileLocation)
	if err != nil {
		return err
	}

	return nil
}
