package sketches

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
	Font  *utils.Font
	Name string
	Color utils.ColorA

}
type BaseSketch struct {
	Sketch
	BaseSlide
}

func SetupSlide(a Slide, w *glfw.Window, f *utils.Font) error {
	a.Setup(w, f)
	return nil
}

func UpdateSlide(a Slide) {
	a.Update()
}

func DrawSlide(a Slide) {
	a.Draw()
}

func CloseSlide(a Slide) {
	a.Close()
}

func HandleKeyboardSketch(s Sketch, k glfw.Key, sc int, a glfw.Action, m glfw.ModifierKey) {
	s.HandleKeyboard(k, sc, a, m)
}

func HandleMousePositionSketch(s Sketch, xpos, ypos float64) {
	s.HandleMousePosition(xpos, ypos)
}

func HandleScrollSketch(s Sketch, xoff, yoff float64) {
	s.HandleScroll(xoff, yoff)
}

