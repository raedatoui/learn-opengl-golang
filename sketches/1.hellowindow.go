package sketches

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type HelloWindow struct {
	window *glfw.Window
}

func (hw *HelloWindow) Setup() error {
	return nil
}

func (hw *HelloWindow) Draw() {

}

func (hw *HelloWindow) Update() {

}

func (hw *HelloWindow) Close() {

}

func (hw *HelloWindow) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		hw.window.SetShouldClose(true)
	}
}

func (hw *HelloWindow) HandleMousePosition(xpos, ypos float64) {

}

func (hw *HelloWindow) HandleScroll(xoff, yoff float64) {

}