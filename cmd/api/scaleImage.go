package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/jpeg"
	"net/http"
	"strconv"

	"github.com/nfnt/resize"
)

func (app *application) scaleImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		app.badRequesetError(w, r, err)
		return
	}

	width := r.PostForm.Get("width")
	height := r.PostForm.Get("height")
	wd, err := strconv.Atoi(width)
	if err != nil {
		app.badRequesetError(w, r, errors.New("width must be a number"))
		return
	}

	ht, err := strconv.Atoi(height)
	if err != nil {
		app.badRequesetError(w, r, errors.New("width must be a number"))
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		app.badRequesetError(w, r, err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		app.internaleServerError(w, r, err)
	}
	file.Close()

	uwd := uint(wd)
	uht := uint(ht)
	m := resize.Resize(uwd, uht, img, resize.Lanczos3)

	var out bytes.Buffer

	err = jpeg.Encode(&out, m, nil)
	if err != nil {
		app.internaleServerError(w, r, err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(out.Bytes())
	if err := app.jsonResponse(w, http.StatusOK, imageBase64); err != nil {
		app.internaleServerError(w, r, err)
	}

}
