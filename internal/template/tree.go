package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"dcastanho.readson/internal/logger"
)

type node interface {
	evaluate(ctx *ASTContext) (string, error)
	next() node
	setNext(n node)
}

type baseNode struct {
	child node
}

func (b *baseNode) next() node {
	return b.child
}

func (b *baseNode) setNext(n node) {
	b.child = n
}

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

type textNode struct {
	baseNode
	text string
}

func (n *textNode) evaluate(ctx *ASTContext) (string, error) {
	thisText := n.text

	logger.DefaultLogger.Node("Text:", thisText)

	result, err := n.withChild(thisText, ctx)
	return result, err
}

type accessNode struct {
	baseNode
	accessPattern string
}

func (n accessNode) evaluate(ctx *ASTContext) (string, error) {

	logger.DefaultLogger.Node("Variable Access:", n.accessPattern)

	s, _, err := ctx.Getter(ctx.Data, n.accessPattern)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Access pattern '%s' is invalid", n.accessPattern))
	}
	result, err := n.withChild(s, ctx)
	return result, err
}

type ifNode struct {
	baseNode
	condition   condition
	trueClause  node
	falseClause node
}

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

type forNode struct {
	baseNode
	itemName  string
	indexName string
	pattern   string
	loop      node
	forType   bool
}

func (n *forNode) rangeFor(ctx *ASTContext, sb *strings.Builder, array []byte) {
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
}

func (n *forNode) propFor(ctx *ASTContext, sb *strings.Builder, object []byte) {

	forEach := func(prop string, val []byte) {

		newGetter := func(data []byte, pattern string) (string, ElementType, error) {

			if pattern == n.indexName {
				return prop, String, nil
			}

			updated, usesItem := strings.CutPrefix(pattern, n.itemName)

			if usesItem {
				return ctx.Getter(val, updated)
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
}

func (n *forNode) evaluate(ctx *ASTContext) (string, error) {
	// sb := strings.Builder{}
	logger.DefaultLogger.Node("For:", n.forType)

	a, _, err := ctx.Getter(ctx.Data, n.pattern)

	if err != nil {
		return "", err
	}

	iterable := []byte(a)
	sb := strings.Builder{}

	if n.forType {
		n.rangeFor(ctx, &sb, iterable)
	} else {
		n.propFor(ctx, &sb, iterable)
	}
	loopString := sb.String()
	return n.withChild(loopString, ctx)
}
