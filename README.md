# typact

[![Go Report Card](https://goreportcard.com/badge/go.l0nax.org/typact)](https://goreportcard.com/report/go.l0nax.org/typact)
[![Godoc Reference](https://pkg.go.dev/badge/go.l0nax.org/typact.svg)](https://pkg.go.dev/go.l0nax.org/typact)
[![Latest Release](https://gitlab.com/l0nax/typact/-/badges/release.svg)](https://gitlab.com/l0nax/typact/-/releases)
[![Coverage](https://gitlab.com/l0nax/typact/badges/master/coverage.svg)](https://gitlab.com/l0nax/typact/-/commits/master)
[![License](https://img.shields.io/gitlab/license/l0nax%2Ftypact)](./LICENSE)

A **zero dependency** type action (typact) library for GoLang.

<br />

**WARNING:** This library has not reached v1 yet!

## Installation

```bash
go get go.l0nax.org/typact
```

## Examples

Examples can be found in the [`examples`](./examples/) directory or at at the official [Godoc](https://pkg.go.dev/go.l0nax.org/typact) page.

## Usage

### Implementing `std.Cloner[T]` for `Option[T]`

**:warning: NOTE:** This feature is marked as `Unstable`!
<br />

Because of some limitations of the GoLang type system, there are some important facts to remember when
using the `Clone` method.

#### Custom Structs: Pointer Receiver

The best way to implement the `Clone` method is by using a pointer receiver like this:
```go
type MyData struct {
  ID        int
  CreatedAt time.Time
  UpdatedAt time.Time
}

// Clone implements the std.Cloner interface.
func (m *MyData) Clone() *MyData {
  return &MyData{
    ID:        m.ID,
    CreatedAt: m.CreatedAt,
    UpdatedAt: m.UpdatedAt,
  }
}
```

This allows you to use your type with and without a pointer.
For example:
```go
var asPtr typact.Option[*MyData]
var asVal typact.Option[MyData]
```

Calling `Clone` on `asPtr` _and_ `asVal` will result in a valid clone. This is because the `Clone`
implementation of the `Option[T]` type checks if a type implements the `std.Cloner` interface with a
pointer receiver.

For now it is not possible to implement the `Clone` method without a pointer receiver and use `Clone`
on a pointer value:
```go
type MyData struct {
  ID        int
  CreatedAt time.Time
  UpdatedAt time.Time
}

// Clone implements the std.Cloner interface.
func (m MyData) Clone() MyData {
  return MyData{
    ID:        m.ID,
    CreatedAt: m.CreatedAt,
    UpdatedAt: m.UpdatedAt,
  }
}

func main() {
  data := typact.Some(&MyData{
    ID:        15,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  })
  data.Clone() // This will panic because *MyData does not implement std.Cloner.
}
```

Thus the best way to support all use-cases is to implement the interface with a pointer receiver.

## Motivation

I've created this library because for one option types are really useful and prevent the _one billion dollar mistake_
with dereferencing nil pointers (at least it reduces the risk).

At my work and within my private projects I often find myself in the position where
- values may be provided (by the programmer, i.e. options)
- values may be `NULL` (i.e. in a database/ JSON/ ...)

Whilst writing my own LSM implementation I frequently ran into the situation where specific checks
are much easier to read when using a declarative approach (see the `cmpop.Ordering` type).
Thus this library is not a _Option type_ only library.

## License

The project is licensed under the [_MIT License_](./LICENSE).
