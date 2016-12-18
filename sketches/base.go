package sketches

import "github.com/go-gl/glfw/v3.2/glfw"

type Sketch interface {
	Setup()
	Update()
	Draw()
	Close()
	HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	HandleMousePosition(xpos, ypos float64)
	HandleScroll(xoff, yoff float64)
}
