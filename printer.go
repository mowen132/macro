// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"bufio"
	"fmt"
	"io"
)

type Printer struct {
	writer *bufio.Writer
	pos    Position
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{
		writer: bufio.NewWriter(w),
		pos:    Position{1, 1},
	}
}

func (p *Printer) PrintToken(token *Token) error {
	switch token.Kind {
	case TokenInt:
		return p.PrintInt(token.Val)

	case TokenFloat:
		return p.PrintFloat(token.Val)

	case TokenString:
		return p.PrintString(token.Val)

	case TokenSymbol:
		return p.PrintSymbol(token.Val)

	case TokenLeftParenthesis:
		return p.PrintLeftParenthesis()

	case TokenRightParenthesis:
		return p.PrintRightParenthesis()

	case TokenLeftSquare:
		return p.PrintLeftSquare()

	case TokenRightSquare:
		return p.PrintRightSquare()

	case TokenLeftCurly:
		return p.PrintLeftCurly()

	case TokenRightCurly:
		return p.PrintRightCurly()

	case TokenQuote:
		return p.PrintQuote()

	case TokenQuasiquote:
		return p.PrintQuasiquote()

	case TokenUnquote:
		return p.PrintUnquote()

	case TokenWhitespace:
		return p.PrintWhitespace(token.Val)

	case TokenComment:
		return p.PrintComment(token.Val)

	case TokenNewline:
		return p.PrintNewline()
	}

	return fmt.Errorf("unknown token: %v", token.Kind)
}

func (p *Printer) PrintInt(val string) error {
	return p.writeString(val)
}

func (p *Printer) PrintFloat(val string) error {
	return p.writeString(val)
}

func (p *Printer) PrintString(val string) error {
	if err := p.writeRune('"'); err != nil {
		return err
	}

	for _, c := range val {
		var err error

		switch c {
		case '"':
			err = p.writeString("\\\"")

		case '\\':
			err = p.writeString("\\\\")

		case '\b':
			err = p.writeString("\\b")

		case '\f':
			err = p.writeString("\\f")

		case '\n':
			err = p.writeString("\\n")

		case '\r':
			err = p.writeString("\\r")

		case '\t':
			err = p.writeString("\\t")

		default:
			err = p.writeRune(c)
		}

		if err != nil {
			return err
		}
	}

	return p.writeRune('"')
}

func (p *Printer) PrintSymbol(val string) error {
	return p.writeString(val)
}

func (p *Printer) PrintLeftParenthesis() error {
	return p.writeRune('(')
}

func (p *Printer) PrintRightParenthesis() error {
	return p.writeRune(')')
}

func (p *Printer) PrintLeftSquare() error {
	return p.writeRune('[')
}

func (p *Printer) PrintRightSquare() error {
	return p.writeRune(']')
}

func (p *Printer) PrintLeftCurly() error {
	return p.writeRune('{')
}

func (p *Printer) PrintRightCurly() error {
	return p.writeRune('}')
}

func (p *Printer) PrintQuote() error {
	return p.writeRune('\'')
}

func (p *Printer) PrintQuasiquote() error {
	return p.writeRune('`')
}

func (p *Printer) PrintUnquote() error {
	return p.writeRune(',')
}

func (p *Printer) PrintWhitespace(val string) error {
	return p.writeString(val)
}

func (p *Printer) PrintComment(val string) error {
	if err := p.writeRune(';'); err != nil {
		return err
	}

	return p.writeString(val)
}

func (p *Printer) PrintNewline() error {
	if _, err := p.writer.WriteRune('\n'); err != nil {
		return err
	}

	p.pos.Line++
	p.pos.Col = 1
	return nil
}

func (p *Printer) Pos() Position {
	return p.pos
}

func (p *Printer) Flush() error {
	return p.writer.Flush()
}

func (p *Printer) writeString(s string) error {
	for _, c := range s {
		if err := p.writeRune(c); err != nil {
			return err
		}
	}

	return nil
}

func (p *Printer) writeRune(c rune) error {
	if _, err := p.writer.WriteRune(c); err != nil {
		return err
	}

	p.pos.Col++
	return nil
}
