package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Stealer"
	app.Usage = "Steal variable and its values and write it to GO"
	app.Version = "1.0.0"

	// flags for option command
	var location string
	var destination string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file",
			Value:       "files/",
			Usage:       "specify files location or directory",
			Destination: &location,
		},
		cli.StringFlag{
			Name:        "folder",
			Value:       "",
			Usage:       "destination to save all the converted files",
			Destination: &destination,
		},
	}

	// default action
	app.Action = func(c *cli.Context) error {
		if location == "" && destination == "" {
			fmt.Println("please specify location and destination")
			return nil
		}

		fileConverter(location, destination)
		return nil
	}

	// Commad to execute
	app.Commands = []cli.Command{

		//first command
		{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "Show load result",
			Action: func(c *cli.Context) error {

				return nil
			},
		},
	}

	app.Run(os.Args)
}
