// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"bytes"
)

func Marshal(val any) ([]byte, error) {
	var b bytes.Buffer
	e := NewEncoder(&b)

	if err := e.Encode(val); err != nil {
		return nil, err
	}

	if err := e.Flush(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func Unmarshal(b []byte) (any, error) {
	d := NewDecoder(bytes.NewReader(b))
	return d.Decode()
}
