package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Image struct {
	Name string `json:"name"`
	Mime string `json:"mime"`
	Data []byte `json:"data"`
}

type WorkerPool struct {
	Images        []Image
	concurrency   int
	imageSendChan chan Image
	imageRecChan  chan Image

	wg sync.WaitGroup
}

func (wp *WorkerPool) worker() {
	for image := range wp.imageSendChan {
		newImg, err := image.ConvertJPG()
		if err == nil {
			wp.imageRecChan <- newImg
		}
		wp.wg.Done()
	}
}

func (wp *WorkerPool) run(w http.ResponseWriter) {
	wp.imageSendChan = make(chan Image, len(wp.Images))
	wp.imageRecChan = make(chan Image, len(wp.Images))

	for i := 0; i < wp.concurrency; i++ {
		go wp.worker()
	}

	wp.wg.Add(len(wp.Images))
	for _, image := range wp.Images {
		wp.imageSendChan <- image
	}

	close(wp.imageSendChan)
	wp.wg.Wait()
	close(wp.imageRecChan)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"images.zip\"")

	zipWriter := zip.NewWriter(w)
	for result := range wp.imageRecChan {
name := strings.TrimSuffix(, ) + ".png"

	}

}

func (image *Image) ConvertJPG() (Image, error) {
	data, decodeErr := jpeg.Decode(bytes.NewReader(image.Data))
	if decodeErr != nil {
		return Image{}, decodeErr

	}

	var buff bytes.Buffer

	encodeErr := png.Encode(&buff, data)
	if encodeErr != nil {
		return Image{}, encodeErr
	}

	newImage := Image{
		Name: image.Name,
		Mime: image.Mime,
		Data: buff.Bytes(),
	}

	return newImage, nil

}

func (app *application) extractImage(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseMultipartForm(10 * 1024 * 1024)
	images := []Image{}
	if parseErr != nil {
		app.internaleServerError(w, r, parseErr)
	}

	nImages := r.PostForm.Get("nImages")

	app.logger.Infof("nImages %s", nImages)

	n, err := strconv.Atoi(nImages)
	if err != nil {
		app.badRequesetError(w, r, err)
		return
	}

	for i := 0; i < n; i++ {

		key := fmt.Sprintf("image%d", i)

		file, header, err := r.FormFile(key)
		if err != nil {
			app.badRequesetError(w, r, err)
			return
		}

		fileName := header.Filename
		fileMime := header.Header.Get("Content-Type")

		defer file.Close()

		data, fileReadErr := io.ReadAll(file)

		if fileReadErr != nil {
			app.badRequesetError(w, r, fileReadErr)
			return
		}

		image := Image{
			Name: fileName,
			Mime: fileMime,
			Data: data,
		}

		images = append(images, image)

		wp := WorkerPool{
			Images:      images,
			concurrency: 5,
		}

		wp.run(w)
		app.logger.Infof("Received %s (%d bytes)", key, len(data))

	}

	if err := app.jsonResponse(w, http.StatusOK, "images extracted successfully"); err != nil {
		app.internaleServerError(w, r, err)
	}

}
