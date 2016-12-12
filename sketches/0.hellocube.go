// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package sketches

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/learn-opengl/utils"
	_ "image/png"
	"log"
)

type HelloCube struct {
	Window              *glfw.Window
	Program             uint32
	Vao, Vbo            uint32
	Texture             uint32
	Angle, PreviousTime float64
	Model               mgl32.Mat4
	ModelUniform        int32
}

func (sketch *HelloCube) Setup() {
	var cubeVertices = []float32{
		//  X, Y, Z, U, V
		// Bottom
		-1.0, -1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,

		// Top
		-1.0, 1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0, 1.0,

		// Front
		-1.0, -1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,

		// Back
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 1.0,

		// Left
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,

		// Right
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
	}

	var vertexShader = `
	#version 330

	uniform mat4 projection;
	uniform mat4 camera;
	uniform mat4 model;

	in vec3 vert;
	in vec2 vertTexCoord;

	out vec2 fragTexCoord;

	void main() {
		fragTexCoord = vertTexCoord;
		gl_Position = projection * camera * model * vec4(vert, 1);
	}
	` + "\x00"

	var fragmentShader = `
	#version 330

	uniform sampler2D tex;

	in vec2 fragTexCoord;

	out vec4 outputColor;

	void main() {
		outputColor = texture(tex, fragTexCoord);
	}
	` + "\x00"
	// Configure the vertex and fragment shaders
	var err error
	sketch.Program, err = utils.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(sketch.Program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(800.0)/600.0, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(sketch.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(sketch.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	sketch.Model = mgl32.Ident4()
	sketch.ModelUniform = gl.GetUniformLocation(sketch.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(sketch.ModelUniform, 1, false, &sketch.Model[0])

	textureUniform := gl.GetUniformLocation(sketch.Program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(sketch.Program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	sketch.Texture, err = utils.NewTexture("sketches/square.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	gl.GenVertexArrays(1, &sketch.Vao)
	gl.BindVertexArray(sketch.Vao)

	gl.GenBuffers(1, &sketch.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, sketch.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(sketch.Program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(sketch.Program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	sketch.Angle = 0.0
	sketch.PreviousTime = glfw.GetTime()
}

func (sketch *HelloCube) Update() {
	time := glfw.GetTime()
	elapsed := time - sketch.PreviousTime
	sketch.PreviousTime = time

	sketch.Angle += elapsed
	sketch.Model = mgl32.HomogRotate3D(float32(sketch.Angle), mgl32.Vec3{0, 1, 0})
}

func (sketch *HelloCube) Draw() {
	// Render
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(sketch.Program)
	gl.UniformMatrix4fv(sketch.ModelUniform, 1, false, &sketch.Model[0])

	gl.BindVertexArray(sketch.Vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

}

func (sketch *HelloCube) Close() {
	gl.DeleteVertexArrays(1, &sketch.Vao)
	gl.DeleteBuffers(1, &sketch.Vbo)
	gl.DeleteBuffers(1, &sketch.Vao)
}

func (sketch *HelloCube) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		sketch.Window.SetShouldClose(true)
	}
}
