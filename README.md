# macro

A lightweight and generic **S-expression (sexp)** library for Go.

`macro` provides tools to **encode and decode S-expressions** without any evaluation or runtime. It is designed for:

- **Serialization/Deserialization** – Store and retrieve structured data in a simple format.
- **DSLs & Preprocessors** – Build lightweight languages or configuration layers on top of S-expressions.

The library offers both **low-level token-based APIs** and **high-level encoder/decoder APIs**, allowing you to choose the level of control you need.

---

## Installation

```bash
go get github.com/mowen132/macro
```

---

## Quick Start

### Encode a Value

```go
package main

import (
    "fmt"
    "github.com/mowen132/macro"
)

func main() {
    data := []any{
        macro.Symbol("define"),
        macro.Symbol("x"),
        42,
    }

    b, err := macro.Marshal(data)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(b))
}
```

Output:

```
(define x 42)
```

### Decode a Value

```go
package main

import (
    "fmt"
    "github.com/mowen132/macro"
)

func main() {
    input := []byte("(define x 42)")

    val, err := macro.Unmarshal(input)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%#v\n", val)
}
```

---

## API Overview

`macro` provides two levels of API:

1. **Low-Level API** – Work directly with tokens (`Scanner`, `Printer`, `Token`, `Position`).
2. **High-Level API** – Encode and decode Go values (`Decoder`, `Encoder`, `Symbol`, `Marshal`, `Unmarshal`).

---

### Low-Level API

#### Scanner

Reads S-expression tokens from an `io.Reader`.

```go
s := macro.NewScanner(strings.NewReader("(foo 123)"))
for {
    tok, err := s.Scan()
    if err != nil {
        log.Fatal(err)
    }
    if tok.Kind == macro.EndToken {
        break
    }
    fmt.Printf("Token: %v (value: %v)\n", tok.Kind, tok.Val)
}
```

#### Printer

Writes tokens to an `io.Writer`.

```go
var b bytes.Buffer
p := macro.NewPrinter(&b)
p.PrintLeftParenthesis()
p.PrintSymbol("foo")
p.PrintWhitespace(" ")
p.PrintInt(123)
p.PrintRightParenthesis()
p.Flush()

fmt.Println(b.String()) // (foo 123)
```

#### Token

Represents a single token:

```go
type Token struct {
    Kind TokenKind
    Val  string
    Pos  Position
}
```

#### Position

Represents a token’s position in the input.

```go
type Position struct {
    Line int
    Col  int
}
```

---

### High-Level API

#### Decoder

Decodes an S-expression stream into Go values.

```go
d := macro.NewDecoder(strings.NewReader("(1 2 3)"))
val, err := d.Decode()
// val will be []any{1, 2, 3}
```

#### Encoder

Encodes Go values into S-expressions.

```go
var b bytes.Buffer
e := macro.NewEncoder(&b)
e.Encode([]any{macro.Symbol("foo"), 123})
e.Flush()
fmt.Println(b.String()) // (foo 123)
```

#### Symbol

Represents a Lisp-like symbol:

```go
type Symbol string
```

#### Marshal / Unmarshal

Convenience functions for one-shot encoding and decoding:

```go
b, _ := macro.Marshal([]any{macro.Symbol("bar"), 99})
val, _ := macro.Unmarshal(b)
```

---

## Roadmap

- [ ] Full grammar specification
- [ ] Extended examples and idiomatic usage patterns
- [ ] Additional integration utilities

---

## License

This project is licensed under the [MIT License](LICENSE).

## Commit Message Convention

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages. Please refer to the official documentation for guidelines on how to format your commits.
