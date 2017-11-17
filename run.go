package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//go:generate go-assets-builder -s /webui/ -p webui -o webui/assets.go webui/template webui/assets

func main() {
	app := NewConoHaIso()
	if err := app.Run(os.Args); err != nil {
		log.Errorf(err.Error())
	}
}
