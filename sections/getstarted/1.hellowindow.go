package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloWindow struct {
	sections.BaseSketch
}

func (hw *HelloWindow) InitGL() error {
	hw.Name = "1. Hello Window"
	return nil
}

func (hw *HelloWindow) Update() {

}

func (hw *HelloWindow) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hw.Color.R, hw.Color.G, hw.Color.B, hw.Color.A)
}

func (hw *HelloWindow) Close() {

}

func (hc *HelloWindow) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {

}

func (hw *HelloWindow) HandleMousePosition(xpos, ypos float64) {

}

func (hw *HelloWindow) HandleScroll(xoff, yoff float64) {

}
