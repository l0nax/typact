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

Benchmarking has also shown that implementing `std.Cloner[T]` with a pointer receiver results in better performance for
all use cases:

<details>
<summary>Benchmark Results</summary>

```
  goos: linux
goarch: amd64
pkg: go.l0nax.org/typact
cpu: AMD Ryzen 9 5900X 12-Core Processor
                                           │  /tmp/base   │
                                           │    sec/op    │
Option_Clone/None-24                         1.915n ±  1%
Option_Clone/String-24                       2.461n ±  1%
Option_Clone/Int64-24                        2.465n ±  1%
Option_Clone/CustomStructPointer-24          64.79n ±  3%
Option_Clone/CustomStructFallback-24         168.2n ±  9%
Option_Clone/ScalarSlice-24                  214.4n ±  4%

Option_Clone/PtrSliceWrapper_PtrRecv-24      1.537µ ±  4%
Option_Clone/SliceWrapper_PtrRecv-24         1.588µ ±  5%
Option_Clone/PtrSliceWrapper_NormalRecv-24   1.943µ ±  5%
Option_Clone/SliceWrapper_NormalRecv-24      1.887µ ± 15%
geomean                                      109.3n

                                           │  /tmp/base   │
                                           │     B/op     │
Option_Clone/None-24                         0.000 ± 0%
Option_Clone/String-24                       0.000 ± 0%
Option_Clone/Int64-24                        0.000 ± 0%
Option_Clone/CustomStructPointer-24          48.00 ± 0%
Option_Clone/CustomStructFallback-24         96.00 ± 0%
Option_Clone/ScalarSlice-24                  112.0 ± 0%

Option_Clone/PtrSliceWrapper_PtrRecv-24      328.0 ± 0%
Option_Clone/SliceWrapper_PtrRecv-24         408.0 ± 0%
Option_Clone/PtrSliceWrapper_NormalRecv-24   424.0 ± 0%
Option_Clone/SliceWrapper_NormalRecv-24      408.0 ± 0%
geomean                                                 ¹
¹ summaries must be >0 to compute geomean

                                           │  /tmp/base   │
                                           │  allocs/op   │
Option_Clone/None-24                         0.000 ± 0%
Option_Clone/String-24                       0.000 ± 0%
Option_Clone/Int64-24                        0.000 ± 0%
Option_Clone/CustomStructPointer-24          1.000 ± 0%
Option_Clone/CustomStructFallback-24         2.000 ± 0%
Option_Clone/ScalarSlice-24                  3.000 ± 0%

Option_Clone/PtrSliceWrapper_PtrRecv-24      11.00 ± 0%
Option_Clone/SliceWrapper_PtrRecv-24         11.00 ± 0%
Option_Clone/PtrSliceWrapper_NormalRecv-24   13.00 ± 0%
Option_Clone/SliceWrapper_NormalRecv-24      11.00 ± 0%
geomean                                                 ¹
¹ summaries must be >0 to compute geomean
```

</details>

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
