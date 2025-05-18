package main

import (
	"net/http"
)

func (app *application) internaleServerError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	wirteJsonError(w, http.StatusInternalServerError, "internal server error")

}
