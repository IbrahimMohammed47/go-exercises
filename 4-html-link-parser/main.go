package main

import (
	"fmt"
	"os"

	"linkparser"
)

func main() {
	// read file
	r, err := os.Open("examples/ex2.html")
	if err != nil {
		panic(err)
	}

	x := linkparser.ParseHtml(r)

	fmt.Println(x)
}
