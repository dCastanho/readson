package parser

import (
	"errors"
)

// TODO Add logging

type Getter func([]byte, string) (string, ElementType, error)

type Template struct {
	top node
}

type ASTContext struct {
	Getter Getter
	Data   []byte
}

func ParseTemplate(templateData []byte) (*Template, error) {
	top, err := Parse("", templateData)
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
