# typact

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
