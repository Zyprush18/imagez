package pkg

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Zyprush18/imagez/utils"
	"github.com/h2non/bimg"
)

type JobChannels struct {
	worker                                   int
	nameFileZip                              chan string
	ImgOriEntity                             chan utils.ImageOri
	ZipEntity                                chan *ZipEntity
	errs                                     chan<- error
	format                                   string
	wgConvert, wgResize, wgToZip, wgCompress sync.WaitGroup
}

type ZipEntity struct {
	name  string
	image []byte
}

func NewJobChannel(worker int, ImgOriEntity chan utils.ImageOri, nameFileZip chan string, errs chan<- error, format string) *JobChannels {
	return &JobChannels{
		worker:       worker,
		ImgOriEntity: ImgOriEntity,
		nameFileZip:  nameFileZip,
		ZipEntity:    make(chan *ZipEntity),
		errs:         errs,
		format:       format,
		wgConvert:    sync.WaitGroup{},
		wgResize:     sync.WaitGroup{},
		wgToZip:      sync.WaitGroup{},
		wgCompress:   sync.WaitGroup{},
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

func (j *JobChannels) ProcessSaveZip(nameZip string, zipFile *zip.Writer) {
	defer j.wgToZip.Done()
	for v := range j.ZipEntity {
		f, err := zipFile.Create(v.name)
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

func (j *JobChannels) ResizeJob(width, height int) {
	nameZip := fmt.Sprintf("./img/%s.zip", rand.Text())
	file, zipFile, err := j.createZip(nameZip)
	if err != nil {
		j.errs <- err
		close(j.errs)
		return
	}

	for i := 0; i < j.worker; i++ {
		j.wgResize.Add(1)
		go j.ProcessResize(width, height, i)
	}

	j.wgToZip.Add(1)
	go j.ProcessSaveZip(nameZip, zipFile)

	go func() {
		j.wgResize.Wait()
		close(j.ZipEntity)
		j.wgToZip.Wait()

		j.nameFileZip <- nameZip

		zipFile.Close()
		file.Close()
		close(j.nameFileZip)
		close(j.errs)
	}()
}

func (j *JobChannels) CompressJob(size int) {
	nameZip := fmt.Sprintf("./img/%s.zip", rand.Text())
	file, zipFile, err := j.createZip(nameZip)
	if err != nil {
		j.errs <- err
		close(j.errs)
		return
	}

	for i := 0; i < j.worker; i++ {
		j.wgCompress.Add(1)
		go j.ProcessCompress(size)
	}

	j.wgToZip.Add(1)
	go j.ProcessSaveZip(nameZip, zipFile)

	go func() {
		j.wgCompress.Wait()
		close(j.ZipEntity)
		j.wgToZip.Wait()

		j.nameFileZip <- nameZip

		zipFile.Close()
		file.Close()
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
		j.ZipEntity <- &ZipEntity{
			name:  fmt.Sprintf("%s-%s.%s", v.Name, rand.Text(), j.format),
			image: newImg,
		}
	}

}

func (j *JobChannels) ProcessResize(width, height, worker int) {
	defer j.wgResize.Done()
	for v := range j.ImgOriEntity {
		format := strings.Trim(filepath.Ext(v.Name), ".")
		switch format {
		case "jpg", "jpeg", "png", "webp", "gif", "avif", "tiff":
			newImg, err := bimg.NewImage(v.Image).Resize(width, height)
			if err != nil {
				j.errs <- err
				continue
			}
			j.ZipEntity <- &ZipEntity{
				name:  fmt.Sprintf("%s-%s.%s", strings.TrimSuffix(v.Name, filepath.Ext(v.Name)), rand.Text(), format),
				image: newImg,
			}
		default:
			j.errs <- fmt.Errorf(utils.UNSUPPORTED_FORMAT)
			continue
		}

	}
}

func (j *JobChannels) ProcessCompress(size int) {
	defer j.wgCompress.Done()
	for v := range j.ImgOriEntity {
		var newImg []byte
		var err error
		quality := 80
		extFile := strings.Trim(filepath.Ext(v.Name), ".")

		for i := 0; i < 8; i++ {
			if quality > 100 {
				quality = 100
			} else if quality < 0 {
				quality = 1
			}

			newImg, err = utils.ProcessBimg(v.Image, quality)
			if err != nil {
				j.errs <- err
				continue
			}

			if len(newImg) >= size-(1024 * 100) && len(newImg) <= size+(1024 * 100) {
				fmt.Println("masukkk")
				break
			}

			if len(newImg) > size {
				quality -= 9
			} else if len(newImg) < size {
				quality += 9
			}
		}

		j.ZipEntity <- &ZipEntity{
			name:  fmt.Sprintf("%s-%s.%s", strings.TrimSuffix(v.Name, filepath.Ext(v.Name)), rand.Text(), extFile),
			image: newImg,
		}
	}
}
