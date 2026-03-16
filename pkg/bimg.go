package pkg

import (
	"fmt"
	"sync"

	"github.com/Zyprush18/imagez/internal/utils"
	"github.com/h2non/bimg"
)

type JobChannels struct {
	worker int
	name   <-chan string
	image  <-chan []byte
	errs   chan<- error
	format string
	wgConvert sync.WaitGroup
}

func NewJobChannel(worker int, image <-chan []byte, name <-chan string, errs chan<- error, format string) *JobChannels {
	return &JobChannels{
		worker: worker,
		image:  image,
		name:   name,
		errs:   errs,
		format: format,
		wgConvert: sync.WaitGroup{},
	}
}

func (j *JobChannels) Extension() bimg.ImageType {
	switch j.format {
	case "jpg", "jpeg":
		return bimg.JPEG
	case "png":
		return bimg.PNG
	case "webp":
		return bimg.WEBP
	case "gif":
		return bimg.GIF
	case "avif":
		return bimg.AVIF
	case "tiff":
		return bimg.TIFF
	default:
		return bimg.UNKNOWN
	}
}


func (j *JobChannels) SaveImage(name string, img []byte) error {
	return bimg.Write("./img/"+name, img)
}

func (j *JobChannels) ConvertJob() {

	ext := j.Extension()
	if ext == bimg.UNKNOWN {
		j.errs <- fmt.Errorf(utils.UNSUPPORTED_FORMAT)
		return
	}
	for i := 0; i < j.worker; i++ {
		j.wgConvert.Add(1)
		go j.ProcessConvert(ext)
	}

	go func() {
		j.wgConvert.Wait()
		close(j.errs)
	}()
}

func (j *JobChannels) ProcessConvert(extension bimg.ImageType) {
	defer j.wgConvert.Done()

	for v := range j.image {
		newImg, err := bimg.NewImage(v).Convert(extension)
		if err != nil {
			j.errs <- err
			continue
		}
		nameFile := fmt.Sprintf("%s.%s", <-j.name, j.format)

		if err := j.SaveImage(nameFile, newImg); err != nil {
			j.errs <- err
			continue
		}
	}
}

