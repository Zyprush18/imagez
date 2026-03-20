package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
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
	Compress(data []*multipart.FileHeader,size int) (string, error)
	Crop(data []*multipart.FileHeader, width, height int) (string, error)
	Watermark(data []*multipart.FileHeader, textWatermark string) (string, error)
	Downloads(w io.Writer,filename string) error
	DeleteFileZip(filename string) error
}

type ImageService struct{}

func NewServiceImage() ImagesService {
	return &ImageService{}
}

func (i *ImageService) getNameFile(filename string) string  {
	return strings.Trim(filename, "./img/")
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

	return i.getNameFile(fileName), nil
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

	return i.getNameFile(fileName), nil
}

func (i *ImageService) Compress(data []*multipart.FileHeader,  size int) (string, error) {
	imgFile := make(chan utils.ImageOri, len(data))
	FileNameZip := make(chan string, 1)
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

		imgFile <- utils.ImageOri{
			Name:  v.Filename,
			Image: file,
		}
	}

	close(imgFile)

	compress := pkg.NewJobChannel(worker, imgFile, FileNameZip, errs, "")
	compress.CompressJob(size)

	for v := range errs {
		if v != nil {
			return "", v
		}
	}

	fileName := <-FileNameZip

	return i.getNameFile(fileName), nil
}



func (i *ImageService) Crop(data []*multipart.FileHeader, width, height int) (string, error) {
	imgFile := make(chan utils.ImageOri, len(data))
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

		imgFile <- utils.ImageOri{
			Name:  v.Filename,
			Image: file,
		}
	}

	close(imgFile)

	crop := pkg.NewJobChannel(worker, imgFile, zipFileName, errs, "")
	crop.CropJob(width, height)

	for v := range errs {
		if v != nil {
			return "", v
		}
	}

	fileName := <-zipFileName

	return i.getNameFile(fileName), nil
}

func (i *ImageService) Watermark(data []*multipart.FileHeader, textWatermark string) (string, error) {
	imgFile := make(chan utils.ImageOri, len(data))
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

		imgFile <- utils.ImageOri{
			Name:  v.Filename,
			Image: file,
		}
	}

	close(imgFile)

	watermark := pkg.NewJobChannel(worker, imgFile, zipFileName, errs, "")
	watermark.WatermarkJob(textWatermark)

	for v := range errs {
		if v != nil {
			return "", v
		}
	}

	fileName := <-zipFileName

	return i.getNameFile(fileName), nil
}


func (i *ImageService) Downloads(w io.Writer,filename string) error {
	pathfile := fmt.Sprintf("./img/%s", filename)
	file, err := os.Open(pathfile)
	if err != nil {
		return  err
	}

	defer file.Close()

	io.Copy(w, file)

	return nil
}


func (i *ImageService) DeleteFileZip(filename string) error {
	pathfile := fmt.Sprintf("./img/%s", filename)
	if err:= os.Remove(pathfile);err != nil {
		return  err
	}

	return nil
}