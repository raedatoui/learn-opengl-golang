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
