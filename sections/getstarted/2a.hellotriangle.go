package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

type HelloTriangle struct {
	sections.BaseSketch
	program  uint32
	vao, vbo uint32
}

func (ht *HelloTriangle) Setup(w *glfw.Window, f *utils.Font) error {
	ht.Window = w
	ht.Font = f
	ht.Color = utils.StepColor(utils.MAG, utils.BLACK, 10, 3)
	ht.Name = "2. Hello Triangles"

	var vertexShader = `
	#version 330 core
	layout (location = 0) in vec3 position;
	void main() {
	  gl_Position = vec4(position.x, position.y, position.z, 1.0);
	}` + "\x00"

	var fragShader = `
	#version 330 core
	out vec4 color;
	void main() {
	  color = vec4(1.0f, 1.0f, 0.2f, 1.0f);
	}` + "\x00"

	var vertices = []float32{
		-0.5, -0.5, 0.0, // Left
		0.5, -0.5, 0.0, // Right
		0.0, 0.5, 0.0, // Top
	}
	var err error
	ht.program, err = utils.BasicProgram(vertexShader, fragShader)
	if err != nil {
		return err
	}
	gl.UseProgram(ht.program)

	gl.GenVertexArrays(1, &ht.vao)
	gl.GenBuffers(1, &ht.vbo)

	gl.BindVertexArray(ht.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, ht.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	//vertAttrib := uint32(gl.GetAttribLocation(ht.program, gl.Str("position\x00")))
	// here we can skip computing the vertAttrib value and use 0 since our shader declares layout = 0 for
	// the uniform
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return nil
}

func (ht *HelloTriangle) Update() {

}

func (ht *HelloTriangle) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(ht.Color.R, ht.Color.G, ht.Color.B, ht.Color.A)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// Draw our first triangle
	gl.UseProgram(ht.program)
	gl.BindVertexArray(ht.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)

	ht.Font.SetColor(0.0, 0.0, 0.0, 1.0)
	ht.Font.Printf(30, 30, 0.5, ht.Name)
}

func (ht *HelloTriangle) Close() {
	gl.DeleteVertexArrays(1, &ht.vao)
	gl.DeleteBuffers(1, &ht.vbo)
	gl.UseProgram(0)
}

func (hc *HelloTriangle) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {

}

func (ht *HelloTriangle) HandleMousePosition(xpos, ypos float64) {

}

func (ht *HelloTriangle) HandleScroll(xoff, yoff float64) {

}
