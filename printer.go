// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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

func (p *Printer) PrintToken(tok *Token) error {
	switch tok.Kind {
	case TokenInt:
		return p.PrintInt(tok.Val)

	case TokenFloat:
		return p.PrintFloat(tok.Val)

	case TokenString:
		return p.PrintString(tok.Val)

	case TokenSymbol:
		return p.PrintSymbol(tok.Val)

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

	case TokenUnquote:
		return p.PrintUnquote()

	case TokenWhitespace:
		return p.PrintWhitespace(tok.Val)

	case TokenComment:
		return p.PrintComment(tok.Val)

	case TokenNewline:
		return p.PrintNewline()
	}

	return fmt.Errorf("unsupported token kind: %v", tok.Kind)
}

func (p *Printer) PrintInt(val string) error {
	return p.writeString(val)
}

func (p *Printer) PrintFloat(val string) error {
	return p.writeString(val)
}

func (p *Printer) PrintString(val string) error {
	var s strings.Builder
	s.Grow(len(val) + 2)
	s.WriteByte('"')

	for _, c := range val {
		switch c {
		case '"':
			s.WriteString("\\\"")

		case '\\':
			s.WriteString("\\\\")

		case '\b':
			s.WriteString("\\b")

		case '\f':
			s.WriteString("\\f")

		case '\n':
			s.WriteString("\\n")

		case '\r':
			s.WriteString("\\r")

		case '\t':
			s.WriteString("\\t")

		default:
			s.WriteRune(c)
		}
	}

	s.WriteByte('"')
	return p.writeString(s.String())
}

func (p *Printer) PrintStringRaw(val string) error {
	var s strings.Builder
	s.Grow(len(val) + 2)
	s.WriteByte('`')
	p.pos.Col++

	for _, c := range val {
		s.WriteRune(c)

		if c == '\n' {
			p.pos.Line++
			p.pos.Col = 1
		} else {
			p.pos.Col++
		}
	}

	s.WriteByte('`')
	p.pos.Col++
	_, err := p.writer.WriteString(s.String())
	return err
}

func (p *Printer) PrintSymbol(val string) error {
	return p.writeString(val)
}

func (p *Printer) PrintLeftParenthesis() error {
	return p.writeByte('(')
}

func (p *Printer) PrintRightParenthesis() error {
	return p.writeByte(')')
}

func (p *Printer) PrintLeftSquare() error {
	return p.writeByte('[')
}

func (p *Printer) PrintRightSquare() error {
	return p.writeByte(']')
}

func (p *Printer) PrintLeftCurly() error {
	return p.writeByte('{')
}

func (p *Printer) PrintRightCurly() error {
	return p.writeByte('}')
}

func (p *Printer) PrintQuote() error {
	return p.writeByte('\'')
}

func (p *Printer) PrintUnquote() error {
	return p.writeByte(',')
}

func (p *Printer) PrintWhitespace(val string) error {
	return p.writeString(val)
}

func (p *Printer) PrintComment(val string) error {
	if err := p.writeByte(';'); err != nil {
		return err
	}

	return p.writeString(val)
}

func (p *Printer) PrintNewline() error {
	p.pos.Line++
	p.pos.Col = 1
	return p.writer.WriteByte('\n')
}

func (p *Printer) Pos() Position {
	return p.pos
}

func (p *Printer) Flush() error {
	return p.writer.Flush()
}

func (p *Printer) writeString(s string) error {
	p.pos.Col += len(s)
	_, err := p.writer.WriteString(s)
	return err
}

func (p *Printer) writeByte(c byte) error {
	p.pos.Col++
	return p.writer.WriteByte(c)
}
