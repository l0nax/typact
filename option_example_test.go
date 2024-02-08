package typact_test

import (
	"fmt"
	"slices"

	"go.l0nax.org/typact"
)

func assertEq[T comparable](src T, expected T) {
	if src != expected {
		panic(fmt.Sprintf("expected %v, got %v", expected, src))
	}
}

func ExampleOption_Insert() {
	opt := typact.None[int]()
	val := opt.Insert(5)

	fmt.Println(*val)
	fmt.Println(opt.Unwrap())

	*val = 3
	fmt.Println(opt.Unwrap())

	// Output:
	// 5
	// 5
	// 3
}

func ExampleOption_Unwrap_some() {
	x := typact.Some("foo")
	fmt.Println(x.Unwrap())

	// Output:
	// foo
}

func ExampleOption_Unwrap_none() {
	// WARN: This function panics!
	x := typact.None[string]()
	fmt.Println(x.Unwrap())
}

func ExampleOption_UnwrapOr() {
	fmt.Println(typact.Some("foo").UnwrapOr("bar"))
	fmt.Println(typact.None[string]().UnwrapOr("world"))

	// Output:
	// foo
	// world
}

func ExampleOption_Inspect() {
	x := typact.Some("foo")
	x.Inspect(func(val string) {
		fmt.Println(val)
	})

	// Does print nothing
	y := typact.None[string]()
	y.Inspect(func(val string) {
		fmt.Println(val)
	})

	// Output:
	// foo
}

func ExampleOption_Deconstruct() {
	x := typact.Some("foo")

	val, ok := x.Deconstruct()
	fmt.Printf("%q, %t\n", val, ok)

	y := typact.None[string]()
	val, ok = y.Deconstruct()
	fmt.Printf("%q, %t\n", val, ok)

	// Output:
	// "foo", true
	// "", false
}

func ExampleOption_GetOrInsertWith() {
	x := typact.None[int]()

	{
		y := x.GetOrInsertWith(func() int { return 2 })
		fmt.Println(*y)

		*y = 20
	}

	fmt.Println(x.Unwrap())

	// Output:
	// 2
	// 20
}

func ExampleOption_IsSomeAnd() {
	x := typact.Some("foo")

	ok := x.IsSomeAnd(func(str string) bool {
		return str == "foo"
	})
	fmt.Printf("%t\n", ok)

	ok = x.IsSomeAnd(func(str string) bool {
		return str == ""
	})
	fmt.Printf("%t\n", ok)

	y := typact.None[string]()
	ok = y.IsSomeAnd(func(str string) bool {
		return str != ""
	})
	fmt.Printf("%t\n", ok)

	// Output:
	// true
	// false
	// false
}

func ExampleOption_Filter_slice() {
	x := []typact.Option[string]{
		typact.Some("foo"),
		typact.None[string](),
		typact.Some("bar"),
		typact.Some("baz"),
		typact.None[string](),
		typact.Some("hello"),
		typact.Some("world"),
	}

	x = slices.DeleteFunc(x, func(val typact.Option[string]) bool {
		return val.IsSomeAnd(IsNotZero[string])
	})
	fmt.Println(x)

	// Output:
	// [foo bar baz hello world]
}

func IsZero[T comparable](val T) bool {
	var zero T
	return val == zero
}

func IsNotZero[T comparable](val T) bool {
	var zero T
	return val != zero
}
