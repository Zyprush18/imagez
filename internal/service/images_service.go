package service

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Zyprush18/imagez/internal/utils"
	"github.com/Zyprush18/imagez/pkg"
)


var worker int = runtime.NumCPU()

type ImagesService interface {
	Convert (data []*multipart.FileHeader, extFormat string) error
}

type ImageService struct {}

func NewServiceImage () ImagesService {
	return &ImageService{}
}


func (i *ImageService) Convert(data []*multipart.FileHeader, extFormat string) error {
	name := make(chan string, len(data))
	img := make(chan []byte, len(data))
	errs := make(chan error, len(data))
	for _, v := range data {
		typeFile := v.Header.Get("Content-Type")
		if err := utils.CheckType(typeFile); err != nil {
			return err
		}

		src, err := v.Open()
		if err != nil {
			return  err
		}

		file, err := io.ReadAll(src)
		if err != nil {
			return err
		}

		src.Close()

		filename :=  strings.TrimSuffix(v.Filename, filepath.Ext(v.Filename))

		name <- filename
		img <- file
	}

	close(name)
	close(img)

	convert := pkg.NewJobChannel(worker, img, name, errs, extFormat)
	convert.ConvertJob()


	for v := range errs {
		if v != nil {
			return v
		}
	}


	return nil
}