package main

import (
	"log"
	"os"
	"tucil/src/internal/viewer"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("3D Viewer")
	if len(os.Args) < 2 {
		println("Valid commands:")
		println("Required argument: <path to file.obj>")
		return
	}
	objfile := os.Args[1]
	sim := &viewer.Simulation{}

	err := sim.GetObj(objfile)
	if err != nil {
		log.Fatal("Error loading OBJ file: ", err)
	}

	if err := ebiten.RunGame(sim); err != nil {
		log.Fatal("Engine crashed: ", err)
	}
}
