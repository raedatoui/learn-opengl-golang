package lighting

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/raedatoui/glutils"
)

type BasicSpecular struct {
	LightingColors
}

func (bc *BasicSpecular) GetHeader() string {
	return "2. Basic Specular Lighting"
}

func (lc *BasicSpecular) getVertices() []float32 {
	return []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
	}
}
func (bc *BasicSpecular) initContainers(vertices []float32) {

	attr := glutils.NewAttributesMap()
	attr.Add(bc.lightingShader.Attributes["position"], 3, 0)
	attr.Add(bc.lightingShader.Attributes["normal"], 3, 3)
	bc.containerVa = glutils.VertexArray{
		Data:       vertices,
		Stride:     6,
		DrawMode:   gl.STATIC_DRAW,
		Normalized: false,
		Attributes: attr,
	}
	bc.containerVa.Setup()

	attr2 := glutils.NewAttributesMap()
	attr2.Add(bc.lampShader.Attributes["position"], 3, 0)
	bc.lightVa = glutils.VertexArray{
		Vbo:        bc.containerVa.Vbo,
		Attributes: attr2,
		DrawMode:   gl.STATIC_DRAW,
		Normalized: false,
		Stride:     6,
	}
	bc.lightVa.Setup()
}

func (bc *BasicSpecular) setLightingUniforms() {
	gl.Uniform3f(bc.lightingShader.Uniforms["objectColor"], 1.0, 0.5, 0.31)
	gl.Uniform3f(bc.lightingShader.Uniforms["lightColor"], 1.0, 0.5, 1.0)
	gl.Uniform3f(bc.lightingShader.Uniforms["lightPos"], bc.lightPos[0], bc.lightPos[1], bc.lightPos[2])
	gl.Uniform3f(bc.lightingShader.Uniforms["viewPos"], bc.camera.Position[0], bc.camera.Position[1], bc.camera.Position[2])
}
func (bc *BasicSpecular) InitGL() error {
	bc.initCamera()
	if err := bc.initShaders(
		"_assets/lighting/2.basic/lighting.vs",
		"_assets/lighting/2.basic/lighting.frag",
		"_assets/lighting/2.basic/lamp.vs",
		"_assets/lighting/2.basic/lamp.frag",
	); err != nil {
		return err
	}
	bc.initContainers(bc.getVertices())
	return nil
}

func (bc *BasicSpecular) Draw() {
	bc.clear()
	gl.UseProgram(bc.lightingShader.Program)
	bc.setLightingUniforms()
	v, p := bc.getCameraTransforms()
	bc.transformShader(bc.lightingShader, v, p)
	bc.drawContainer()

	// Also draw the lamp object, again binding the appropriate shader
	gl.UseProgram(bc.lampShader.Program)
	// Set matrices
	bc.transformShader(bc.lampShader, v, p)
	bc.drawLamp()
}

//func (bc *BasicSpecular) Draw() {
//	// Clear the colorbuffer
//	gl.ClearColor(bc.Color32.R, bc.Color32.G, bc.Color32.B, bc.Color32.A)
//	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
//
//	// Use corresponding shader when setting uniforms/drawing objects
//	gl.UseProgram(bc.lightingShader)
//	objectColorLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("objectColor\x00"))
//	lightColorLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("lightColor\x00"))
//	lightPosLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("lightPos\x00"))
//	viewPosLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("viewPos\x00"))
//	gl.Uniform3f(objectColorLoc, 1.0, 0.5, 0.31)
//	gl.Uniform3f(lightColorLoc, 1.0, 0.5, 1.0)
//	gl.Uniform3f(lightPosLoc, bc.lightPos[0], bc.lightPos[1], bc.lightPos[2])
//	gl.Uniform3f(viewPosLoc, bc.camera.Position[0], bc.camera.Position[1], bc.camera.Position[2])
//
//	// Create camera transformations
//	view := bc.camera.GetViewMatrix()
//	projection := mgl32.Perspective(float32(bc.camera.Zoom), sections.RATIO, 0.1, 100.0)
//	// Get the uniform locations
//	modelLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("model\x00"))
//	viewLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("view\x00"))
//	projLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("projection\x00"))
//	// Pass the matrices to the shader
//	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
//	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])
//
//	// Draw the container (using container's vertex attributes)
//	gl.BindVertexArray(bc.containerVAO)
//	model := mgl32.Translate3D(0, 0, 0.0)
//	angle := float32(glfw.GetTime())
//	model = model.Mul4(mgl32.HomogRotate3D(angle, bc.rotationAxis))
//	gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
//	gl.DrawArrays(gl.TRIANGLES, 0, 36)
//	gl.BindVertexArray(0)
//
//	// Also draw the lamp object, again binding the appropriate shader
//	gl.UseProgram(bc.lampShader)
//	// Get location objects for the matrices on the lamp shader (these could be different on a different shader)
//	modelLoc = gl.GetUniformLocation(bc.lampShader, gl.Str("model\x00"))
//	viewLoc = gl.GetUniformLocation(bc.lampShader, gl.Str("view\x00"))
//	projLoc = gl.GetUniformLocation(bc.lampShader, gl.Str("projection\x00"))
//	// Set matrices
//	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
//	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])
//
//	// Get location objects for the matrices on the lamp shader (these could be different on a different shader)
//	model2 := bc.lightPositionMat.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2)) // Make it a smaller cube
//	gl.UniformMatrix4fv(modelLoc, 1, false, &model2[0])
//	// Draw the light object (using light's vertex attributes)
//	gl.BindVertexArray(bc.lightVAO)
//	gl.DrawArrays(gl.TRIANGLES, 0, 36)
//	gl.BindVertexArray(0)
//}
