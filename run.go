package main

import (
	"os"
)

func main() {
	app := NewConoHaIso()
	app.Run(os.Args)
}
