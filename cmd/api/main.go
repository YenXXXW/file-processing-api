package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {

	config := config{
		addr: ":8080",
	}

	logger, _ := zap.NewProduction()
	sugarLogger := logger.Sugar()

	app := &application{
		config: config,
		logger: sugarLogger,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
