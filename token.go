// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"fmt"
)

type Token struct {
	Kind TokenKind
	Val  string
	Pos  Position
}

type TokenKind int

const (
	TokenBOF TokenKind = iota
	TokenInt
	TokenFloat
	TokenString
	TokenSymbol
	TokenLeftParenthesis
	TokenRightParenthesis
	TokenLeftSquare
	TokenRightSquare
	TokenLeftCurly
	TokenRightCurly
	TokenQuote
	TokenQuasiquote
	TokenUnquote
	TokenWhitespace
	TokenComment
	TokenNewline
	TokenEOF
)

var BOFToken = &Token{Kind: TokenBOF}

func NewToken(kind TokenKind, val string, pos Position) *Token {
	return &Token{
		Kind: kind,
		Val:  val,
		Pos:  pos,
	}
}

func (t *Token) String() string {
	switch t.Kind {
	case TokenBOF:
		return fmt.Sprintf("BOF")

	case TokenInt:
		return fmt.Sprintf("INT %v %s", t.Pos, t.Val)

	case TokenFloat:
		return fmt.Sprintf("FLT %v %s", t.Pos, t.Val)

	case TokenString:
		return fmt.Sprintf("STR %v %q", t.Pos, t.Val)

	case TokenSymbol:
		return fmt.Sprintf("SYM %v %s", t.Pos, t.Val)

	case TokenLeftParenthesis:
		return fmt.Sprintf("LPA %v", t.Pos)

	case TokenRightParenthesis:
		return fmt.Sprintf("RPA %v", t.Pos)

	case TokenLeftSquare:
		return fmt.Sprintf("LSQ %v", t.Pos)

	case TokenRightSquare:
		return fmt.Sprintf("RSQ %v", t.Pos)

	case TokenLeftCurly:
		return fmt.Sprintf("LCU %v", t.Pos)

	case TokenRightCurly:
		return fmt.Sprintf("RCU %v", t.Pos)

	case TokenQuote:
		return fmt.Sprintf("QUO %v", t.Pos)

	case TokenQuasiquote:
		return fmt.Sprintf("QQU %v", t.Pos)

	case TokenUnquote:
		return fmt.Sprintf("UNQ %v", t.Pos)

	case TokenWhitespace:
		return fmt.Sprintf("WHI %v %q", t.Pos, t.Val)

	case TokenComment:
		return fmt.Sprintf("CMT %v %q", t.Pos, t.Val)

	case TokenNewline:
		return fmt.Sprintf("NEW %v", t.Pos)

	case TokenEOF:
		return fmt.Sprintf("EOF %v", t.Pos)

	default:
		return fmt.Sprintf("UNK %v %q", t.Pos, t.Val)
	}
}
