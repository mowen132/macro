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
	TokenUnquote
	TokenWhitespace
	TokenComment
	TokenNewline
	TokenEnd
)

func (t *Token) String() string {
	prefix := t.Pos.String()

	switch t.Kind {
	case TokenInt:
		return fmt.Sprintf("INT %s %v", prefix, t.Val)

	case TokenFloat:
		return fmt.Sprintf("FLT %s %v", prefix, t.Val)

	case TokenString:
		return fmt.Sprintf("STR %s %q", prefix, t.Val)

	case TokenSymbol:
		return fmt.Sprintf("SYM %s %q", prefix, t.Val)

	case TokenLeftParenthesis:
		return "LPA " + prefix

	case TokenRightParenthesis:
		return "RPA " + prefix

	case TokenLeftSquare:
		return "LSQ " + prefix

	case TokenRightSquare:
		return "RSQ " + prefix

	case TokenLeftCurly:
		return "LCU " + prefix

	case TokenRightCurly:
		return "RCU " + prefix

	case TokenQuote:
		return "QUO " + prefix

	case TokenUnquote:
		return "UNQ " + prefix

	case TokenWhitespace:
		return fmt.Sprintf("WHI %s %q", prefix, t.Val)

	case TokenComment:
		return fmt.Sprintf("CMT %s %q", prefix, t.Val)

	case TokenNewline:
		return "NEW " + prefix

	case TokenEnd:
		return "END " + prefix

	default:
		return fmt.Sprintf("UNK %s %v", prefix, t.Val)
	}
}
