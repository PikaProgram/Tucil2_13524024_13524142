package cli

import (
	"os"
	"strconv"
	"tucil/src/internal/model"
	"tucil/src/internal/obj"
)

func Convert(args []string) {
	if len(args) < 2 {
		println("Usage: convert <input.obj> <maxDepth> [output.obj]")
		return
	}

	inputFile, maxDepthStr := args[0], args[1]

	outputFile := "output.obj"

	if len(args) > 2 {
		outputFile = args[2]
	}

	maxDepth, err := strconv.Atoi(maxDepthStr)
	if err != nil {
		println("Invalid maxDepth value:", maxDepthStr)
		return
	}

	basepath, err := os.Getwd()
	if err != nil {
		println("Error getting working directory:", err)
		return
	}

	inputFile, outputFile = basepath+"/"+inputFile, basepath+"/"+outputFile

	res, err := obj.ReadFile(inputFile)
	if err != nil {
		println("Error reading input file:", err)
		return
	}

	object, err := obj.ParseOBJ(res)
	if err != nil {
		println("Error parsing OBJ content:", err)
		return
	}

	box, err := object.GetBoundingBox()
	if err != nil {
		println("Error getting bounding box:", err)
		return
	}

	rootCube, err := box.GetRootCube()
	if err != nil {
		println("Error getting root cube:", err)
		return
	}

	cubes, err := rootCube.SubDivideCube(maxDepth, object)
	if err != nil {
		println("Error subdividing cube:", err)
		return
	}

	resultObject := model.CubesToOBJ(cubes)

	err = obj.WriteOBJToFile(outputFile, resultObject)
	if err != nil {
		println("Error writing output file:", err)
		return
	}

	println("Conversion completed successfully. Output written to", outputFile)

}
