// Copyright (c) 2026 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

type Node interface {
	GetPos() Position
	SetPos(Position)
}

type IntNode struct {
	Val string
	Pos Position
}

func NewIntNode(val string, pos Position) *IntNode {
	return &IntNode{
		Val: val,
		Pos: pos,
	}
}

func (n *IntNode) GetPos() Position {
	return n.Pos
}

func (n *IntNode) SetPos(pos Position) {
	n.Pos = pos
}

type FloatNode struct {
	Val string
	Pos Position
}

func NewFloatNode(val string, pos Position) *FloatNode {
	return &FloatNode{
		Val: val,
		Pos: pos,
	}
}

func (n *FloatNode) GetPos() Position {
	return n.Pos
}

func (n *FloatNode) SetPos(pos Position) {
	n.Pos = pos
}

type StringNode struct {
	Val string
	Pos Position
}

func NewStringNode(val string, pos Position) *StringNode {
	return &StringNode{
		Val: val,
		Pos: pos,
	}
}

func (n *StringNode) GetPos() Position {
	return n.Pos
}

func (n *StringNode) SetPos(pos Position) {
	n.Pos = pos
}

type SymbolNode struct {
	Val string
	Pos Position
}

func NewSymbolNode(val string, pos Position) *SymbolNode {
	return &SymbolNode{
		Val: val,
		Pos: pos,
	}
}

func (n *SymbolNode) GetPos() Position {
	return n.Pos
}

func (n *SymbolNode) SetPos(pos Position) {
	n.Pos = pos
}

type CallNode struct {
	Head any
	Args []any
	Pos  Position
}

func NewCallNode(head any, args []any, pos Position) *CallNode {
	return &CallNode{
		Head: head,
		Args: args,
		Pos:  pos,
	}
}

func (n *CallNode) GetPos() Position {
	return n.Pos
}

func (n *CallNode) SetPos(pos Position) {
	n.Pos = pos
}

type ListNode struct {
	Elems []any
	Pos   Position
}

func NewListNode(elems []any, pos Position) *ListNode {
	return &ListNode{
		Elems: elems,
		Pos:   pos,
	}
}

func (n *ListNode) GetPos() Position {
	return n.Pos
}

func (n *ListNode) SetPos(pos Position) {
	n.Pos = pos
}

type DictNode struct {
	Pairs []*KeyValPair
	Pos   Position
}

type KeyValPair struct {
	Key any
	Val any
}

func NewDictNode(pairs []*KeyValPair, pos Position) *DictNode {
	return &DictNode{
		Pairs: pairs,
		Pos:   pos,
	}
}

func NewKeyValPair(key any, val any) *KeyValPair {
	return &KeyValPair{
		Key: key,
		Val: val,
	}
}

func (n *DictNode) GetPos() Position {
	return n.Pos
}

func (n *DictNode) SetPos(pos Position) {
	n.Pos = pos
}

type QuoteNode struct {
	Arg any
	Pos Position
}

func NewQuoteNode(arg any, pos Position) *QuoteNode {
	return &QuoteNode{
		Arg: arg,
		Pos: pos,
	}
}

func (n *QuoteNode) GetPos() Position {
	return n.Pos
}

func (n *QuoteNode) SetPos(pos Position) {
	n.Pos = pos
}

type QuasiquoteNode struct {
	Arg any
	Pos Position
}

func NewQuasiquoteNode(arg any, pos Position) *QuasiquoteNode {
	return &QuasiquoteNode{
		Arg: arg,
		Pos: pos,
	}
}

func (n *QuasiquoteNode) GetPos() Position {
	return n.Pos
}

func (n *QuasiquoteNode) SetPos(pos Position) {
	n.Pos = pos
}

type UnquoteNode struct {
	Arg any
	Pos Position
}

func NewUnquoteNode(arg any, pos Position) *UnquoteNode {
	return &UnquoteNode{
		Arg: arg,
		Pos: pos,
	}
}

func (n *UnquoteNode) GetPos() Position {
	return n.Pos
}

func (n *UnquoteNode) SetPos(pos Position) {
	n.Pos = pos
}

type EOFNode struct {
	Pos Position
}

func NewEOFNode(pos Position) *EOFNode {
	return &EOFNode{
		Pos: pos,
	}
}

func (n *EOFNode) GetPos() Position {
	return n.Pos
}

func (n *EOFNode) SetPos(pos Position) {
	n.Pos = pos
}

func MapPos(expr any, fn func(Position) Position) {
	if node, ok := expr.(Node); ok {
		node.SetPos(fn(node.GetPos()))

		switch val := node.(type) {
		case *CallNode:
			MapPos(val.Head, fn)

			for _, arg := range val.Args {
				MapPos(arg, fn)
			}

		case *ListNode:
			for _, elem := range val.Elems {
				MapPos(elem, fn)
			}

		case *DictNode:
			for _, pair := range val.Pairs {
				MapPos(pair.Key, fn)
				MapPos(pair.Val, fn)
			}

		case *QuoteNode:
			MapPos(val.Arg, fn)

		case *QuasiquoteNode:
			MapPos(val.Arg, fn)

		case *UnquoteNode:
			MapPos(val.Arg, fn)
		}
	}
}
