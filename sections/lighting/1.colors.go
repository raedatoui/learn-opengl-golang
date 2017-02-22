package lighting

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type LightingColors struct {
	sections.BaseSketch
	lightingShader, lampShader glutils.Shader
	containerVa, lightVa       glutils.VertexArray
	lastX                      float64
	lastY                      float64
	firstMouse                 bool
	deltaTime, lastFrame       float64
	camera                     glutils.Camera
	lightPos                   mgl32.Vec3
	w, a, s, d                 bool
	rotationAxis               mgl32.Vec3
	translationMat             mgl32.Mat4
	lightPositionMat           mgl32.Mat4
	scaleMat                   mgl32.Mat4
}

func (lc *LightingColors) GetHeader() string {
	return "1. Colors"
}

func (lc *LightingColors) initShaders(v1, f1, v2, f2 string) error {
	if sh, err := glutils.NewShader(v1, f1, ""); err != nil {
		return err
	} else {
		lc.lightingShader = sh
	}
	if sh, err := glutils.NewShader(v2, f2, ""); err != nil {
		return err
	} else {
		lc.lampShader = sh
	}
	return nil
}
func (lc *LightingColors) initCamera() {
	// Camera
	lc.camera = glutils.NewCamera(
		mgl32.Vec3{0.0, 0.0, 3.0},
		mgl32.Vec3{0.0, 1.0, 3.0},
		glutils.YAW, glutils.PITCH,
	)
	lc.lastX = sections.WIDTH / 2.0
	lc.lastY = sections.HEIGHT / 2.0

	// Light attributes
	lc.lightPos = mgl32.Vec3{1.2, 1.0, 2.0}
	lc.lightPositionMat = mgl32.Translate3D(lc.lightPos[0], lc.lightPos[1], lc.lightPos[2])

	// Deltatime
	lc.deltaTime = 0.0 // Time between current frame and last frame
	lc.lastFrame = 0.0 // Time of last frame

	lc.translationMat = mgl32.Translate3D(0, 0, 0.0)
	lc.scaleMat = mgl32.Scale3D(0.2, 0.2, 0.2)
	lc.rotationAxis = mgl32.Vec3{1.0, 0.3, 0.5}.Normalize()
}
func (lc *LightingColors) getVertices() []float32 {
	return []float32{
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
}
func (lc *LightingColors) initContainers(vertices []float32) {

	attr := glutils.NewAttributesMap()
	attr.Add(lc.lightingShader.Attributes["position"], 3, 0)
	lc.containerVa = glutils.VertexArray{
		Data:       vertices,
		Stride:     3,
		DrawMode:   gl.STATIC_DRAW,
		Normalized: false,
		Attributes: attr,
	}
	lc.containerVa.Setup()

	attr2 := glutils.NewAttributesMap()
	attr2.Add(lc.lampShader.Attributes["position"], 3, 0)
	lc.lightVa = glutils.VertexArray{
		Vbo:        lc.containerVa.Vbo,
		Attributes: attr2,
		DrawMode:   gl.STATIC_DRAW,
		Normalized: false,
		Stride:     3,
	}
	lc.lightVa.Setup()
}
func (lc *LightingColors) InitGL() error {
	lc.initCamera()
	if err := lc.initShaders(
		"_assets/lighting/1.colors/colors.vs",
		"_assets/lighting/1.colors/colors.frag",
		"_assets/lighting/1.colors/lamp.vs",
		"_assets/lighting/1.colors/lamp.frag",
	); err != nil {
		return err
	}
	vertices := lc.getVertices()
	lc.initContainers(vertices)
	return nil
}

func (lc *LightingColors) Update() {
	// Set frame time
	currentFrame := glfw.GetTime()
	lc.deltaTime = currentFrame - lc.lastFrame
	lc.lastFrame = currentFrame
	if lc.w {
		lc.camera.ProcessKeyboard(glutils.FORWARD, lc.deltaTime)
	}
	if lc.s {
		lc.camera.ProcessKeyboard(glutils.BACKWARD, lc.deltaTime)
	}
	if lc.a {
		lc.camera.ProcessKeyboard(glutils.LEFT, lc.deltaTime)
	}
	if lc.d {
		lc.camera.ProcessKeyboard(glutils.RIGHT, lc.deltaTime)
	}
}

func (lc *LightingColors) clear() {
	// Clear the colorbuffer
	gl.ClearColor(lc.Color32.R, lc.Color32.G, lc.Color32.B, lc.Color32.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
func (lc *LightingColors) setLightingUniforms() {
	// Use corresponding shader when setting uniforms/drawing objects
	gl.Uniform3f(lc.lightingShader.Uniforms["objectColor"], 1.0, 0.5, 0.31)
	gl.Uniform3f(lc.lightingShader.Uniforms["lightColor"], 1.0, 0.5, 1.0)
}
func (lc *LightingColors) getCameraTransforms() (mgl32.Mat4, mgl32.Mat4) {
	// Create camera transformations
	view := lc.camera.GetViewMatrix()
	projection := mgl32.Perspective(float32(lc.camera.Zoom), sections.Ratio, 0.1, 100.0)
	return view, projection
}
func (lc *LightingColors) transformShader(shader glutils.Shader, view, projection mgl32.Mat4) {
	// Pass the matrices to the shader
	gl.UniformMatrix4fv(shader.Uniforms["view"], 1, false, &view[0])
	gl.UniformMatrix4fv(shader.Uniforms["projection"], 1, false, &projection[0])
}
func (lc *LightingColors) drawContainer() {
	// Draw the container (using container's vertex attributes)
	angle := float32(glfw.GetTime())
	model := lc.translationMat.Mul4(mgl32.HomogRotate3D(angle, lc.rotationAxis))
	gl.UniformMatrix4fv(lc.lightingShader.Uniforms["model"], 1, false, &model[0])

	gl.BindVertexArray(lc.containerVa.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)
}
func (lc *LightingColors) drawLamp() {
	//model2 = model2.Mul4(mgl32.Translate3D(lc.lightPos[0], lc.lightPos[1], lc.lightPos[2]))
	model := lc.lightPositionMat.Mul4(lc.scaleMat) // Make it a smaller cube
	gl.UniformMatrix4fv(lc.lampShader.Uniforms["model"], 1, false, &model[0])
	// Draw the light object (using light's vertex attributes)
	gl.BindVertexArray(lc.lightVa.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)
}
func (lc *LightingColors) Draw() {
	lc.clear()
	gl.UseProgram(lc.lightingShader.Program)
	lc.setLightingUniforms()
	v, p := lc.getCameraTransforms()
	lc.transformShader(lc.lightingShader, v, p)
	lc.drawContainer()

	// Also draw the lamp object, again binding the appropriate shader
	gl.UseProgram(lc.lampShader.Program)
	// Set matrices
	lc.transformShader(lc.lampShader, v, p)

	lc.drawLamp()
}

func (lc *LightingColors) Close() {
	lc.lampShader.Delete()
	lc.lightingShader.Delete()
	lc.lightVa.Delete()
	lc.containerVa.Delete()
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
