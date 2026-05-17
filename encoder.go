// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"io"
)

type Encoder struct {
	printer *Printer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		printer: NewPrinter(w),
	}
}

func (e *Encoder) Encode(expr Node) error {
	switch node := expr.(type) {
	case *NodeInt:
		if err := e.printer.PrintInt(node.Val); err != nil {
			return err
		}

	case *NodeFloat:
		if err := e.printer.PrintFloat(node.Val); err != nil {
			return err
		}

	case *NodeString:
		if err := e.printer.PrintString(node.Val); err != nil {
			return err
		}

	case *NodeSymbol:
		if err := e.printer.PrintSymbol(node.Val); err != nil {
			return err
		}

	case *NodeCall:
		if err := e.printer.PrintLeftParenthesis(); err != nil {
			return err
		}

		if err := e.Encode(node.Head); err != nil {
			return err
		}

		for _, arg := range node.Args {
			if err := e.printer.PrintWhitespace(" "); err != nil {
				return err
			}

			if err := e.Encode(arg); err != nil {
				return err
			}
		}

		if err := e.printer.PrintRightParenthesis(); err != nil {
			return err
		}

	case *NodeList:
		if err := e.printer.PrintLeftSquare(); err != nil {
			return err
		}

		for i, elem := range node.Elems {
			if i > 0 {
				if err := e.printer.PrintWhitespace(" "); err != nil {
					return err
				}
			}

			if err := e.Encode(elem); err != nil {
				return err
			}
		}

		if err := e.printer.PrintRightSquare(); err != nil {
			return err
		}

	case *NodeDict:
		if err := e.printer.PrintLeftCurly(); err != nil {
			return err
		}

		for i, pair := range node.Pairs {
			if i > 0 {
				if err := e.printer.PrintWhitespace(" "); err != nil {
					return err
				}
			}

			if err := e.Encode(pair.Key); err != nil {
				return err
			}

			if err := e.printer.PrintWhitespace(" "); err != nil {
				return err
			}

			if err := e.Encode(pair.Val); err != nil {
				return err
			}
		}

		if err := e.printer.PrintRightCurly(); err != nil {
			return err
		}

	case *NodeQuote:
		if err := e.printer.PrintQuote(); err != nil {
			return err
		}

		if err := e.Encode(node.Val); err != nil {
			return err
		}

	case *NodeQuasiquote:
		if err := e.printer.PrintQuasiquote(); err != nil {
			return err
		}

		if err := e.Encode(node.Val); err != nil {
			return err
		}

	case *NodeUnquote:
		if err := e.printer.PrintUnquote(); err != nil {
			return err
		}

		if err := e.Encode(node.Val); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) Flush() error {
	return e.printer.Flush()
}
