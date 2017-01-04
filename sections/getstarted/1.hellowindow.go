package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloWindow struct {
	sections.BaseSketch
}

func (hw *HelloWindow) InitGL() error {
	hw.Name = "1. Hello Window"
	return nil
}

func (hw *HelloWindow) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hw.Color32.R, hw.Color32.G, hw.Color32.B, hw.Color32.A)
}

func (hw *HelloWindow) Close() {

}
