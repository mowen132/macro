// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	bof = -1
	eof = -2
)

type Scanner struct {
	reader *bufio.Reader
	char   rune
	pos    Position
	buf    strings.Builder
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(r),
		char:   bof,
		pos:    Position{1, 0},
	}
}

func (s *Scanner) Scan() (*Token, error) {
	switch s.char {
	case bof:
		if err := s.read(); err != nil {
			return nil, err
		}

		return s.Scan()

	case '+', '-':
		pos := s.pos

		if err := s.consume(); err != nil {
			return nil, err
		}

		switch s.char {
		case '0':
			return s.scanZero(pos)

		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return s.scanDigit(pos)

		case '!', '#', '$', '%', '&', '*', '+', '-', '/', ':',
			'<', '=', '>', '?', '@', 'A', 'B', 'C', 'D', 'E',
			'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
			'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y',
			'Z', '\\', '^', '_', 'a', 'b', 'c', 'd', 'e', 'f',
			'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
			'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'|', '~':

			return s.scanSymbol(pos)

		case '.':
			return s.scanDot(pos)

		case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
			return &Token{SymbolToken, s.extract(), pos}, nil

		default:
			if unicode.IsLetter(s.char) {
				return s.scanSymbol(pos)
			}

			return nil, s.errorUnexpectedf("%q in symbol", s.char)
		}

	case '0':
		return s.scanZero(s.pos)

	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.scanDigit(s.pos)

	case '"':
		pos := s.pos

		for {
			if err := s.read(); err != nil {
				return nil, err
			}

			switch s.char {
			case '"':
				if err := s.read(); err != nil {
					return nil, err
				}

				switch s.char {
				case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
					return &Token{StringToken, s.extract(), pos}, nil

				default:
					return nil, s.errorUnexpectedf("%q after closing '\"'", s.char)
				}

			case '\\':
				if err := s.read(); err != nil {
					return nil, err
				}

				switch s.char {
				case '"':
					s.buf.WriteRune('"')

				case '\\':
					s.buf.WriteRune('\\')

				case 'b':
					s.buf.WriteRune('\b')

				case 'f':
					s.buf.WriteRune('\f')

				case 'n':
					s.buf.WriteRune('\n')

				case 'r':
					s.buf.WriteRune('\r')

				case 't':
					s.buf.WriteRune('\t')

				case eof:
					return nil, s.errorUnexpectedf("eof in escape sequence")

				default:
					return nil, s.errorUnexpectedf("%q in escape sequence", s.char)
				}

			case '\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\a', '\b', '\n',
				'\v', '\f', '\r', '\x0e', '\x0f', '\x10', '\x11', '\x12', '\x13', '\x14',
				'\x15', '\x16', '\x17', '\x18', '\x19', '\x1a', '\x1b', '\x1c', '\x1d', '\x1e',
				'\x1f', '\x7f':

				return nil, s.errorUnexpectedf("%q in string", s.char)

			case eof:
				return nil, s.errorUnexpectedf("eof in string")

			default:
				s.buf.WriteRune(s.char)
			}
		}

	case '`':
		pos := s.pos

		for {
			if err := s.read(); err != nil {
				return nil, err
			}

			switch s.char {
			case '`':
				if err := s.read(); err != nil {
					return nil, err
				}

				switch s.char {
				case '`':
					s.buf.WriteRune('`')

				case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
					return &Token{StringToken, s.extract(), pos}, nil

				default:
					return nil, s.errorUnexpectedf("%q after closing '`'", s.char)
				}

			case '\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\a', '\b', '\v',
				'\f', '\r', '\x0e', '\x0f', '\x10', '\x11', '\x12', '\x13', '\x14', '\x15',
				'\x16', '\x17', '\x18', '\x19', '\x1a', '\x1b', '\x1c', '\x1d', '\x1e', '\x1f',
				'\x7f':

				return nil, s.errorUnexpectedf("%q in raw string", s.char)

			case eof:
				return nil, s.errorUnexpectedf("eof in raw string")

			default:
				s.buf.WriteRune(s.char)
			}
		}

	case '!', '#', '$', '%', '&', '*', '/', ':', '<', '=',
		'>', '?', '@', 'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q',
		'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '\\',
		'^', '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
		'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
		's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '|', '~':

		return s.scanSymbol(s.pos)

	case '.':
		return s.scanDot(s.pos)

	case '(':
		return s.scanSingle(LeftParenthesisToken)

	case ')':
		return s.scanSingleTerm(RightParenthesisToken)

	case '[':
		return s.scanSingle(LeftSquareToken)

	case ']':
		return s.scanSingleTerm(RightSquareToken)

	case '{':
		return s.scanSingle(LeftCurlyToken)

	case '}':
		return s.scanSingleTerm(RightCurlyToken)

	case '\'':
		return s.scanSingle(QuoteToken)

	case ',':
		return s.scanSingle(UnquoteToken)

	case '\t', ' ':
		pos := s.pos

		for {
			if err := s.consume(); err != nil {
				return nil, err
			}

			switch s.char {
			case '\t', ' ':
				continue

			default:
				return &Token{WhitespaceToken, s.extract(), pos}, nil
			}
		}

	case ';':
		pos := s.pos

		for {
			if err := s.read(); err != nil {
				return nil, err
			}

			switch s.char {
			case '\n', '\r', eof:
				return &Token{CommentToken, s.extract(), pos}, nil

			case '\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\a', '\b', '\v',
				'\f', '\x0e', '\x0f', '\x10', '\x11', '\x12', '\x13', '\x14', '\x15', '\x16',
				'\x17', '\x18', '\x19', '\x1a', '\x1b', '\x1c', '\x1d', '\x1e', '\x1f', '\x7f':

				return nil, s.errorUnexpectedf("%q in comment", s.char)

			default:
				s.buf.WriteRune(s.char)
			}
		}

	case '\n':
		return s.scanSingle(NewlineToken)

	case '\r':
		if err := s.read(); err != nil {
			return nil, err
		}

		switch s.char {
		case '\n':
			return s.scanSingle(NewlineToken)

		case eof:
			return nil, s.errorUnexpectedf("eof after '\r'")

		default:
			return nil, s.errorUnexpectedf("%q after '\r'", s.char)
		}

	case eof:
		return &Token{EndToken, nil, s.pos}, nil

	default:
		if unicode.IsLetter(s.char) {
			return s.scanSymbol(s.pos)
		}

		return nil, s.errorUnexpectedf("%q", s.char)
	}
}

