package sketches

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

type HelloShaders struct {
	Window   *glfw.Window
	vao, vbo uint32
	shader  uint32
}


func (hs *HelloShaders) Setup() error {
	var err error
	hs.shader, err = utils.Shader(
		"sketches/_assets/3.shaders/basic.vs",
		"sketches/_assets/3.shaders/basic.frag", "")

	if err != nil {
		return err
	}

	var vertices = []float32{
		// Positions      // Colors
		0.5, -0.5, 0.0,   1.0, 0.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0,  0.0, 1.0, 0.0, // Bottom Left
		0.0, 0.5, 0.0,    0.0, 0.0, 1.0, // Top
	}
	gl.GenVertexArrays(1, &hs.vao)
	gl.GenBuffers(1, &hs.vbo)

	gl.BindVertexArray(hs.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, hs.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) * utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	// position uniform
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6 * utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	//color uniform
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6 * utils.GL_FLOAT32_SIZE, gl.PtrOffset(3 * utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return nil
}

func (hs *HelloShaders) Update() {

}

func (hs *HelloShaders) Draw() {
	// Draw the triangle
	gl.UseProgram(hs.shader)
	gl.BindVertexArray(hs.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (hs *HelloShaders) Close() {
	gl.DeleteVertexArrays(1, &hs.vao)
	gl.DeleteBuffers(1, &hs.vbo)
	gl.UseProgram(0)
}

func (hs *HelloShaders) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		hs.Window.SetShouldClose(true)
	}
}

func (hs *HelloShaders) HandleMousePosition(xpos, ypos float64) {

}

func (hs *HelloShaders) HandleScroll(xoff, yoff float64) {

}