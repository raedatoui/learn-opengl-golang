package sketches

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type HelloWindow struct {
	Window *glfw.Window
}

func (sketch *HelloWindow) Setup() error {
	return nil
}

func (sketch *HelloWindow) Draw() {

}

func (sketch *HelloWindow) Update() {

}

func (sketch *HelloWindow) Close() {

}

func (sketch *HelloWindow) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		sketch.Window.SetShouldClose(true)
	}
}

func (sketch *HelloWindow) HandleMousePosition(xpos, ypos float64) {

}

func (sketch *HelloWindow) HandleScroll(xoff, yoff float64) {

}