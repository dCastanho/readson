package parser

import "errors"

type Getter func([]byte, string) (string, error)

type Template struct {
	top node
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

func ApplyTemplate(template *Template, data []byte, getter Getter) (string, error) {
	return template.top.evaluate(data, getter)
}
