package parser

type Getter func([]byte, string) (string, error)

type Template struct {
	top node
}

func ParseTemplate(templateData []byte) Template {
	top, err := Parse("", templateData)
	if err != nil {
		// TODO handle parse error
	}

	actual, ok := top.(node)
	if !ok {
		// TODO handle incorrect type error
	}

	return Template{top: actual}
}

func ApplyTemplate(template Template, data []byte, getter Getter) string {
	return template.top.evaluate(data, getter)
}
