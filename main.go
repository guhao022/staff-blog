package main

import (
	"staff/listener/blog"
	"staff/tools/env"

	"github.com/num5/axiom"
)

func main() {
	_, err := env.Load()
	if err != nil {
		panic(err)
	}

	b := axiom.New()
	b.AddAdapter(axiom.NewShell(b))
	b.Register(&blog.BlogListener{})

	b.Start()
}
