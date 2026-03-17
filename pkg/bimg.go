package pkg

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"os"
	"sync"

	"github.com/Zyprush18/imagez/utils"
	"github.com/h2non/bimg"
)

type JobChannels struct {
	worker      int
	nameFileZip chan string
	ImgOriEntity chan utils.ImageOri
	ZipEntity   chan ZipEntity
	errs        chan<- error
	format      string
	wgConvert   sync.WaitGroup
	wgToZip     sync.WaitGroup
	MxToZip     sync.Mutex
}

type ZipEntity struct {
	name string
	image []byte
}

func NewJobChannel(worker int, ImgOriEntity chan utils.ImageOri, nameFileZip chan string, errs chan<- error, format string) *JobChannels {
	return &JobChannels{
		worker:      worker,
		ImgOriEntity: ImgOriEntity,
		nameFileZip: nameFileZip,
		ZipEntity:   make(chan ZipEntity),
		errs:        errs,
		format:      format,
		wgConvert:   sync.WaitGroup{},
		wgToZip:     sync.WaitGroup{},
		MxToZip:     sync.Mutex{},
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

func (j *JobChannels) createZip(NameZip string) (file *os.File, zipWriter *zip.Writer, err error) {
	file, err = os.Create(NameZip)
	if err != nil {
		return nil, nil, err
	}
	return file, zip.NewWriter(file), nil
}

func (j *JobChannels) ConvertJob() {
	ext := j.Extension()

	if ext == bimg.UNKNOWN {
		j.errs <- fmt.Errorf(utils.UNSUPPORTED_FORMAT)
		close(j.errs)
		return
	}
	nameZip := fmt.Sprintf("./img/%s.zip", rand.Text())
	file, zipFile, err := j.createZip(nameZip)
	if err != nil {
		j.errs <- err
		close(j.errs)
		return
	}

	for i := 0; i < j.worker; i++ {
		j.wgConvert.Add(1)
		go j.ProcessConvert(ext)
	}

	j.wgToZip.Add(1)
	go j.ProcessSaveZip(nameZip, zipFile)

	go func() {
		j.wgConvert.Wait()
		close(j.ZipEntity)
		j.wgToZip.Wait()

		zipFile.Close()
		file.Close()

		j.nameFileZip <- nameZip
		close(j.nameFileZip)
		close(j.errs)
	}()
}

func (j *JobChannels) ProcessConvert(extension bimg.ImageType) {
	defer j.wgConvert.Done()
	for v := range j.ImgOriEntity {
		newImg, err := bimg.NewImage(v.Image).Convert(extension)
		if err != nil {
			j.errs <- err
			continue
		}
		j.ZipEntity <- ZipEntity{
			name: fmt.Sprintf("%s-%s", v.Name, rand.Text()),
			image: newImg,
		}
	}

}

func (j *JobChannels) ProcessSaveZip(nameZip string, zipFile *zip.Writer) {
	defer j.wgToZip.Done()
	for v := range j.ZipEntity {
		nameFile := fmt.Sprintf("imagez-%s.%s", v.name, j.format)
		f, err := zipFile.Create(nameFile)
		if err != nil {
			j.errs <- err
			continue
		}

		if _, err := f.Write(v.image); err != nil {
			j.errs <- err
			continue
		}
	}
}
