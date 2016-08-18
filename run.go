package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {
	app := NewConoHaIso()
	if err := app.Run(os.Args); err != nil {
		log.Errorf(err.Error())
	}
}
