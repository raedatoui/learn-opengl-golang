package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

type HelloTriangleC struct {
	sections.BaseSketch
	program  uint32
	vao, vbo uint32
	program2  uint32
	vao2, vbo2 uint32
}

func (ht *HelloTriangleC) InitGL() error {
	ht.Name = "2c. Hello 2 Triangles"

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


	var err error
	ht.program, err = utils.BasicProgram(vertexShader, fragShader)
	if err != nil {
		return err
	}

	ht.program2, err = utils.BasicProgram(vertexShader2, fragShader2)
	if err != nil {
		return err
	}

	gl.GenVertexArrays(1, &ht.vao)
	gl.GenVertexArrays(1, &ht.vao2)

	gl.GenBuffers(1, &ht.vbo)
	gl.GenBuffers(1, &ht.vbo2)

	gl.BindVertexArray(ht.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, ht.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	//vertAttrib := uint32(gl.GetAttribLocation(ht.program, gl.Str("position\x00")))
	// here we can skip computing the vertAttrib value and use 0 since our shader declares layout = 0 for
	// the uniform
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	gl.BindVertexArray(ht.vao2)

	gl.BindBuffer(gl.ARRAY_BUFFER, ht.vbo2)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices2)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices2), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	return nil
}

func (ht *HelloTriangleC) Draw() {
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

func (ht *HelloTriangleC) Close() {
	gl.DeleteVertexArrays(1, &ht.vao)
	gl.DeleteBuffers(1, &ht.vbo)
	gl.UseProgram(0)
}
