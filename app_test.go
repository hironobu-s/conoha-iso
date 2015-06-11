package main

import (
	"flag"
	"github.com/k0kubun/pp"
	"testing"
)

func TestNoRegion(t *testing.T) {
	app := NewConoHaIso()

	f := flag.NewFlagSet("test", 0)

	for _, command := range []string{"list", "download", "eject", "insert"} {
		test := []string{"conoha-iso", command}
		f.Parse(test)

		err := app.Run(f.Args())
		if err == nil {
			pp.Printf("err:%v\n", err)
			t.Errorf("No region specified. Test should be fail in this case.")
		}
	}
}

func TestNoDownloadUrl(t *testing.T) {
	app := NewConoHaIso()

	f := flag.NewFlagSet("test", 0)
	test := []string{"conoha-iso", "download"}
	f.Parse(test)

	err := app.Run(f.Args())
	if err == nil {
		pp.Printf("err:%v\n", err)
		t.Errorf("No download url specified. Test should be fail in this case.")
	}
}
