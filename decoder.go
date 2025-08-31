// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"fmt"
	"io"
)

type scopeType int

const (
	docScope scopeType = iota
	listScope
	listLiteralScope
	dictLiteralScope
	quoteScope
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
	return d.decode(docScope)
}

func (d *Decoder) decode(scope scopeType) (any, error) {
	for {
		tok, err := d.scanner.Scan()

		if err != nil {
			return nil, err
		}

		switch tok.Kind {
		case IntToken:
			return tok.AsInt(), nil

		case FloatToken:
			return tok.AsFloat(), nil

		case StringToken:
			return tok.AsString(), nil

		case SymbolToken:
			return Symbol(tok.AsString()), nil

		case LeftParenthesisToken:
			return d.decodeList(listScope, []any{})

		case RightParenthesisToken:
			return nil, d.checkEndDelimiter(scope == listScope, tok.Pos, ")")

		case LeftSquareToken:
			return d.decodeList(listLiteralScope, []any{Symbol("list")})

		case RightSquareToken:
			return nil, d.checkEndDelimiter(scope == listLiteralScope, tok.Pos, "]")

		case LeftCurlyToken:
			return d.decodeList(dictLiteralScope, []any{Symbol("dict")})

		case RightCurlyToken:
			return nil, d.checkEndDelimiter(scope == dictLiteralScope, tok.Pos, "}")

		case QuoteToken:
			return d.decodeQuoted("quote")

		case UnquoteToken:
			return d.decodeQuoted("unquote")

		case EndToken:
			return nil, d.checkEndDelimiter(scope == docScope, tok.Pos, "eof")
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
	val, err := d.decode(quoteScope)

	if err != nil {
		return nil, err
	}

	return []any{Symbol(name), val}, nil
}
