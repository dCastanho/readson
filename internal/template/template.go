package parser

import (
	"errors"
)

type Getter func([]byte, string) (string, ElementType, error)

type ArrayEach func(data []byte, forEach func(curr []byte, dataType ElementType)) error

type ObjectEach func(data []byte, forEach func(prop string, val []byte, dataType ElementType)) error

type Template struct {
	top node
}

type ASTContext struct {
	Getter     Getter
	ObjectEach ObjectEach
	ArrayEach  ArrayEach
	Data       []byte
}

func ParseTemplate(templateName string) (*Template, error) {
	top, err := ParseFile(templateName)
	if err != nil {
		return nil, err
	}

	actual, ok := top.(node)
	if !ok {
		return nil, errors.New("Incorrect syntax somewhere") // Not great, but this error should not happen.
	}

	return &Template{top: actual}, nil
}

func ApplyTemplate(template *Template, ctx *ASTContext) (string, error) {
	return template.top.evaluate(ctx)
}
