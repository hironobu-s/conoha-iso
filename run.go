package main

import (
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

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "api-username, u",
			Value: "",
			Usage: "API Username",
		},
		cli.StringFlag{
			Name:  "api-password, p",
			Value: "",
			Usage: "API Password",
		},
		cli.StringFlag{
			Name:  "api-tenant-id, t",
			Value: "",
			Usage: "API TenantId",
		},
		cli.StringFlag{
			Name:  "region, r",
			Value: "tyo1",
			Usage: "Region name that ISO image will be uploaded.",
		},
	}

	app.Commands = []cli.Command{
		list(flags),
		download(flags),
	}

}

func list(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List ISO Images.",
		Flags:   flags,
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
		},
	}
	return cmd
}

func download(flags []cli.Flag) cli.Command {
	cmd := cli.Command{
		Name:    "download",
		Aliases: []string{"l"},
		Usage:   "Download ISO Image from the specific server.",
		Flags:   flags,
		Action: func(c *cli.Context) {
			if len(c.Args()) == 0 {
				return
			}

			url := c.Args()[0]

			ident, err := auth(c)
			if err != nil {
				log.Errorf("%s", err)
				return
			}

			var compute *command.Compute
			compute = command.NewCompute()
			compute.Identity = ident

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
