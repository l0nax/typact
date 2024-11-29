# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).


## 0.3.0 (2024-11-29)

### Added (9 changes)
- Add `CloneWith` method to `Option[T]` type
- Add `IsZero` method to `Option[T]` to support [yaml](https://pkg.go.dev/gopkg.in/yaml.v3\#Marshal) `null` value representation
- Add `Take` method to Option type
- Add `std/randx` random helper package
- Add `xslices.FillValues` function to fill a value pattern
- Add `xslices.Fill` function to fill a slice very fast
- Add experimental `iterop` package to ease working with iterators
- Add new experimental package `immutable` with `List` type
- Add new package `exp/xslices` which provides additional helper to the `slices` std package

### Fixed (1 change)
- Fix `UnmarshalText` and `MarshalText` to correctly handle scalar types

### Other (2 changes)
- Initial release
- The `Option[T].Clone()` method now fallbacks to pointer-receiver calling on custom types, if available instead of panicing

