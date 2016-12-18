package sketches

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/learn-opengl/utils"
)
var keys map[glfw.Key]bool
var lastX float64 = 400
var lastY float64 = 300
var firstMouse bool = true

type HelloCamera struct {
	Window             *glfw.Window
	Shader             uint32
	Vao, Vbo, Ebo      uint32
	Texture1, Texture2 uint32
	Transform          mgl32.Mat4
	CubePositions      []mgl32.Vec3
	Camera             utils.Camera
	DeltaTime, LastFrame float64
}

func (sketch *HelloCamera) Setup() {
	var err error
	sketch.Shader, err = utils.Shader("sketches/assets/6.coordinates/coordinate.vs",
		"sketches/assets/6.coordinates/coordinate.frag", "")
	if err != nil {
		panic(err)
	}
	gl.UseProgram(sketch.Shader)

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

	sketch.CubePositions = []mgl32.Vec3{
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{2.0, 5.0, -15.0},
		mgl32.Vec3{-1.5, -2.2, -2.5},
		mgl32.Vec3{-3.8, -2.0, -12.3},
		mgl32.Vec3{2.4, -0.4, -3.5},
		mgl32.Vec3{-1.7, 3.0, -7.5},
		mgl32.Vec3{1.3, -2.0, -2.5},
		mgl32.Vec3{1.5, 2.0, -2.5},
		mgl32.Vec3{1.5, 0.2, -1.5},
		mgl32.Vec3{-1.3, 1.0, -1.5},
	}

	// ====================
	// Camera
	// ====================
	sketch.Camera = utils.NewCamera(
		mgl32.Vec3{0.0, 0.0, 3.0},
		mgl32.Vec3{0.0, 1.0, 3.0},
		utils.YAW, utils.PITCH,
	)

	gl.GenVertexArrays(1, &sketch.Vao)
	gl.GenBuffers(1, &sketch.Vbo)
	gl.GenBuffers(1, &sketch.Ebo)

	gl.BindVertexArray(sketch.Vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, sketch.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5 * utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// TexCoord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 5 * utils.GL_FLOAT32_SIZE, gl.PtrOffset(3*utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0) // Unbind VAO

	// ====================
	// Texture 1
	// ====================
	gl.GenTextures(1, &sketch.Texture1)
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	rgba, err := utils.ImageToPixelData("sketches/assets/images/container.png")
	if err != nil {
		panic(err)
	}
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	// ====================
	// Texture 2
	// ====================
	gl.GenTextures(1, &sketch.Texture2)
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture2)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	rgba, err = utils.ImageToPixelData("sketches/assets/images/awesomeface.png")
	if err != nil {
		panic(err)
	}
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	gl.Enable(gl.DEPTH_TEST)

	keys = make(map[glfw.Key]bool)
}

func (sketch *HelloCamera) Update() {
	// Set frame time
	currentFrame := glfw.GetTime()
	sketch.DeltaTime = currentFrame - sketch.LastFrame
	sketch.LastFrame = currentFrame
	if keys[glfw.KeyW] {
		sketch.Camera.ProcessKeyboard(utils.FORWARD, float32(sketch.DeltaTime))
	}
	if keys[glfw.KeyS] {
		sketch.Camera.ProcessKeyboard(utils.BACKWARD, float32(sketch.DeltaTime))
	}
	if keys[glfw.KeyA] {
		sketch.Camera.ProcessKeyboard(utils.LEFT, float32(sketch.DeltaTime))
	}
	if keys[glfw.KeyD] {
		sketch.Camera.ProcessKeyboard(utils.RIGHT, float32(sketch.DeltaTime))
	}
}

func (sketch *HelloCamera) Draw() {
	// Render
	// Clear the colorbuffer
	gl.ClearColor(0.8, 0.3, 0.8, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Bind Textures using texture units
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture1)
	loc1 := gl.GetUniformLocation(sketch.Shader, gl.Str("ourTexture1\x00"))
	gl.Uniform1i(loc1, 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture2)
	loc2 := gl.GetUniformLocation(sketch.Shader, gl.Str("ourTexture2\x00"))
	gl.Uniform1i(loc2, 1)

	// Activate shader
	gl.UseProgram(sketch.Shader)

	// Create camera transformations
	view := sketch.Camera.GetViewMatrix()
	projection := mgl32.Perspective(sketch.Camera.Zoom, 800.0/600.0, 0.1, 1000.0)

	// Get their uniform location
	modelLoc := gl.GetUniformLocation(sketch.Shader, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(sketch.Shader,  gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(sketch.Shader, gl.Str("projection\x00"))
	// Pass the matrices to the shader
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	// Note: currently we set the projection matrix each frame,
	// but since the projection matrix rarely changes it's often best practice to set it outside the main loop only once.
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	// Draw container
	gl.BindVertexArray(sketch.Vao)

	for i := 0; i < 10; i++ {
		// Calculate the model matrix for each object and pass it to shader before drawing
		model := mgl32.Translate3D(
			sketch.CubePositions[i][0],
			sketch.CubePositions[i][1],
			sketch.CubePositions[i][2])
		//angle := 20.0 * float32(i)
		//if i % 3 == 0 {
		angle := float32(glfw.GetTime()) * float32(i+1)
		//}  // Every 3rd iteration (including the first) we set the angle using GLFW's time function.

		model = model.Mul4(mgl32.HomogRotate3D(angle, mgl32.Vec3{1.0, 0.3, 0.5}))
		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
	gl.BindVertexArray(0)
}

func (sketch *HelloCamera) Close() {
	gl.DeleteVertexArrays(1, &sketch.Vao)
	gl.DeleteBuffers(1, &sketch.Vbo)
	gl.DeleteBuffers(1, &sketch.Ebo)
	gl.UseProgram(0)
}

func (sketch *HelloCamera) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		sketch.Window.SetShouldClose(true)
	}
    if action == glfw.Press {
		keys[key] = true
	} else if action == glfw.Release {
		keys[key] = false
	}
}

func (sketch *HelloCamera) HandleMousePosition(xpos, ypos float64) {
	if firstMouse {
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}

    xoffset := xpos - lastX
    yoffset := lastY - ypos  // Reversed since y-coordinates go from bottom to left

    lastX = xpos
    lastY = ypos

    sketch.Camera.ProcessMouseMovement(float32(xoffset), float32(yoffset), true)
}

func (sketch *HelloCamera) HandleScroll(xoff, yoff float64) {
	sketch.Camera.ProcessMouseScroll(float32(yoff))
}
