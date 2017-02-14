package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloTransformations struct {
	sections.BaseSketch
	shader             uint32
	vao, vbo, ebo      uint32
	texture1, texture2 uint32
	translationMat     mgl32.Mat4
	rotationAxis       mgl32.Vec3
}

func (ht *HelloTransformations) InitGL() error {
	ht.Name = "5. Transformations"

	ht.translationMat = mgl32.Translate3D(0.5, -0.5, 0.0)
	ht.rotationAxis = mgl32.Vec3{0.0, 0.0, 1.0}.Normalize()

	var err error
	ht.shader, err = glutils.Shader("_assets/getting_started/5.transformations/transform.vs",
		"_assets/getting_started/5.transformations/transform.frag", "")
	if err != nil {
		return err
	}
	gl.UseProgram(ht.shader)

	vertices := []float32{
		// Positions      // Colors       // Texture Coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // Top Right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // Top Left
	}

	indices := []uint32{ // Note that we start from 0!
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}

	gl.GenVertexArrays(1, &ht.vao)
	gl.GenBuffers(1, &ht.vbo)
	gl.GenBuffers(1, &ht.ebo)

	gl.BindVertexArray(ht.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, ht.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*glutils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ht.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*glutils.GL_FLOAT32_SIZE, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*glutils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*glutils.GL_FLOAT32_SIZE, gl.PtrOffset(3*glutils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(1)
	// TexCoord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*glutils.GL_FLOAT32_SIZE, gl.PtrOffset(6*glutils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0) // Unbind VAO

	// Texture 1
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/container.png"); err != nil {
		return err
	} else {
		ht.texture1 = tex
	}

	// Texture 2
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/awesomeface.png"); err != nil {
		return err
	} else {
		ht.texture2 = tex
	}

	return nil
}

func (ht *HelloTransformations) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(ht.Color32.R, ht.Color32.G, ht.Color32.B, ht.Color32.A)

	// Bind Textures using texture units
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture1)
	loc1 := gl.GetUniformLocation(ht.shader, gl.Str("ourTexture1\x00"))
	gl.Uniform1i(loc1, 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture2)
	loc2 := gl.GetUniformLocation(ht.shader, gl.Str("ourTexture2\x00"))
	gl.Uniform1i(loc2, 1)

	// Activate shader
	gl.UseProgram(ht.shader)

	// rotate
	transform := ht.translationMat.Mul4(mgl32.HomogRotate3D(float32(glfw.GetTime()), ht.rotationAxis))
	transformLoc := gl.GetUniformLocation(ht.shader, gl.Str("transform\x00"))
	// here we create a pointer from the first element of the matrix?
	// read up and update this comm
	gl.UniformMatrix4fv(transformLoc, 1, false, &transform[0])

	// Draw container
	gl.BindVertexArray(ht.vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)

}

func (ht *HelloTransformations) Close() {
	gl.DeleteVertexArrays(1, &ht.vao)
	gl.DeleteBuffers(1, &ht.vbo)
	gl.DeleteBuffers(1, &ht.ebo)
	gl.DeleteProgram(ht.shader)
}
