package utils

import "errors"

const (
	UNSUPPORTED_FORMAT string = "unsupported format"
	UNSUPPORTED_TYPE   string = "unsupported type"
)

func CheckType(typefile string) error {
	typeImage := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/tiff",
		"image/webp",
		"image/avif",
		"image/svg+xml",
		"image/x-icon",
	}
	for _, v := range typeImage {
		if typefile == v {
			return nil
		}
	}
	return errors.New(UNSUPPORTED_TYPE)
}
