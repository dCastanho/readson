package main

import (
	"bufio"
	"errors"
	"fmt"
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

// TODO math
// TODO tests?
// TODO functions in conditions
// TODO Documentation
// TODO check correctness of current grammar - can it be reduced/optimized? <-
// TODO make blocks nicer to use
// TODO make errors nicer if possible
// TODO better logging when possible

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
			&cli.StringFlag{
				Name:    "functions",
				Aliases: []string{"f"},
				// Value: "template",
				Usage: "javascript `FILE` with user custom functions",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				// Value: "template",
				Usage: "Print logs to standard out",
				// Required: true,
			},
			&cli.BoolFlag{
				Name:    "keep",
				Aliases: []string{"k"},
				// Value: "template",
				Usage: "Keep the pre-processed template",
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

			functionsFile := cCtx.String("functions")

			if functionsFile != "" {
				fmt.Println(functionsFile)
				err := SetupFunctions(functionsFile)

				if err != nil {
					panic(err.Error())
				}
			}

			logger.DeployLogger(cCtx.Bool("verbose"), os.Stdout)

			out, err := processTemplateFile(templFile)

			if err != nil {
				panic(err.Error())
			}

			OneTemplate(jsonFile, out, filePattern, output)

			if !cCtx.Bool("keep") {
				err = os.Remove(out)
				if err != nil {
					panic(err.Error())
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func SetupFunctions(jsFile string) error {
	dat, err := os.ReadFile(jsFile)
	if err != nil {
		return err
	}

	text := string(dat)
	err = md.SetupUserFunctions(text)
	return err
}

func replaceName(defines *map[string]string, line string) string {

	curr := line
	for key, val := range *defines {
		curr = strings.Replace(curr, "$"+key+"$", val, -1)
	}
	return curr
}
func processTemplateFile(inputFilePath string) (string, error) {
	processedDir, processedName := filepath.Split(inputFilePath)
	outputFilePath := filepath.Join(processedDir, "processed"+processedName)

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

	defines := make(map[string]string)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if !found && strings.HasPrefix(line, "$$$ ") {
			parts := strings.SplitN(line[4:], " ", 2)
			if len(parts) == 2 {
				name := parts[0]
				newText := parts[1]
				defines[name] = newText
			}
		} else {
			found = true
			line = replaceName(&defines, line)
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

		os.Remove(filename)
		os.WriteFile(filename, []byte(res), 0644)

		i++
		ctx = iterator()
	}

}
