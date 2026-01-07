package main

import (
	"log"
	"suscord_ws/internal/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	if err = app.Run(); err != nil {
		log.Printf("%+v\n", err)
	}

	if err := app.Stop(); err != nil {
		log.Printf("%+v\n", err)
	}
}
