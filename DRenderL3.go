package main

import (
	"dr"
	"fmt"
	"math"
	"utils"
)

func main() {
	width := 1000.0
	height := 1000.0
	dRender := dr.NewDRender(int(width), int(height))

	model := utils.LoadModelFromFileWithDiffuse("./obj/african_head.obj", "./obj/african_head_diffuse.tga")
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
		for v := face.Vertexs.Front(); v != nil && i < 3; v = v.Next() {
			vIndex := v.Value
			vertexs[i] = dr.VertexTrans(model.Vertexs[vIndex.(int)])
			vertexs_for_intentsity[i] = dr.VertexTrans(model.Vertexs[vIndex.(int)])
			vertexs[i].X = math.Round((vertexs[i].X + 1.0) * width / 2.0)
			vertexs[i].Y = math.Round((vertexs[i].Y + 1.0) * height / 2.0)
			i++
		}

		// intensity
		v1 := vertexs_for_intentsity[1].Minus(vertexs_for_intentsity[0])
		v2 := vertexs_for_intentsity[2].Minus(vertexs_for_intentsity[0])
		intensity := dr.Cross(v1.ToVec3f(), v2.ToVec3f())
		intensity = intensity.Normalize()
		lightDirect := dr.Vec3f{0, 0, 1.0}
		lightStrength := intensity.DotProduct(lightDirect)
		// fmt.Println(lightStrength)
		lightStrength = math.Min(lightStrength, 1.0)
		lightStrength = math.Max(lightStrength, 0.0)
		dRender.FillTriangleWithTexture(vertexs[0], vertexs[1], vertexs[2], &zBuffer, vertexs_for_texture, model.Texture, lightStrength)
	}

	// in one line
	// dRender.LoadModelForRenderWithTexture("./obj/african_head.obj", "./obj/african_head_diffuse.tga", width, height, dr.Vec3f{0, 0, 1})

	dRender.SavePNG("./lecture3.png")

}
