package utils

import "math/rand"

type ColorA struct {
	R, G, B, A float32
}

func RandColor() ColorA {
	return ColorA{rand.Float32(),rand.Float32(), rand.Float32(), 1.0 }
}
