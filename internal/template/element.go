package parser

import (
	"errors"
	"strconv"
	"strings"
)

// ElementType defines the type of an element
type ElementType uint8

const Boolean ElementType = 0
const String ElementType = 1
const Number ElementType = 2
const Array ElementType = 3
const Object ElementType = 4
const NotExists ElementType = 5

// An element represents a value in the tree - it can be any value. It can refer to functions,
// variable accesses or constants.
type element interface {

	// Returns the value, given a context ctx (in case variable accesses are necessary) and
	// its element type.
	value(ctx *ASTContext) (any, ElementType, error)

	// stringValue returns a string representation of the value, particularly useful during
	// node evaluation
	stringValue(ctx *ASTContext) (string, error)
}

// convertType converts a text type definition (array, bool, object, number or string)
// into its matching ElementType
func convertType(text string) ElementType {
	if text == "array" {
		return Array
	} else if text == "object" {
		return Object
	} else if text == "number" {
		return Number
	} else if text == "string" {
		return String
	} else if text == "bool" {
		return Boolean
	}
	return NotExists
}

// textToElem takes plain text constants and parses it into its correct element
func textToElem(text string) (any, ElementType, error) {
	var tpe ElementType
	if text[0] == '"' {
		tpe = String
		text = text[1 : len(text)-1]
		return text, tpe, nil
	} else if text == "true" {
		return true, Boolean, nil
	} else if text == "false" {
		return false, Boolean, nil
	} else {
		fl, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, NotExists, err
		}

		tpe = Number
		return fl, Number, nil
	}

}

// anyElemToString converts the value returned by an element (any), into a string
// when the element is of type String, Boolean or Number.
func anyElemToString(elem any, elementType ElementType) (string, error) {
	if elementType == String {
		s := elem.(string)
		return s, nil
	} else if elementType == Boolean {
		b := elem.(bool)
		return strconv.FormatBool(b), nil
	} else if elementType == Number {
		fl := elem.(float64)
		s := strconv.FormatFloat(fl, 'f', 2, 64)
		return s, nil
	}
	return elem.(string), nil
}

// typedStringToElem takes a string that has already been classified and converts it to its
// correct .go representation (ex: boolean become true/false)
func typedStringToElem(text string, elemenType ElementType) (any, error) {
	if elemenType == String {
		return text, nil
	} else if elemenType == Boolean {
		if text == "true" {
			return true, nil
		} else if text == "false" {
			return false, nil
		}
	} else if elemenType == Number {
		fl, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, err
		}
		return fl, nil
	}
	return text, nil
}

// compare takes two elements of any kind and compares the two, returning 0 if equal, -1 if e < other or
// 1 if e > other.
//
// ctx is the necessary context to execute variable accesses and actual represents the type of the elements
// being compared - which are assumed to be of the same type
func compare(e element, other element, ctx *ASTContext, actual ElementType) (int8, error) {

	otherVal, _, err := other.value(ctx)

	if err != nil {
		return 0, err
	}

	thisVal, _, err := e.value(ctx)

	if err != nil {
		return 0, err
	}

	switch actual {
	case Boolean:
		first, _ := thisVal.(bool)
		second, _ := otherVal.(bool)

		if first == second {
			return 0, nil
		} else {
			return -1, nil
		}
	case String:
		first, _ := thisVal.(string)
		second, _ := otherVal.(string)

		return int8(strings.Compare(first, second)), nil

	case Number:
		return numberCompare(thisVal, otherVal), nil
	}
	return 0, errors.New("Unsuported comparison")
}

// accessElements represent variable accesses to the given data
type accessElement struct {

	// pattern that represents the path to the variable (ex: name->sub)
	pattern string
}

// name of the access element returns the variable access pattern
func (e accessElement) name() string {
	return e.pattern
}

// getPattern searches for a given pattern in a context ctx, returning it as a string
// Returns an error if it does not exist
func getPattern(pattern string, ctx *ASTContext) (string, ElementType, error) {
	elem, ttype, err := ctx.Getter(ctx.Data, pattern)
	if err != nil {
		return "", NotExists, err
	}

	return elem, ttype, err
}

// The value of access element is the value of its variable in the data. Returns error
// if a variable matching the pattern does not exist.
func (e accessElement) value(ctx *ASTContext) (any, ElementType, error) {
	elem, tpe, err := getPattern(e.pattern, ctx)
	if err != nil {
		return nil, NotExists, err
	}

	toAny, err := typedStringToElem(elem, tpe)

	if err != nil {
		return nil, NotExists, err
	}

	return toAny, tpe, nil
}

// stringValue finds and returns the variable matching the pattern, as a string
// Improves efficency when directly accessing for a string, as it avoids the need for type conversions
func (e accessElement) stringValue(ctx *ASTContext) (string, error) {
	elem, _, err := getPattern(e.pattern, ctx)
	if err != nil {
		return "", err
	}

	return elem, err
}

// constantElement represents constants defined in the template, in functions or comparisons for example.
// Constants can be bools, numbers or strings
type constantElement struct {
	// constant, saved as a string representation of its value
	constant string
}

// stringValue is  textual representation of the constant (how it was extracted from the text)
func (e constantElement) stringValue(ctx *ASTContext) (string, error) {
	return e.constant, nil
}

// value of the constant, already parsed in case of floats or bools
func (e constantElement) value(ctx *ASTContext) (any, ElementType, error) {
	return textToElem(e.constant)
}

// numberCompares compares to numbers in the form of any variables
// if a > b it returns 1
// if a < b it returns -1
// if a = b it returns 0
func numberCompare(a, b any) int8 {
	actualA, _ := a.(float64)
	actualB, _ := b.(float64)

	if actualA == actualB {
		return 0
	} else if actualA > actualB {
		return 1
	} else {
		return -1
	}

}
