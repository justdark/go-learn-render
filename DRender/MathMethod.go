package dr

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

type Vec3f struct {
	X float64
	Y float64
	Z float64
}
type Vector3f [3]float64
type Vector4f [4]float64
type Matrix44f [4][4]float64

func IdentityMatrix44f() Matrix44f {
	return Matrix44f{[4]float64{1, 0, 0, 0}, [4]float64{0, 1, 0, 0}, [4]float64{0, 0, 1, 0}, [4]float64{0, 0, 0, 1}}
}

func ZeroMatrix44f() Matrix44f {
	return Matrix44f{[4]float64{0, 0, 0, 0}, [4]float64{0, 0, 0, 0}, [4]float64{0, 0, 0, 0}, [4]float64{0, 0, 0, 0}}
}

func (m *Matrix44f) PVector(v Vector4f) Vector4f {
	rst := Vector4f{0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			rst[i] += m[i][j] * v[j]
		}
	}
	return rst
}

func (m *Matrix44f) PMatrix(v Matrix44f) Matrix44f {
	rst := ZeroMatrix44f()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				rst[i][j] += m[i][k] * v[k][j]
			}
		}
	}
	return rst
}

func (m *Matrix44f) Transpose() Matrix44f {
	rst := ZeroMatrix44f()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			rst[i][j] += m[j][i]
		}
	}
	return rst
}

func (m Matrix44f) Augment(right Matrix44f) [4][8]float64 {
	result := [4][8]float64{}
	for r, row := range m {
		for c := range row {
			result[r][c] = m[r][c]
		}
		cols := len(m[0])
		for c := range right[0] {
			result[r][cols+c] = right[r][c]
		}
	}
	return result
}
func (m Matrix44f) Invert() Matrix44f {
	data := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			data[i*4+j] = m[i][j]
		}
	}
	dense := mat.NewDense(4, 4, data)

	m2 := dense.T().T()
	dense.Inverse(m2)
	rst := ZeroMatrix44f()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			rst[i][j] = dense.At(i, j)
		}
	}
	return rst
}

func (v Matrix44f) P() *Matrix44f {
	return &v
}

func (v Vec3f) P() *Vec3f {
	return &v
}

func (v Vector3f) P() *Vector3f {
	return &v
}

func (v Vector4f) P() *Vector4f {
	return &v
}

func (v *Vec3f) DotProduct(v2 Vec3f) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v *Vec3f) ProductFloat(f float64) Vec3f {
	return Vec3f{v.X * f, v.Y * f, v.Z * f}
}

func (v *Vec3f) ToVector3fArray() Vector3f {
	return Vector3f{v.X, v.Y, v.Z}
}

func (v *Vec3f) ToVector4fArray() Vector4f {
	return Vector4f{v.X, v.Y, v.Z, 1.0}
}

func (v *Vec3f) Minus(v2 Vec3f) Vec3f {
	return Vec3f{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v *Vec3f) Plus(v2 Vec3f) Vec3f {
	return Vec3f{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v *Vec3f) Normalize() Vec3f {
	base := v.Norm()
	return Vec3f{v.X / base, v.Y / base, v.Z / base}
}

func (v *Vec3f) Norm() float64 {
	base := math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
	return base
}

func (v *Vector4f) ToVector3f() Vector3f {
	return Vector3f{v[0] / v[3], v[1] / v[3], v[2] / v[3]}
}

func (v *Vector3f) ToVec3f() Vec3f {
	return Vec3f{v[0], v[1], v[2]}
}
func (v *Vector3f) ToVertex() Vertex {
	return Vertex{v[0], v[1], v[2]}
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
