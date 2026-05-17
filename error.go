// Copyright (c) 2026 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"fmt"
)

type ParseError struct {
	Message string
	Pos     Position
}

func NewParseError(message string, pos Position) *ParseError {
	return &ParseError{
		Message: message,
		Pos:     pos,
	}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s %s", e.Pos, e.Message)
}
