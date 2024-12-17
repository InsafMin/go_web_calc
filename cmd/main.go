package main

import (
	"github.com/InsafMin/go_web_calc/internal/application"
	"log"
)

func main() {
	app := application.New()
	//err := app.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	err := app.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
