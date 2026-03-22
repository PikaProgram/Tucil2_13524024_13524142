package model

import "errors"

type Vertex struct {
	X, Y, Z float64
}

type Face struct {
	Vertexes [3]int
}

type Object struct {
	Faces    []Face
	Vertexes []Vertex
}

type Box struct {
	Min, Max, Center Vertex
}

type Cube struct {
	Box
	Size float64
}

func (o *Object) GetBoundingBox() (Box, error) {
	if len(o.Vertexes) == 0 {
		return Box{}, errors.New("object has no vertexes")
	}

	min := o.Vertexes[0]
	max := o.Vertexes[0]

	for _, v := range o.Vertexes {
		for _, c := range []struct {
			Min *float64
			Max *float64
			Val float64
		}{
			{Min: &min.X, Max: &max.X, Val: v.X},
			{Min: &min.Y, Max: &max.Y, Val: v.Y},
			{Min: &min.Z, Max: &max.Z, Val: v.Z},
		} {
			if c.Val < *c.Min {
				*c.Min = c.Val
			}
			if c.Val > *c.Max {
				*c.Max = c.Val
			}
		}
	}

	center := Vertex{
		X: (min.X + max.X) / 2,
		Y: (min.Y + max.Y) / 2,
		Z: (min.Z + max.Z) / 2,
	}

	return Box{
		Min:    min,
		Max:    max,
		Center: center,
	}, nil
}

func (b *Box) GetRootCube() (Cube, error) {
	size := max(max(b.Max.X-b.Min.X, b.Max.Y-b.Min.Y), b.Max.Z-b.Min.Z)
	return Cube{
		Box: Box{
			Min: Vertex{
				X: b.Center.X - size/2,
				Y: b.Center.Y - size/2,
				Z: b.Center.Z - size/2,
			},
			Max: Vertex{
				X: b.Center.X + size/2,
				Y: b.Center.Y + size/2,
				Z: b.Center.Z + size/2,
			},
			Center: b.Center,
		},
		Size: size,
	}, nil
}