func (s *Scanner) scanZero(pos Position) (*Token, error) {
	if err := s.consume(); err != nil {
		return nil, err
	}

	switch s.char {
	case '.':
		return s.scanDecimal(pos)

	case 'e', 'E':
		return s.scanExponent(pos)

	case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
		s.buf.Reset()
		return &Token{IntToken, 0, pos}, nil

	default:
		return nil, s.errorUnexpectedf("%q after '0'", s.char)
	}
}

func (s *Scanner) scanDigit(pos Position) (*Token, error) {
	for {
		if err := s.consume(); err != nil {
			return nil, err
		}

		switch s.char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue

		case '.':
			return s.scanDecimal(pos)

		case 'e', 'E':
			return s.scanExponent(pos)

		case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
			val, _ := strconv.Atoi(s.extract())
			return &Token{IntToken, val, pos}, nil

		default:
			return nil, s.errorUnexpectedf("%q after digit", s.char)
		}
	}
}

func (s *Scanner) scanDecimal(pos Position) (*Token, error) {
	if err := s.consume(); err != nil {
		return nil, err
	}

	switch s.char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		for {
			if err := s.consume(); err != nil {
				return nil, err
			}

			switch s.char {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				continue

			case 'e', 'E':
				return s.scanExponent(pos)

			case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
				val, _ := strconv.ParseFloat(s.extract(), 64)
				return &Token{FloatToken, val, pos}, nil

			default:
				return nil, s.errorUnexpectedf("%q in decimal", s.char)
			}
		}

	case eof:
		return nil, s.errorUnexpectedf("eof after '.'")

	default:
		return nil, s.errorUnexpectedf("%q after '.'", s.char)
	}
}

