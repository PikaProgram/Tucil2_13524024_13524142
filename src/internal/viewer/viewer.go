package viewer

import (
	"image/color"
	"log"
	"sync"
	"tucil/src/internal/model"
	"tucil/src/internal/obj"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Simulation struct {
	obj model.Object

	distance     float64
	rotationXRad float64
	rotationYRad float64

	prevMouseX int
	prevMouseY int
}

func (g *Simulation) GetObj(filename string) error {

	object, err := obj.ParseOBJ(filename)
	if err != nil {
		return err
	}
	g.obj = *object
	g.distance = 100.0
	return nil
}

func (g *Simulation) Update() error {
	// dy is positive when scrolling up, negative when scrolling down
	_, dy := ebiten.Wheel()
	g.distance -= dy

	// Mouse drag (rotation)
	cx, cy := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		deltaX := cx - g.prevMouseX
		deltaY := cy - g.prevMouseY

		g.rotationYRad += float64(deltaX) * 0.01
		g.rotationXRad += float64(deltaY) * 0.01
	}

	g.prevMouseX = cx
	g.prevMouseY = cy

	return nil
}

func (g *Simulation) Draw(screen *ebiten.Image) {
	green := color.RGBA{0, 255, 0, 255}
	matrix := Identity()
	matrix = matrix.Multiply(Perspective(0.1, 1000.0))
	matrix = matrix.Multiply(Translation(0, 0, -g.distance))
	matrix = matrix.Multiply(RotationX(g.rotationXRad))
	matrix = matrix.Multiply(RotationY(g.rotationYRad))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, face := range g.obj.Faces {
		wg.Add(1)
		go func(m Matrix4, f model.Face) {
			defer wg.Done()
			var v1, v2, v3 model.Vertex
			m.MultiplyVector(g.obj.Vertexes[f.Vertexes[0]], &v1)
			m.MultiplyVector(g.obj.Vertexes[f.Vertexes[1]], &v2)
			m.MultiplyVector(g.obj.Vertexes[f.Vertexes[2]], &v3)
			x1, y1 := 300+(v1.X*300), 300-(v1.Y*300)
			x2, y2 := 300+(v2.X*300), 300-(v2.Y*300)
			x3, y3 := 300+(v3.X*300), 300-(v3.Y*300)

			//Safeguard for excessive distance
			if !(x1 > -5000 && x1 < 5000 && y1 > -5000 && y1 < 5000 &&
				x2 > -5000 && x2 < 5000 && y2 > -5000 && y2 < 5000 &&
				x3 > -5000 && x3 < 5000 && y3 > -5000 && y3 < 5000) {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 2, green, true)
			vector.StrokeLine(screen, float32(x2), float32(y2), float32(x3), float32(y3), 2, green, true)
			vector.StrokeLine(screen, float32(x3), float32(y3), float32(x1), float32(y1), 2, green, true)
		}(matrix, face)

	}
	wg.Wait()
}

// Must have this function for the ebiten library
func (g *Simulation) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 600
}

func main() {
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("3D viewer")

	Simulation := &Simulation{}

	if err := ebiten.RunGame(Simulation); err != nil {
		log.Fatal(err)
	}
}
