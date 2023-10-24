package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/robertkrimen/otto"
)

var VM *otto.Otto // All functions assume it has been initialized

type userFunc struct {
	name       string
	parameters []element
}

// Returns the value, given a context ctx (in case variable accesses are necessary) and
// its element type.
// value(ctx *ASTContext) (any, ElementType, error)

// // stringValue returns a string representation of the value, particularly useful during
// // node evaluation
// stringValue(ctx *ASTContext) (string, error)

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

		parameter, _, err := par.value(ctx)

		if err != nil {
			return nil, err
		}

		sb.WriteString(fmt.Sprintf("\"%s\"", parameter))
	}

	sb.WriteRune(')')
	funcCall := sb.String()

	result, err := VM.Run(funcCall)
	return &result, err
}

func ottoToElemType(v *otto.Value) ElementType {
	if v.IsBoolean() {
		return Boolean
	} else if v.IsString() {
		return String
	} else if v.IsNumber() {
		return Number
	} else if v.IsObject() {
		return Object
	}
	return NotExists

}

// Returns the type of element that it is
func (f userFunc) typeof(ctx *ASTContext) ElementType {
	v, err := f.call(ctx)
	if err != nil {
		return NotExists
	}

	return ottoToElemType(v)
}

// Returns the value, given a context ctx (in case variable accesses are necessary) and
// its element type.
func (f userFunc) value(ctx *ASTContext) (any, ElementType, error) {
	v, err := f.call(ctx)
	if err != nil {
		return nil, NotExists, err
	}

	return v.String(), ottoToElemType(v), nil
}

// Returns the value, given a context ctx (in case variable accesses are necessary) and
// its element type.
func (f userFunc) stringValue(ctx *ASTContext) (string, error) {
	v, _, e := f.value(ctx)
	return v.(string), e
}
