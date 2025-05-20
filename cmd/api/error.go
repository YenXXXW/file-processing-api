package main

import (
	"net/http"
)

func (app *application) internaleServerError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusInternalServerError, "internal server error")

}

func (app *application) badRequesetError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusBadRequest, err.Error())
}
