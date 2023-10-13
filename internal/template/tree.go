package parser

import (
	"errors"
)

func Execute(t any, data []byte, getter Getter) (string, error) {

	if top, ok := t.(node); ok {
		return top.evaluate(data, getter), nil
	} else {
		return "", errors.New("Not a node")
	}

}

type node interface {
	evaluate(data []byte, getter Getter) string
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

func (b *baseNode) withChild(def string, data []byte, getter Getter) string {
	var res string
	if b.child == nil {
		res = def
	} else {
		childText := b.child.evaluate(data, getter)
		res = def + childText
	}
	return res
}

type textNode struct {
	baseNode
	text string
}

func (n *textNode) evaluate(data []byte, getter Getter) string {
	println("evaluating text", n.text)
	thisText := n.text
	result := n.withChild(thisText, data, getter)
	return result
}

type accessNode struct {
	baseNode
	accessPattern string
}

func (n accessNode) evaluate(data []byte, getter Getter) string {
	println("evaluating access", n.accessPattern)
	s, err := getter(data, n.accessPattern)
	if err != nil {
		// TODO handle incorrect access
	}
	result := n.withChild(s, data, getter)
	return result
}

type ifNode struct {
	baseNode
	condition   string
	trueClause  node
	falseClause node
}

func (n *ifNode) evaluate(data []byte, getter Getter) string {
	println("evaluating if", n.condition)
	// result := evalCondition(n.condition, data, getter)
	// TODO: implement condition and false clause
	temp := n.trueClause.evaluate(data, getter)
	result := n.withChild(temp, data, getter)
	return result
}
