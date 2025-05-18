package main

import "net/http"

func (app *application) healthCheckHanlder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("okay! working fine"))

}
