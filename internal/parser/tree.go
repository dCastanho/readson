package parser

import (
	md "dcastanho.readson/internal/template"
)

type node interface {
	evaluate(data []byte, getter md.Getter) string
	next() node
}

type baseNode struct {
	child node
}

func (b baseNode) next() node {
	return b.child
}

type textNode struct {
	baseNode
	text string
}

func (n textNode) evaluate(data []byte, getter md.Getter) string {
	return n.text
}

type accessNode struct {
	baseNode
	accessPattern string
}

func (n accessNode) evaluate(data []byte, getter md.Getter) string {
	s, err := getter(data, n.accessPattern)
	if err != nil {
		// TODO handle incorrect access
	}
	return s
}

type ifNode struct {
	baseNode
	condition   string
	trueClause  node
	falseClause node
}

func evalCondition(cond string, data []byte, getter md.Getter) bool {
	// strings.Split(cond, "=")
	return true
}

func (n ifNode) evaluate(data []byte, getter md.Getter) string {
	result := evalCondition(n.condition, data, getter)

	if result {
		return n.trueClause.evaluate(data, getter)
	} else {
		return n.falseClause.evaluate(data, getter)
	}

}
