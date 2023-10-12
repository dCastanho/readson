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
}

func ParseTemplate(templateText string) Template {

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

	if sb.Len() > 0 {
		components = append(components, block{text: sb.String(), isExpression: false})
	}

	return Template{blocks: components}
}

// TODO allow for more complex patterns of template (if, for, ternary operator)

func ApplyTemplate(template Template, values []byte, getter Getter) string {
	var sb strings.Builder
	for _, current_block := range template.blocks {
		text := current_block.text
		println("block", text)
		if current_block.isExpression { // is a variable
			str, _ := getter(values, text)
			sb.WriteString(str)
		} else {
			sb.WriteString(text)
		}
	}

	return sb.String()
}
