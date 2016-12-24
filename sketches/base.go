package sketches

import "github.com/go-gl/glfw/v3.2/glfw"

type Sketch interface {
	Setup() error
	Update()
	Draw()
	Close()
	HandleKeyboard(k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey)
	HandleMousePosition(xpos, ypos float64)
	HandleScroll(xoff, yoff float64)
}

type BaseSketch struct {
	Sketch
	Window *glfw.Window
}

func (s *BaseSketch) HandleKeyboard(k glfw.Key, sc int, a glfw.Action, m glfw.ModifierKey) {
	if k == glfw.KeyEscape && a == glfw.Press {
		s.Window.SetShouldClose(true)
	}
}
