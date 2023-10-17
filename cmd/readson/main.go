package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"dcastanho.readson/internal/files"
	"dcastanho.readson/internal/logger"
	md "dcastanho.readson/internal/template"
	"github.com/urfave/cli/v2"
)

// TODO tests?
// TODO Capitalize/Title functions
// TODO Defines
// TODO Documentation

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
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"p"},
				// Value: "template",
				Usage: "`PATTERN` to assign each file a name",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				// Value: "template",
				Usage: "Output to a single file `OUTPUT`",
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
			filePattern := cCtx.String("name")
			output := cCtx.String("output")

			if output != "" && filePattern != "" {
				panic("Cannot assign both a name pattern and an output file")
			}

			logger.DeployLogger(cCtx.Bool("verbose"), os.Stdout)

			out, err := processTemplateFile(templFile)

			if err != nil {
				panic(err.Error())
			}

			OneTemplate(jsonFile, out, filePattern, output)

			err = os.Remove(out)
			if err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func replaceName(template string, name string, newText string) string {
	return strings.Replace(template, "$"+name+"$", newText, -1)
}
func processTemplateFile(inputFilePath string) (string, error) {
	processedDir, processedName := filepath.Split(inputFilePath)
	outputFilePath := filepath.Join(processedDir, "processed_"+processedName)

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return "", err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var name, newText string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "$$$ ") {
			parts := strings.SplitN(line[4:], " ", 2)
			if len(parts) == 2 {
				name = parts[0]
				newText = parts[1]
			}
		} else {
			line = replaceName(line, name, newText)
			_, _ = fmt.Fprintln(outputFile, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return outputFilePath, err
}

func OneTemplate(pattern string, templateFile string, filePattern string, output string) {
	ext := filepath.Ext(templateFile)
	templ, err := md.ParseTemplate(templateFile)

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

		var filename string

		if filePattern == "" && output == "" {
			filename = files.FileName(pattern) + strconv.Itoa(i) + ext
		} else if pattern != "" {
			dir, pattern := filepath.Split(filePattern)
			filename, _, err = ctx.Getter(ctx.Data, pattern)
			filename = dir + filename + ext
			if err != nil {
				panic(err.Error())
			}
		} else if output != "" {
			filename = output + ext
		}
		os.WriteFile(filename, []byte(res), fs.FileMode(os.O_CREATE))
		i++
		ctx = iterator()
	}

}
