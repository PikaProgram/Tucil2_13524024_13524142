package cli

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"tucil/src/internal/model"
	"tucil/src/internal/obj"
)

func Convert(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: convert <input.obj> <maxDepth> [output.obj]")
		return
	}

	startTime := time.Now()

	inputFile, maxDepthStr := args[0], args[1]
	outputFile := "output.obj"

	if len(args) > 2 {
		outputFile = args[2]
	}

	maxDepth, err := strconv.Atoi(maxDepthStr)
	if err != nil {
		fmt.Println("Invalid maxDepth value:", maxDepthStr)
		return
	}

	basepath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err.Error())
		return
	}

	inputFile, outputFile = basepath+"/"+inputFile, basepath+"/"+outputFile

	res, err := obj.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err.Error())
		return
	}

	object, err := obj.ParseOBJ(res)
	if err != nil {
		fmt.Println("Error parsing OBJ content:", err.Error())
		return
	}

	box, err := object.GetBoundingBox()
	if err != nil {
		fmt.Println("Error getting bounding box:", err.Error())
		return
	}

	rootCube, err := box.GetRootCube()
	if err != nil {
		fmt.Println("Error getting root cube:", err.Error())
		return
	}

	stats := &model.OctreeCount{
		NodesFormed: make(map[int]int),
		NodesPruned: make(map[int]int),
	}

	cubes, err := rootCube.SubDivideCube(0, maxDepth, object, stats)
	if err != nil {
		fmt.Println("Error subdividing cube:", err.Error())
		return
	}

	resultObject := model.CubesToOBJ(cubes)

	err = obj.WriteOBJToFile(outputFile, resultObject)
	if err != nil {
		fmt.Println("Error writing output file:", err.Error())
		return
	}

	duration := time.Since(startTime)

	fmt.Println("\n=== Voxelization Report ===")
	fmt.Printf("Banyaknya voxel yang terbentuk: %d\n", len(cubes))
	fmt.Printf("Banyaknya vertex yang terbentuk: %d\n", len(resultObject.Vertexes))
	fmt.Printf("Banyaknya faces yang terbentuk: %d\n", len(resultObject.Faces))

	fmt.Println("\nStatistik node octree yang terbentuk:")
	for i := 1; i <= maxDepth; i++ {
		fmt.Printf("%d: %d\n", i, stats.NodesFormed[i])
	}

	fmt.Println("\nStatistik node yang tidak perlu ditelusuri:")
	for i := 1; i <= maxDepth; i++ {
		fmt.Printf("%d: %d\n", i, stats.NodesPruned[i])
	}

	fmt.Printf("\nKedalaman octree: %d\n", maxDepth)
	fmt.Printf("Lama waktu program berjalan: %v\n", duration)
	fmt.Printf("Path dimana file .obj disimpan: %s\n", outputFile)
	fmt.Println("===========================\n")
}
