// Copyright (c) 2026 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

type Node interface {
	Pos() *Position
}

type BaseNode struct {
	pos Position
}

func (n *BaseNode) Pos() *Position {
	return &n.pos
}

type IntNode struct {
	BaseNode
	Val string
}

func NewIntNode(val string, pos Position) *IntNode {
	return &IntNode{
		BaseNode: BaseNode{pos},
		Val:      val,
	}
}

type FloatNode struct {
	BaseNode
	Val string
}

func NewFloatNode(val string, pos Position) *FloatNode {
	return &FloatNode{
		BaseNode: BaseNode{pos},
		Val:      val,
	}
}

type StringNode struct {
	BaseNode
	Val string
}

func NewStringNode(val string, pos Position) *StringNode {
	return &StringNode{
		BaseNode: BaseNode{pos},
		Val:      val,
	}
}

type SymbolNode struct {
	BaseNode
	Val string
}

func NewSymbolNode(val string, pos Position) *SymbolNode {
	return &SymbolNode{
		BaseNode: BaseNode{pos},
		Val:      val,
	}
}

type CallNode struct {
	BaseNode
	Head Node
	Args []Node
}

func NewCallNode(head Node, args []Node, pos Position) *CallNode {
	return &CallNode{
		BaseNode: BaseNode{pos},
		Head:     head,
		Args:     args,
	}
}

type ListNode struct {
	BaseNode
	Elems []Node
}

func NewListNode(elems []Node, pos Position) *ListNode {
	return &ListNode{
		BaseNode: BaseNode{pos},
		Elems:    elems,
	}
}

type DictNode struct {
	BaseNode
	Pairs []*KeyValPair
}

type KeyValPair struct {
	Key Node
	Val Node
}

func NewDictNode(pairs []*KeyValPair, pos Position) *DictNode {
	return &DictNode{
		BaseNode: BaseNode{pos},
		Pairs:    pairs,
	}
}

func NewKeyValPair(key Node, val Node) *KeyValPair {
	return &KeyValPair{
		Key: key,
		Val: val,
	}
}

type QuoteNode struct {
	BaseNode
	Arg Node
}

func NewQuoteNode(arg Node, pos Position) *QuoteNode {
	return &QuoteNode{
		BaseNode: BaseNode{pos},
		Arg:      arg,
	}
}

type QuasiquoteNode struct {
	BaseNode
	Arg Node
}

func NewQuasiquoteNode(arg Node, pos Position) *QuasiquoteNode {
	return &QuasiquoteNode{
		BaseNode: BaseNode{pos},
		Arg:      arg,
	}
}

type UnquoteNode struct {
	BaseNode
	Arg Node
}

func NewUnquoteNode(arg Node, pos Position) *UnquoteNode {
	return &UnquoteNode{
		BaseNode: BaseNode{pos},
		Arg:      arg,
	}
}

type EOFNode struct {
	BaseNode
}

func NewEOFNode(pos Position) *EOFNode {
	return &EOFNode{
		BaseNode: BaseNode{pos},
	}
}

func WalkPos(expr Node, fn func(*Position)) {
	fn(expr.Pos())

	switch val := expr.(type) {
	case *CallNode:
		WalkPos(val.Head, fn)

		for _, arg := range val.Args {
			WalkPos(arg, fn)
		}

	case *ListNode:
		for _, elem := range val.Elems {
			WalkPos(elem, fn)
		}

	case *DictNode:
		for _, pair := range val.Pairs {
			WalkPos(pair.Key, fn)
			WalkPos(pair.Val, fn)
		}

	case *QuoteNode:
		WalkPos(val.Arg, fn)

	case *QuasiquoteNode:
		WalkPos(val.Arg, fn)

	case *UnquoteNode:
		WalkPos(val.Arg, fn)
	}
}
