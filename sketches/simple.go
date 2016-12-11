package sketches

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type SimpleKetch struct {
	Window *glfw.Window
}

func (sketch SimpleKetch) Setup() {

}

func (sketch SimpleKetch) Draw() {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (sketch SimpleKetch) Update() {

}

func (Sketch SimpleKetch) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

}