func (s *Scanner) scanExponent(pos Position) (*Token, error) {
	if err := s.consume(); err != nil {
		return nil, err
	}

	switch s.char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if err := s.consume(); err != nil {
			return nil, err
		}

	case '+', '-':
		if err := s.consume(); err != nil {
			return nil, err
		}

		switch s.char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if err := s.consume(); err != nil {
				return nil, err
			}

		case eof:
			return nil, s.errorUnexpectedf("eof after exponent sign")

		default:
			return nil, s.errorUnexpectedf("%q after exponent sign", s.char)
		}

	case eof:
		return nil, s.errorUnexpectedf("eof after exponent")

	default:
		return nil, s.errorUnexpectedf("%q after exponent", s.char)
	}

	for {
		switch s.char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if err := s.consume(); err != nil {
				return nil, err
			}

		case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
			val, _ := strconv.ParseFloat(s.extract(), 64)
			return &Token{FloatToken, val, pos}, nil

		default:
			return nil, s.errorUnexpectedf("%q in exponent", s.char)
		}
	}
}

func (s *Scanner) scanSymbol(pos Position) (*Token, error) {
	for {
		if err := s.consume(); err != nil {
			return nil, err
		}

		switch s.char {
		case '!', '#', '$', '%', '&', '*', '+', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', '<', '=', '>', '?', '@', 'A', 'B', 'C', 'D',
			'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
			'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
			'Y', 'Z', '\\', '^', '_', 'a', 'b', 'c', 'd', 'e',
			'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
			'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
			'z', '|', '~':

			continue

		case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
			return &Token{SymbolToken, s.extract(), pos}, nil

		default:
			if unicode.IsLetter(s.char) || unicode.IsDigit(s.char) {
				continue
			}

			return nil, s.errorUnexpectedf("%q in symbol", s.char)
		}
	}
}

func (s *Scanner) scanDot(pos Position) (*Token, error) {
	if err := s.consume(); err != nil {
		return nil, err
	}

	switch s.char {
	case '!', '#', '$', '%', '&', '*', '+', '-', '.', '/',
		':', '<', '=', '>', '?', '@', 'A', 'B', 'C', 'D',
		'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
		'Y', 'Z', '\\', '^', '_', 'a', 'b', 'c', 'd', 'e',
		'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
		'z', '|', '~':

		return s.scanSymbol(pos)

	case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
		return &Token{SymbolToken, s.extract(), pos}, nil

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return nil, s.errorUnexpectedf("digit after '.'")

	default:
		if unicode.IsLetter(s.char) {
			return s.scanSymbol(pos)
		}

		return nil, s.errorUnexpectedf("%q in symbol", s.char)
	}
}

func (s *Scanner) scanSingle(kind TokenKind) (*Token, error) {
	pos := s.pos

	if err := s.read(); err != nil {
		return nil, err
	}

	return &Token{kind, nil, pos}, nil
}

func (s *Scanner) scanSingleTerm(kind TokenKind) (*Token, error) {
	char := s.char
	pos := s.pos

	if err := s.read(); err != nil {
		return nil, err
	}

	switch s.char {
	case ')', ']', '}', ' ', '\t', ';', '\n', '\r', eof:
		return &Token{kind, nil, pos}, nil

	default:
		return nil, s.errorUnexpectedf("%q after %q", s.char, char)
	}
}

func (s *Scanner) read() error {
	c, _, err := s.reader.ReadRune()

	if err != nil {
		if err != io.EOF {
			return err
		}

		c = eof
	}

	s.char = c

	if c == '\n' {
		s.pos.Line++
		s.pos.Col = 0
	} else {
		s.pos.Col++
	}

	return nil
}

func (s *Scanner) consume() error {
	s.buf.WriteRune(s.char)
	return s.read()
}

func (s *Scanner) extract() string {
	val := s.buf.String()
	s.buf.Reset()
	return val
}

func (s *Scanner) errorUnexpectedf(format string, args ...any) error {
	return fmt.Errorf("%s unexpected %s", s.pos, fmt.Sprintf(format, args...))
}
