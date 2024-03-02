package main

import (
	"localdothost-dns/app"
	"log"
	"os"
)

func main() {
	app := app.Init()
	err := app.Run(os.Args)

	if(err != nil) {
		log.Fatal(err)
	}
}