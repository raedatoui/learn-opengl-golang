package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

type HelloWindow struct {
	sections.BaseSketch
}

func (hw *HelloWindow) Setup(w *glfw.Window, f *utils.Font) error {
	hw.Window = w
	hw.Font = f
	hw.Color = utils.RandColor()
	hw.Name = "1. Hello Window"
	return nil
}

func (hw *HelloWindow) Update() {

}

func (hw *HelloWindow) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hw.Color.R, hw.Color.G, hw.Color.B, hw.Color.A)
	hw.Font.SetColor(0.0, 0.0, 0.0, 1.0)
	hw.Font.Printf(30, 30, 0.5, hw.Name)
}

func (hw *HelloWindow) Close() {

}

func (hc *HelloWindow) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {

}

func (hw *HelloWindow) HandleMousePosition(xpos, ypos float64) {

}

func (hw *HelloWindow) HandleScroll(xoff, yoff float64) {

}
