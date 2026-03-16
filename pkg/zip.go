package pkg

import (
	"archive/zip"
	"os"
)

func CreateZip(name string) (*os.File,*zip.Writer, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, nil, err
	}

	return file, zip.NewWriter(file), nil
}

func SaveZip(name string, files []byte, file *zip.Writer) error {
	f, err := file.Create(name)
	if err != nil {
		return err
	}

	if _, err := f.Write(files); err != nil {
		return err
	}

	return nil
}
