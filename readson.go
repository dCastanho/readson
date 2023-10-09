package main

import (
	md "dcastanho.readson/template"
	"dcastanho.readson/template/expressions"
)

func ApplyJSON(json []byte, template []byte) string {
	templ := md.ParseTemplate(string(template), expressions.JSONParserGetter)
	result := md.ApplyTemplate(templ, json)
	return result
}
