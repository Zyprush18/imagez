package utils

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	UNSUPPORTED_FORMAT string = "unsupported format"
	UNSUPPORTED_TYPE   string = "unsupported type"
	INVALID_CROP_PARAMETERS string = "invalid crop parameters"
)

func CheckErrFileIsExist(filename string, err error) bool {
	if err.Error() == fmt.Sprintf("open ./img/%s: no such file or directory", filename) {
		return  true
	}

	return false
}


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

func ConvertToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}