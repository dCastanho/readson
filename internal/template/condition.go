package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// TODO Add ( ) logic and negation

type ElementType uint8

const Boolean ElementType = 0
const String ElementType = 1
const Number ElementType = 2
const Invalid ElementType = 3

type operatorType uint8

const lt operatorType = 0
const lte operatorType = 1
const gt operatorType = 2
const gte operatorType = 3
const eq operatorType = 4
const dif operatorType = 5

func convertOperator(operator string) operatorType {
	switch operator {
	case "=":
		return eq
	case ">=":
		return gte
	case "<=":
		return lte
	case ">":
		return gt
	case "<":
		return lt
	case "!=":
		return dif
	}
	return 0
}

type condition interface {
	eval(ctx *ASTContext) (bool, error)
}

type singleCondition struct {
	element element
}

type andCondition struct {
	left  condition
	right condition
}

func (c andCondition) eval(ctx *ASTContext) (bool, error) {
	left, err := c.left.eval(ctx)

	if err != nil {
		return false, err
	}

	right, err := c.right.eval(ctx)

	if err != nil {
		return false, err
	}

	println("evaluating and", left, right)

	return left && right, nil
}

type orCondition struct {
	left  condition
	right condition
}

func (c orCondition) eval(ctx *ASTContext) (bool, error) {
	left, err := c.left.eval(ctx)

	if err != nil {
		return false, err
	}

	right, err := c.right.eval(ctx)

	if err != nil {
		return false, err
	}

	return left || right, nil
}

type operatorCondition struct {
	left     element
	right    element
	operator operatorType
}

func sameTypes(a element, b element, ctx *ASTContext) (bool, ElementType, ElementType) {
	fmt.Println(a, b)
	aType := a.typeof(ctx)
	bTyte := b.typeof(ctx)
	return aType == bTyte, aType, bTyte
}

func (c operatorCondition) eval(ctx *ASTContext) (bool, error) {

	areSameType, actual, bType := sameTypes(c.left, c.right, ctx)

	if !areSameType {
		return false, errors.New(fmt.Sprintf("Cannot compare of different types left:%d right:%d", actual, bType))
	}

	result := compare(c.left, c.right, ctx, actual)

	booleanCompareError := errors.New("Booleans cannot be compared this way")

	switch c.operator {
	case eq:
		return result == 0, nil
	case lt:
		if actual == Boolean {
			return false, booleanCompareError
		}

		return result < 0, nil
	case lte:
		if actual == Boolean {
			return false, booleanCompareError
		}

		return result <= 0, nil
	case gt:
		if actual == Boolean {
			return false, booleanCompareError
		}

		return result > 0, nil
	case gte:
		if actual == Boolean {
			return false, booleanCompareError
		}

		return result >= 0, nil
	case dif:
		if actual == Boolean {
			return false, booleanCompareError
		}

		return result != 0, nil
	}
	return false, errors.New("Invalid operator value")

}

type element interface {
	typeof(ctx *ASTContext) ElementType
	value(ctx *ASTContext) any
}

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
			return nil, Invalid, err
		}

		tpe = Number
		return fl, Number, nil
	}
}

func gottedToElem(text string, elemenType ElementType) (any, error) {
	if elemenType == String {
		return text, nil
	} else if text == "true" {
		return true, nil
	} else if text == "false" {
		return false, nil
	} else {
		fl, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, err
		}
		return fl, nil
	}
}

// Assumes they are comparable
func compare(e element, other element, ctx *ASTContext, actual ElementType) int8 {

	switch actual {
	case Boolean:
		first, _ := e.value(ctx).(bool)
		second, _ := other.value(ctx).(bool)

		if first == second {
			return 0
		} else {
			return -1
		}
	case String:
		first, _ := e.value(ctx).(string)
		second, _ := other.value(ctx).(string)

		return int8(strings.Compare(first, second))

	case Number:
		return numberCompare(e.value(ctx), other.value(ctx))
	}
	return 0
}

type accessElement struct {
	pattern string
}

// TODO make access elements functional

func (e accessElement) typeof(ctx *ASTContext) ElementType {

	_, ttype, err := getPattern(e.pattern, ctx)
	if err != nil {
		//TODO error handling in typeof
	}

	return ttype
}

func getPattern(pattern string, ctx *ASTContext) (any, ElementType, error) {
	elem, ttype, err := ctx.Getter(ctx.Data, pattern)
	if err != nil {
		return nil, Invalid, err
	}

	toAny, err := gottedToElem(elem, ttype)
	return toAny, ttype, err
}

func (e accessElement) value(ctx *ASTContext) any {
	toAny, _, err := getPattern(e.pattern, ctx)
	if err != nil {
		//TODO error handling in typeof
	}

	return toAny
}

type constantElement struct {
	constant  any
	constType ElementType
}

func (e constantElement) typeof(ctx *ASTContext) ElementType {
	return e.constType
}

func (e constantElement) value(ctx *ASTContext) any {
	return e.constant
}

func numberCompare(a, b any) int8 {
	actualA, _ := a.(float64)
	actualB, _ := b.(float64)

	println(actualA, actualB)

	if actualA == actualB {
		return 0
	} else if actualA > actualB {
		return 1
	} else {
		return -1
	}

}
