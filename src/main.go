package main

import (
	"os"
	"tucil/src/internal/obj"
)

func main() {
	args := os.Args[1:]

	res, err := obj.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	_, err = obj.ParseOBJ(res)
	if err != nil {
		panic(err)
	}
}
