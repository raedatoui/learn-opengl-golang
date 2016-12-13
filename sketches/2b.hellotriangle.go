package sketches

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl/utils"
)

type HelloSquare struct {
	Window        *glfw.Window
	Program       uint32
	Vao, Vbo, Ebo uint32
}

func (sketch *HelloSquare) Setup() {
	var vertexShader2 = `
	#version 330 core
	in vec3 vert;
	void main() {
		gl_Position = vec4(vert.x, vert.y, vert.z, 1.0);
	}` + "\x00"

	var fragShader2 = `
	#version 330 core
	out vec4 color;
	void main() {
		color = vec4(1.0f, 1.0f, 0.2f, 1.0f);
	}` + "\x00"

	var vertices = []float32{
		0.5,  0.5, 0.0,  // Top Right
		0.5, -0.5, 0.0,  // Bottom Right
		-0.5, -0.5, 0.0, // Bottom Left
		-0.5, 0.5, 0.0,   // Top Left
	}

	var indices = []uint32{
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}

	var err error
	sketch.Program, err = utils.BasicProgram(vertexShader2, fragShader2)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(sketch.Program)

	gl.GenVertexArrays(1, &sketch.Vao)
	gl.BindVertexArray(sketch.Vao)

	gl.GenBuffers(1, &sketch.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, sketch.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)* utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &sketch.Ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, sketch.Ebo)
	// seems like 4 works best here for the size of uint32
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices) * 4, gl.Ptr(indices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(sketch.Program, gl.Str("vert\x00")))
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3 * utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vertAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.Disable(gl.DEPTH_TEST)
}

func (sketch *HelloSquare) Update() {

}

func (sketch *HelloSquare) Draw() {
	gl.ClearColor(0.5, 0.5, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(sketch.Program)
	gl.BindVertexArray(sketch.Vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func (sketch *HelloSquare) Close() {
	gl.DeleteVertexArrays(1, &sketch.Vao)
	gl.DeleteBuffers(1, &sketch.Vbo)
	gl.DeleteBuffers(1, &sketch.Ebo)
	gl.UseProgram(0)
}

func (sketch *HelloSquare) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		sketch.Window.SetShouldClose(true)
	}
}
