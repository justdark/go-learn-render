package dr

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"utils"

	"github.com/fogleman/gg"
)

type DRender struct {
	context *gg.Context
}

type Vertex struct {
	X float64
	Y float64
	Z float64
}

func (v *Vertex) IntegerLize() {
	v.X = math.Round(v.X)
	v.Y = math.Round(v.Y)
}

func (v Vertex) P() *Vertex {
	return &v
}

func (v *Vertex) Minus(v2 Vertex) Vertex {
	return Vertex{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}
func (v *Vertex) Add(v2 Vertex) Vertex {
	return Vertex{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}
func (v *Vertex) ToVec3f() Vec3f {
	return Vec3f{v.X, v.Y, v.Z}
}
func (v *Vertex) CrossProduct(v2 Vertex) Vec3f {
	return Cross(v.ToVec3f(), v2.ToVec3f())
}
func NewDRender(width int, height int) DRender {
	ctx := gg.NewContext(width, height)
	return DRender{
		context: ctx,
	}
}
func (dRender *DRender) InvertY() {
	dRender.context.InvertY()
}

func (dRender *DRender) PrintSample() { fmt.Println("xxx") }

func (dRender *DRender) DrawOnPixel(x int, y int, color color.Color) {
	dRender.context.SetColor(color)
	dRender.context.SetPixel(x, y)
	dRender.context.Fill()
}

func (dRender *DRender) DrawOnPixelInvertedY(x int, y int, color color.Color) {
	dRender.context.SetColor(color)
	dRender.context.SetPixel(x, dRender.context.Height()-y)
	dRender.context.Fill()
}

func (dRender *DRender) SavePNG(path string) { dRender.context.SavePNG(path) }

func (dRender *DRender) DrawLineByVertex(v1 Vertex, v2 Vertex, color color.Color) {
	dRender.DrawLine(int(v1.X), int(v1.Y), int(v2.X), int(v2.Y), color)
}
func (dRender *DRender) DrawLine(x0 int, y0 int, x1 int, y1 int, color color.Color) {
	steep := false
	if math.Abs(float64(x0-x1)) < math.Abs(float64(y0-y1)) {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		steep = true
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := y1 - y0
	derror := math.Abs(float64(dy) / float64(dx))
	error := 0.0
	y := y0
	yi := 1
	if y0 > y1 {
		yi = -1
	}
	for x := x0; x < x1; x++ {
		if steep {
			dRender.DrawOnPixelInvertedY(y, x, color)
		} else {
			dRender.DrawOnPixelInvertedY(x, y, color)
		}
		error += derror
		if error > 0.5 {
			y += yi
			error -= 1.0
		}
	}
}

func (dRender *DRender) DrawTriangle(v1 Vertex, v2 Vertex, v3 Vertex, color color.Color) {
	dRender.DrawLineByVertex(v1, v2, color)
	dRender.DrawLineByVertex(v1, v3, color)
	dRender.DrawLineByVertex(v2, v3, color)
}

func BoundingBoxFind(vertexs []Vertex) (Vec3f, Vec3f) {
	bboxmin := Vec3f{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	bboxmax := Vec3f{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	for i := 0; i < len(vertexs); i++ {
		bboxmin.X = math.Min(bboxmin.X, vertexs[i].X)
		bboxmin.Y = math.Min(bboxmin.Y, vertexs[i].Y)
		bboxmin.Z = math.Min(bboxmin.Z, vertexs[i].Z)
		bboxmax.X = math.Max(bboxmax.X, vertexs[i].X)
		bboxmax.Y = math.Max(bboxmax.Y, vertexs[i].Y)
		bboxmax.Z = math.Max(bboxmax.Z, vertexs[i].Z)
	}
	return bboxmin, bboxmax
}

func BoundingBoxWithLimit(vertexs []Vertex, min Vec3f, max Vec3f) (Vec3f, Vec3f) {
	bboxmin := Vec3f{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	bboxmax := Vec3f{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	for i := 0; i < len(vertexs); i++ {
		bboxmin.X = math.Min(bboxmin.X, vertexs[i].X)
		bboxmin.Y = math.Min(bboxmin.Y, vertexs[i].Y)
		bboxmin.Z = math.Min(bboxmin.Z, vertexs[i].Z)
		bboxmax.X = math.Max(bboxmax.X, vertexs[i].X)
		bboxmax.Y = math.Max(bboxmax.Y, vertexs[i].Y)
		bboxmax.Z = math.Max(bboxmax.Z, vertexs[i].Z)
	}

	bboxmin.X = math.Max(bboxmin.X, min.X)
	bboxmin.Y = math.Max(bboxmin.Y, min.Y)
	bboxmin.Z = math.Max(bboxmin.Z, min.Z)
	bboxmax.X = math.Min(bboxmax.X, max.X)
	bboxmax.Y = math.Min(bboxmax.Y, max.Y)
	bboxmax.Z = math.Min(bboxmax.Z, max.Z)

	return bboxmin, bboxmax
}

const triangleThresh = 0.0

func (dRender *DRender) FillTriangle(v1 Vertex, v2 Vertex, v3 Vertex, color color.Color) {
	vertexs := []Vertex{v1, v2, v3}
	bbmin, bbmax := BoundingBoxFind(vertexs)
	// dRender.DrawTriangle(v1, v2, v3, color)

	for x := bbmin.X; x < bbmax.X; x++ {
		for y := bbmin.Y; y < bbmax.Y; y++ {
			u := Barycentric2D(vertexs, Vertex{x, y, 0})
			if u.X < triangleThresh || u.Y < triangleThresh || u.Z < triangleThresh {
				continue
			}
			dRender.DrawOnPixelInvertedY(int(x), int(y), color)
		}
	}
}
func (dRender *DRender) FillTriangleWithZBuffer(v1 Vertex, v2 Vertex, v3 Vertex, zBuffer *[]float64, color color.Color) {
	vertexs := []Vertex{v1, v2, v3}
	bbmin, bbmax := BoundingBoxWithLimit(vertexs, Vec3f{0, 0, 0}, Vec3f{float64(dRender.context.Width()), float64(dRender.context.Height()), math.MaxFloat64})
	// dRender.DrawTriangle(v1, v2, v3, color)

	for x := int(bbmin.X); x < int(bbmax.X); x++ {
		for y := int(bbmin.Y); y < int(bbmax.Y); y++ {
			u := Barycentric2D(vertexs, Vertex{float64(x), float64(y), 0})
			if u.X < triangleThresh || u.Y < triangleThresh || u.Z < triangleThresh {
				continue
			}
			z := v1.Z*u.X + v2.Z*u.Y + v3.Z*u.Z
			if (*zBuffer)[int(x+dRender.context.Width()*y)] < z {
				(*zBuffer)[int(x+dRender.context.Width()*y)] = z
				dRender.DrawOnPixelInvertedY(int(x), int(y), color)
			}
		}
	}
}

func (dRender *DRender) FillTriangleWithTexture(v1 Vertex, v2 Vertex, v3 Vertex, zBuffer *[]float64, vts [3]Vertex, texture image.Image, lightStrength float64) {
	vertexs := []Vertex{v1, v2, v3}
	bbmin, bbmax := BoundingBoxWithLimit(vertexs, Vec3f{0, 0, 0}, Vec3f{float64(dRender.context.Width()), float64(dRender.context.Height()), math.MaxFloat64})
	// dRender.DrawTriangle(v1, v2, v3, color)
	for x := int(bbmin.X); x < int(bbmax.X); x++ {
		for y := int(bbmin.Y); y < int(bbmax.Y); y++ {
			u := Barycentric2D(vertexs, Vertex{float64(x), float64(y), 0})
			if u.X < triangleThresh || u.Y < triangleThresh || u.Z < triangleThresh {
				continue
			}
			z := v1.Z*u.X + v2.Z*u.Y + v3.Z*u.Z
			vx := vts[0].X*u.X + vts[1].X*u.Y + vts[2].X*u.Z
			vy := vts[0].Y*u.X + vts[1].Y*u.Y + vts[2].Y*u.Z
			colorRaw := texture.At(int(vx), int(vy))
			r, g, b, _ := colorRaw.RGBA()

			// newColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			newColor := color.RGBA{uint8(float64(r>>8) * lightStrength), uint8(float64(g>>8) * lightStrength), uint8(float64(b>>8) * lightStrength), 255}
			// fmt.Println(colorRaw, newColor, r, g, b, uint8(float64(r)*lightStrength), uint8(float64(g)*lightStrength), uint8(float64(b)*lightStrength))
			// fmt.Println(x, y, int(x)+dRender.context.Width()*int(y))
			if (*zBuffer)[int(x)+dRender.context.Width()*int(y)] < z {
				(*zBuffer)[int(x)+dRender.context.Width()*int(y)] = z
				dRender.DrawOnPixelInvertedY(int(x), int(y), newColor)
			}
		}
	}
}

func (dRender *DRender) FillTriangleWithTextureAndGouraudshading(v1 Vertex, v2 Vertex, v3 Vertex, zBuffer *[]float64, vts [3]Vertex, intensities [3]float64, texture image.Image) {
	vertexs := []Vertex{v1, v2, v3}
	bbmin, bbmax := BoundingBoxWithLimit(vertexs, Vec3f{0, 0, 0}, Vec3f{float64(dRender.context.Width()), float64(dRender.context.Height()), math.MaxFloat64})
	// dRender.DrawTriangle(v1, v2, v3, color)
	for x := int(bbmin.X); x < int(bbmax.X); x++ {
		for y := int(bbmin.Y); y < int(bbmax.Y); y++ {
			u := Barycentric2D(vertexs, Vertex{float64(x), float64(y), 0})
			if u.X < triangleThresh || u.Y < triangleThresh || u.Z < triangleThresh {
				continue
			}
			z := v1.Z*u.X + v2.Z*u.Y + v3.Z*u.Z
			vx := vts[0].X*u.X + vts[1].X*u.Y + vts[2].X*u.Z
			vy := vts[0].Y*u.X + vts[1].Y*u.Y + vts[2].Y*u.Z
			colorRaw := texture.At(int(vx), int(vy))
			r, g, b, _ := colorRaw.RGBA()

			// newColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			lightStrength := intensities[0]*u.X + intensities[1]*u.Y + intensities[2]*u.Z
			lightStrength = math.Max(lightStrength, 0.0)
			lightStrength = math.Min(lightStrength, 1.0)
			newColor := color.RGBA{uint8(float64(r>>8) * lightStrength), uint8(float64(g>>8) * lightStrength), uint8(float64(b>>8) * lightStrength), 255}
			// fmt.Println(colorRaw, newColor, r, g, b, uint8(float64(r)*lightStrength), uint8(float64(g)*lightStrength), uint8(float64(b)*lightStrength))
			// fmt.Println(x, y, int(x)+dRender.context.Width()*int(y))
			if (*zBuffer)[int(x)+dRender.context.Width()*int(y)] < z {
				(*zBuffer)[int(x)+dRender.context.Width()*int(y)] = z
				dRender.DrawOnPixelInvertedY(int(x), int(y), newColor)
			}
		}
	}
}

func (dRender *DRender) FillTriangleWithTextureAndNormMapping(v1 Vertex, v2 Vertex, v3 Vertex, zBuffer *[]float64, vts [3]Vertex, texture image.Image, normMapping image.Image, lightDirect Vec3f) {
	vertexs := []Vertex{v1, v2, v3}
	bbmin, bbmax := BoundingBoxWithLimit(vertexs, Vec3f{0, 0, 0}, Vec3f{float64(dRender.context.Width()), float64(dRender.context.Height()), math.MaxFloat64})
	// dRender.DrawTriangle(v1, v2, v3, color)
	for x := int(bbmin.X); x < int(bbmax.X); x++ {
		for y := int(bbmin.Y); y < int(bbmax.Y); y++ {
			u := Barycentric2D(vertexs, Vertex{float64(x), float64(y), 0})
			if u.X < triangleThresh || u.Y < triangleThresh || u.Z < triangleThresh {
				continue
			}
			z := v1.Z*u.X + v2.Z*u.Y + v3.Z*u.Z
			vx := vts[0].X*u.X + vts[1].X*u.Y + vts[2].X*u.Z
			vy := vts[0].Y*u.X + vts[1].Y*u.Y + vts[2].Y*u.Z

			colorRaw := texture.At(int(vx), int(vy))
			r, g, b, _ := colorRaw.RGBA()

			nx, ny, nz, _ := normMapping.At(int(vx), int(vy)).RGBA()
			norm := Vec3f{float64(nz>>8)/255.0*2.0 - 1.0, float64(ny>>8)/255.0*2.0 - 1.0, float64(nx>>8)/255*2.0 - 1.0}.P().Normalize()

			// newColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			lightStrength := norm.DotProduct(lightDirect.Normalize())
			lightStrength = math.Max(lightStrength, 0.0)
			newColor := color.RGBA{uint8(float64(r>>8) * lightStrength), uint8(float64(g>>8) * lightStrength), uint8(float64(b>>8) * lightStrength), 255}
			// fmt.Println(colorRaw, newColor, r, g, b, uint8(float64(r)*lightStrength), uint8(float64(g)*lightStrength), uint8(float64(b)*lightStrength))
			// fmt.Println(x, y, int(x)+dRender.context.Width()*int(y))
			if (*zBuffer)[int(x)+dRender.context.Width()*int(y)] < z {
				(*zBuffer)[int(x)+dRender.context.Width()*int(y)] = z
				dRender.DrawOnPixelInvertedY(int(x), int(y), newColor)
			}
		}
	}
}

func (dRender *DRender) FillTriangleWithTextureAndNormMappingAndSpecMapping(v1 Vertex, v2 Vertex, v3 Vertex, zBuffer *[]float64, vts [3]Vertex, vns [3]Vertex, model utils.DModel, lightDirect Vec3f, m Matrix44f, mit Matrix44f) {
	vertexs := []Vertex{v1, v2, v3}
	bbmin, bbmax := BoundingBoxWithLimit(vertexs, Vec3f{0, 0, 0}, Vec3f{float64(dRender.context.Width()), float64(dRender.context.Height()), math.MaxFloat64})
	// dRender.DrawTriangle(v1, v2, v3, color)
	for x := int(bbmin.X); x < int(bbmax.X); x++ {
		for y := int(bbmin.Y); y < int(bbmax.Y); y++ {
			u := Barycentric2D(vertexs, Vertex{float64(int(x)), float64(int(y)), 0})
			if u.X < triangleThresh || u.Y < triangleThresh || u.Z < triangleThresh {
				continue
			}
			z := v1.Z*u.X + v2.Z*u.Y + v3.Z*u.Z
			vx := vts[0].X*u.X + vts[1].X*u.Y + vts[2].X*u.Z
			vy := vts[0].Y*u.X + vts[1].Y*u.Y + vts[2].Y*u.Z
			vx = float64(int(vx * 1024))
			vy = float64(int(vy * 1024))

			colorRaw := model.Texture.At(int(vx), int(vy))
			r, g, b, _ := colorRaw.RGBA()

			nx, ny, nz, _ := model.NormalMapping.At(int(vx), int(vy)).RGBA()
			norm := Vec3f{float64(nz)/65536.0*2.0 - 1.0, float64(ny)/65536.0*2.0 - 1.0, float64(nx)/65536*2.0 - 1.0}.P().Normalize()
			lightDirect = lightDirect.Normalize()
			// norm = mit.PVector(norm.ToVector4fArray()).P().ToVector3f().P().ToVec3f().P().Normalize()
			// lightDirect = m.PVector(lightDirect.ToVector4fArray()).P().ToVector3f().P().ToVec3f().P().Normalize()

			specR, specG, specB, specA := model.SpecularMapping.At(int(vx), int(vy)).RGBA()
			// newColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}

			lightStrength := norm.DotProduct(lightDirect)
			rv := norm.ProductFloat(lightStrength * 2.0).P().Minus(lightDirect).P().Normalize() // reflected light
			lightStrength = math.Max(lightStrength, 0.0)
			specValue := math.Pow(math.Max(rv.Z, 0.0), 12*(1.0-math.Min(1.0, float64(specR)/65536)))
			// if spec <= 257 {
			// 	specValue = 0.0
			// }
			// fmt.Println(spec / 256)

			// fmt.Println(specValue)
			// specValue = 0.0
			// r8 := uint8(math.Min(255, 5+float64(r>>8)*(lightStrength+0.6*specValue)))
			// g8 := uint8(math.Min(255, 5+float64(g>>8)*(lightStrength+0.6*specValue)))
			// b8 := uint8(math.Min(255, 5+float64(b>>8)*(lightStrength+0.6*specValue)))

			r8 := uint8(math.Min(255, float64(r>>8)*(1.2*lightStrength+2.6*specValue)))
			g8 := uint8(math.Min(255, float64(g>>8)*(1.2*lightStrength+2.6*specValue)))
			b8 := uint8(math.Min(255, float64(b>>8)*(1.2*lightStrength+2.6*specValue)))

			// r8 = uint8(math.Min(255, float64(specB>>8)))
			// g8 = uint8(math.Min(255, float64(specB>>8)))
			// b8 = uint8(math.Min(255, float64(specB>>8)))

			// r8 = uint8(math.Min(255, specValue*255))
			// g8 = uint8(math.Min(255, specValue*255))
			// b8 = uint8(math.Min(255, specValue*255))

			// r8 = uint8(math.Min(255, lightStrength*255))
			// g8 = uint8(math.Min(255, lightStrength*255))
			// b8 = uint8(math.Min(255, lightStrength*255))
			// r8 = uint8(math.Min(255, math.Max(rv.Z, 0.0)*255))
			// g8 = uint8(math.Min(255, math.Max(rv.Z, 0.0)*255))
			// b8 = uint8(math.Min(255, math.Max(rv.Z, 0.0)*255))
			specB = specB + specG + specR + specA
			newColor := color.RGBA{r8, g8, b8, 255}
			// if float64(spec>>8)/255.0 > 0.5 {
			// 	fmt.Println(model.SpecularMapping.At(int(vx), int(vy)).RGBA())
			// 	fmt.Println(lightStrength, specValue)
			// 	fmt.Println(rv)
			// 	fmt.Println(norm)
			// 	fmt.Println("==============")
			// }
			// fmt.Println(colorRaw, newColor, r, g, b, uint8(float64(r)*lightStrength), uint8(float64(g)*lightStrength), uint8(float64(b)*lightStrength))
			// fmt.Println(x, y, int(x)+dRender.context.Width()*int(y))
			if (*zBuffer)[int(x)+dRender.context.Width()*int(y)] < z {
				(*zBuffer)[int(x)+dRender.context.Width()*int(y)] = z
				dRender.DrawOnPixelInvertedY(int(x), int(y), newColor)
			}
		}
	}
}

func VertexTrans(v utils.Vertex) Vertex {
	return Vertex{v.X, v.Y, v.Z}
}

func (dRender *DRender) LoadModelForRender(modelPath string, width float64, height float64, lightDirect Vec3f) {
	model := utils.LoadModelFromFile(modelPath)
	model.SortFacesByZ()
	// add sort by Z to better performance

	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		vertexs := [3]Vertex{}
		vertexs_for_intentsity := [3]Vertex{}
		i := 0
		for v := face.Vertexs.Front(); v != nil && i < 3; v = v.Next() {
			vIndex := v.Value
			vertexs[i] = VertexTrans(model.Vertexs[vIndex.(int)])
			vertexs_for_intentsity[i] = VertexTrans(model.Vertexs[vIndex.(int)])
			vertexs[i].X = math.Round((vertexs[i].X + 1.0) * float64(width) / 2.0)
			vertexs[i].Y = math.Round((vertexs[i].Y + 1.0) * float64(height) / 2.0)
			i++
		}
		// dRender.FillTriangle(vertexs[0], vertexs[1], vertexs[2], color.RGBA{0, 255, 0, 255})
		// random color
		// dRender.FillTriangle(vertexs[0], vertexs[1], vertexs[2], color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))})

		// intensity
		v1 := vertexs_for_intentsity[1].Minus(vertexs_for_intentsity[0])
		v2 := vertexs_for_intentsity[2].Minus(vertexs_for_intentsity[0])
		intensity := Cross(v1.ToVec3f(), v2.ToVec3f())
		intensity = intensity.Normalize()
		lightStrength := intensity.DotProduct(lightDirect)
		// fmt.Println(lightStrength)
		lightStrength = math.Min(lightStrength, 1.0)
		if lightStrength > 0 {
			dRender.FillTriangle(vertexs[0], vertexs[1], vertexs[2], color.RGBA{uint8(255 * lightStrength), uint8(255 * lightStrength), uint8(255 * lightStrength), 255})
		}
	}
}

func (dRender *DRender) LoadModelForRenderWithTexture(modelPath string, texturePath string, width float64, height float64, lightDirect Vec3f) {
	model := utils.LoadModelFromFileWithDiffuse(modelPath, texturePath)

	zBuffer := make([]float64, int(width*height))
	for i := range zBuffer {
		zBuffer[i] = -math.MaxFloat64
	}
	fmt.Println(model.Texture.Bounds().Max)
	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		vertexs := [3]Vertex{}
		vertexs_for_intentsity := [3]Vertex{}
		vertexs_for_texture := [3]Vertex{}
		textureVertex := face.VTexture
		for i := 0; i < 3; i++ {
			vertexs_for_texture[i] = VertexTrans(model.VTexture[textureVertex[i]])
			vertexs_for_texture[i].X *= float64(model.Texture.Bounds().Max.X)
			vertexs_for_texture[i].Y = float64(model.Texture.Bounds().Max.Y) - vertexs_for_texture[i].Y*float64(model.Texture.Bounds().Max.Y)
		}
		// fmt.Println(vertexs_for_texture)

		i := 0
		for v := face.Vertexs.Front(); v != nil && i < 3; v = v.Next() {
			vIndex := v.Value
			vertexs[i] = VertexTrans(model.Vertexs[vIndex.(int)])
			vertexs_for_intentsity[i] = VertexTrans(model.Vertexs[vIndex.(int)])
			vertexs[i].X = math.Round((vertexs[i].X + 1.0) * width / 2.0)
			vertexs[i].Y = math.Round((vertexs[i].Y + 1.0) * height / 2.0)
			i++
		}

		// intensity
		v1 := vertexs_for_intentsity[1].Minus(vertexs_for_intentsity[0])
		v2 := vertexs_for_intentsity[2].Minus(vertexs_for_intentsity[0])
		intensity := Cross(v1.ToVec3f(), v2.ToVec3f())
		intensity = intensity.Normalize()
		lightDirect := Vec3f{0, 0, 1}
		lightStrength := intensity.DotProduct(lightDirect)
		// fmt.Println(lightStrength)
		lightStrength = math.Min(lightStrength, 1.0)
		lightStrength = math.Max(lightStrength, 0.0)
		dRender.FillTriangleWithTexture(vertexs[0], vertexs[1], vertexs[2], &zBuffer, vertexs_for_texture, model.Texture, lightStrength)

	}
}

func (dRender *DRender) LoadModelForRenderWithTextureAndCamera(modelPath string, texturePath string, width float64, height float64, lightDirect Vec3f, camera Vec3f) {
	model := utils.LoadModelFromFileWithDiffuse(modelPath, texturePath)

	zBuffer := make([]float64, int(width*height))
	for i := range zBuffer {
		zBuffer[i] = -math.MaxFloat64
	}
	fmt.Println(model.Texture.Bounds().Max)
	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		vertexs := [3]Vertex{}
		vertexs_for_intentsity := [3]Vertex{}
		vertexs_for_texture := [3]Vertex{}
		textureVertex := face.VTexture
		for i := 0; i < 3; i++ {
			vertexs_for_texture[i] = VertexTrans(model.VTexture[textureVertex[i]])
			vertexs_for_texture[i].X *= float64(model.Texture.Bounds().Max.X)
			vertexs_for_texture[i].Y = float64(model.Texture.Bounds().Max.Y) - vertexs_for_texture[i].Y*float64(model.Texture.Bounds().Max.Y)
		}
		// fmt.Println(vertexs_for_texture)

		i := 0
		for v := face.Vertexs.Front(); v != nil && i < 3; v = v.Next() {
			vIndex := v.Value
			vertexs[i] = VertexTrans(model.Vertexs[vIndex.(int)])
			vertexs_for_intentsity[i] = vertexs[i]
			vertexs[i].X = math.Round((vertexs[i].X + 1.0) * width / 2.0)
			vertexs[i].Y = math.Round((vertexs[i].Y + 1.0) * height / 2.0)
			i++
		}

		// intensity
		v1 := vertexs_for_intentsity[1].Minus(vertexs_for_intentsity[0])
		v2 := vertexs_for_intentsity[2].Minus(vertexs_for_intentsity[0])
		intensity := Cross(v1.ToVec3f(), v2.ToVec3f())
		intensity = intensity.Normalize()
		lightDirect := Vec3f{0, 0, 1}
		lightStrength := intensity.DotProduct(lightDirect)
		// fmt.Println(lightStrength)
		lightStrength = math.Min(lightStrength, 1.0)
		lightStrength = math.Max(lightStrength, 0.0)
		dRender.FillTriangleWithTexture(vertexs[0], vertexs[1], vertexs[2], &zBuffer, vertexs_for_texture, model.Texture, lightStrength)

	}
}
