package sections

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

// Slide is the most basic slide. it has to setup, update, draw and cloae
type Slide interface {
	Setup(w *glfw.Window, f *utils.Font) error
	Update()
	Draw()
	Close()
}

// Sketch is an interactive Slide and process user interactions
type Sketch interface {
	Slide
	HandleKeyboard(k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey)
	HandleMousePosition(xpos, ypos float64)
	HandleScroll(xoff, yoff float64)
}

// BaseSlide is the base implementation of Slide with the min required fields
type BaseSlide struct {
	Slide
	Window *glfw.Window
	Font   *utils.Font
	Name   string
	Color  utils.ColorA
}
type BaseSketch struct {
	Sketch
	BaseSlide
}
