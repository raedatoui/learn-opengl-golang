package utils

import (
	"reflect"
	"github.com/go-gl/mathgl/mgl32"
	"fmt"
	"strconv"
)

// checks if o implements i
// i must be an interface and not an instance of a struct.
// example: (*MyInterface)(nil)
func Implements(o, i interface{}) bool {
	ot := reflect.ValueOf(o).Type()
	bs := reflect.TypeOf(i).Elem()
	return ot.Implements(bs)
}

// compares 2 objects and determines if they have the same type.
func IsType(o, i interface{}) bool {
	a := reflect.ValueOf(o).Type()
	b := reflect.ValueOf(i).Type()
	return a == b
}

func PrintMat4(m mgl32.Mat4) {
	fmt.Printf("%s\n%s\n%s\n%s\n-------\n",
		ftos(m[0:4]),
		ftos(m[4:8]),
		ftos(m[8:12]),
		ftos(m[12:16]),
	)
}

func ftos(f []float32) string {
    out := ""
	for i := range f {
        out += strconv.FormatFloat(float64(f[i]), 'f', 2, 32) + ", "
    }
	return out
}