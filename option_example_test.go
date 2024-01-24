package typact_test

import (
	"fmt"

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
