// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"fmt"
	"io"
	"strconv"
)

type Encoder struct {
	printer *Printer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		printer: NewPrinter(w),
	}
}

func (e *Encoder) Encode(val any) error {
	var err error

	switch v := val.(type) {
	case int:
		err = e.printer.PrintInt(strconv.Itoa(v))

	case float64:
		err = e.printer.PrintFloat(strconv.FormatFloat(v, 'g', -1, 64))

	case string:
		err = e.printer.PrintString(v)

	case Symbol:
		err = e.printer.PrintSymbol(string(v))

	case []any:
		err = e.encodeList(v)

	default:
		err = fmt.Errorf("unsupported type: %T", val)
	}

	return err
}

func (e *Encoder) encodeList(list []any) error {
	p := e.printer

	if len(list) > 0 {
		if v, ok := list[0].(Symbol); ok {
			switch v {
			case "list":
				return e.encodeDelimitedList(list[1:], p.PrintLeftSquare, p.PrintRightSquare)

			case "dict":
				return e.encodeDelimitedList(list[1:], p.PrintLeftCurly, p.PrintRightCurly)

			case "quote":
				return e.encodeQuoted(list[1:], p.PrintQuote, "quote")

			case "unquote":
				return e.encodeQuoted(list[1:], p.PrintUnquote, "unquote")
			}
		}
	}

	return e.encodeDelimitedList(list, p.PrintLeftParenthesis, p.PrintRightParenthesis)
}

func (e *Encoder) encodeDelimitedList(list []any, printLeft, printRight func() error) error {
	if err := printLeft(); err != nil {
		return err
	}

	for i, v := range list {
		if i > 0 {
			if err := e.printer.PrintWhitespace(" "); err != nil {
				return err
			}
		}

		if err := e.Encode(v); err != nil {
			return err
		}
	}

	return printRight()
}

func (e *Encoder) encodeQuoted(list []any, print func() error, name string) error {
	if len(list) != 1 {
		return fmt.Errorf("wrong number of arguments in %s", name)
	}

	if err := print(); err != nil {
		return err
	}

	return e.Encode(list[0])
}

func (e *Encoder) Flush() error {
	return e.printer.Flush()
}
