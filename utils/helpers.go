package utils

import (
	"reflect"
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