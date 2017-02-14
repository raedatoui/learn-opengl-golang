package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloTriangle struct {
	sections.BaseSketch
	program  uint32
	vao, vbo uint32
}

func (ht *HelloTriangle) createBuffers(vertices []float32) (uint32, uint32) {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*glutils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	//vertAttrib := uint32(gl.GetAttribLocation(ht.program, gl.Str("position\x00")))
	// here we can skip computing the vertAttrib value and use 0 since our shader declares layout = 0 for
	// the uniform
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*glutils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)
	return vao, vbo
}

func (ht *HelloTriangle) InitGL() error {
	ht.Name = "2a. Hello Triangle"

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

	var err error
	ht.program, err = glutils.BasicProgram(vertexShader, fragShader)
	if err != nil {
		return err
	}

	var vertices = []float32{
		-0.5, -0.5, 0.0, // Left
		0.5, -0.5, 0.0, // Right
		0.0, 0.5, 0.0, // Top
	}
	ht.vao, ht.vbo = ht.createBuffers(vertices)
	return nil
}

func (ht *HelloTriangle) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(ht.Color32.R, ht.Color32.G, ht.Color32.B, ht.Color32.A)

	// Draw our first triangle
	gl.UseProgram(ht.program)
	gl.BindVertexArray(ht.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (ht *HelloTriangle) Close() {
	gl.DeleteVertexArrays(1, &ht.vao)
	gl.DeleteBuffers(1, &ht.vbo)
	gl.DeleteProgram(ht.program)
}

type TriangleEx1 struct {
	HelloTriangle
	ebo         uint32
	currentMode int32
}

func (hs *TriangleEx1) InitGL() error {
	hs.Name = "2b. Triangle Ex1"
	var vertexShader = `
	#version 330 core
	in vec3 vert;
	void main() {
		gl_Position = vec4(vert.x, vert.y, vert.z, 1.0);
	}` + "\x00"

	var fragShader = `
	#version 330 core
	out vec4 color;
	void main() {
		color = vec4(1.0f, 1.0f, 0.2f, 1.0f);
	}` + "\x00"

	var err error
	hs.program, err = glutils.BasicProgram(vertexShader, fragShader)
	if err != nil {
		return err
	}

	var vertices = []float32{
		0.5, 0.5, 0.0, // Top Right
		0.5, -0.5, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, // Top Left
	}

	var indices = []uint32{
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}

	gl.GenVertexArrays(1, &hs.vao)
	gl.BindVertexArray(hs.vao)

	gl.GenBuffers(1, &hs.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, hs.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*glutils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &hs.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, hs.ebo)
	// seems like 4 works best here for the size of uint32
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(hs.program, gl.Str("vert\x00")))
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*glutils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vertAttrib)

	gl.BindVertexArray(0)

	return nil
}

func (hs *TriangleEx1) Draw() {
	gl.GetIntegerv(gl.POLYGON_MODE, &hs.currentMode)

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hs.Color32.R, hs.Color32.G, hs.Color32.B, hs.Color32.A)

	gl.UseProgram(hs.program)
	gl.BindVertexArray(hs.vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
	gl.PolygonMode(gl.FRONT_AND_BACK, uint32(hs.currentMode))
}

func (hs *TriangleEx1) Close() {
	gl.DeleteVertexArrays(1, &hs.vao)
	gl.DeleteBuffers(1, &hs.vbo)
	gl.DeleteBuffers(1, &hs.ebo)
	gl.UseProgram(0)
}

func (hs *TriangleEx1) GetSubHeader() string {
	return "the square always uses GL_LINE for the polygon mode"
}

type TriangleEx2 struct {
	HelloTriangle
	program2   uint32
	vao2, vbo2 uint32
}

func (ht *TriangleEx2) InitGL() error {
	ht.Name = "2c. Triangle Ex2"

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
	  color = vec4(1.0f, 1.0f, 0.2f, 0.9f);
	}` + "\x00"

	var vertexShader2 = `
	#version 330 core
	layout (location = 0) in vec3 position;
	void main() {
	  gl_Position = vec4(position.x/2.0, position.y/2.0, position.z/2.0, 1.0);
	}` + "\x00"

	var fragShader2 = `
	#version 330 core
	out vec4 color;
	void main() {
	  color = vec4(0.8f, 0.2f, 1.0f, 0.8f);
	}` + "\x00"

	var err error
	ht.program, err = glutils.BasicProgram(vertexShader, fragShader)
	if err != nil {
		return err
	}

	ht.program2, err = glutils.BasicProgram(vertexShader2, fragShader2)
	if err != nil {
		return err
	}

	var vertices = []float32{
		0.5, -0.5, 0.0, // Right
		-0.5, -0.5, 0.0, // Left
		0.0, 0.5, 0.0, // Top
	}

	var vertices2 = []float32{
		1.0, -1.0, 0.0, // Left
		1.0, 0.0, 0.0, // Right
		0.0, -1.0, 0.0, // Top
	}

	ht.vao, ht.vbo = ht.createBuffers(vertices)
	ht.vao2, ht.vbo2 = ht.createBuffers(vertices2)

	return nil
}

func (ht *TriangleEx2) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(ht.Color32.R, ht.Color32.G, ht.Color32.B, ht.Color32.A)

	// Draw our first triangle
	gl.UseProgram(ht.program)
	gl.BindVertexArray(ht.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)

	gl.UseProgram(ht.program2)
	gl.BindVertexArray(ht.vao2)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}
