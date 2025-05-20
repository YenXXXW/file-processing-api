package main

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type User struct {
	Age          int `json:"age"`
	Fnlwgt       int `json:"fmlwgt"`
	EducationNum int `json:"education_num"`
	CapitalGain  int `json:"capital_gain"`
	CapitalLoss  int `json:"capital_lost"`
	HoursPerWeek int `json:"hours_per_week"`
	IncomeLevel  int `json:"income_level"`
}

func (app *application) readCsv(w http.ResponseWriter, r *http.Request) {

	reader := csv.NewReader(r.Body)
	defer r.Body.Close()

	_, err := reader.Read()
	if err != nil {
		app.internaleServerError(w, r, err)
		return
	}

	for {

		var age, fnlwgt, education_num, captail_gain, captail_loss, hours_per_week, income_level int
		attributes := []int{age, fnlwgt, education_num, captail_gain, captail_loss, hours_per_week, income_level}
		record, err := reader.Read()

		for i := range record {
			attributes[i], err = strconv.Atoi(record[i])

			if err != nil {
				app.badRequesetError(w, r, errors.New("all values must be  number"))
				return
			}

		}

		if err == io.EOF {
			break
		}

		if err != nil {
			app.internaleServerError(w, r, err)
			return
		}

		app.logger.Info(record)
	}

	if err := app.jsonResponse(w, http.StatusOK, "csv readed successfully"); err != nil {

		app.internaleServerError(w, r, err)
	}

}
