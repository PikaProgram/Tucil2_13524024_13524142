package viewer

import (
	"image/color"
	"log"
	"math"
	"runtime"
	"sync"
	"tucil/src/internal/model"
	"tucil/src/internal/obj"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Simulation struct {
	obj model.Object

	distance     float64
	rotationXRad float64
	rotationYRad float64

	prevMouseX     int
	prevMouseY     int
	cullingEnabled bool
}

func (g *Simulation) GetObj(filename string) error {

	object, err := obj.ParseOBJ(filename)
	if err != nil {
		return err
	}
	g.obj = *object
	g.distance = 35.0
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
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.cullingEnabled = !g.cullingEnabled
	}
	g.prevMouseX = cx
	g.prevMouseY = cy

	return nil
}

func (g *Simulation) Draw(screen *ebiten.Image) {
	matrix := Identity()
	matrix = matrix.Multiply(Perspective(0.1, 1000.0))
	matrix = matrix.Multiply(Translation(0, 0, -g.distance))
	matrix = matrix.Multiply(RotationX(g.rotationXRad))
	matrix = matrix.Multiply(RotationY(g.rotationYRad))

	var wg sync.WaitGroup
	var mu sync.Mutex
	numWorkers := runtime.NumCPU()

	if len(g.obj.Faces) == 0 {
		return
	}
	chunkSize := (len(g.obj.Faces) + numWorkers - 1) / numWorkers

	whiteSubImage := ebiten.NewImage(1, 1)
	whiteSubImage.Fill(color.White)

	var batchedVertices [][]ebiten.Vertex
	var batchedIndices [][]uint16
	currentV := make([]ebiten.Vertex, 0, 65532)
	currentI := make([]uint16, 0, 65532)
	var baseIndex uint16 = 0

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(g.obj.Faces) {
			end = len(g.obj.Faces)
		}

		// Skip if there are more CPU cores than faces
		if start >= len(g.obj.Faces) {
			continue
		}

		wg.Add(1)

		go func(start, end int, m Matrix4) {
			defer wg.Done()

			for j := start; j < end; j++ {
				f := g.obj.Faces[j]

				var v1, v2, v3 model.Vertex
				m.MultiplyVector(g.obj.Vertexes[f.Vertexes[0]], &v1)
				m.MultiplyVector(g.obj.Vertexes[f.Vertexes[1]], &v2)
				m.MultiplyVector(g.obj.Vertexes[f.Vertexes[2]], &v3)
				x1, y1 := 300+(v1.X*300), 300-(v1.Y*300)
				x2, y2 := 300+(v2.X*300), 300-(v2.Y*300)
				x3, y3 := 300+(v3.X*300), 300-(v3.Y*300)

				if g.cullingEnabled {
					//This culling algorithm is VERY funny and VERY buggy...
					rotY := math.Mod(math.Abs(g.rotationYRad), 2*math.Pi)
					isLookingAtBack := rotY > math.Pi/2 && rotY < 3*math.Pi/2

					cross := (x2-x1)*(y3-y1) - (y2-y1)*(x3-x1)

					if isLookingAtBack {
						if cross < 0 {
							continue
						}
					} else {
						if cross > 0 {
							continue
						}
					}
				}

				if !(x1 > -5000 && x1 < 5000 && y1 > -5000 && y1 < 5000 &&
					x2 > -5000 && x2 < 5000 && y2 > -5000 && y2 < 5000 &&
					x3 > -5000 && x3 < 5000 && y3 > -5000 && y3 < 5000) {
					continue
				}

				mu.Lock()

				if baseIndex >= 65532 {
					batchedVertices = append(batchedVertices, currentV)
					batchedIndices = append(batchedIndices, currentI)
					currentV = make([]ebiten.Vertex, 0, 65532)
					currentI = make([]uint16, 0, 65532)
					baseIndex = 0
				}

				currentV = append(currentV,
					ebiten.Vertex{DstX: float32(x1), DstY: float32(y1), ColorR: 0, ColorG: 1, ColorB: 0, ColorA: 1},
					ebiten.Vertex{DstX: float32(x2), DstY: float32(y2), ColorR: 0, ColorG: 1, ColorB: 0, ColorA: 1},
					ebiten.Vertex{DstX: float32(x3), DstY: float32(y3), ColorR: 0, ColorG: 1, ColorB: 0, ColorA: 1},
				)
				currentI = append(currentI, baseIndex, baseIndex+1, baseIndex+2)
				baseIndex += 3
				mu.Unlock()
			}
		}(start, end, matrix)

	}
	wg.Wait()

	if len(currentV) > 0 {
		batchedVertices = append(batchedVertices, currentV)
		batchedIndices = append(batchedIndices, currentI)
	}

	for i := range batchedVertices {
		screen.DrawTriangles(batchedVertices[i], batchedIndices[i], whiteSubImage, nil)
	}
	status := "OFF"
	if g.cullingEnabled {
		status = "ON (visual glitches in some angles but improves performance)"
	}
	ebitenutil.DebugPrint(screen, "Backface Culling (Press C): "+status+"\nScroll down to zoom out\nScroll up to zoom in\nClick and drag to rotate")
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
