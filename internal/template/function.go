package parser

import (
	"errors"
	"fmt"
	"strings"

	"dcastanho.readson/internal/logger"
	"github.com/robertkrimen/otto"
)

var VM *otto.Otto // All functions assume it has been initialized

type userFunc struct {
	baseNode
	name       string
	parameters []node
}

func (f userFunc) call(ctx *ASTContext) (*otto.Value, error) {

	if VM == nil {
		return nil, errors.New("Javascript environment not started, you most likely forgot to supply a file")
	}

	sb := strings.Builder{}
	sb.WriteString("result = ")
	sb.WriteString(f.name)
	sb.WriteRune('(')

	for i, par := range f.parameters {
		if i > 0 {
			sb.WriteRune(',')
		}

		parameter, err := par.evaluate(ctx)

		if err != nil {
			return nil, err
		}

		sb.WriteString(fmt.Sprintf("\"%s\"", parameter))
	}

	sb.WriteRune(')')
	funcCall := sb.String()
	fmt.Println(funcCall)

	fmt.Println(VM)

	result, err := VM.Run(funcCall)
	fmt.Println(result, err)
	return &result, err
}

func (f userFunc) evaluate(ctx *ASTContext) (string, error) {

	logger.DefaultLogger.Node("Function: " + f.name)

	res, err := f.call(ctx)

	if err != nil {
		return "", err
	}

	return f.withChild((*res).String(), ctx)
}

func (f userFunc) eval(ctx *ASTContext) (bool, error) {
	res, err := f.call(ctx)

	if err != nil {
		return false, err
	}

	return res.ToBoolean()
}
