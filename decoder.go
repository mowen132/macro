// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"fmt"
	"io"
)

const tokenBOF TokenKind = -1

type Decoder struct {
	scanner *Scanner
	token   *Token
	depth   int
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scanner: NewScanner(r),
		token:   &Token{Kind: tokenBOF},
		depth:   0,
	}
}

func (d *Decoder) Decode() (Node, error) {
	switch d.token.Kind {
	case tokenBOF:
		if err := d.scan(); err != nil {
			return nil, err
		}

		return d.Decode()

	case TokenInt:
		token := d.token

		if err := d.scan(); err != nil {
			return nil, err
		}

		return &NodeInt{
			Val: token.Val,
			Pos: token.Pos,
		}, nil

	case TokenFloat:
		token := d.token

		if err := d.scan(); err != nil {
			return nil, err
		}

		return &NodeFloat{
			Val: token.Val,
			Pos: token.Pos,
		}, nil

	case TokenString:
		token := d.token

		if err := d.scan(); err != nil {
			return nil, err
		}

		return &NodeString{
			Val: token.Val,
			Pos: token.Pos,
		}, nil

	case TokenSymbol:
		token := d.token

		if err := d.scan(); err != nil {
			return nil, err
		}

		return &NodeSymbol{
			Val: token.Val,
			Pos: token.Pos,
		}, nil

	case TokenLeftParenthesis:
		d.depth++
		pos := d.token.Pos

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

		return &NodeCall{
			Head: head,
			Args: args,
			Pos:  pos,
		}, nil

	case TokenLeftSquare:
		d.depth++
		pos := d.token.Pos

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

		return &NodeList{
			Elems: elems,
			Pos:   pos,
		}, nil

	case TokenLeftCurly:
		d.depth++
		pos := d.token.Pos

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

			pairs = append(pairs, &KeyValPair{
				Key: key,
				Val: val,
			})
		}

		if err := d.scan(); err != nil {
			return nil, err
		}

		d.depth--

		return &NodeDict{
			Pairs: pairs,
			Pos:   pos,
		}, nil

	case TokenQuote:
		d.depth++
		pos := d.token.Pos

		if err := d.scan(); err != nil {
			return nil, err
		}

		val, err := d.Decode()

		if err != nil {
			return nil, err
		}

		d.depth--

		return &NodeQuote{
			Val: val,
			Pos: pos,
		}, nil

	case TokenQuasiquote:
		d.depth++
		pos := d.token.Pos

		if err := d.scan(); err != nil {
			return nil, err
		}

		val, err := d.Decode()

		if err != nil {
			return nil, err
		}

		d.depth--

		return &NodeQuasiquote{
			Val: val,
			Pos: pos,
		}, nil

	case TokenUnquote:
		d.depth++
		pos := d.token.Pos

		if err := d.scan(); err != nil {
			return nil, err
		}

		val, err := d.Decode()

		if err != nil {
			return nil, err
		}

		d.depth--

		return &NodeUnquote{
			Val: val,
			Pos: pos,
		}, nil

	case TokenEOF:
		if d.depth == 0 {
			return &NodeEOF{
				Pos: d.token.Pos,
			}, nil
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

func (d *Decoder) unexpected() *ParseError {
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
