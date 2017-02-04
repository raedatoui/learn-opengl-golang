package lighting

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

type LightingColors struct {
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

func (lc *LightingColors) InitGL() error {
	lc.Name = "1. Colors"

	// Camera
	lc.camera = utils.NewCamera(
		mgl32.Vec3{0.0, 0.0, 3.0},
		mgl32.Vec3{0.0, 1.0, 3.0},
		utils.YAW, utils.PITCH,
	)
	lc.lastX = utils.WIDTH / 2.0
	lc.lastY = utils.HEIGHT / 2.0

	// Light attributes
	lc.lightPos = mgl32.Vec3{1.2, 1.0, 2.0}

	// Deltatime
	lc.deltaTime = 0.0 // Time between current frame and last frame
	lc.lastFrame = 0.0 // Time of last frame

	if sh, err := utils.Shader(
		"_assets/lighting/1.colors/colors.vs",
		"_assets/lighting/1.colors/colors.frag", ""); err != nil {
		return err
	} else {
		lc.lightingShader = sh
	}
	if sh, err := utils.Shader(
		"_assets/lighting/1.colors/lamp.vs",
		"_assets/lighting/1.colors/lamp.frag", ""); err != nil {
		return err
	} else {
		lc.lampShader = sh
	}

	vertices := []float32{
		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,

		-0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, -0.5, 0.5,

		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, 0.5,
		-0.5, 0.5, 0.5,

		0.5, 0.5, 0.5,
		0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,

		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		-0.5, -0.5, 0.5,
		-0.5, -0.5, -0.5,

		-0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
	}

	// First, set the container's VAO (and VBO)
	gl.GenVertexArrays(1, &lc.containerVAO)
	gl.GenBuffers(1, &lc.vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, lc.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(lc.containerVAO)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)

	// Then, we set the light's VAO (VBO stays the same. After all, the vertices are the same for the light object (also a 3D cube))
	gl.GenVertexArrays(1, &lc.lightVAO)
	gl.BindVertexArray(lc.lightVAO)
	// We only need to bind to the VBO (to link it with glVertexAttribPointer), no need to fill it; the VBO's data already contains all we need.
	gl.BindBuffer(gl.ARRAY_BUFFER, lc.vbo)
	// Set the vertex attributes (only position data for the lamp))
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)

	return nil
}

func (lc *LightingColors) Update() {
	// Set frame time
	currentFrame := glfw.GetTime()
	lc.deltaTime = currentFrame - lc.lastFrame
	lc.lastFrame = currentFrame
	if lc.w {
		lc.camera.ProcessKeyboard(utils.FORWARD, lc.deltaTime)
	}
	if lc.s {
		lc.camera.ProcessKeyboard(utils.BACKWARD, lc.deltaTime)
	}
	if lc.a {
		lc.camera.ProcessKeyboard(utils.LEFT, lc.deltaTime)
	}
	if lc.d {
		lc.camera.ProcessKeyboard(utils.RIGHT, lc.deltaTime)
	}
}

func (lc *LightingColors) Draw() {
	// Clear the colorbuffer
	gl.ClearColor(lc.Color32.R, lc.Color32.G, lc.Color32.B, lc.Color32.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Use cooresponding shader when setting uniforms/drawing objects
	gl.UseProgram(lc.lightingShader)
	objectColorLoc := gl.GetUniformLocation(lc.lightingShader, gl.Str("objectColor\x00"))
	lightColorLoc := gl.GetUniformLocation(lc.lightingShader, gl.Str("lightColor\x00"))
	gl.Uniform3f(objectColorLoc, 1.0, 0.5, 0.31)
	gl.Uniform3f(lightColorLoc, 1.0, 0.5, 1.0)

	// Create camera transformations
	view := lc.camera.GetViewMatrix()
	projection := mgl32.Perspective(float32(lc.camera.Zoom), utils.RATIO, 0.1, 100.0)

	// Get the uniform locations
	modelLoc := gl.GetUniformLocation(lc.lightingShader, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(lc.lightingShader, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(lc.lightingShader, gl.Str("projection\x00"))
	// Pass the matrices to the shader
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	// Draw the container (using container's vertex attributes)
	gl.BindVertexArray(lc.containerVAO)
	model := mgl32.Translate3D(0, 0, 0.0)
	angle := float32(glfw.GetTime())
	model = model.Mul4(mgl32.HomogRotate3D(angle, mgl32.Vec3{1.0, 0.3, 0.5}))
	gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

	// Also draw the lamp object, again binding the appropriate shader
	gl.UseProgram(lc.lampShader)
	// Get location objects for the matrices on the lamp shader (these could be different on a different shader)
	modelLoc = gl.GetUniformLocation(lc.lampShader, gl.Str("model\x00"))
	viewLoc = gl.GetUniformLocation(lc.lampShader, gl.Str("view\x00"))
	projLoc = gl.GetUniformLocation(lc.lampShader, gl.Str("projection\x00"))
	// Set matrices
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	//model2 = model2.Mul4(mgl32.Translate3D(lc.lightPos[0], lc.lightPos[1], lc.lightPos[2]))
	model2 := mgl32.Translate3D(lc.lightPos[0], lc.lightPos[1], lc.lightPos[2])
	model2 = model2.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2)) // Make it a smaller cube
	gl.UniformMatrix4fv(modelLoc, 1, false, &model2[0])
	// Draw the light object (using light's vertex attributes)
	gl.BindVertexArray(lc.lightVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

}

func (lc *LightingColors) Close() {
	gl.DeleteVertexArrays(1, &lc.lightVAO)
	gl.DeleteVertexArrays(1, &lc.containerVAO)
	gl.DeleteBuffers(1, &lc.vbo)
	gl.DeleteProgram(lc.lightingShader)
	gl.DeleteProgram(lc.lampShader)
}

func (lc *LightingColors) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey, keys map[glfw.Key]bool) {
	lc.w = keys[glfw.KeyW]
	lc.a = keys[glfw.KeyA]
	lc.s = keys[glfw.KeyS]
	lc.d = keys[glfw.KeyD]
}

func (lc *LightingColors) HandleMousePosition(xpos, ypos float64) {
	if lc.firstMouse {
		lc.lastX = xpos
		lc.lastY = ypos
		lc.firstMouse = false
	}

	xoffset := xpos - lc.lastX
	yoffset := lc.lastY - ypos // Reversed since y-coordinates go from bottom to left

	lc.lastX = xpos
	lc.lastY = ypos

	lc.camera.ProcessMouseMovement(xoffset, yoffset, true)
}

func (lc *LightingColors) HandleScroll(xoff, yoff float64) {
	lc.camera.ProcessMouseScroll(yoff)
}
