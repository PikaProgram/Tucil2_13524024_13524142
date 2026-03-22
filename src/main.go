package main

import (
	"os"
	"tucil/src/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		println("Valid commands:")
		println("  convert <input.obj> <maxDepth> [output.obj]")
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "convert":
		cli.Convert(args)
	default:
		println("Unknown command:", cmd)
	}
}
