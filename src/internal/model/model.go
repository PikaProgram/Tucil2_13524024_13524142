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

func (c *Cube) DivideCube() ([8]Cube, error) {
	var cubes [8]Cube
	half := c.Size / 2
	for i := 0; i < 8; i++ {
		offset := Vertex{
			X: half * float64(i&1),
			Y: half * float64((i>>1)&1),
			Z: half * float64((i>>2)&1),
		}
		cubes[i] = Cube{
			Box: Box{
				Min: Vertex{
					X: c.Min.X + offset.X,
					Y: c.Min.Y + offset.Y,
					Z: c.Min.Z + offset.Z,
				},
				Max: Vertex{
					X: c.Min.X + offset.X + half,
					Y: c.Min.Y + offset.Y + half,
					Z: c.Min.Z + offset.Z + half,
				},
				Center: Vertex{
					X: c.Center.X + offset.X - half/2,
					Y: c.Center.Y + offset.Y - half/2,
					Z: c.Center.Z + offset.Z - half/2,
				},
			},
			Size: half,
		}
	}
	return cubes, nil
}

func (c *Cube) SubDivideCube(depth int, originalObject *Object) ([]Cube, error) {
	if depth <= 0 {
		return []Cube{*c}, nil
	}

	cubes, err := c.DivideCube()
	if err != nil {
		return nil, err
	}

	var result []Cube
	for _, cube := range cubes {
		if cube.IntersectsObject(originalObject) {
			subCubes, err := cube.SubDivideCube(depth-1, originalObject)
			if err != nil {
				return nil, err
			}
			result = append(result, subCubes...)
		}
	}

	return result, nil
}

func (c *Cube) IntersectsObject(object *Object) bool {
	for _, face := range object.Faces {
		v1 := object.Vertexes[face.Vertexes[0]]
		v2 := object.Vertexes[face.Vertexes[1]]
		v3 := object.Vertexes[face.Vertexes[2]]

		if c.intersectTriangle(v1, v2, v3) {
			return true
		}
	}
	return false
}

func (c *Cube) intersectTriangle(v1, v2, v3 Vertex) bool {
	bMin, bMax := c.Min, c.Max
	tMin := Vertex{
		X: min(min(v1.X, v2.X), v3.X),
		Y: min(min(v1.Y, v2.Y), v3.Y),
		Z: min(min(v1.Z, v2.Z), v3.Z),
	}
	tMax := Vertex{
		X: max(max(v1.X, v2.X), v3.X),
		Y: max(max(v1.Y, v2.Y), v3.Y),
		Z: max(max(v1.Z, v2.Z), v3.Z),
	}

	return !(tMax.X < bMin.X || tMin.X > bMax.X ||
		tMax.Y < bMin.Y || tMin.Y > bMax.Y ||
		tMax.Z < bMin.Z || tMin.Z > bMax.Z)
}

func CubesToOBJ(cubes []Cube) *Object {
	obj := &Object{
		Faces:    []Face{},
		Vertexes: []Vertex{},
	}
	for _, cube := range cubes {
		v1 := Vertex{X: cube.Min.X, Y: cube.Min.Y, Z: cube.Min.Z}
		v2 := Vertex{X: cube.Max.X, Y: cube.Min.Y, Z: cube.Min.Z}
		v3 := Vertex{X: cube.Max.X, Y: cube.Max.Y, Z: cube.Min.Z}
		v4 := Vertex{X: cube.Min.X, Y: cube.Max.Y, Z: cube.Min.Z}
		v5 := Vertex{X: cube.Min.X, Y: cube.Min.Y, Z: cube.Max.Z}
		v6 := Vertex{X: cube.Max.X, Y: cube.Min.Y, Z: cube.Max.Z}
		v7 := Vertex{X: cube.Max.X, Y: cube.Max.Y, Z: cube.Max.Z}
		v8 := Vertex{X: cube.Min.X, Y: cube.Max.Y, Z: cube.Max.Z}

		baseIndex := len(obj.Vertexes)
		obj.Vertexes = append(obj.Vertexes, v1, v2, v3, v4, v5, v6, v7, v8)

		for _, face := range [][4]int{
			{0, 1, 2, 3}, // Bottom
			{4, 5, 6, 7}, // Top
			{0, 1, 5, 4}, // Front
			{2, 3, 7, 6}, // Back
			{1, 2, 6, 5}, // Right
			{0, 3, 7, 4}, // Left
		} {
			obj.Faces = append(obj.Faces,
				Face{Vertexes: [3]int{baseIndex + face[0], baseIndex + face[1], baseIndex + face[2]}},
				Face{Vertexes: [3]int{baseIndex + face[0], baseIndex + face[2], baseIndex + face[3]}},
			)
		}
	}
	return obj
}
