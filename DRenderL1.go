package main

import (
	"dr"
	"fmt"
	"image/color"
	"math"
	"utils"
)

func main() {
	width := 1000.0
	height := 1000.0
	dRender := dr.NewDRender(int(width), int(height))
	dRender.PrintSample()
	// dRender.DrawOnPixel(15, 20, color.RGBA{255, 0, 0, 255})
	// dRender.DrawLine(30, 30, 300, 300, color.RGBA{255, 0, 0, 255})
	// dRender.DrawLine(30, 30, 200, 300, color.RGBA{0, 255, 0, 255})

	fmt.Println(math.Round(0.51))

	model := utils.LoadModelFromFile("./obj/african_head.obj")
	fmt.Println(len(model.Faces))

	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		for v := face.Vertexs.Front(); v != nil; v = v.Next() {
			nvIndex := face.Vertexs.Front().Value
			if v.Next() != nil {
				nvIndex = v.Next().Value
			}
			vIndex := v.Value
			v0 := model.Vertexs[vIndex.(int)]
			v1 := model.Vertexs[nvIndex.(int)]
			x0 := (v0.X + 1.0) * width / 2.0
			y0 := (v0.Y + 1.0) * height / 2.0
			x1 := (v1.X + 1.0) * width / 2.0
			y1 := (v1.Y + 1.0) * height / 2.0
			dRender.DrawLine(int(x0), int(y0), int(x1), int(y1), color.RGBA{255, 255, 255, 255})
		}
	}
	dRender.SavePNG("./lecture1.png")
}
