// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package macro

import (
	"bytes"
)

func Marshal(expr any) ([]byte, error) {
	var b bytes.Buffer
	encoder := NewEncoder(&b)

	if err := encoder.Encode(expr); err != nil {
		return nil, err
	}

	if err := encoder.Flush(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func Unmarshal(b []byte) (any, error) {
	return NewDecoder(bytes.NewReader(b)).Decode()
}
