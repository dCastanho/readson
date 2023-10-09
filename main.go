package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "readson",
		Usage: "Turn JSON files into readable Markdown ones!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "templ",
				Aliases: []string{"t"},
				// Value: "template",
				Usage:    "File path to the `TEMPLATE`",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			jsonFile := cCtx.Args().Get(0)
			if jsonFile == "" {
				return errors.New("Missing JSON file path")
			}
			templFile := cCtx.String("templ")
			OneToOne(jsonFile, templFile)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func OneToOne(jsonFile string, templateFile string) {
	datJSON, err := os.ReadFile(jsonFile)
	if err != nil {
		print("Error ")
		println(err.Error())
	}

	datTempl, err := os.ReadFile(templateFile)
	if err != nil {
		print("Error ")
		println(err.Error())
	}

	res := ApplyJSON(datJSON, datTempl)

	os.WriteFile("test", []byte(res), 0644)
}
