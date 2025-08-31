// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"fmt"
)

type Token struct {
	Kind TokenKind
	Val  any
	Pos  Position
}

type TokenKind int

const (
	IntToken TokenKind = iota
	FloatToken
	StringToken
	SymbolToken
	LeftParenthesisToken
	RightParenthesisToken
	LeftSquareToken
	RightSquareToken
	LeftCurlyToken
	RightCurlyToken
	QuoteToken
	UnquoteToken
	WhitespaceToken
	CommentToken
	NewlineToken
	EndToken
)

func (t *Token) AsInt() int {
	if v, ok := t.Val.(int); ok {
		return v
	}

	return 0
}

func (t *Token) AsFloat() float64 {
	if v, ok := t.Val.(float64); ok {
		return v
	}

	return 0.0
}

func (t *Token) AsString() string {
	if v, ok := t.Val.(string); ok {
		return v
	}

	return ""
}

func (t *Token) String() string {
	prefix := t.Pos.String()

	switch t.Kind {
	case IntToken:
		return fmt.Sprintf("INT %s %v", prefix, t.Val)

	case FloatToken:
		return fmt.Sprintf("FLT %s %v", prefix, t.Val)

	case StringToken:
		return fmt.Sprintf("STR %s %q", prefix, t.Val)

	case SymbolToken:
		return fmt.Sprintf("SYM %s %q", prefix, t.Val)

	case LeftParenthesisToken:
		return "LPA " + prefix

	case RightParenthesisToken:
		return "RPA " + prefix

	case LeftSquareToken:
		return "LSQ " + prefix

	case RightSquareToken:
		return "RSQ " + prefix

	case LeftCurlyToken:
		return "LCU " + prefix

	case RightCurlyToken:
		return "RCU " + prefix

	case QuoteToken:
		return "QUO " + prefix

	case UnquoteToken:
		return "UNQ " + prefix

	case WhitespaceToken:
		return fmt.Sprintf("WHI %s %q", prefix, t.Val)

	case CommentToken:
		return fmt.Sprintf("CMT %s %q", prefix, t.Val)

	case NewlineToken:
		return "NEW " + prefix

	case EndToken:
		return "END " + prefix

	default:
		return fmt.Sprintf("UNK %s %v", prefix, t.Val)
	}
}
