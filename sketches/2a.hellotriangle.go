package sketches

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl/utils"
)

type HelloTriangle struct {
	Window        *glfw.Window
	Program       uint32
	Vao, Vbo      uint32
}

func (sketch *HelloTriangle) Setup() {

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
	sketch.Program, err = utils.BasicProgram(vertexShader, fragShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(sketch.Program)

	gl.GenVertexArrays(1, &sketch.Vao)
	gl.GenBuffers(1, &sketch.Vbo)

	gl.BindVertexArray(sketch.Vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, sketch.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)* utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	//vertAttrib := uint32(gl.GetAttribLocation(sketch.Program, gl.Str("position\x00")))
	// here we can skip computing the vertAttrib value and use 0 since our shader declares layout = 0 for
	// the uniform
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3 * utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	// switch to 2d mode
	gl.Disable(gl.DEPTH_TEST)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
}

func (sketch *HelloTriangle) Update() {

}

func (sketch *HelloTriangle) Draw() {
	gl.ClearColor(0.5, 0.5, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	// Draw our first triangle
	gl.UseProgram(sketch.Program)
	gl.BindVertexArray(sketch.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func (sketch *HelloTriangle) Close() {
	gl.DeleteVertexArrays(1, &sketch.Vao)
	gl.DeleteBuffers(1, &sketch.Vbo)
	gl.UseProgram(0)
}

func (sketch *HelloTriangle) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		sketch.Window.SetShouldClose(true)
	}
}

func (sketch *HelloTriangle) HandleMousePosition(xpos, ypos float64) {

}

func (sketch *HelloTriangle) HandleScroll(xoff, yoff float64) {

}