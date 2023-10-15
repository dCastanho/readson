package parser

import (
	"errors"
	"fmt"

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
		return "", nil
	}

	var clause string

	if result {
		clause, err = n.trueClause.evaluate(ctx)
	} else {
		clause, err = n.falseClause.evaluate(ctx)
	}

	if err != nil {
		return "", nil
	}

	res, err := n.withChild(clause, ctx)
	return res, err

}
