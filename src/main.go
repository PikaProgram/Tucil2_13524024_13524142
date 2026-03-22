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

	object, err := obj.ParseOBJ(res)
	if err != nil {
		panic(err)
	}

	box, err := object.GetBoundingBox()
	if err != nil {
		panic(err)
	}
	println("Bounding Box:")
	println("Min:", box.Min.X, box.Min.Y, box.Min.Z)
	println("Max:", box.Max.X, box.Max.Y, box.Max.Z)
	println("Center:", box.Center.X, box.Center.Y, box.Center.Z)

	rootCube, err := box.GetRootCube()
	if err != nil {
		panic(err)
	}
	println("Root Cube:")
	println("Min:", rootCube.Min.X, rootCube.Min.Y, rootCube.Min.Z)
	println("Max:", rootCube.Max.X, rootCube.Max.Y, rootCube.Max.Z)
	println("Center:", rootCube.Center.X, rootCube.Center.Y, rootCube.Center.Z)
	println("Size:", rootCube.Size)
}
