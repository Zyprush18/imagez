package service

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Zyprush18/imagez/pkg"
	"github.com/Zyprush18/imagez/utils"
)

var worker int = runtime.NumCPU()

type ImagesService interface {
	Convert(data []*multipart.FileHeader, extFormat string) (string, error)
	Resize(data []*multipart.FileHeader, width, height int) (string, error)
}

type ImageService struct{}

func NewServiceImage() ImagesService {
	return &ImageService{}
}

func (i *ImageService) Convert(data []*multipart.FileHeader, extFormat string) (string, error) {
	FileData := make(chan utils.ImageOri, len(data))
	errs := make(chan error, worker)
	nameFileZip := make(chan string, 1)
	for _, v := range data {
		typeFile := v.Header.Get("Content-Type")
		if err := utils.CheckType(typeFile); err != nil {
			return "", err
		}

		src, err := v.Open()
		if err != nil {
			return "", err
		}

		file, err := io.ReadAll(src)
		if err != nil {
			return "", err
		}

		src.Close()

		FileData <- utils.ImageOri{
			Name:  strings.TrimSuffix(v.Filename, filepath.Ext(v.Filename)),
			Image: file,
		}

	}

	close(FileData)

	convert := pkg.NewJobChannel(worker, FileData, nameFileZip, errs, extFormat)
	convert.ConvertJob()

	for v := range errs {
		if v != nil {
			return "", v
		}
	}

	fileName := <-nameFileZip

	return fileName, nil
}

func (i *ImageService) Resize(data []*multipart.FileHeader, width, height int) (string, error) {
	ImageData := make(chan utils.ImageOri, len(data))
	zipFileName := make(chan string, 1)
	errs := make(chan error, worker)

	for _, v := range data {
		typeFile := v.Header.Get("Content-Type")
		if err := utils.CheckType(typeFile); err != nil {
			return "", err
		}

		src, err := v.Open()
		if err != nil {
			return "", err
		}

		file, err := io.ReadAll(src)
		if err != nil {
			return "", err
		}

		src.Close()

		ImageData <- utils.ImageOri{
			Name:  v.Filename,
			Image: file,
		}
	}

	close(ImageData)

	resize := pkg.NewJobChannel(worker, ImageData, zipFileName, errs, "")
	resize.ResizeJob(width, height)

	for v := range errs {
		if v != nil {
			return "", v
		}
	}

	fileName := <-zipFileName

	return fileName, nil
}
