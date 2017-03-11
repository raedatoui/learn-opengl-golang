package modelloading

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"path"
	"fmt"
)

type ModelLoading struct {
	sections.BaseSketch
	shader               glutils.Shader
	model                glutils.Model
	camera               glutils.Camera
	deltaTime, lastFrame float64
	lastX, lastY         float64
	firstMouse           bool
	w, a, s, d           bool
}

func (ml *ModelLoading) InitGL() error {
	ml.firstMouse = false
	ml.camera = glutils.NewCamera(
		mgl32.Vec3{0.0, 0.0, 3.0},
		mgl32.Vec3{0.0, 1.0, 3.0},
		glutils.YAW, glutils.PITCH,
	)
	ml.Name = "3. Model Loading"
	// Setup and compile our shaders
	ml.shader, _ = glutils.NewShader("_assets/model_loading/shader.vs",
		"_assets/model_loading/shader.frag", "")
	// Load models
	ml.model, _ = glutils.NewModel("_assets/objects/nanosuit/", "nanosuit.obj", false)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	return nil
}

func (ml *ModelLoading) Update() {
	// Set frame time
	currentFrame := glfw.GetTime()
	ml.deltaTime = currentFrame - ml.lastFrame
	ml.lastFrame = currentFrame
	if ml.w {
		ml.camera.ProcessKeyboard(glutils.FORWARD, ml.deltaTime)
	}
	if ml.s {
		ml.camera.ProcessKeyboard(glutils.BACKWARD, ml.deltaTime)
	}
	if ml.a {
		ml.camera.ProcessKeyboard(glutils.LEFT, ml.deltaTime)
	}
	if ml.d {
		ml.camera.ProcessKeyboard(glutils.RIGHT, ml.deltaTime)
	}
}

func (ml *ModelLoading) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(ml.Color32.R, ml.Color32.G, ml.Color32.B, ml.Color32.A)

	gl.UseProgram(ml.shader.Program)

	// Transformation matrices
	projection := mgl32.Perspective(float32(ml.camera.Zoom), sections.Ratio, 0.1, 100.0)
	view := ml.camera.GetViewMatrix()

	gl.UniformMatrix4fv(ml.shader.Uniforms["view"], 1, false, &view[0])
	gl.UniformMatrix4fv(ml.shader.Uniforms["projection"], 1, false, &projection[0])

	// Draw the loaded model
	model := mgl32.Translate3D(0, -1.75, 0.0)        // Translate it down a bit so it's at the center of the scene
	model = model.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2)) // It's a bit too big for our scene, so scale it down

	gl.UniformMatrix4fv(ml.shader.Uniforms["model"], 1, false, &model[0])
	ml.model.Draw(ml.shader.Program)
}

func (lc *ModelLoading) HandleMousePosition(xpos, ypos float64) {
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

func (ml *ModelLoading) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey, keys map[glfw.Key]bool) {
	ml.w = keys[glfw.KeyW]
	ml.a = keys[glfw.KeyA]
	ml.s = keys[glfw.KeyS]
	ml.d = keys[glfw.KeyD]
}

func (ml *ModelLoading) HandleScroll(xoff, yoff float64) {
	ml.camera.ProcessMouseScroll(yoff)
}

func (ml *ModelLoading) HandleFiles(names []string) {
	f := path.Base(names[0])
	dir := path.Dir(names[0]) + "/"
	fmt.Println(f, dir)
	ml.model.Dispose()
	ml.model, _ = glutils.NewModel(dir, f, false)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
}

func (ml *ModelLoading) Close() {
	ml.model.Dispose()
	gl.UseProgram(0)
}
