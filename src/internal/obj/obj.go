package obj

type Vertice struct {
	X, Y, Z float64
}

type Face struct {
	Vertices []Vertice
}

type Object struct {
	Vertices []Vertice
	Faces    []Face
}
