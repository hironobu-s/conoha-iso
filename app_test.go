package main

import (
	"flag"
	"fmt"
	"testing"
)

func TestNoRegion(t *testing.T) {
	app := NewConoHaIso()

	f := flag.NewFlagSet("test", 0)

	for _, command := range []string{"list"} {
		test := []string{"conoha-iso", command}
		f.Parse(test)

		err := app.Run(f.Args())
		if err == nil {
			fmt.Printf("err: %v\n", err)
		}
	}
}
