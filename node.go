// Copyright (c) 2026 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

type Node interface {
	GetPos() *Position
}

type NodeInt struct {
	Val string
	Pos Position
}

func (n *NodeInt) GetPos() *Position {
	return &n.Pos
}

type NodeFloat struct {
	Val string
	Pos Position
}

func (n *NodeFloat) GetPos() *Position {
	return &n.Pos
}

type NodeString struct {
	Val string
	Pos Position
}

func (n *NodeString) GetPos() *Position {
	return &n.Pos
}

type NodeSymbol struct {
	Val string
	Pos Position
}

func (n *NodeSymbol) GetPos() *Position {
	return &n.Pos
}

type NodeCall struct {
	Head Node
	Args []Node
	Pos  Position
}

func (n *NodeCall) GetPos() *Position {
	return &n.Pos
}

type NodeList struct {
	Elems []Node
	Pos   Position
}

func (n *NodeList) GetPos() *Position {
	return &n.Pos
}

type NodeDict struct {
	Pairs []*KeyValPair
	Pos   Position
}

type KeyValPair struct {
	Key Node
	Val Node
}

func (n *NodeDict) GetPos() *Position {
	return &n.Pos
}

type NodeQuote struct {
	Val Node
	Pos Position
}

func (n *NodeQuote) GetPos() *Position {
	return &n.Pos
}

type NodeQuasiquote struct {
	Val Node
	Pos Position
}

func (n *NodeQuasiquote) GetPos() *Position {
	return &n.Pos
}

type NodeUnquote struct {
	Val Node
	Pos Position
}

func (n *NodeUnquote) GetPos() *Position {
	return &n.Pos
}

type NodeEOF struct {
	Pos Position
}

func (n *NodeEOF) GetPos() *Position {
	return &n.Pos
}

func Walk(expr Node, walkFn func(node Node)) {
	walkFn(expr)

	switch node := expr.(type) {
	case *NodeCall:
		Walk(node.Head, walkFn)

		for _, arg := range node.Args {
			Walk(arg, walkFn)
		}

	case *NodeList:
		for _, elem := range node.Elems {
			Walk(elem, walkFn)
		}

	case *NodeDict:
		for _, pair := range node.Pairs {
			Walk(pair.Key, walkFn)
			Walk(pair.Val, walkFn)
		}
	}
}
