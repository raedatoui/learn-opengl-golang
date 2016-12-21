package sketches

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl/utils"
)

type HelloTextures struct {
	Window             *glfw.Window
	Shader             uint32
	Vao, Vbo, Ebo      uint32
	Texture1, Texture2 uint32
}

func (sketch *HelloTextures) Setup() error {
	var err error
	sketch.Shader, err = utils.Shader("sketches/_assets/4.textures/texture.vs",
		"sketches/_assets/4.textures/texture.frag", "")
	if err != nil {
		panic(err)
	}
	gl.UseProgram(sketch.Shader)

	vertices := []float32{
		// Positions      // Colors       // Texture Coords
		 0.5,  0.5, 0.0,  1.0, 0.0, 0.0,  1.0, 1.0, // Top Right
		 0.5, -0.5, 0.0,  0.0, 1.0, 0.0,  1.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0,  0.0, 0.0, 1.0,  0.0, 0.0, // Bottom Left
		-0.5,  0.5, 0.0,  1.0, 1.0, 0.0,  0.0, 1.0, // Top Left
	}

	indices := []uint32{ // Note that we start from 0!
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}

	gl.GenVertexArrays(1, &sketch.Vao)
	gl.GenBuffers(1, &sketch.Vbo)
	gl.GenBuffers(1, &sketch.Ebo)

	gl.BindVertexArray(sketch.Vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, sketch.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, sketch.Ebo)
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
	gl.GenTextures(1, &sketch.Texture1)
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	rgba, err := utils.ImageToPixelData("sketches/_assets/images/container.png")
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
	gl.GenTextures(1, &sketch.Texture2)
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture2)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	rgba, err = utils.ImageToPixelData("sketches/_assets/images/awesomeface.png")
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

func (sketch *HelloTextures) Update() {

}

func (sketch *HelloTextures) Draw() {
	// Bind Textures using texture units
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture1)
	loc1 := gl.GetUniformLocation(sketch.Shader, gl.Str("ourTexture1\x00"))
	gl.Uniform1i(loc1, 0)

	gl.ActiveTexture(gl.TEXTURE1);
	gl.BindTexture(gl.TEXTURE_2D, sketch.Texture2)
	loc2 := gl.GetUniformLocation(sketch.Shader, gl.Str("ourTexture2\x00"))
	gl.Uniform1i(loc2, 1);

	// Activate shader
	gl.UseProgram(sketch.Shader)

	// Draw container
	gl.BindVertexArray(sketch.Vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func (sketch *HelloTextures) Close() {
	gl.DeleteVertexArrays(1, &sketch.Vao)
	gl.DeleteBuffers(1, &sketch.Vbo)
	gl.DeleteBuffers(1, &sketch.Ebo)
	gl.UseProgram(0)
}

func (sketch *HelloTextures) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		sketch.Window.SetShouldClose(true)
	}
}

func (sketch *HelloTextures) HandleMousePosition(xpos, ypos float64) {

}

func (sketch *HelloTextures) HandleScroll(xoff, yoff float64) {

}