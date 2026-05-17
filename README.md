# macro

A lightweight and generic **S-expression** library for Go.

`macro` provides tools to **encode and decode S-expressions** without any evaluation or runtime. It is designed for:

- Serialization/Deserialization
- DSLs & Preprocessors

The library offers both **low-level token-based APIs** and **high-level AST-based APIs**, allowing you to choose the level of control you need.

---

## Installation

```bash
go get github.com/mowen132/macro
```

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/mowen132/macro"
)

func main() {
    ast, err := macro.Unmarshal([]byte("(define x 42)"))
    if err != nil {
        panic(err)
    }

    /* Put additional logic here! */

    b, err := macro.Marshal(ast)
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

---

## API Overview

`macro` provides two levels of API:

1. **Low-Level API** – Work directly with tokens (`Scanner`, `Printer`, `Token`, `Position`).
2. **High-Level API** – Encode and decode abstract syntax trees (`Decoder`, `Encoder`, `Marshal`, `Unmarshal`).

---

### Low-Level API

#### Scanner

Reads S-expression tokens from an `io.Reader`.

```go
s := macro.NewScanner(strings.NewReader("(foo 123)"))
for {
    tok, err := s.Scan()
    if err != nil {
        panic(err)
    }
    if tok.Kind == macro.TokenEOF {
        break
    }
    fmt.Println(tok)
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
p.PrintInt("123")
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

Represents a position:

```go
type Position struct {
    Line int
    Col  int
}
```

---

### High-Level API

#### Decoder

Decodes an S-expression into an abstract syntax tree.

```go
d := macro.NewDecoder(strings.NewReader("(1 2 3)"))
ast, err := d.Decode()
```

#### Encoder

Encodes an abstract syntax tree into an S-expressions.

```go
var b bytes.Buffer
e := macro.NewEncoder(&b)
e.Encode(ast)
e.Flush()
fmt.Println(b.String())
```

#### Marshal / Unmarshal

Convenience functions for one-shot encoding and decoding:

```go
b, err := macro.Marshal(ast)
```

```go
ast, err := macro.Unmarshal(b)
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
