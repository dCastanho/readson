package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"dcastanho.readson/internal/files"
	parser "dcastanho.readson/internal/parser"
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
				Usage: "File path to the `TEMPLATE`",
				// Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			jsonFile := cCtx.Args().Get(0)
			if jsonFile == "" {
				return errors.New("Missing JSON file path")
			}
			// getNext := GetX(path)

			ret, err := parser.ParseFile("test.md")

			println(ret)
			if err != nil {
				println(err.Error())
			}
			// templFile := cCtx.String("templ")
			// // OneTemplate(jsonFile, templFile)

			// OneTemplate(jsonFile, templFile)

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
	templ := md.ParseTemplate(string(datTempl))

	iterator := files.GetData(pattern)
	curr, get := iterator()

	i := 0

	for curr != nil {
		res := md.ApplyTemplate(templ, curr, get)
		filename := files.FileName(pattern) + strconv.Itoa(i) + ext
		println(filename)
		os.WriteFile(filename, []byte(res), 0064)
		i++
		curr, get = iterator()
	}

}
