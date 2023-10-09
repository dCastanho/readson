package md

import (
	"strings"
)

type block struct {
	text         string
	isExpression bool
}

type Getter func([]byte, string) (string, error)

type Template struct {
	blocks []block
	get    Getter
}

func ParseTemplate(templateText string, getter Getter) Template {

	const expressionByte byte = 36

	var components []block = make([]block, 0)

	var isExp bool = templateText[0] == expressionByte
	var sb strings.Builder

	for i := 0; i < len(templateText); i++ {
		if templateText[i] == expressionByte {
			components = append(components, block{text: sb.String(), isExpression: isExp})
			sb = strings.Builder{}
			isExp = !isExp
		} else {
			sb.WriteByte(templateText[i])
		}
	}

	return Template{blocks: components, get: getter}
}

func ApplyTemplate(template Template, values []byte) string {
	var sb strings.Builder
	for _, current_block := range template.blocks {
		text := current_block.text
		if current_block.isExpression { // is a variable
			str, _ := template.get(values, text)
			sb.WriteString(str)
		} else {
			sb.WriteString(text)
		}
	}

	return sb.String()
}
