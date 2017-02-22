package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloCamera struct {
	HelloCoordinates
	camera               glutils.Camera
	deltaTime, lastFrame float64
	w, a, s, d           bool
	lastX, lastY         float64
	firstMouse           bool
}

func (hc *HelloCamera) GetHeader() string {
	return "7. Camera (use WSDA and mouse)"
}

func (hc *HelloCamera) InitGL() error {
	hc.camera = glutils.NewCamera(
		mgl32.Vec3{0.0, 0.0, 3.0},
		mgl32.Vec3{0.0, 1.0, 3.0},
		glutils.YAW, glutils.PITCH,
	)
	hc.lastX = sections.WIDTH / 2
	hc.lastY = sections.HEIGHT / 2
	hc.firstMouse = true

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

func (hc *HelloCamera) Update() {
	// Set frame time
	currentFrame := glfw.GetTime()
	hc.deltaTime = currentFrame - hc.lastFrame
	hc.lastFrame = currentFrame
	if hc.w {
		hc.camera.ProcessKeyboard(glutils.FORWARD, hc.deltaTime)
	}
	if hc.s {
		hc.camera.ProcessKeyboard(glutils.BACKWARD, hc.deltaTime)
	}
	if hc.a {
		hc.camera.ProcessKeyboard(glutils.LEFT, hc.deltaTime)
	}
	if hc.d {
		hc.camera.ProcessKeyboard(glutils.RIGHT, hc.deltaTime)
	}
}

func (hc *HelloCamera) setTransformations() {
	// Create transformations
	view := hc.camera.GetViewMatrix()
	projection := mgl32.Perspective(float32(hc.camera.Zoom), sections.Ratio, 0.1, 1000.0)

	// Pass the matrices to the shader
	gl.UniformMatrix4fv(hc.shader.Uniforms["view"], 1, false, &view[0])
	// Note: currently we set the projection matrix each frame,
	// but since the projection matrix rarely changes it's often best practice to set it outside the main loop only once.
	gl.UniformMatrix4fv(hc.shader.Uniforms["projection"], 1, false, &projection[0])
}

func (hc *HelloCamera) Draw() {
	hc.clear()
	// Activate shader
	gl.UseProgram(hc.shader.Program)
	hc.setTextures()
	hc.setTransformations()
	hc.renderVertexArray()
}

func (hc *HelloCamera) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey, keys map[glfw.Key]bool) {
	hc.w = keys[glfw.KeyW]
	hc.a = keys[glfw.KeyA]
	hc.s = keys[glfw.KeyS]
	hc.d = keys[glfw.KeyD]
}

func (hc *HelloCamera) HandleMousePosition(xpos, ypos float64) {
	if hc.firstMouse {
		hc.lastX = xpos
		hc.lastY = ypos
		hc.firstMouse = false
	}

	xoffset := xpos - hc.lastX
	yoffset := hc.lastY - ypos // Reversed since y-coordinates go from bottom to left

	hc.lastX = xpos
	hc.lastY = ypos

	hc.camera.ProcessMouseMovement(xoffset, yoffset, true)
}

func (hc *HelloCamera) HandleScroll(xoff, yoff float64) {
	hc.camera.ProcessMouseScroll(yoff)
}
