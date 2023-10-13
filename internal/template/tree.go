package parser

import (
	"errors"
	"fmt"
)

type node interface {
	evaluate(data []byte, getter Getter) (string, error)
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

func (b *baseNode) withChild(def string, data []byte, getter Getter) (string, error) {

	if b.child == nil {
		return def, nil
	} else {
		childText, err := b.child.evaluate(data, getter)
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

func (n *textNode) evaluate(data []byte, getter Getter) (string, error) {
	thisText := n.text
	result, err := n.withChild(thisText, data, getter)
	return result, err
}

type accessNode struct {
	baseNode
	accessPattern string
}

func (n accessNode) evaluate(data []byte, getter Getter) (string, error) {
	println("evaluating access", n.accessPattern)
	s, err := getter(data, n.accessPattern)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Access pattern %s is invalid", n.accessPattern))
	}
	result, err := n.withChild(s, data, getter)
	return result, err
}

type ifNode struct {
	baseNode
	condition   string
	trueClause  node
	falseClause node
}

func evalCondition(condition string, data []byte, getter Getter) (bool, error) {
	// TODO: implement condition evaluation
	return false, nil
}

func (n *ifNode) evaluate(data []byte, getter Getter) (string, error) {
	println("evaluating if")

	result, err := evalCondition(n.condition, data, getter)

	if err != nil {
		return "", nil
	}

	var clause string

	if result {
		clause, err = n.trueClause.evaluate(data, getter)
	} else {
		clause, err = n.falseClause.evaluate(data, getter)
	}

	if err != nil {
		return "", nil
	}

	res, err := n.withChild(clause, data, getter)
	return res, err

}
