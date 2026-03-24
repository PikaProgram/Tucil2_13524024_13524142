package model

import (
	"errors"
	"math"
)

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

type OctreeCount struct {
	NodesFormed map[int]int
	NodesPruned map[int]int
}

func (o *Object) GetBoundingBox() (Box, error) {
	if len(o.Vertexes) == 0 {
		return Box{}, errors.New("object has no vertexes")
	}
	min := o.Vertexes[0]
	max := o.Vertexes[0]

	for _, v := range o.Vertexes {
		if min.X > v.X {
			min.X = v.X
		}
		if min.Y > v.Y {
			min.Y = v.Y
		}
		if min.Z > v.Z {
			min.Z = v.Z
		}

		if max.X < v.X {
			max.X = v.X
		}
		if max.Y < v.Y {
			max.Y = v.Y
		}
		if max.Z < v.Z {
			max.Z = v.Z
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

func (c *Cube) SubDivideCube(currentDepth int, maxDepth int, originalObject *Object, stats *OctreeCount) ([]Cube, error) {
	stats.NodesFormed[currentDepth]++

	if currentDepth >= maxDepth {
		return []Cube{*c}, nil
	}

	cubes, err := c.DivideCube()
	if err != nil {
		return nil, err
	}

	var result []Cube
	for _, cube := range cubes {
		if cube.IntersectsObject(originalObject) {
			subCubes, err := cube.SubDivideCube(currentDepth+1, maxDepth, originalObject, stats)
			if err != nil {
				return nil, err
			}
			result = append(result, subCubes...)
		} else {
			stats.NodesPruned[currentDepth+1]++
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

	// AABB check
	if tMax.X < bMin.X || tMin.X > bMax.X ||
		tMax.Y < bMin.Y || tMin.Y > bMax.Y ||
		tMax.Z < bMin.Z || tMin.Z > bMax.Z {
		return false
	}

	half := Vertex{
		X: (bMax.X - bMin.X) / 2,
		Y: (bMax.Y - bMin.Y) / 2,
		Z: (bMax.Z - bMin.Z) / 2,
	}

	center := Vertex{
		X: (bMin.X + bMax.X) / 2,
		Y: (bMin.Y + bMax.Y) / 2,
		Z: (bMin.Z + bMax.Z) / 2,
	}

	v1 = Vertex{X: v1.X - center.X, Y: v1.Y - center.Y, Z: v1.Z - center.Z}
	v2 = Vertex{X: v2.X - center.X, Y: v2.Y - center.Y, Z: v2.Z - center.Z}
	v3 = Vertex{X: v3.X - center.X, Y: v3.Y - center.Y, Z: v3.Z - center.Z}

	e0 := Vertex{X: v2.X - v1.X, Y: v2.Y - v1.Y, Z: v2.Z - v1.Z}
	e1 := Vertex{X: v3.X - v2.X, Y: v3.Y - v2.Y, Z: v3.Z - v2.Z}
	e2 := Vertex{X: v1.X - v3.X, Y: v1.Y - v3.Y, Z: v1.Z - v3.Z}

	axes := []Vertex{
		{X: 0, Y: -e0.Z, Z: e0.Y},
		{X: 0, Y: -e1.Z, Z: e1.Y},
		{X: 0, Y: -e2.Z, Z: e2.Y},
		{X: e0.Z, Y: 0, Z: -e0.X},
		{X: e1.Z, Y: 0, Z: -e1.X},
		{X: e2.Z, Y: 0, Z: -e2.X},
		{X: -e0.Y, Y: e0.X, Z: 0},
		{X: -e1.Y, Y: e1.X, Z: 0},
		{X: -e2.Y, Y: e2.X, Z: 0},
	}

	// 9 Axis test
	for _, axis := range axes {
		// Skip near-zero so math doesn't get mad
		if math.Abs(axis.X) < 1e-12 && math.Abs(axis.Y) < 1e-12 && math.Abs(axis.Z) < 1e-12 {
			continue
		}

		p1 := v1.X*axis.X + v1.Y*axis.Y + v1.Z*axis.Z
		p2 := v2.X*axis.X + v2.Y*axis.Y + v2.Z*axis.Z
		p3 := v3.X*axis.X + v3.Y*axis.Y + v3.Z*axis.Z

		r := half.X*math.Abs(axis.X) + half.Y*math.Abs(axis.Y) + half.Z*math.Abs(axis.Z)

		if max(max(p1, p2), p3) < -r || min(min(p1, p2), p3) > r {
			return false
		}
	}

	// AABB test for triangle vertices against cube half-size
	if (max(max(v1.X, v2.X), v3.X) < -half.X || min(min(v1.X, v2.X), v3.X) > half.X) ||
		(max(max(v1.Y, v2.Y), v3.Y) < -half.Y || min(min(v1.Y, v2.Y), v3.Y) > half.Y) ||
		(max(max(v1.Z, v2.Z), v3.Z) < -half.Z || min(min(v1.Z, v2.Z), v3.Z) > half.Z) {
		return false
	}

	// Triangle normal test
	cross := Vertex{
		X: e0.Y*e1.Z - e0.Z*e1.Y,
		Y: e0.Z*e1.X - e0.X*e1.Z,
		Z: e0.X*e1.Y - e0.Y*e1.X,
	}

	// Skip near-zero so math doesn't get mad, part 2 (degen triangles kekw)
	if cross.X*cross.X+cross.Y*cross.Y+cross.Z*cross.Z < 1e-12 {
		return false
	}

	q1 := cross.X*v1.X + cross.Y*v1.Y + cross.Z*v1.Z
	q2 := cross.X*v2.X + cross.Y*v2.Y + cross.Z*v2.Z
	q3 := cross.X*v3.X + cross.Y*v3.Y + cross.Z*v3.Z

	r := half.X*math.Abs(cross.X) + half.Y*math.Abs(cross.Y) + half.Z*math.Abs(cross.Z)

	if max(max(q1, q2), q3) < -r || min(min(q1, q2), q3) > r {
		return false
	}

	return true
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
