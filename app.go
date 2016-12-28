package main

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/hironobu-s/conoha-iso/command"
	"github.com/hironobu-s/conoha-iso/webui"
	"github.com/urfave/cli"
)

const (
	APP_VERSION            = "0.4"
	DEFAULT_LISTEN_ADDRESS = "127.0.0.1:6543"
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

func (app *ConoHaIso) setup() {
	app.Name = "conoha-iso"
	app.Usage = "This app allow you to manage the ISO images on ConoHa."
	app.Version = APP_VERSION // Version should be updated by hand at each release.

	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "api-username, u",
			Value:  "",
			Usage:  "API Username",
			EnvVar: "OS_USERNAME,CONOHA_USERNAME",
		},
		cli.StringFlag{
			Name:   "api-password, p",
			Value:  "",
			Usage:  "API Password",
			EnvVar: "OS_PASSWORD,CONOHA_PASSWORD",
		},
		cli.StringFlag{
			Name:   "api-tenant-id, t",
			Value:  "",
			Usage:  "API TenantId",
			EnvVar: "OS_TENANT_ID,CONOHA_TENANT_ID",
		},
		cli.StringFlag{
			Name:   "api-tenant-name, n",
			Value:  "",
			Usage:  "API TenantName",
			EnvVar: "OS_TENANT_NAME,CONOHA_TENANT_NAME",
		},
		cli.StringFlag{
			Name:   "region, r",
			Value:  "tyo1",
			Usage:  "Region name that ISO image will be uploaded. Allowed values are tyo1, sin1 or sjc1.",
			EnvVar: "OS_REGION_NAME,CONOHA_REGION",
		},
	}

	app.Commands = []cli.Command{
		app.list(flags),
		app.download(flags),
		app.insert(flags),
		app.eject(flags),
		app.server(flags),
	}
}

func (app *ConoHaIso) afterAction(context *cli.Context) error {
	return app.lastError
}

func (app *ConoHaIso) server(flags []cli.Flag) cli.Command {
	flags = append(flags, cli.StringFlag{
		Name:  "listen,l",
		Value: DEFAULT_LISTEN_ADDRESS,
		Usage: "Listen address, of the form <host:port>",
	})

	cmd := cli.Command{
		Name:  "server",
		Usage: "Run the server.",
		Flags: flags,
		After: app.afterAction,
		Action: func(c *cli.Context) error {
			ident, err := app.auth(c)
			if err != nil {
				return err
			}

			if err = webui.RunServer(c.String("listen"), ident); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func (app *ConoHaIso) list(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:  "list",
		Usage: "List ISO Images.",
		Flags: flags,
		After: app.afterAction,
		Action: func(c *cli.Context) error {
			ident, err := app.auth(c)
			if err != nil {
				app.lastError = err
				return nil
			}

			var compute *command.Compute
			compute = command.NewCompute(ident)

			isos, err := compute.Isos()
			if err != nil {
				app.lastError = err
				return nil
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
			return nil
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
		Action: func(c *cli.Context) error {
			ident, err := app.auth(c)
			if err != nil {
				app.lastError = err
				return nil
			}

			var compute *command.Compute
			compute = command.NewCompute(ident)

			err = compute.InsertIntractive()
			if err != nil {
				app.lastError = err
				return nil
			}
			log.Info("ISO file was inserted and changed boot device.")
			return nil
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
		Action: func(c *cli.Context) error {
			ident, err := app.auth(c)
			if err != nil {
				app.lastError = err
				return nil
			}

			var compute *command.Compute
			compute = command.NewCompute(ident)

			err = compute.EjectIntractive()
			if err != nil {
				app.lastError = err
				return nil
			}
			log.Info("ISO file was ejected.")
			return nil
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
				return errors.New("ISO file url is required.")
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			ident, err := app.auth(c)
			if err != nil {
				app.lastError = err
				return nil
			}

			var compute *command.Compute
			compute = command.NewCompute(ident)

			if err = compute.Download(c.String("url")); err != nil {
				app.lastError = err
				return nil
			}

			log.Info("Download request was accepted.")
			return nil
		},
	}
	return cmd
}

func (app *ConoHaIso) auth(c *cli.Context) (*command.Identity, error) {

	ident := command.NewIdentity()

	requires := map[string]*string{
		"api-username": &ident.ApiUsername,
		"api-password": &ident.ApiPassword,
	}

	for name, v := range requires {
		if c.String(name) != "" {
			*v = c.String(name)
		} else {
			return nil, fmt.Errorf("Parameter \"%s\" is required.", name)
		}
	}

	// TenantIdとTenantNameはどちらか必須
	if c.String("api-tenant-name") == "" && c.String("api-tenant-id") == "" {
		return nil, fmt.Errorf("Ethier \"api-tenant-id\" or \"api-tenant-name\" is required.")
	}

	ident.ApiTenantName = c.String("api-tenant-name")
	ident.ApiTenantId = c.String("api-tenant-id")

	if c.String("region") == "" {
		ident.Region = "tyo1"
	} else {
		ident.Region = c.String("region")
	}

	if err := ident.Auth(); err != nil {
		return nil, err
	}
	return ident, nil
}
