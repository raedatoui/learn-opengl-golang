package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloCoordinates struct {
	sections.BaseSketch
	shader             glutils.Shader
	va                 glutils.VertexArray
	texture1, texture2 uint32
	transform          mgl32.Mat4
	cubePositions      []mgl32.Mat4
	rotationAxis       mgl32.Vec3
}

func (hc *HelloCoordinates) GetHeader() string {
	return "6. Coordinate Systems"
}

func (hc *HelloCoordinates) createShader() error {
	var err error
	hc.shader, err = glutils.NewShader(
		"_assets/getting_started/6.coordinates/coordinate.vs",
		"_assets/getting_started/6.coordinates/coordinate.frag", "")
	if err != nil {
		return err
	}
	return nil
}
func (hc *HelloCoordinates) createBuffers() error {
	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	attr := glutils.NewAttributesMap()
	attr.Add(hc.shader.Attributes["position"], 3, 0)
	attr.Add(hc.shader.Attributes["texCoord"], 2, 3)

	hc.va = glutils.VertexArray{
		Data:       vertices,
		Stride:     5,
		DrawMode:   gl.STATIC_DRAW,
		Normalized: false,
		Attributes: attr,
	}
	hc.va.Setup()

	hc.rotationAxis = mgl32.Vec3{1.0, 0.3, 0.5}.Normalize()

	hc.cubePositions = []mgl32.Mat4{
		mgl32.Translate3D(0.0, 0.0, 0.0),
		mgl32.Translate3D(2.0, 5.0, -15.0),
		mgl32.Translate3D(-1.5, -2.2, -2.5),
		mgl32.Translate3D(-3.8, -2.0, -12.3),
		mgl32.Translate3D(2.4, -0.4, -3.5),
		mgl32.Translate3D(-1.7, 3.0, -7.5),
		mgl32.Translate3D(1.3, -2.0, -2.5),
		mgl32.Translate3D(1.5, 2.0, -2.5),
		mgl32.Translate3D(1.5, 0.2, -1.5),
		mgl32.Translate3D(-1.3, 1.0, -1.5),
	}
	return nil
}
func (hc *HelloCoordinates) createTextures() error {
	// Texture 1
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/container.png"); err != nil {
		return err
	} else {
		hc.texture1 = tex
	}

	// Texture 2
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/awesomeface.png"); err != nil {
		return err
	} else {
		hc.texture2 = tex
	}
	return nil
}
func (hc *HelloCoordinates) InitGL() error {
	if err := hc.createShader(); err != nil {
		return err
	}
	if err := hc.createBuffers(); err != nil {
		return err
	}
	if err := hc.createTextures(); err != nil {
		return err
	}
	return nil
}

func (hc *HelloCoordinates) clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hc.Color32.R, hc.Color32.G, hc.Color32.B, hc.Color32.A)
}
func (hc *HelloCoordinates) setTextures() {
	// Bind Textures using texture units
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, hc.texture1)
	gl.Uniform1i(hc.shader.Uniforms["ourTexture1"], 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, hc.texture2)
	gl.Uniform1i(hc.shader.Uniforms["ourTexture2"], 1)
}
func (hc *HelloCoordinates) setTransformations() {
	// Create transformations
	view := mgl32.Translate3D(0.0, 0.0, -3.0)
	projection := mgl32.Perspective(45.0, sections.RATIO, 0.1, 100.0)
	// Pass the matrices to the shader
	gl.UniformMatrix4fv(hc.shader.Uniforms["view"], 1, false, &view[0])
	// Note: currently we set the projection matrix each frame,
	// but since the projection matrix rarely changes it's often best practice to set it outside the main loop only once.
	gl.UniformMatrix4fv(hc.shader.Uniforms["projection"], 1, false, &projection[0])
}
func (hc *HelloCoordinates) renderVertexArray() {
	// Draw container
	gl.BindVertexArray(hc.va.Vao)
	for i := 0; i < 10; i++ {
		// Calculate the model matrix for each object and pass it to shader before drawing
		model := hc.cubePositions[i]

		angle := float32(glfw.GetTime()) * float32(i+1)

		model = model.Mul4(mgl32.HomogRotate3D(angle, hc.rotationAxis))
		gl.UniformMatrix4fv(hc.shader.Uniforms["model"], 1, false, &model[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
	gl.BindVertexArray(0)
}
func (hc *HelloCoordinates) Draw() {
	hc.clear()
	// Activate shader
	gl.UseProgram(hc.shader.Program)
	hc.setTextures()
	hc.setTransformations()
	hc.renderVertexArray()
}

func (hc *HelloCoordinates) Close() {
	hc.shader.Delete()
	hc.va.Delete()
}
