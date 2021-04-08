package main

import (
	"dr"
	"fmt"
	"image/color"
	"math"
)

func main() {
	dRender := dr.NewDRender(500, 400)
	dRender.PrintSample()
	dRender.DrawOnPixel(15, 20, color.RGBA{255, 0, 0, 255})
	dRender.DrawLine(30, 30, 300, 300, color.RGBA{255, 0, 0, 255})
	dRender.DrawLine(30, 30, 200, 300, color.RGBA{0, 255, 0, 255})
	dRender.SavePNG("./out.png")
	fmt.Println(math.Round(0.51))
	fmt.Println(math.Round(0.5))
	fmt.Println(math.Round(0.49))
}
