package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloTransformations struct {
	sections.BaseSketch
	shader             uint32
	vao, vbo, ebo      uint32
	texture1, texture2 uint32
	transform          mgl32.Mat4
}

func (ht *HelloTransformations) Setup(w *glfw.Window, f *utils.Font) error {
	ht.Window = w
	ht.Font = f
	ht.Name = "5. Transformations"
	ht.Color = utils.RandColor()

	var err error
	ht.shader, err = utils.Shader("_assets/5.transformations/transform.vs",
		"_assets/5.transformations/transform.frag", "")
	if err != nil {
		return err
	}
	gl.UseProgram(ht.shader)

	vertices := []float32{
		// Positions      // Colors       // Texture Coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // Top Right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // Top Left
	}

	indices := []uint32{ // Note that we start from 0!
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}

	gl.GenVertexArrays(1, &ht.vao)
	gl.GenBuffers(1, &ht.vbo)
	gl.GenBuffers(1, &ht.ebo)

	gl.BindVertexArray(ht.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, ht.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ht.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*utils.GL_FLOAT32_SIZE, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*utils.GL_FLOAT32_SIZE, gl.PtrOffset(3*utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(1)
	// TexCoord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*utils.GL_FLOAT32_SIZE, gl.PtrOffset(6*utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0) // Unbind VAO

	// ====================
	// Texture 1
	// ====================
	gl.GenTextures(1, &ht.texture1)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	rgba, err := utils.ImageToPixelData("_assets/images/container.png")
	if err != nil {
		return err
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
	gl.GenTextures(1, &ht.texture2)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture2)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	rgba, err = utils.ImageToPixelData("_assets/images/awesomeface.png")
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

	return nil
}

func (ht *HelloTransformations) Update() {

}

func (ht *HelloTransformations) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(ht.Color.R, ht.Color.G, ht.Color.B, ht.Color.A)

	// Bind Textures using texture units
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture1)
	loc1 := gl.GetUniformLocation(ht.shader, gl.Str("ourTexture1\x00"))
	gl.Uniform1i(loc1, 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture2)
	loc2 := gl.GetUniformLocation(ht.shader, gl.Str("ourTexture2\x00"))
	gl.Uniform1i(loc2, 1)

	// Activate shader
	gl.UseProgram(ht.shader)

	// create transform
	ht.transform = mgl32.Translate3D(0.5, -0.5, 0.0)
	// rotate
	ht.transform = ht.transform.Mul4(mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0.0, 0.0, 1.0}))
	transformLoc := gl.GetUniformLocation(ht.shader, gl.Str("transform\x00"))
	// here we create a pointer from the first element of the matrix?
	// read up and update this comm
	gl.UniformMatrix4fv(transformLoc, 1, false, &ht.transform[0])

	// Draw container
	gl.BindVertexArray(ht.vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)

	ht.Font.SetColor(0.0, 0.0, 0.0, 1.0)
	ht.Font.Printf(30, 30, 0.5, ht.Name)
}

func (ht *HelloTransformations) Close() {
	gl.DeleteVertexArrays(1, &ht.vao)
	gl.DeleteBuffers(1, &ht.vbo)
	gl.DeleteBuffers(1, &ht.ebo)
	gl.UseProgram(0)
}

func (hc *HelloTransformations) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {

}

func (ht *HelloTransformations) HandleMousePosition(xpos, ypos float64) {

}

func (ht *HelloTransformations) HandleScroll(xoff, yoff float64) {

}
