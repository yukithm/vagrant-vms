package main

import "github.com/codegangsta/cli"

const (
	// ApplicationName is the name of this application.
	ApplicationName = "vagrant-vms"

	// Version is the version number of this application.
	Version = "0.1.0"
)

func main() {
	app := cli.NewApp()
	app.Name = ApplicationName
	app.Version = Version
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "vms",
			Usage:  "show provider's VM list (currently supported VirtualBox only)",
			Action: vmsCommand,
		},
		cli.Command{
			Name:  "boxes",
			Usage: "show vagrant boxes",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "vertical",
					Usage: "print vertically",
				},
			},
			Action: boxesCommand,
		},
		cli.Command{
			Name:  "status",
			Usage: "show vagrant global status",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "vertical",
					Usage: "print vertically",
				},
			},
			Action: statusCommand,
		},
	}
	app.RunAndExitOnError()
}
