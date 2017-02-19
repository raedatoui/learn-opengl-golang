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
	shader             glutils.Shader
	va glutils.VertexArray
	texture1, texture2 uint32
	translationMat     mgl32.Mat4
	rotationAxis       mgl32.Vec3
}

func (ht *HelloTransformations) InitGL() error {
	ht.Name = "5. Transformations"

	ht.translationMat = mgl32.Translate3D(0.5, -0.5, 0.0)
	ht.rotationAxis = mgl32.Vec3{0.0, 0.0, 1.0}.Normalize()

	var err error
	ht.shader, err = glutils.NewShader(
		"_assets/getting_started/5.transformations/transform.vs",
		"_assets/getting_started/5.transformations/transform.frag", "")
	if err != nil {
		return err
	}

	vertices := []float32{
		// Positions      // Colors       // Texture Coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // Top Right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // Top Left
	}

	attr := make(glutils.AttributesMap)
	attr[ht.shader.Attributes["position"]] = [2]int{3, 0}
	attr[ht.shader.Attributes["texCoord"]] = [2]int{2, 6}

	indices := []uint32{ // Note that we start from 0!
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}
	ht.va = glutils.VertexArray{
		Data: vertices,
		Indices: indices,
		DrawMode: gl.STATIC_DRAW,
		Normalized: false,
		Stride: 8,
		Attributes: attr,
	}
	ht.va.Setup()

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
	loc1 := gl.GetUniformLocation(ht.shader.Program, gl.Str("ourTexture1\x00"))
	gl.Uniform1i(loc1, 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture2)
	loc2 := gl.GetUniformLocation(ht.shader.Program, gl.Str("ourTexture2\x00"))
	gl.Uniform1i(loc2, 1)

	// Activate shader
	gl.UseProgram(ht.shader.Program)

	// rotate
	transform := ht.translationMat.Mul4(mgl32.HomogRotate3D(float32(glfw.GetTime()), ht.rotationAxis))
	transformLoc := gl.GetUniformLocation(ht.shader.Program, gl.Str("transform\x00"))
	// here we create a pointer from the first element of the matrix?
	// read up and update this comm
	gl.UniformMatrix4fv(transformLoc, 1, false, &transform[0])

	// Draw container
	gl.BindVertexArray(ht.va.Vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)

}

func (ht *HelloTransformations) Close() {
	ht.shader.Delete()
	ht.va.Delete()
}
