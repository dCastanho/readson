package main

import (
	"errors"
	"log"
	"os"

	md "dcastanho.readson/template"
	"dcastanho.readson/template/expressions"
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
				Usage: "File path to the `TEMPLATE`",
				// Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			jsonFile := cCtx.Args().Get(0)
			if jsonFile == "" {
				return errors.New("Missing JSON file path")
			}
			files := GetFiles(jsonFile)
			templFile := cCtx.String("templ")
			// OneTemplate(jsonFile, templFile)

			OneTemplate(files, templFile)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func OneTemplate(jsonFiles []string, templateFile string) {
	datTempl, err := os.ReadFile(templateFile)
	if err != nil {
		print("Error ")
		println(err.Error())
	}
	templ := md.ParseTemplate(string(datTempl), expressions.JSONParserGetter)

	for _, file := range jsonFiles {
		WriteOne(file, templ)
	}

}

func WriteOne(jsonFile string, template md.Template) {

	dat, _ := os.ReadFile(jsonFile) // check for file has been done before
	res := []byte(md.ApplyTemplate(template, dat))

	os.WriteFile(FileName(jsonFile)+".md", res, 0064)

}
