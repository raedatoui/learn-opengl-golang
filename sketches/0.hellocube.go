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

func (hc *HelloCube) Setup() error {
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
	hc.Program, err = utils.BasicProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}

	gl.UseProgram(hc.Program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(800.0)/600.0, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(hc.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(hc.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	hc.Model = mgl32.Ident4()
	hc.ModelUniform = gl.GetUniformLocation(hc.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(hc.ModelUniform, 1, false, &hc.Model[0])

	textureUniform := gl.GetUniformLocation(hc.Program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(hc.Program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	hc.Texture, err = utils.NewTexture("sketches/_assets/0.cube/square.png")
	if err != nil {
		return err
	}

	// Configure the vertex data
	gl.GenVertexArrays(1, &hc.Vao)
	gl.BindVertexArray(hc.Vao)

	gl.GenBuffers(1, &hc.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, hc.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(hc.Program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(hc.Program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Configure global settings

	hc.Angle = 0.0
	hc.PreviousTime = glfw.GetTime()
	return nil
}

func (hc *HelloCube) Update() {
	time := glfw.GetTime()
	elapsed := time - hc.PreviousTime
	hc.PreviousTime = time

	hc.Angle += elapsed
	hc.Model = mgl32.HomogRotate3D(float32(hc.Angle), mgl32.Vec3{0, 1, 0})
}

func (hc *HelloCube) Draw() {
	gl.UseProgram(hc.Program)
	gl.UniformMatrix4fv(hc.ModelUniform, 1, false, &hc.Model[0])

	gl.BindVertexArray(hc.Vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, hc.Texture)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

}

func (hc *HelloCube) Close() {
	gl.DeleteVertexArrays(1, &hc.Vao)
	gl.DeleteBuffers(1, &hc.Vbo)
	gl.DeleteBuffers(1, &hc.Vao)
	gl.UseProgram(0)
}

func (hc *HelloCube) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		hc.Window.SetShouldClose(true)
	}
}

func (hc *HelloCube) HandleMousePosition(xpos, ypos float64) {

}

func (hc *HelloCube) HandleScroll(xoff, yoff float64) {

}