package lighting

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

type BasicSpecular struct {
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

func (bc *BasicSpecular) InitGL() error {
	bc.Name = "2. Basic Specular Lighting"

	// Camera
	bc.camera = utils.NewCamera(
		mgl32.Vec3{0.0, 0.0, 3.0},
		mgl32.Vec3{0.0, 1.0, 3.0},
		utils.YAW, utils.PITCH,
	)
	bc.lastX = utils.WIDTH / 2.0
	bc.lastY = utils.HEIGHT / 2.0
	// Light attributes
	bc.lightPos = mgl32.Vec3{1.2, 1.0, 2.0}

	// Deltatime
	bc.deltaTime = 0.0 // Time between current frame and last frame
	bc.lastFrame = 0.0 // Time of last frame

	if sh, err := utils.Shader("_assets/9.basic_lighting/lighting.vs", "_assets/9.basic_lighting/lighting.frag", ""); err != nil {
		return err
	} else {
		bc.lightingShader = sh
	}
	if sh, err := utils.Shader("_assets/9.basic_lighting/lamp.vs", "_assets/9.basic_lighting/lamp.frag", ""); err != nil {
		return err
	} else {
		bc.lampShader = sh
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
	gl.GenVertexArrays(1, &bc.containerVAO)
	gl.GenBuffers(1, &bc.vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, bc.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(bc.containerVAO)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// Normal attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(3*utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	// Then, we set the light's VAO (VBO stays the same. After all, the vertices are the same for the light object (also a 3D cube))
	gl.GenVertexArrays(1, &bc.lightVAO)
	gl.BindVertexArray(bc.lightVAO)
	// We only need to bind to the VBO (to link it with glVertexAttribPointer), no need to fill it; the VBO's data already contains all we need.
	gl.BindBuffer(gl.ARRAY_BUFFER, bc.vbo)
	// Set the vertex attributes (only position data for the lamp))
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	return nil
}

func (bc *BasicSpecular) Update() {
	// Set frame time
	currentFrame := glfw.GetTime()
	bc.deltaTime = currentFrame - bc.lastFrame
	bc.lastFrame = currentFrame
	if bc.w {
		bc.camera.ProcessKeyboard(utils.FORWARD, float32(bc.deltaTime))
	}
	if bc.s {
		bc.camera.ProcessKeyboard(utils.BACKWARD, float32(bc.deltaTime))
	}
	if bc.a {
		bc.camera.ProcessKeyboard(utils.LEFT, float32(bc.deltaTime))
	}
	if bc.d {
		bc.camera.ProcessKeyboard(utils.RIGHT, float32(bc.deltaTime))
	}
}

func (bc *BasicSpecular) Draw() {
	// Clear the colorbuffer
	gl.ClearColor(bc.Color32.R, bc.Color32.G, bc.Color32.B, bc.Color32.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Use corresponding shader when setting uniforms/drawing objects
	gl.UseProgram(bc.lightingShader)
	objectColorLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("objectColor\x00"))
	lightColorLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("lightColor\x00"))
	lightPosLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("lightPos\x00"))
	viewPosLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("viewPos\x00"))
	gl.Uniform3f(objectColorLoc, 1.0, 0.5, 0.31)
	gl.Uniform3f(lightColorLoc, 1.0, 0.5, 1.0)
	gl.Uniform3f(lightPosLoc, bc.lightPos[0], bc.lightPos[1], bc.lightPos[2])
	gl.Uniform3f(viewPosLoc, bc.camera.Position[0], bc.camera.Position[1], bc.camera.Position[2])

	// Create camera transformations
	view := bc.camera.GetViewMatrix()
	projection := mgl32.Perspective(bc.camera.Zoom, float32(utils.WIDTH)/float32(utils.HEIGHT), 0.1, 100.0)
	// Get the uniform locations
	modelLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(bc.lightingShader, gl.Str("projection\x00"))
	// Pass the matrices to the shader
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	// Draw the container (using container's vertex attributes)
	gl.BindVertexArray(bc.containerVAO)
	model := mgl32.Translate3D(0, 0, 0.0)
	angle := float32(glfw.GetTime())
	model = model.Mul4(mgl32.HomogRotate3D(angle, mgl32.Vec3{1.0, 0.3, 0.5}))
	gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

	// Also draw the lamp object, again binding the appropriate shader
	gl.UseProgram(bc.lampShader)
	// Get location objects for the matrices on the lamp shader (these could be different on a different shader)
	modelLoc = gl.GetUniformLocation(bc.lampShader, gl.Str("model\x00"))
	viewLoc = gl.GetUniformLocation(bc.lampShader, gl.Str("view\x00"))
	projLoc = gl.GetUniformLocation(bc.lampShader, gl.Str("projection\x00"))
	// Set matrices
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	// Get location objects for the matrices on the lamp shader (these could be different on a different shader)
	model2 := mgl32.Translate3D(bc.lightPos[0], bc.lightPos[1], bc.lightPos[2])
	model2 = model2.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2)) // Make it a smaller cube
	gl.UniformMatrix4fv(modelLoc, 1, false, &model2[0])
	// Draw the light object (using light's vertex attributes)
	gl.BindVertexArray(bc.lightVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)
}

func (bc *BasicSpecular) Close() {

}

func (lc *BasicSpecular) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey, keys map[glfw.Key]bool) {
	lc.w = keys[glfw.KeyW]
	lc.a = keys[glfw.KeyA]
	lc.s = keys[glfw.KeyS]
	lc.d = keys[glfw.KeyD]
}

func (bc *BasicSpecular) HandleMousePosition(xpos, ypos float64) {
	if bc.firstMouse {
		bc.lastX = xpos
		bc.lastY = ypos
		bc.firstMouse = false
	}

	xoffset := xpos - bc.lastX
	yoffset := bc.lastY - ypos // Reversed since y-coordinates go from bottom to left

	bc.lastX = xpos
	bc.lastY = ypos

	bc.camera.ProcessMouseMovement(float32(xoffset), float32(yoffset), true)
}
