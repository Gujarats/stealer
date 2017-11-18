package main

import (
	"fmt"
	"os"

	"github.com/Gujarats/logger"
	"github.com/Gujarats/stealer"
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
			Name:        "location",
			Value:       "",
			Usage:       "specify files location or directory for coverting it to go",
			Destination: &location,
		},
		cli.StringFlag{
			Name:        "destination",
			Value:       "",
			Usage:       "destination to save all the converted files",
			Destination: &destination,
		},
	}

	// default action
	app.Action = func(c *cli.Context) {
		if location == "" {
			cli.ShowAppHelp(c)
			message := logger.GetColorFormat(logger.Yellow, logger.Faint, "location got =  \"\" and must not empty")
			err := logger.GetColorFormat(logger.Cyan, logger.Faint, "Error :: ")
			fmt.Println(err, message)
			return
		}

		stealer.Convert(location, destination)
	}

	app.Run(os.Args)
}
