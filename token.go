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
	TokenInt TokenKind = iota
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

func (t *Token) String() string {
	pos := t.Pos.String()

	switch t.Kind {
	case TokenInt:
		return fmt.Sprintf("INT %s %v", pos, t.Val)

	case TokenFloat:
		return fmt.Sprintf("FLT %s %v", pos, t.Val)

	case TokenString:
		return fmt.Sprintf("STR %s %q", pos, t.Val)

	case TokenSymbol:
		return fmt.Sprintf("SYM %s %q", pos, t.Val)

	case TokenLeftParenthesis:
		return "LPA " + pos

	case TokenRightParenthesis:
		return "RPA " + pos

	case TokenLeftSquare:
		return "LSQ " + pos

	case TokenRightSquare:
		return "RSQ " + pos

	case TokenLeftCurly:
		return "LCU " + pos

	case TokenRightCurly:
		return "RCU " + pos

	case TokenQuote:
		return "QUO " + pos

	case TokenQuasiquote:
		return "QQU " + pos

	case TokenUnquote:
		return "UNQ " + pos

	case TokenWhitespace:
		return fmt.Sprintf("WHI %s %q", pos, t.Val)

	case TokenComment:
		return fmt.Sprintf("CMT %s %q", pos, t.Val)

	case TokenNewline:
		return "NEW " + pos

	case TokenEOF:
		return "EOF " + pos

	default:
		return fmt.Sprintf("UNK %s %v", pos, t.Val)
	}
}
