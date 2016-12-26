package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

type HelloShaders struct {
	sections.BaseSketch
	vao, vbo uint32
	shader   uint32
}

func (hs *HelloShaders) Setup(w *glfw.Window, f *utils.Font) error {
	hs.Window = w
	hs.Font = f
	hs.Color = utils.RandColor()
	hs.Name = "3. Shaders"

	var err error
	hs.shader, err = utils.Shader(
		"_assets/3.shaders/basic.vs",
		"_assets/3.shaders/basic.frag", "")

	if err != nil {
		return err
	}

	var vertices = []float32{
		// Positions      // Colors
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // Bottom Left
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // Top
	}
	gl.GenVertexArrays(1, &hs.vao)
	gl.GenBuffers(1, &hs.vbo)

	gl.BindVertexArray(hs.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, hs.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	// position uniform
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	//color uniform
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(3*utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return nil
}

func (hs *HelloShaders) Update() {

}

func (hs *HelloShaders) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hs.Color.R, hs.Color.G, hs.Color.B, hs.Color.A)

	// Draw the triangle
	gl.UseProgram(hs.shader)
	gl.BindVertexArray(hs.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)

	hs.Font.SetColor(0.0, 0.0, 0.0, 1.0)
	hs.Font.Printf(30, 30, 0.5, hs.Name)
}

func (hs *HelloShaders) Close() {
	gl.DeleteVertexArrays(1, &hs.vao)
	gl.DeleteBuffers(1, &hs.vbo)
	gl.UseProgram(0)
}

func (hc *HelloShaders) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {

}

func (hs *HelloShaders) HandleMousePosition(xpos, ypos float64) {

}

func (hs *HelloShaders) HandleScroll(xoff, yoff float64) {

}
