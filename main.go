package main

import (
	"os"

	md "dcastanho.readson/template"
	"dcastanho.readson/template/expressions"
)

func main() {

	dat, err := os.ReadFile("test.json")
	if err != nil {
		print("Error ")
		println(err.Error())
		return
	}

	templ, err := os.ReadFile("test.md")
	if err != nil {
		print("Error ")
		println(err.Error())
		return
	}

	t_text := string(templ)

	temp := md.ParseTemplate(t_text, expressions.JSONParserGetter)
	print(md.ApplyTemplate(temp, dat))
}
