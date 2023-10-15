package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"dcastanho.readson/internal/files"
	"dcastanho.readson/internal/logger"
	md "dcastanho.readson/internal/template"
	"github.com/urfave/cli/v2"
)

// TODO allow for custom file names with json expressions

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
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				// Value: "template",
				Usage: "Print logs to standard out",
				// Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			jsonFile := cCtx.Args().Get(0)
			if jsonFile == "" {
				return errors.New("Missing JSON file path")
			}
			// getNext := GetX(path)
			templFile := cCtx.String("templ")
			logger.DeployLogger(cCtx.Bool("verbose"), os.Stdout)
			OneTemplate(jsonFile, templFile)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func OneTemplate(pattern string, templateFile string) {
	datTempl, err := os.ReadFile(templateFile)
	ext := filepath.Ext(templateFile)
	if err != nil {
		print("Error ")
		println(err.Error())
	}
	templ, err := md.ParseTemplate(datTempl)

	if err != nil {
		panic(err)
	}

	iterator := files.GetData(pattern)
	ctx := iterator()

	i := 0

	for ctx != nil {
		res, err := md.ApplyTemplate(templ, ctx)

		if err != nil {
			panic(err)
		}

		filename := files.FileName(pattern) + strconv.Itoa(i) + ext
		os.WriteFile(filename, []byte(res), fs.FileMode(os.O_CREATE))
		i++
		ctx = iterator()
	}

}
