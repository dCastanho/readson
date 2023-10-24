package parser

import (
	"strconv"
	"strings"

	"dcastanho.readson/internal/logger"
)

// node is a node in the AST tree, evaluate returns the interpreted
// text (replacements and flows executed) of itself and all its 'next's
type node interface {
	evaluate(ctx *ASTContext) (string, error)
	next() node
	setNext(n node)
}

// baseNode serves to provide the next functionality to all nodes, due to
// being present regardless of its type
type baseNode struct {
	// child is the node to evaluate next
	child node
}

// returns the next node to evaluate
func (b *baseNode) next() node {
	return b.child
}

// setNext sets the next node to evaluate
func (b *baseNode) setNext(n node) {
	b.child = n
}

// withChild receives a string and if the node has a next, evaluates it and
// concatenates the two strings. This proceess is common across all nodes -
// they evaluate themselves, and then the next one, returning the concatenation.
func (b *baseNode) withChild(def string, ctx *ASTContext) (string, error) {

	if b.child == nil {
		return def, nil
	} else {
		childText, err := b.child.evaluate(ctx)
		if err != nil {
			return "", err
		}
		return def + childText, nil
	}
}

// textNode is a node which represents plain text, to be printed in the same way
// as it was in the provided template
type textNode struct {
	baseNode

	// text in the template
	text string
}

// evaluate of a textNode returns the stored text + the evaluated text of its proceeding nodes
func (n *textNode) evaluate(ctx *ASTContext) (string, error) {
	thisText := n.text

	logger.DefaultLogger.Node("Text:", thisText)

	result, err := n.withChild(thisText, ctx)
	return result, err
}

// accessNode represents the access to a variable - be it a function, data access or a constant
type accessNode struct {
	baseNode

	// accessorElement is the element generated from the template
	// it is from where the string of this node is extracted
	accessorElement element
}

// evaluate of an accessNode fetches the string value of the variable + the evaluated text of its
// proceeding nodes
func (n accessNode) evaluate(ctx *ASTContext) (string, error) {

	logger.DefaultLogger.Node("Acessor Block:")

	v, err := n.accessorElement.stringValue(ctx)

	if err != nil {
		return "", err
	}

	result, err := n.withChild(v, ctx)
	return result, err
}

// ifNode represents an if in the text, which has a condition and two children nodes.
// If the condition evaluates to true, the trueClause child is evaluated.
// If the condition evaluates to false, the falseClause child is evaluated.
type ifNode struct {
	baseNode

	// condition of the if
	condition condition

	// trueClause is evaluated if condition
	trueClause node

	// falseClause is evaluated if !condition
	falseClause node
}

// evaluate on ifNode checks the condition and evaluates the correct child,
// depending on the result
func (n *ifNode) evaluate(ctx *ASTContext) (string, error) {

	logger.DefaultLogger.Node("If")

	var err error

	result, err := n.condition.eval(ctx)

	if err != nil {
		return "", err
	}

	var clause string

	if result {
		clause, err = n.trueClause.evaluate(ctx)

	} else if n.falseClause != nil {
		clause, err = n.falseClause.evaluate(ctx)
	}

	if err != nil {
		return "", err
	}

	res, err := n.withChild(clause, ctx)
	return res, err

}

// forNode represents the execution of a loop in the template. It always has two associated variables,
// which vary depending on the type of loop. It is an iterative loop, which evaluates its child for every
// member of the element it is iterating
type forNode struct {
	baseNode
	itemName  string
	indexName string
	pattern   element
	loop      node
	forType   bool
}

// rangeFor is the execution of a for node according to the iteration of array (represented by a slice of bytes)
// it returns the evaluation of the child node for each element of the array
func (n *forNode) rangeFor(ctx *ASTContext, array []byte) string {
	sb := strings.Builder{}
	i := 1

	forEach := func(curr []byte, dataType ElementType) {

		newGetter := func(data []byte, pattern string) (string, ElementType, error) {

			if pattern == n.indexName {
				return strconv.Itoa(i), Number, nil
			}

			updated, usesItem := strings.CutPrefix(pattern, n.itemName)

			if usesItem {
				if updated != "" {
					g, t, err := ctx.Getter(curr, updated)
					return g, t, err
				} else {
					return string(curr), dataType, nil
				}
			} else {
				return ctx.Getter(data, pattern)
			}

		}

		s, err := n.loop.evaluate(&ASTContext{Data: ctx.Data, Getter: newGetter, ArrayEach: ctx.ArrayEach, ObjectEach: ctx.ObjectEach})
		i++

		if err != nil {
			panic(err.Error())
		}

		sb.WriteString(s)
	}

	ctx.ArrayEach(array, forEach)
	return sb.String()
}

// propFor is the execution of a forNode in the context of property iteration of an object.
// It evaluates the loop node for each property of the given object and returns the result.
func (n *forNode) propFor(ctx *ASTContext, object []byte) string {
	sb := strings.Builder{}

	forEach := func(prop string, val []byte, dataType ElementType) {

		newGetter := func(data []byte, pattern string) (string, ElementType, error) {

			if pattern == n.indexName {
				return prop, String, nil
			}

			updated, usesItem := strings.CutPrefix(pattern, n.itemName)

			if usesItem {
				if updated != "" {
					return ctx.Getter(val, updated)
				} else {
					return string(val), dataType, nil
				}
			} else {
				return ctx.Getter(data, pattern)
			}

		}

		s, err := n.loop.evaluate(&ASTContext{Data: ctx.Data, Getter: newGetter, ArrayEach: ctx.ArrayEach, ObjectEach: ctx.ObjectEach})

		if err != nil {
			panic(err.Error())
		}

		sb.WriteString(s)
	}

	ctx.ObjectEach(object, forEach)
	return sb.String()
}

// evaluate on a forNode checks which kind of for it is (range vs props) and
// performs the necessary loop, evaluating its loop node for each element in the iterable
// and returning the concatenation.
func (n *forNode) evaluate(ctx *ASTContext) (string, error) {
	logger.DefaultLogger.Node("For:", n.pattern)

	a, err := n.pattern.stringValue(ctx)

	if err != nil {
		return "", err
	}

	iterable := []byte(a)
	var loopString string

	if n.forType {
		loopString = n.rangeFor(ctx, iterable)
	} else {
		loopString = n.propFor(ctx, iterable)
	}
	return n.withChild(loopString, ctx)
}
