package option_test

import (
	"fmt"

	"go.l0nax.org/typact"
	"go.l0nax.org/typact/std/option"
)

func ExampleMap_some() {
	x := typact.Some("Hello, World!")
	ll := option.Map(x, func(s string) int { return len(s) })
	fmt.Printf("%t, %v\n", x.IsSome(), ll.Unwrap())

	// Output:
	// true, 13
}

func ExampleMap_none() {
	x := typact.None[string]()
	ll := option.Map(x, func(s string) int { return len(s) })
	fmt.Printf("%t, %v\n", x.IsSome(), ll.UnwrapOr(-1))

	// Output:
	// false, -1
}

func ExampleMapOr_some() {
	x := typact.Some("Hello, World!")
	ll := option.MapOr(x, func(s string) int { return len(s) }, 100)
	fmt.Println(ll)

	// Output:
	// 13
}

func ExampleMapOr_none() {
	x := typact.None[string]()
	ll := option.MapOr(x, func(s string) int { return len(s) }, 100)
	fmt.Println(ll)

	// Output:
	// 100
}

func ExampleMapOrElse_some() {
	x := typact.Some("Hello, World!")
	ll := option.MapOrElse(x,
		func(s string) int { return len(s) },
		func() int { return 100 },
	)
	fmt.Println(ll)

	// Output:
	// 13
}

func ExampleMapOrElse_none() {
	x := typact.None[string]()
	ll := option.MapOrElse(x,
		func(s string) int { return len(s) },
		func() int { return 100 },
	)
	fmt.Println(ll)

	// Output:
	// 100
}
