package main

import (
	"log"

	"github.com/EmptyInsid/db_gui/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Error while run app: %v\n", err)
	}
}
