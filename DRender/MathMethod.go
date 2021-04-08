package dr

import (
	"math"
)

type Vec3f struct {
	X float64
	Y float64
	Z float64
}

type Vector4f [4]float64
type Matrix44f [4][4]float64

func IdentityMatrix44f() Matrix44f {
	return Matrix44f{[4]float64{1, 0, 0, 0}, [4]float64{0, 1, 0, 0}, [4]float64{0, 0, 1, 0}, [4]float64{1, 1, 1, 1}}
}

func (m *Matrix44f) pVector(v Vector4f) Vector4f {
	rst := Vector4f{0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			rst[i] += m[i][j] * v[j]
		}
	}
	return rst
}

func (v *Vec3f) DotProduct(v2 Vec3f) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v *Vec3f) Normalize() Vec3f {
	base := math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
	return Vec3f{v.X / base, v.Y / base, v.Z / base}
}

func Cross(a Vec3f, b Vec3f) Vec3f {
	return Vec3f{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
}

func Barycentric2D(pts []Vertex, point Vertex) Vec3f {
	u := Cross(Vec3f{pts[2].X - pts[0].X, pts[1].X - pts[0].X, pts[0].X - point.X}, Vec3f{pts[2].Y - pts[0].Y, pts[1].Y - pts[0].Y, pts[0].Y - point.Y})
	if math.Abs(u.Z) <= 1.0 {
		return Vec3f{-1, 1, 1}
	}
	return Vec3f{1.0 - (u.X+u.Y)/u.Z, u.Y / u.Z, u.X / u.Z}
}
