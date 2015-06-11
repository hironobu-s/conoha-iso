package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func main() {
	app := NewConoHaIso()
	if err := app.Run(os.Args); err != nil {
		log.Errorf(err.Error())
	}
}
