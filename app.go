package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/hironobu-s/conoha-iso/command"
	"io/ioutil"
)

type ConoHaIso struct {
	lastError error

	*cli.App
}

func NewConoHaIso() *ConoHaIso {
	app := &ConoHaIso{
		App: cli.NewApp(),
	}
	app.setup()
	return app
}

func (app *ConoHaIso) version() string {
	// os.Getwd()
	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		return ""
	} else {
		println(version)
		return string(version)
	}
}

func (app *ConoHaIso) setup() {
	app.Name = "conoha-iso"
	app.Usage = "This app allow you to manage ISO images on ConoHa."
	app.Version = app.version()

	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "api-username, u",
			Value:  "",
			Usage:  "API Username",
			EnvVar: "CONOHA_USERNAME",
		},
		cli.StringFlag{
			Name:   "api-password, p",
			Value:  "",
			Usage:  "API Password",
			EnvVar: "CONOHA_PASSWORD",
		},
		cli.StringFlag{
			Name:   "api-tenant-id, t",
			Value:  "",
			Usage:  "API TenantId",
			EnvVar: "CONOHA_TENANT_ID",
		},
		cli.StringFlag{
			Name:   "region, r",
			Value:  "",
			Usage:  "Region name that ISO image will be uploaded. Allowed values are tyo1, sin1 or sjc1.",
			EnvVar: "CONOHA_REGION",
		},
	}

	app.Commands = []cli.Command{
		app.list(flags),
		app.download(flags),
		app.insert(flags),
		app.eject(flags),
	}
}

func (app *ConoHaIso) afterAction(context *cli.Context) error {
	return app.lastError
}

func (app *ConoHaIso) list(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:  "list",
		Usage: "List ISO Images.",
		Flags: flags,
		After: app.afterAction,
		Action: func(c *cli.Context) {
			ident, err := app.auth(c)
			if err != nil {
				app.lastError = err
				return
			}

			var compute *command.Compute
			compute = command.NewCompute(ident)

			isos, err := compute.List()
			if err != nil {
				app.lastError = err
				return
			}

			for i, iso := range isos.IsoImages {
				fmt.Printf("[Image%d]\n", i+1)
				fmt.Printf("%-6s %s\n", "Name:", iso.Name)
				fmt.Printf("%-6s %s\n", "Url:", iso.Url)
				fmt.Printf("%-6s %s\n", "Path:", iso.Path)
				fmt.Printf("%-6s %s\n", "Ctime:", iso.Ctime)
				fmt.Printf("%-6s %d\n", "Size:", iso.Size)

				if i != len(isos.IsoImages)-1 {
					println()
				}
			}
			if len(isos.IsoImages) == 0 {
				println("No ISO images.")
			}
		},
	}
	return cmd
}

func (app *ConoHaIso) insert(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:  "insert",
		Usage: "Insert an ISO images to the VPS.",
		Flags: flags,
		After: app.afterAction,
		Action: func(c *cli.Context) {
			ident, err := app.auth(c)
			if err != nil {
				app.lastError = err
				return
			}

			var compute *command.Compute
			compute = command.NewCompute(ident)

			err = compute.Insert()
			if err != nil {
				app.lastError = err
				return
			}
			log.Info("ISO file was inserted and changed boot device.")
		},
	}
	return cmd
}

func (app *ConoHaIso) eject(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:  "eject",
		Usage: "Eject an ISO image from the VPS.",
		Flags: flags,
		After: app.afterAction,
		Action: func(c *cli.Context) {
			ident, err := app.auth(c)
			if err != nil {
				app.lastError = err
				return
			}

			var compute *command.Compute
			compute = command.NewCompute(ident)

			err = compute.Eject()
			if err != nil {
				app.lastError = err
				return
			}
			log.Info("ISO file was ejected.")
		},
	}
	return cmd
}

func (app *ConoHaIso) download(flags []cli.Flag) cli.Command {

	flags = append(flags, cli.StringFlag{
		Name:  "url, i",
		Value: "",
		Usage: "ISO file url.",
	})

	cmd := cli.Command{
		Name:  "download",
		Usage: "Download ISO file from the FTP/HTTP server.",
		Flags: flags,
		After: app.afterAction,
		Before: func(c *cli.Context) error {
			if c.String("url") == "" {
				return errors.New("ISO file url required.")
			}
			return nil
		},
		Action: func(c *cli.Context) {

			ident, err := app.auth(c)
			if err != nil {
				app.lastError = err
				return
			}

			var compute *command.Compute
			compute = command.NewCompute(ident)

			if err = compute.Download(c.String("url")); err != nil {
				app.lastError = err
				return
			}

			log.Info("Download request was accepted.")
		},
	}
	return cmd
}

func (app *ConoHaIso) auth(c *cli.Context) (*command.Identity, error) {

	ident := command.NewIdentity()

	requires := map[string]*string{
		"api-username":  &ident.ApiUsername,
		"api-password":  &ident.ApiPassword,
		"api-tenant-id": &ident.ApiTenantId,
	}

	for name, v := range requires {
		if c.String(name) != "" {
			*v = c.String(name)
		} else {
			return nil, fmt.Errorf("Parameter \"%s\" is required.", name)
		}
	}

	if c.String("region") == "" {
		return nil, fmt.Errorf("Region shoud be required.")
	} else {
		ident.Region = c.String("region")
	}

	if err := ident.Auth(); err != nil {
		return nil, err
	}
	return ident, nil
}
