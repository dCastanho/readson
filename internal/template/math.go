package parser

import (
	"errors"
	"fmt"
	"strconv"
)

func textToOperator(text byte) mathOperator {
	switch text {
	case '+':
		return addition
	case '-':
		return subtraction
	case '*':
		return multiplication
	case '/':
		return division
	default:
		return invalid
	}
}

type mathOperator = uint8

const addition mathOperator = 0
const subtraction mathOperator = 1
const division mathOperator = 2
const multiplication mathOperator = 3
const invalid mathOperator = 4

// expression represents a mathmatical expression between two elements
type expression struct {

	// element on the left of the operator
	left element

	// element on the right of the operator
	right element

	// operation is the kind of operation this expression performs (i.e. addition)
	operation mathOperator
}

// given an element, getNumber evaluates it and returns its conversion to a float,
// if that element does in fact represent a number, otherwise it gives an error.
func getNumber(e element, ctx *ASTContext) (float64, error) {
	val, tpe, err := e.value(ctx)

	fmt.Println(val, tpe, err)

	if tpe != Number {
		return 0, errors.New("Invalid mathmatical expression, can only operate over numbers")
	}

	if err != nil {
		return 0, err
	}

	return val.(float64), nil
}

// performMath performs the appropriate mathmatical operation, defined by the operation argument
// and using left and right arguments in their correct order.
func performMath(left float64, right float64, operation mathOperator) float64 {
	switch operation {
	case addition:
		return left + right
	case subtraction:
		return left - right
	case division:
		return left / right
	case multiplication:
		return left * right
	}
	return 0
}

// value of an expression evaluates the left side and checks that it is a number
// then does the samne for the right side and finally performs the mathmatical operation.
func (e expression) value(ctx *ASTContext) (any, ElementType, error) {
	leftNumber, err := getNumber(e.left, ctx)

	if err != nil {
		return nil, NotExists, err
	}

	rightNumber, err := getNumber(e.right, ctx)

	if err != nil {
		return nil, NotExists, err
	}

	result := performMath(leftNumber, rightNumber, e.operation)

	return result, Number, nil
}

func (e expression) stringValue(ctx *ASTContext) (string, error) {
	v, _, err := e.value(ctx)

	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(v.(float64), 'f', -1, 64), nil

}
