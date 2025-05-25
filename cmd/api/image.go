package main

import (
	"net/http"
)

func (app *application) extractImage(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseMultipartForm(10 * 1024 * 1024)
	if parseErr != nil {
		app.internaleServerError(w, r, parseErr)
	}

	nImages := r.PostForm.Get("nImages")

	app.logger.Infof("nImages %s", nImages)

	if err := app.jsonResponse(w, http.StatusOK, "image exreacted successfully"); err != nil {

		app.internaleServerError(w, r, err)
	}

}
