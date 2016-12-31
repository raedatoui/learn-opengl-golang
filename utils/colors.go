package utils

import (
	"fmt"
	"math"
	"math/rand"
)

type Color struct {
	R, G, B, A float64
}

type Color32 struct {
	R, G, B, A float32
}

func (c *Color) To32() Color32 {
	return Color32{
		R: float32(c.R), G: float32(c.G),
		B: float32(c.B), A: float32(c.A)}
}

func RandColor() Color {
	return Color{rand.Float64(), rand.Float64(), rand.Float64(), 1.0}
}

func StepColor(c1, c2 Color, t, i int) Color {
	factorStep := 1 / (float64(t) - 1.0)
	c := interpolateColor(c1, c2, factorStep*float64(i))
	return Color{
		R: float64(c.R),
		G: float64(c.G),
		B: float64(c.B),
	}
}

func interpolateColor(c1, c2 Color, factor float64) Color {
	result := new(Color)
	result.R = c1.R + factor*(c2.R-c1.R)
	result.G = c1.G + factor*(c2.G-c1.G)
	result.B = c1.B + factor*(c2.B-c1.B)
	return *result
}

func Rgb2Hex(c Color) string {
	rgb := []int{
		int(round(c.R*255, 0.5, 0)),
		int(round(c.G*255, 0.5, 0)),
		int(round(c.B*255, 0.5, 0)),
	}
	t := (1 << 24) + (rgb[0] << 16) + (rgb[1] << 8) + rgb[2]
	s := fmt.Sprintf("%x", t)
	return "#" + s
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

var Magenta = Color{R: 236.0 / 255.0, G: 0, B: 140.0 / 255.0}
var Black = Color{R: 0, G: 0, B: 0}
var White = Color{R: 1.0, G: 1.0, B: 1.0}
