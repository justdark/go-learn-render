package main

import (
	"dr"
	"utils"
)

func VertexTrans(v utils.Vertex) dr.Vertex {
	return dr.Vertex{v.X, v.Y, v.Z}
}
func main() {
	width := 1000.0
	height := 1000.0
	dRender := dr.NewDRender(int(width), int(height))
	// dRender.DrawTriangle(dr.Vertex{20, 30, 0}, dr.Vertex{800, 300, 0}, dr.Vertex{100, 30, 0}, color.RGBA{0, 255, 0, 255})
	// dRender.FillTriangle(dr.Vertex{20, 30, 0}, dr.Vertex{800, 300, 0}, dr.Vertex{100, 30, 0}, color.RGBA{0, 255, 0, 255})
	// dRender.FillTriangle(dr.Vertex{231, 350, 0}, dr.Vertex{500, 200, 0}, dr.Vertex{300, 300, 0}, color.RGBA{0, 255, 0, 255})
	// dRender.SavePNG("./out.png")

	// model := utils.LoadModelFromFile("./obj/african_head.obj")
	// model.SortFacesByZ()
	// fmt.Println(len(model.Faces))
	// // add sort by Z to better performance

	// for i := 0; i < len(model.Faces); i++ {
	// 	face := model.Faces[i]
	// 	vertexs := [3]dr.Vertex{}
	// 	vertexs_for_intentsity := [3]dr.Vertex{}
	// 	i := 0
	// 	for v := face.Vertexs.Front(); v != nil && i < 3; v = v.Next() {
	// 		vIndex := v.Value
	// 		vertexs[i] = VertexTrans(model.Vertexs[vIndex.(int)])
	// 		vertexs_for_intentsity[i] = VertexTrans(model.Vertexs[vIndex.(int)])
	// 		vertexs[i].X = math.Round((vertexs[i].X + 1.0) * width / 2.0)
	// 		vertexs[i].Y = math.Round((vertexs[i].Y + 1.0) * height / 2.0)
	// 		i++
	// 	}
	// 	// dRender.FillTriangle(vertexs[0], vertexs[1], vertexs[2], color.RGBA{0, 255, 0, 255})
	// 	// random color
	// 	// dRender.FillTriangle(vertexs[0], vertexs[1], vertexs[2], color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))})

	// 	// intensity
	// 	v1 := vertexs_for_intentsity[1].Minus(vertexs_for_intentsity[0])
	// 	v2 := vertexs_for_intentsity[2].Minus(vertexs_for_intentsity[0])
	// 	intensity := dr.Cross(v1.ToVec3f(), v2.ToVec3f())
	// 	intensity = intensity.Normalize()
	// 	lightDirect := dr.Vec3f{0, 0, 1}
	// 	lightStrength := intensity.DotProduct(lightDirect)
	// 	// fmt.Println(lightStrength)
	// 	lightStrength = math.Min(lightStrength, 1.0)
	// 	if lightStrength > 0 {
	// 		dRender.FillTriangle(vertexs[0], vertexs[1], vertexs[2], color.RGBA{uint8(255 * lightStrength), uint8(255 * lightStrength), uint8(255 * lightStrength), 255})
	// 	}
	// }

	// in one line
	dRender.LoadModelForRender("./obj/african_head.obj", width, height, dr.Vec3f{0, 0, 1})
	dRender.SavePNG("./lecture2.png")
}
