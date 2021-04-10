package main

import (
	"dr"
	"fmt"
	"math"
	"utils"
)

func lookAt(eye dr.Vec3f, center dr.Vec3f, up dr.Vec3f) dr.Matrix44f {
	z := eye.Minus(center).P().Normalize()
	x := dr.Cross(up, z).P().Normalize()
	y := dr.Cross(z, x).P().Normalize()
	centerArr := center.ToVector3fArray()

	az := z.ToVector3fArray()
	ax := x.ToVector3fArray()
	ay := y.ToVector3fArray()

	Minv := dr.IdentityMatrix44f()
	Tr := dr.IdentityMatrix44f()
	for i := 0; i < 3; i++ {
		Minv[0][i] = ax[i]
		Minv[1][i] = ay[i]
		Minv[2][i] = az[i]
		Tr[i][3] = -centerArr[i]
	}
	ModelView := Minv.PMatrix(Tr)
	return ModelView
}
func genViewPoint(x, y, w, h, depth int) dr.Matrix44f {
	m := dr.IdentityMatrix44f()
	m[0][3] = float64(x) + float64(w)/2.0
	m[1][3] = float64(y) + float64(h)/2.0
	m[2][3] = float64(depth) / 2.0

	m[0][0] = float64(w) / 2.0
	m[1][1] = float64(h) / 2.0
	m[2][2] = float64(depth) / 2.0
	return m
}

func main() {
	width := 1000.0
	height := 1000.0
	depth := 255
	eye := dr.Vec3f{1, 1, 4}
	center := dr.Vec3f{0, 0, 0}
	lightDirect := dr.Vec3f{1, 1, 1}.P().Normalize()
	viewPoint := genViewPoint(int(width)/8, int(height)/8, int(width)*3/4, int(height)*3/4, depth)
	modelView := lookAt(eye, center, dr.Vec3f{0, 1, 0})
	dRender := dr.NewDRender(int(width), int(height))
	projectMatrix := dr.IdentityMatrix44f()
	projectMatrix[3][2] = -1.0 / eye.Minus(center).P().Norm()

	finalTransMatrix := viewPoint.PMatrix(projectMatrix).P().PMatrix(modelView)

	model := utils.LoadModelFromFileWithAll("./obj/african_head.obj", "./obj/african_head_diffuse.tga", "./obj/african_head_nm.tga", "./obj/african_head_spec.tga")
	// fmt.Println(len(model.Faces))
	// add sort by Z to better performance
	zBuffer := make([]float64, int(width*height))
	for i := range zBuffer {
		zBuffer[i] = -math.MaxFloat64
	}
	fmt.Println(model.Texture.Bounds().Max)
	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		vertexs := [3]dr.Vertex{}
		vertexs_for_intentsity := [3]dr.Vertex{}
		vertexs_for_texture := [3]dr.Vertex{}
		textureVertex := face.VTexture
		for i := 0; i < 3; i++ {
			vertexs_for_texture[i] = dr.VertexTrans(model.VTexture[textureVertex[i]])
			vertexs_for_texture[i].X *= float64(model.Texture.Bounds().Max.X)
			vertexs_for_texture[i].Y = float64(model.Texture.Bounds().Max.Y) - vertexs_for_texture[i].Y*float64(model.Texture.Bounds().Max.Y)
		}
		// fmt.Println(vertexs_for_texture)

		i := 0

		// fmt.Println(projectMatrix)
		// finalTransMatrix = finalTransMatrix.PMatrix(projectMatrix)
		for v := face.Vertexs.Front(); v != nil && i < 3; v = v.Next() {
			vIndex := v.Value
			v := dr.VertexTrans(model.Vertexs[vIndex.(int)])
			vertexs[i] = finalTransMatrix.PVector(v.ToVec3f().P().ToVector4fArray()).P().ToVector3f().P().ToVertex() // trans at here
			vertexs[i].IntegerLize()
			vertexs_for_intentsity[i] = finalTransMatrix.PVector(v.ToVec3f().P().ToVector4fArray()).P().ToVector3f().P().ToVertex()
			i++
		}

		// fmt.Println(intensities)

		//6.0 paint without norm mapping src
		// dRender.FillTriangleWithTextureAndNormMapping(vertexs[0], vertexs[1], vertexs[2], &zBuffer, vertexs_for_texture, model.Texture, model.NormalMapping, lightDirect)

		//6.1 paint with spec
		dRender.FillTriangleWithTextureAndNormMappingAndSpecMapping(vertexs[0], vertexs[1], vertexs[2], &zBuffer, vertexs_for_texture, model, lightDirect)
	}

	// dRender.SavePNG("./lecture6.0.png")
	dRender.SavePNG("./lecture6.2.png")

}
