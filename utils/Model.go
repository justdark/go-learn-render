package utils

import (
	"bufio"
	"container/list"
	"image"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ftrvxmtrx/tga"
)

type Vertex struct {
	X float64
	Y float64
	Z float64
}

type Face struct {
	Vertexs  list.List
	Z        float64
	VTexture [3]int
}

type DModel struct {
	Vertexs  []Vertex
	Faces    []Face
	VTexture []Vertex
	VNormal  []Vertex
	Texture  image.Image
}

func (dModel *DModel) SortFacesByZ() {
	sort.SliceStable(dModel.Faces, func(i, j int) bool {
		return dModel.Faces[i].Z < dModel.Faces[j].Z
	})
}

func str2float64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func extractVertexId(str string, index int) int {
	i, _ := strconv.ParseInt(strings.Split(str, "/")[index], 10, 32)
	return int(i)
}

func LoadModelFromFile(path string) DModel {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)
	vlist := []Vertex{}
	vtlist := []Vertex{}
	vnlist := []Vertex{}
	flist := []Face{}
	for s.Scan() {
		if strings.Index(s.Text(), "v ") == 0 {
			sp := strings.Split(s.Text(), " ")
			newV := Vertex{str2float64(sp[1]), str2float64(sp[2]), str2float64(sp[3])}
			vlist = append(vlist, newV)
			// vlist.PushBack(newV)
		} else if strings.Index(s.Text(), "f ") == 0 {
			sp := strings.Split(s.Text(), " ")
			fvs := list.New()
			fvts := [3]int{}
			fz := 0.0
			if len(sp) == 4 {
				for i := 1; i <= 3; i++ {
					index := extractVertexId(sp[i], 0) - 1
					fvs.PushBack(index)
					fvts[i-1] = extractVertexId(sp[i], 1) - 1
					fz += vlist[index].Z
				}
			}
			f := Face{*fvs, fz, fvts}
			flist = append(flist, f)
		} else if strings.Index(s.Text(), "vt ") == 0 {
			sp := strings.Split(s.Text(), " ")
			newV := Vertex{str2float64(sp[2]), str2float64(sp[3]), str2float64(sp[4])}
			vtlist = append(vtlist, newV)
		} else if strings.Index(s.Text(), "vn ") == 0 {
			sp := strings.Split(s.Text(), " ")
			newV := Vertex{str2float64(sp[2]), str2float64(sp[3]), str2float64(sp[4])}
			vnlist = append(vnlist, newV)
		}
	}
	return DModel{vlist, flist, vtlist, vnlist, nil}
}

func LoadModelFromFileWithDiffuse(path string, diffusePath string) DModel {
	model := LoadModelFromFile(path)
	file, _ := os.Open(diffusePath)
	model.Texture, _ = tga.Decode(file)
	return model
}
