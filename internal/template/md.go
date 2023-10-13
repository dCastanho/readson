package md

import (
	"strings"
)

type block struct {
	text         string
	isExpression bool
}

type Getter func([]byte, string) (string, error)

type templateIter func() (block, bool)

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

func iterateTemplate(template Template) templateIter {
	i := 0

	return func() (block, bool) {
		if len(template.blocks) == i {
			return block{text: "", isExpression: false}, false
		} else {
			dat := template.blocks[i]
			i++
			return dat, true
		}
	}

}

func ApplyTemplate(template Template, values []byte, getter Getter) string {
	var sb strings.Builder

	iterator := iterateTemplate(template)
	current_block, going := iterator()

	for going {
		text := current_block.text
		println("block", text)
		if current_block.isExpression { // is a variable
			str, _ := getter(values, text)
			sb.WriteString(str)
		} else {
			sb.WriteString(text)
		}

		current_block, going = iterator()
	}

	return sb.String()
}

func resolveExpression(expression string, values []byte, getter Getter, iterator templateIter) string {
	if expression[0] == '"' {
		return expression
	} else if expression[:2] == "if" {
		// TODO handle IF
	} else if expression[:3] == "for" {
		// TODO handle for
	}

	val, _ := getter(values, strings.Trim(expression, ""))
	return val
}

func evaluteCondition(condition string, values []byte, getter Getter) bool {
	return false
}

// func handleIf(expression string, values []byte, getter Getter, iterator templateIter) string {
// 	condition := expression[2:]

// 	evalToTrue, _ := iterator()

// 	trueBlocks := make([]block, 0)

// 	done := false
// 	// accumulate block to eval in true case
// 	for  !done {
// 		// 	TODO iconrrect, need a stack to check if this else/end refers to this if.
// 		if evalToTrue.isExpression {
// 			if evalToTrue.text == "else" {

// 			} else if evalToTrue.text == "end" {

// 			} else

// 		} else {
// 			trueBlocks = append(trueBlocks, evalToTrue)
// 			// add to accumulator
// 		}
// 	}

// }
