package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloTextures struct {
	sections.BaseSketch
	shader             glutils.Shader
	va                 glutils.VertexArray
	texture1, texture2 uint32
}

func (ht *HelloTextures) getShaders() []string {
	return []string{"_assets/getting_started/4.textures/texture.vs",
		"_assets/getting_started/4.textures/texture.frag"}
}
func (ht *HelloTextures) getVertices() []float32 {
	return []float32{
		// Positions      // Colors       // Texture Coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // Top Right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // Top Left
	}
}
func (ht *HelloTextures) createBuffers(vertices []float32) {
	indices := []uint32{ // Note that we start from 0!
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}
	attr := glutils.NewAttributesMap()
	attr.Add(ht.shader.Attributes["position"], 3, 0)
	attr.Add(ht.shader.Attributes["color"], 3, 2)
	attr.Add(ht.shader.Attributes["texCoord"], 2, 6)

	ht.va = glutils.VertexArray{
		Data:       vertices,
		Indices:    indices,
		Stride:     8,
		Normalized: false,
		DrawMode:   gl.STATIC_DRAW,
		Attributes: attr,
	}
	ht.va.Setup()
}
func (ht *HelloTextures) InitGL() error {
	ht.Name = "4. Textures"

	var err error
	shaders := ht.getShaders()
	ht.shader, err = glutils.NewShader(shaders[0], shaders[1], "")
	if err != nil {
		return err
	}

	ht.createBuffers(ht.getVertices())

	// ====================
	// Texture 1
	// ====================
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/container.png"); err != nil {
		return err
	} else {
		ht.texture1 = tex
	}

	// ====================
	// Texture 2
	// ====================
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/awesomeface.png"); err != nil {
		return err
	} else {
		ht.texture2 = tex
	}

	return nil
}

func (ht *HelloTextures) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(ht.Color32.R, ht.Color32.G, ht.Color32.B, ht.Color32.A)

	// Bind Textures using texture units
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture1)
	gl.Uniform1i(ht.shader.Uniforms["ourTexture1"], 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture2)
	gl.Uniform1i(ht.shader.Uniforms["ourTexture2"], 1)

	// Activate shader
	gl.UseProgram(ht.shader.Program)

	// Draw container
	gl.BindVertexArray(ht.va.Vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func (ht *HelloTextures) Close() {
	ht.shader.Delete()
	ht.va.Delete()
}

type TexturesEx1 struct {
	HelloTextures
}

func (ht *TexturesEx1) getShaders() []string {
	return []string{"_assets/getting_started/4.textures/texture.vs",
		"_assets/getting_started/4.textures/textureex1.frag"}
}
func (ht *TexturesEx1) InitGL() error {
	ht.Name = "4a. Textures Ex1"

	var err error
	shaders := ht.getShaders()
	ht.shader, err = glutils.NewShader(shaders[0], shaders[1], "")
	if err != nil {
		return err
	}

	ht.createBuffers(ht.getVertices())

	// ====================
	// Texture 1
	// ====================
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/container.png"); err != nil {
		return err
	} else {
		ht.texture1 = tex
	}

	// ====================
	// Texture 2
	// ====================
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/awesomeface.png"); err != nil {
		return err
	} else {
		ht.texture2 = tex
	}

	return nil
}

func (ht *TexturesEx1) GetSubHeader() string {
	return "flip the happy face in the frag shader"
}

type TexturesEx2 struct {
	TexturesEx1
}

func (ht *TexturesEx2) getVertices() []float32 {
	return []float32{
		// Positions      // Colors       // Texture Coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 2.0, 2.0, // Top Right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 2.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 2.0, // Top Left
	}
}
func (ht *TexturesEx2) InitGL() error {
	ht.Name = "4b. Textures Ex2"

	var err error
	shaders := ht.getShaders()
	ht.shader, err = glutils.NewShader(shaders[0], shaders[1], "")
	if err != nil {
		return err
	}

	ht.createBuffers(ht.getVertices())

	// ====================
	// Texture 1
	// ====================
	if tex, err := glutils.NewTexture(gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.NEAREST, gl.NEAREST, "_assets/images/container.png"); err != nil {
		return err
	} else {
		ht.texture1 = tex
	}

	// ====================
	// Texture 2
	// ====================
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.NEAREST, gl.NEAREST, "_assets/images/awesomeface.png"); err != nil {
		return err
	} else {
		ht.texture2 = tex
	}

	return nil
}

