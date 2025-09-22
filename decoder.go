// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"fmt"
	"io"
	"strconv"
)

type scopeType int

const (
	scopeDoc scopeType = iota
	scopeList
	scopeListLiteral
	scopeDictLiteral
	scopeQuote
)

type Decoder struct {
	scanner *Scanner
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scanner: NewScanner(r),
	}
}

func (d *Decoder) Decode() (any, error) {
	return d.decode(scopeDoc)
}

func (d *Decoder) decode(scope scopeType) (any, error) {
	for {
		tok, err := d.scanner.Scan()

		if err != nil {
			return nil, err
		}

		switch tok.Kind {
		case TokenInt:
			val, err := strconv.Atoi(tok.Val)
			return val, err

		case TokenFloat:
			val, err := strconv.ParseFloat(tok.Val, 64)
			return val, err

		case TokenString:
			return tok.Val, nil

		case TokenSymbol:
			return Symbol(tok.Val), nil

		case TokenLeftParenthesis:
			return d.decodeList(scopeList, []any{})

		case TokenRightParenthesis:
			return nil, d.checkEndDelimiter(scope == scopeList, tok.Pos, ")")

		case TokenLeftSquare:
			return d.decodeList(scopeListLiteral, []any{Symbol("list")})

		case TokenRightSquare:
			return nil, d.checkEndDelimiter(scope == scopeListLiteral, tok.Pos, "]")

		case TokenLeftCurly:
			return d.decodeList(scopeDictLiteral, []any{Symbol("dict")})

		case TokenRightCurly:
			return nil, d.checkEndDelimiter(scope == scopeDictLiteral, tok.Pos, "}")

		case TokenQuote:
			return d.decodeQuoted("quote")

		case TokenQuasiquote:
			return d.decodeQuoted("quasiquote")

		case TokenUnquote:
			return d.decodeQuoted("unquote")

		case TokenEnd:
			return nil, d.checkEndDelimiter(scope == scopeDoc, tok.Pos, "eof")
		}
	}
}

func (d *Decoder) decodeList(scope scopeType, list []any) ([]any, error) {
	for {
		val, err := d.decode(scope)

		if err != nil {
			return nil, err
		}

		if val != nil {
			list = append(list, val)
		} else {
			return list, nil
		}
	}
}

func (d *Decoder) checkEndDelimiter(expected bool, pos Position, delim string) error {
	if !expected {
		return fmt.Errorf("%s unexpected %s", pos, delim)
	}

	return nil
}

func (d *Decoder) decodeQuoted(name string) (any, error) {
	val, err := d.decode(scopeQuote)

	if err != nil {
		return nil, err
	}

	return []any{Symbol(name), val}, nil
}
