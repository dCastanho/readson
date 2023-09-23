package md

import "strings"

type block struct {
	text       string
	expression bool
}

type Template []block

func ParseTemplate(templateText string) Template {

	const expressionByte byte = 36

	var components Template = make([]block, 0)

	var isExp bool = templateText[0] == expressionByte
	var sb strings.Builder

	for i := 0; i < len(templateText); i++ {
		if templateText[i] == expressionByte {
			components = append(components, block{text: sb.String(), expression: isExp})
			sb = strings.Builder{}
			isExp = !isExp
		} else {
			sb.WriteByte(templateText[i])
		}
	}

	return components
}

func ApplyTemplate(template Template, values map[string]string) string {
	var sb strings.Builder
	for i := 0; i < len(template); i++ {
		if template[i].expression { // is a variable
			sb.WriteString(values[template[i].text])
		} else {
			sb.WriteString(template[i].text)
		}
	}

	return sb.String()
}
