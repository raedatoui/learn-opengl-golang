package lighting

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"math"
)

type Materials struct {
	sections.BaseSketch
	lightingShader, lampShader  uint32
	vbo, containerVAO, lightVAO uint32
	lastX                       float64
	lastY                       float64
	firstMouse                  bool
	deltaTime, lastFrame        float64
	camera                      utils.Camera
	lightPos                    mgl32.Vec3
	w, a, s, d                  bool
}

func (m *Materials) InitGL() error {
	m.Name = "2. Basic Specular Lighting"

	m.firstMouse = true

	// Camera
	m.camera = utils.NewCamera(
		mgl32.Vec3{0.0, 0.0, 3.0},
		mgl32.Vec3{0.0, 1.0, 3.0},
		utils.YAW, utils.PITCH,
	)
	m.lastX = utils.WIDTH / 2.0
	m.lastY = utils.HEIGHT / 2.0
	// Light attributes
	m.lightPos = mgl32.Vec3{1.2, 1.0, 2.0}

	// Deltatime
	m.deltaTime = 0.0 // Time between current frame and last frame
	m.lastFrame = 0.0 // Time of last frame

	if sh, err := utils.Shader(
		"_assets/lighting/3.materials/materials.vs",
		"_assets/lighting/3.materials/materials.frag", ""); err != nil {
		return err
	} else {
		m.lightingShader = sh
	}
	if sh, err := utils.Shader(
		"_assets/lighting/3.materials/lamp.vs",
		"_assets/lighting/3.materials/lamp.frag", ""); err != nil {
		return err
	} else {
		m.lampShader = sh
	}

	vertices := []float32{
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

	// First, set the container's VAO (and VBO)
	gl.GenVertexArrays(1, &m.containerVAO)
	gl.GenBuffers(1, &m.vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(m.containerVAO)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// Normal attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(3*utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	// Then, we set the light's VAO (VBO stays the same. After all, the vertices are the same for the light object (also a 3D cube))
	gl.GenVertexArrays(1, &m.lightVAO)
	gl.BindVertexArray(m.lightVAO)
	// We only need to bind to the VBO (to link it with glVertexAttribPointer), no need to fill it; the VBO's data already contains all we need.
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	// Set the vertex attributes (only position data for the lamp))
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)

	return nil
}

func (m *Materials) Update() {
	// Set frame time
	currentFrame := glfw.GetTime()
	m.deltaTime = currentFrame - m.lastFrame
	m.lastFrame = currentFrame
	if m.w {
		m.camera.ProcessKeyboard(utils.FORWARD, m.deltaTime)
	}
	if m.s {
		m.camera.ProcessKeyboard(utils.BACKWARD, m.deltaTime)
	}
	if m.a {
		m.camera.ProcessKeyboard(utils.LEFT, m.deltaTime)
	}
	if m.d {
		m.camera.ProcessKeyboard(utils.RIGHT, m.deltaTime)
	}
}

func (m *Materials) Draw() {
	gl.ClearColor(m.Color32.R, m.Color32.G, m.Color32.B, m.Color32.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Use cooresponding shader when setting uniforms/drawing objects
	gl.UseProgram(m.lightingShader)

	lightPosLoc := gl.GetUniformLocation(m.lightingShader, gl.Str("light.position\x00"))
	viewPosLoc  := gl.GetUniformLocation(m.lightingShader, gl.Str("viewPos\x00"))
	gl.Uniform3f(lightPosLoc, m.lightPos.X(), m.lightPos.Y(), m.lightPos.Z())
	gl.Uniform3f(viewPosLoc, m.camera.Position.X(), m.camera.Position.Y(), m.camera.Position.Z())
	// Set lights properties
	lightColor := mgl32.Vec3{
		float32(math.Sin(glfw.GetTime() * 2.0)),
		float32(math.Sin(glfw.GetTime() * 0.7)),
		float32(math.Sin(glfw.GetTime() * 1.3)),
	}

	 // Decrease the influence
	diffuseColor := mgl32.Vec3{
		lightColor.X() * 0.5,
		lightColor.Y() * 0.5,
		lightColor.Z() * 0.5,
	}
	// Low influence
	ambientColor := mgl32.Vec3{
		diffuseColor.X() * 0.2,
		diffuseColor.Y() * 0.2,
		diffuseColor.Z() * 0.2,
	}

	gl.Uniform3f(gl.GetUniformLocation(
		m.lightingShader, gl.Str("light.ambient\x00")),  ambientColor.X(), ambientColor.Y(), ambientColor.Z())
	gl.Uniform3f(
		gl.GetUniformLocation(m.lightingShader, gl.Str("light.diffuse\x00")),
		diffuseColor.X(), diffuseColor.Y(), diffuseColor.Z())

	gl.Uniform3f(
		gl.GetUniformLocation(m.lightingShader, gl.Str("light.specular\x00")),
		1.0, 1.0, 1.0)
	// Set material properties
	gl.Uniform3f(
		gl.GetUniformLocation(m.lightingShader, gl.Str("material.ambient\x00")),
		1.0, 0.5, 0.31)
	gl.Uniform3f(
		gl.GetUniformLocation(m.lightingShader, gl.Str("material.diffuse\x00")),
		1.0, 0.5, 0.31)
	gl.Uniform3f(
		gl.GetUniformLocation(m.lightingShader, gl.Str("material.specular\x00")),
		0.5, 0.5, 0.5) // Specular doesn't have full effect on this object's material
	gl.Uniform1f(gl.GetUniformLocation(m.lightingShader, gl.Str("material.shininess\x00")), 32.0)

	// Create camera transformations
	view := m.camera.GetViewMatrix()
	projection := mgl32.Perspective(float32(m.camera.Zoom), float32(utils.WIDTH)/float32(utils.HEIGHT), 0.1, 100.0)
	// Get the uniform locations
	modelLoc := gl.GetUniformLocation(m.lightingShader, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(m.lightingShader, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(m.lightingShader, gl.Str("projection\x00"))
	// Pass the matrices to the shader
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	// Draw the container (using container's vertex attributes)
	gl.BindVertexArray(m.containerVAO)
	model := mgl32.Translate3D(0, 0, 0.0)
	angle := float32(glfw.GetTime())
	model = model.Mul4(mgl32.HomogRotate3D(angle, mgl32.Vec3{1.0, 0.3, 0.5}))
	gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

	// Also draw the lamp object, again binding the appropriate shader
	gl.UseProgram(m.lampShader)
	// Get location objects for the matrices on the lamp shader (these could be different on a different shader)
	modelLoc = gl.GetUniformLocation(m.lampShader, gl.Str("model\x00"))
	viewLoc = gl.GetUniformLocation(m.lampShader, gl.Str("view\x00"))
	projLoc = gl.GetUniformLocation(m.lampShader, gl.Str("projection\x00"))
	// Set matrices
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	// Get location objects for the matrices on the lamp shader (these could be different on a different shader)
	model2 := mgl32.Translate3D(m.lightPos[0], m.lightPos[1], m.lightPos[2])
	model2 = model2.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2)) // Make it a smaller cube
	gl.UniformMatrix4fv(modelLoc, 1, false, &model2[0])
	// Draw the light object (using light's vertex attributes)
	gl.BindVertexArray(m.lightVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

}

func (m *Materials) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey, keys map[glfw.Key]bool) {
	m.w = keys[glfw.KeyW]
	m.a = keys[glfw.KeyA]
	m.s = keys[glfw.KeyS]
	m.d = keys[glfw.KeyD]
}

func (m *Materials) HandleMousePosition(xpos, ypos float64) {
	if m.firstMouse {
		m.lastX = xpos
		m.lastY = ypos
		m.firstMouse = false
	}

	xoffset := xpos - m.lastX
	yoffset := m.lastY - ypos // Reversed since y-coordinates go from bottom to left

	m.lastX = xpos
	m.lastY = ypos

	m.camera.ProcessMouseMovement(xoffset, yoffset, true)
}

func (m *Materials) HandleScroll(xoff, yoff float64) {
	m.camera.ProcessMouseScroll(yoff)
}

func (m *Materials) Close() {
	gl.DeleteVertexArrays(1, &m.lightVAO)
	gl.DeleteVertexArrays(1, &m.containerVAO)
	gl.DeleteBuffers(1, &m.vbo)
	gl.DeleteProgram(m.lightingShader)
	gl.DeleteProgram(m.lampShader)
}