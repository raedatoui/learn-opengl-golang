package sketches

import "github.com/go-gl/glfw/v3.2/glfw"

type Sketch interface {
	Setup()
	Update()
	Draw()
	HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
}
