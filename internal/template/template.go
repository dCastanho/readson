// Parser handles all the interactions with string templates.
// Parses and executes the template according to provided functions
// using the given data to replace in the proper locations
package parser

import (
	"errors"

	"github.com/robertkrimen/otto"
)

// Getter defines how to fetch a variable name from the data.
// The variable comes in the form of a string
// data is the set of bytes which represent the variables applicable to the template (extracted from a json file for example)
// pattern is string representing the pattern of the variable to access and must follow these rules:
//
//	"name": returns the variable with the name, returns error if it does not exist
//	"name[i]": returns the index i of the variable, if it is not an array returns error
//	"name->sub": returns the variable sub, that is a subproperty of name, error if either do not exist
type Getter func(data []byte, pattern string) (string, ElementType, error)

// ArrayEach defines to iterate over an array.
// data is the bytes of the array
// forEach is the function that is executed for each element, where curr is the bytes of the element and dataType is its type
//
// Should return an error if data is not an array
type ArrayEach func(data []byte, forEach func(curr []byte, dataType ElementType)) error

// ObjectEach defines how to iterate over the properties of an key:value object
// data is the bytes of the object
// forEach is the function that is execute for each property of the pair, where prop is the name of the property, val
// are the bytes of the stored value, and dataType is its element
//
// Should return an error if data is not an object
type ObjectEach func(data []byte, forEach func(prop string, val []byte, dataType ElementType)) error

// Struct that represents a parsed text Template
type Template struct {
	top node
}

// ASTContext contains all the necessary definitions to execute a template
type ASTContext struct {
	// Getter will serve to get variables from data
	Getter Getter

	// ObjectEach will serve to iterate over properties of objects
	ObjectEach ObjectEach

	// ArrayEach will serve to iterate over elements of arrays
	ArrayEach ArrayEach

	// Data is the data where variables to insert into the template exist
	Data []byte
}

// ParseTemplate takes a filename and parses the template into a Template struct
// if an error happens in parsing, it is returned.
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

// ApplyTemplate takes a context and a parsed template and performs the necessary replacements.
// The functions defined in ctx will be used to replace where needed sections of the parsed template.
// It returns the string of the final template - with all the replacements performed, or
func ApplyTemplate(template *Template, ctx *ASTContext) (string, error) {
	s, err := template.top.evaluate(ctx)
	return s, err
}

// TODO Slap VM in ctx?
// SetupUserFunctions sets up user provided functions in Javascript.
// It receives a string with all the functions properly defined in Javascript syntax and
// creates the necessary environment to run them in. If there is an error in the Javscript code,
// it is returned.
func SetupUserFunctions(text string) error {
	VM = otto.New()

	_, err := VM.Run(text)
	return err
}
