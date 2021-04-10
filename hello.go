package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func main() {
	/* 这是我的第一个简单的程序 */
	fmt.Println("Hello, World!")
	data := []float64{1, 2, -1, -3}
	dense := mat.NewDense(2, 2, data)
	// data := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	// dense := mat.NewDense(4, 4, data)
	m2 := dense.T().T()
	dense.Inverse(m2)
	fmt.Println(dense.T()
}
