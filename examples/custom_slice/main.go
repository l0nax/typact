package main

import (
	"fmt"
	"time"

	"go.l0nax.org/typact"
	"go.l0nax.org/typact/std"
)

type MyData struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Clone implements [std.Cloner].
func (m *MyData) Clone() *MyData {
	return &MyData{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

var _ std.Cloner[*MyData] = (*MyData)(nil)

// MyDataSlice is a helper type implementing [std.Cloner]
// for [MyData].
type MyDataSlice []*MyData

func (m MyDataSlice) Clone() MyDataSlice {
	cpy := make(MyDataSlice, len(m))
	for i := 0; i < len(m); i++ {
		cpy[i] = m[i].Clone()
	}

	return cpy
}

func main() {
	raw := []*MyData{
		{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	opt := typact.Some(MyDataSlice(raw))

	cpy := opt.Clone().Unwrap()
	fmt.Printf("%+v\n", cpy)
	fmt.Printf("At index 0: %+v", cpy[0])
}