func (ht *TexturesEx2) GetSubHeader() string {
	return "scale tex coord by 2 and set nearest filter on tex to see individual pixels"
}

type TexturesEx3 struct {
	TexturesEx1
}

func (ht *TexturesEx3) getVertices() []float32 {
	return []float32{
		// Positions     // Colors        // Texture Coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.55, 0.55, // Top Right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 0.55, 0.45, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.45, 0.45, // Bottom Left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.45, 0.55, // Top Left
	}
}
func (ht *TexturesEx3) InitGL() error {
	ht.Name = "4c. Textures Ex3"

	var err error
	shaders := ht.getShaders()
	ht.shader, err = glutils.NewShader(shaders[0], shaders[1], "")
	if err != nil {
		return err
	}
	ht.createBuffers(ht.getVertices())

	// Texture 1
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.NEAREST, gl.NEAREST, "_assets/images/container.png"); err != nil {
		return err
	} else {
		ht.texture1 = tex
	}

	// Texture 2
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.NEAREST, gl.NEAREST, "_assets/images/awesomeface.png"); err != nil {
		return err
	} else {
		ht.texture2 = tex
	}

	// mixvalue uniform
	return nil
}

type TexturesEx4 struct {
	HelloTextures
	mixValue float32
}

func (ht *TexturesEx3) GetSubHeader() string {
	return "scale tex down and test various wrapping"
}

func (ht *TexturesEx4) getShaders() []string {
	return []string{"_assets/getting_started/4.textures/texture.vs",
		"_assets/getting_started/4.textures/textureex4.frag"}
}
func (ht *TexturesEx4) InitGL() error {
	ht.Name = "4d. Textures Ex4"

	var err error
	shaders := ht.getShaders()
	ht.shader, err = glutils.NewShader(shaders[0], shaders[1], "")
	if err != nil {
		return err
	}
	ht.createBuffers(ht.getVertices())

	// Texture 1
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.NEAREST, gl.NEAREST, "_assets/images/container.png"); err != nil {
		return err
	} else {
		ht.texture1 = tex
	}

	// Texture 2
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.NEAREST, gl.NEAREST, "_assets/images/awesomeface.png"); err != nil {
		return err
	} else {
		ht.texture2 = tex
	}

	return nil
}

func (ht *TexturesEx4) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey, keys map[glfw.Key]bool) {
	if keys[glfw.KeyUp] {
		ht.mixValue += 0.1
		if ht.mixValue >= 1.0 {
			ht.mixValue = 1.0
		}
	}
	if keys[glfw.KeyDown] {
		ht.mixValue -= 0.1
		if ht.mixValue <= 0.0 {
			ht.mixValue = 0.0
		}
	}
}

func (ht *TexturesEx4) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(ht.Color32.R, ht.Color32.G, ht.Color32.B, ht.Color32.A)

	// Activate shader
	gl.UseProgram(ht.shader.Program)

	// Bind Textures using texture units
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture1)
	gl.Uniform1i(ht.shader.Uniforms["ourTexture1"], 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ht.texture2)
	gl.Uniform1i(ht.shader.Uniforms["ourTexture2"], 1)

	gl.Uniform1f(ht.shader.Uniforms["mixValue"], ht.mixValue)

	// Draw container
	gl.BindVertexArray(ht.va.Vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func (ht *TexturesEx4) GetSubHeader() string {
	return "cross fade tex using up/dwn arrows, setting opacity as frag uniform "
}
