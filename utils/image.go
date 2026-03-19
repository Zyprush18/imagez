package utils

import "github.com/h2non/bimg"

type ImageOri struct {
	Name string
	Image []byte
}

func ProcessBimg(img []byte, quality int) ([]byte, error)  {
	return  bimg.NewImage(img).Process(bimg.Options{Quality: quality})
}