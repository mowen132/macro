// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"fmt"
	"io"
)

type Decoder struct {
	scanner *Scanner
	token   *Token
	depth   int
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scanner: NewScanner(r),
		token:   BOFToken,
		depth:   0,
	}
}

func (d *Decoder) Decode() (Node, error) {
	switch d.token.Kind {
	case TokenBOF:
		if err := d.scan(); err != nil {
			return nil, err
		}

		return d.Decode()

	case TokenInt:
		token := d.token

		if err := d.scan(); err != nil {
			return nil, err
		}

		return NewIntNode(token.Val, token.Pos), nil

	case TokenFloat:
		token := d.token

		if err := d.scan(); err != nil {
			return nil, err
		}

		return NewFloatNode(token.Val, token.Pos), nil

	case TokenString:
		token := d.token

		if err := d.scan(); err != nil {
			return nil, err
		}

		return NewStringNode(token.Val, token.Pos), nil

	case TokenSymbol:
		token := d.token

		if err := d.scan(); err != nil {
			return nil, err
		}

		return NewSymbolNode(token.Val, token.Pos), nil

	case TokenLeftParenthesis:
		pos := d.token.Pos
		d.depth++

		if err := d.scan(); err != nil {
			return nil, err
		}

		head, err := d.Decode()

		if err != nil {
			return nil, err
		}

		args := make([]Node, 0)

		for d.token.Kind != TokenRightParenthesis {
			arg, err := d.Decode()

			if err != nil {
				return nil, err
			}

			args = append(args, arg)
		}

		if err := d.scan(); err != nil {
			return nil, err
		}

		d.depth--
		return NewCallNode(head, args, pos), nil

	case TokenLeftSquare:
		pos := d.token.Pos
		d.depth++

		if err := d.scan(); err != nil {
			return nil, err
		}

		elems := make([]Node, 0)

		for d.token.Kind != TokenRightSquare {
			elem, err := d.Decode()

			if err != nil {
				return nil, err
			}

			elems = append(elems, elem)
		}

		if err := d.scan(); err != nil {
			return nil, err
		}

		d.depth--
		return NewListNode(elems, pos), nil

	case TokenLeftCurly:
		pos := d.token.Pos
		d.depth++

		if err := d.scan(); err != nil {
			return nil, err
		}

		pairs := make([]*KeyValPair, 0)

		for d.token.Kind != TokenRightCurly {
			key, err := d.Decode()

			if err != nil {
				return nil, err
			}

			val, err := d.Decode()

			if err != nil {
				return nil, err
			}

			pairs = append(pairs, NewKeyValPair(key, val))
		}

		if err := d.scan(); err != nil {
			return nil, err
		}

		d.depth--
		return NewDictNode(pairs, pos), nil

	case TokenQuote:
		pos := d.token.Pos
		d.depth++

		if err := d.scan(); err != nil {
			return nil, err
		}

		arg, err := d.Decode()

		if err != nil {
			return nil, err
		}

		d.depth--
		return NewQuoteNode(arg, pos), nil

	case TokenQuasiquote:
		pos := d.token.Pos
		d.depth++

		if err := d.scan(); err != nil {
			return nil, err
		}

		arg, err := d.Decode()

		if err != nil {
			return nil, err
		}

		d.depth--
		return NewQuasiquoteNode(arg, pos), nil

	case TokenUnquote:
		pos := d.token.Pos
		d.depth++

		if err := d.scan(); err != nil {
			return nil, err
		}

		arg, err := d.Decode()

		if err != nil {
			return nil, err
		}

		d.depth--
		return NewUnquoteNode(arg, pos), nil

	case TokenEOF:
		if d.depth == 0 {
			return NewEOFNode(d.token.Pos), nil
		}
	}

	return nil, d.unexpected()
}

func (d *Decoder) scan() error {
	for {
		token, err := d.scanner.Scan()

		if err != nil {
			return err
		}

		switch token.Kind {
		case TokenWhitespace, TokenComment, TokenNewline:
		default:
			d.token = token
			return nil
		}
	}
}

func (d *Decoder) unexpected() error {
	return NewParseError(
		fmt.Sprintf("unexpected %s", describe(d.token.Kind)),
		d.token.Pos,
	)
}

func describe(kind TokenKind) string {
	switch kind {
	case TokenInt:
		return "int"

	case TokenFloat:
		return "float"

	case TokenString:
		return "string"

	case TokenSymbol:
		return "symbol"

	case TokenLeftParenthesis:
		return "("

	case TokenRightParenthesis:
		return ")"

	case TokenLeftSquare:
		return "["

	case TokenRightSquare:
		return "]"

	case TokenLeftCurly:
		return "{"

	case TokenRightCurly:
		return "}"

	case TokenQuote:
		return "'"

	case TokenQuasiquote:
		return "`"

	case TokenUnquote:
		return ","

	case TokenEOF:
		return "eof"
	}

	return "unknown"
}
