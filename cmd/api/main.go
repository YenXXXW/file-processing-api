package main

import (
	"log"

	"github.com/yenxxxw/image-processing-api/internal/modifiedLogger"
)

func main() {

	config := config{
		addr: ":8080",
	}

	sugar := modifiedLogger.InitLogger()

	app := &application{
		config: config,
		logger: sugar,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
