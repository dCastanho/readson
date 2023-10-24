package parser

import (
	"errors"
)

// operatorType defines the operation comparison conditions (>, <, =, etc.)
type operatorType uint8

// lt means the comparison is a lesser than (a < b)
const lt operatorType = 0

// lte means the comparison is a lesser or equal (a <= b)
const lte operatorType = 1

// gt means the comparison is a greater than (a > b)
const gt operatorType = 2

// gt means the comparison is a greater or equal (a >= b)
const gte operatorType = 3

// eq means the comparison is an equal than (a = b)
const eq operatorType = 4

// dif means the comparison is a difference (a != b)
const dif operatorType = 5

// convertOperator takes a string operator and returns the correct type
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

// condition interface defines conditions to be used in If nodes of the AST
type condition interface {

	// eval evaluates the condition and returns the boolean
	eval(ctx *ASTContext) (bool, error)
}

//! Necessary?

// condition that is associated with a single element, returning its value if it is a boolean
type singleCondition struct {
	element element
}

// singleCondition eval checks if its element is a boolean and returns it if so, error otherwise
func (c singleCondition) eval(ctx *ASTContext) (bool, error) {
	v, tpe, err := c.element.value(ctx)

	if err != nil {
		return false, err
	}

	if tpe != Boolean {
		return false, errors.New("Condition must be a boolean")
	}

	return v.(bool), err
}

// andCondition represents an AND operation between two other conditions
type andCondition struct {
	left  condition
	right condition
}

// eval on andCondition evaluates the left condition and returns false if it is false (thus not evaluating right)
// if left is true then it evaluates right and checks the result.
func (c andCondition) eval(ctx *ASTContext) (bool, error) {
	left, err := c.left.eval(ctx)

	if err != nil {
		return false, err
	}

	if !left {
		return false, nil
	}

	right, err := c.right.eval(ctx)

	if err != nil {
		return false, err
	}

	return left && right, nil
}

// orCondition represents an OR operation between conditions
type orCondition struct {
	left  condition
	right condition
}

// eval on orCondition evaluates the left condition and returns true if it is true (thus not evaluating right)
// if left is false then it evaluates right and checks the result.
func (c orCondition) eval(ctx *ASTContext) (bool, error) {
	left, err := c.left.eval(ctx)

	if err != nil {
		return false, err
	}

	if left {
		return true, nil
	}

	right, err := c.right.eval(ctx)

	if err != nil {
		return false, err
	}

	return left || right, nil
}

// negatedCondition is the negation of a condition
type negatedCondition struct {
	toNegate condition
}

// eval on negatedCondition checks the result of its child condition and returns the opposite
func (c negatedCondition) eval(ctx *ASTContext) (bool, error) {
	v, err := c.toNegate.eval(ctx)

	return !v, err
}

// operatorCondition represents a boolean comparison between elements (such as >, <, or =).
type operatorCondition struct {
	left     element
	right    element
	operator operatorType
}

// eval on operatorCondition checks whether the comparison is possible and performs it.
// returns an error if the elements are of different types or if types other than numbers
// are compared with a non-equal operator
func (c operatorCondition) eval(ctx *ASTContext) (bool, error) {

	v1, tpe1, err := c.left.value(ctx)

	if err != nil {
		return false, err
	}

	v2, tpe2, err := c.left.value(ctx)

	if err != nil {
		return false, err
	}

	if tpe1 != tpe2 {
		return false, errors.New("Values must be of the same type")
	}

	if c.operator != eq && tpe1 != Number {
		return false, errors.New("Only numbers can have non-equal comparisons")
	}

	switch c.operator {
	case eq:
		return v1 == v2, nil
	case lt:
		r := numberCompare(v1, v2)
		return r < 0, nil
	case lte:
		r := numberCompare(v1, v2)
		return r <= 0, nil
	case gt:
		r := numberCompare(v1, v2)
		return r > 0, nil
	case gte:
		r := numberCompare(v1, v2)
		return r >= 0, nil

	}

	return false, errors.New("Invalid operator value")

}

// existsCondition is a condition that checks whether a given element exists,
// meant as a variable verification prior to an access if there is a chance
// that variable might not exist
type existsCondition struct {
	element element
}

// eval on existsCondition returns true if the element exists, false otherwise.
// constants will always return true.
func (n existsCondition) eval(ctx *ASTContext) (bool, error) {
	_, tpe, err := n.element.value(ctx)
	return tpe != NotExists && err != nil, nil
}

// ofType checks whether a given element is of a certain type.
// Meant as a prior check to variable accesses or function use
type ofType struct {
	element element
	typeOf  ElementType
}

// eval on ofType fetches the value of the variable and checks its type
func (n ofType) eval(ctx *ASTContext) (bool, error) {
	_, tpe, err := n.element.value(ctx)
	if err != nil {
		return false, err
	}

	return tpe == n.typeOf, nil
}
