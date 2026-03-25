package viewer

import (
	"math"
	"tucil/src/internal/model"
)

type Matrix4 [4][4]float64

func Identity() Matrix4 {
	return Matrix4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m1 Matrix4) Multiply(m2 Matrix4) Matrix4 {
	var out Matrix4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			out[i][j] = m1[i][0]*m2[0][j] +
				m1[i][1]*m2[1][j] +
				m1[i][2]*m2[2][j] +
				m1[i][3]*m2[3][j]
		}
	}
	return out
}

func (m *Matrix4) MultiplyVector(v model.Vertex, out *model.Vertex) {
	x := v.X*m[0][0] + v.Y*m[0][1] + v.Z*m[0][2] + m[0][3]
	y := v.X*m[1][0] + v.Y*m[1][1] + v.Z*m[1][2] + m[1][3]
	z := v.X*m[2][0] + v.Y*m[2][1] + v.Z*m[2][2] + m[2][3]
	w := v.X*m[3][0] + v.Y*m[3][1] + v.Z*m[3][2] + m[3][3]

	if w != 0 && w != 1 {
		x /= w
		y /= w
		z /= w
	}

	out.X = x
	out.Y = y
	out.Z = z
}

// Creates a matrix that rotates a point around the X-axis
func RotationX(angleRad float64) Matrix4 {
	c := math.Cos(angleRad)
	s := math.Sin(angleRad)
	return Matrix4{
		{1, 0, 0, 0},
		{0, c, -s, 0},
		{0, s, c, 0},
		{0, 0, 0, 1},
	}
}

// Creates a matrix that rotates a point around the Y-axis
func RotationY(angleRad float64) Matrix4 {
	c := math.Cos(angleRad)
	s := math.Sin(angleRad)
	return Matrix4{
		{c, 0, s, 0},
		{0, 1, 0, 0},
		{-s, 0, c, 0},
		{0, 0, 0, 1},
	}
}

// Used for zoom. Not an actual zoom, but enough for this assignment ig
func Translation(tx, ty, tz float64) Matrix4 {
	return Matrix4{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1},
	}
}

func Perspective(near, far float64) Matrix4 {
	fovRad := math.Pi / 2
	f := 1.0 / math.Tan(fovRad/2.0)
	rangeInv := 1.0 / (near - far)

	return Matrix4{
		{f, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, (far + near) * rangeInv, (2 * far * near) * rangeInv},
		{0, 0, -1, 0},
	}
}
