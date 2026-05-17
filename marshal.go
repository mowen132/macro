// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"bytes"
)

func Marshal(node Node) ([]byte, error) {
	var b bytes.Buffer
	e := NewEncoder(&b)

	if err := e.Encode(node); err != nil {
		return nil, err
	}

	if err := e.Flush(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func Unmarshal(b []byte) (Node, error) {
	return NewDecoder(bytes.NewReader(b)).Decode()
}
