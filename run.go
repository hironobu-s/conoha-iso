package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/hironobu-s/conoha-iso/command"
	"os"
)

func main() {
	app := cli.NewApp()
	setup(app)
	app.Run(os.Args)
}

func setup(app *cli.App) {
	app.Name = "conoha-iso"
	app.Usage = "This app allow you to manage ISO images on ConoHa."
	app.Version = "0.1"

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
		list(flags),
		download(flags),
		insert(flags),
		eject(flags),
	}
}

func list(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:  "list",
		Usage: "List ISO Images.",
		Flags: flags,

		Action: func(c *cli.Context) {
			ident, err := auth(c)
			if err != nil {
				log.Errorf("%s", err)
				return
			}

			var compute *command.Compute
			compute = command.NewCompute()
			compute.Identity = ident

			isos, err := compute.List()
			if err != nil {
				log.Errorf("%s", err)
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

func insert(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:  "insert",
		Usage: "Insert an ISO images to the VPS.",
		Flags: flags,
		Action: func(c *cli.Context) {
			ident, err := auth(c)
			if err != nil {
				log.Errorf("%s", err)
				return
			}

			var compute *command.Compute
			compute = command.NewCompute()
			compute.Identity = ident

			err = compute.Insert()
			if err != nil {
				log.Errorf("%s", err)
				return
			}
			log.Info("ISO file was inserted and changed boot device.")
		},
	}
	return cmd
}

func eject(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:  "eject",
		Usage: "Eject an ISO image from the VPS.",
		Flags: flags,
		Action: func(c *cli.Context) {
			ident, err := auth(c)
			if err != nil {
				log.Errorf("%s", err)
				return
			}

			var compute *command.Compute
			compute = command.NewCompute()
			compute.Identity = ident

			err = compute.Eject()
			if err != nil {
				log.Errorf("%s", err)
				return
			}
			log.Info("ISO file was ejected.")
		},
	}
	return cmd
}

func download(flags []cli.Flag) cli.Command {

	flags = append(flags, cli.StringFlag{
		Name:  "url, i",
		Value: "",
		Usage: "ISO image url.",
	})

	cmd := cli.Command{
		Name:  "download",
		Usage: "Download ISO image from the FTP/HTTP server.",
		Flags: flags,
		Before: func(c *cli.Context) error {
			if c.String("url") == "" {
				log.Errorf("%s", "ISO image url required.")
				return errors.New("ISO image url required.")
			}
			return nil
		},
		Action: func(c *cli.Context) {

			ident, err := auth(c)
			if err != nil {
				log.Errorf("%s", err)
				return
			}

			var compute *command.Compute
			compute = command.NewCompute()
			compute.Identity = ident

			url := c.String("url")

			if err = compute.Download(url); err != nil {
				log.Errorf("%s", err)
				return
			}

			log.Info("A download request was accepted.")
		},
	}
	return cmd
}

func auth(c *cli.Context) (*command.Identity, error) {

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
	}

	switch c.String("region") {
	case "tyo1":
		ident.Region = command.TYO1
	case "sin1":
		ident.Region = command.SIN1
	case "sjc1":
		ident.Region = command.SJC1
	default:
		return nil, fmt.Errorf("Undefined region \"%s\"", c.String("region"))
	}

	if err := ident.Auth(); err != nil {
		return nil, err
	}
	return ident, nil
}
