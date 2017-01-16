package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"math"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type HelloShaders struct {
	sections.BaseSketch
	vao, vbo uint32
	shader   uint32
}

func (hs *HelloShaders) InitGL() error {
	hs.Name = "3a. Shaders"

	var err error
	hs.shader, err = utils.Shader(
		"_assets/3.shaders/basic.vs",
		"_assets/3.shaders/basic.frag", "")

	if err != nil {
		return err
	}

	var vertices = []float32{
		// Positions      // Colors
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // Bottom Left
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // Top
	}
	gl.GenVertexArrays(1, &hs.vao)
	gl.GenBuffers(1, &hs.vbo)

	gl.BindVertexArray(hs.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, hs.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	// position uniform
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	//color uniform
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(3*utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return nil
}

func (hs *HelloShaders) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hs.Color32.R, hs.Color32.G, hs.Color32.B, hs.Color32.A)

	// Draw the triangle
	gl.UseProgram(hs.shader)
	gl.BindVertexArray(hs.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (hs *HelloShaders) Close() {
	gl.DeleteVertexArrays(1, &hs.vao)
	gl.DeleteBuffers(1, &hs.vbo)
	gl.UseProgram(0)
}

type HelloShaderUniform struct {
	sections.BaseSketch
	vao, vbo uint32
	shader   uint32
	timeValue float64
	greenValue float32
	vertexColorLocation int32 //=
}

func (hs *HelloShaderUniform) InitGL() error {
	hs.Name = "3b. Shaders"

	var err error
	hs.shader, err = utils.Shader(
		"_assets/3.shaders/basic.vs",
		"_assets/3.shaders/uniform.frag", "")

	if err != nil {
		return err
	}
	hs.vertexColorLocation = gl.GetUniformLocation(hs.shader, gl.Str("ourColor\x00"))

	var vertices = []float32{
		// Positions      // Colors
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // Bottom Left
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // Top
	}
	gl.GenVertexArrays(1, &hs.vao)
	gl.GenBuffers(1, &hs.vbo)

	gl.BindVertexArray(hs.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, hs.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	// position uniform
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	//color uniform
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*utils.GL_FLOAT32_SIZE, gl.PtrOffset(3*utils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return nil
}

func (hs *HelloShaderUniform) Update() {
	hs.timeValue = glfw.GetTime()
	hs.greenValue = float32(math.Sin(hs.timeValue) / 2) + 0.5
}
func (hs *HelloShaderUniform) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hs.Color32.R, hs.Color32.G, hs.Color32.B, hs.Color32.A)

	// Draw the triangle
	gl.UseProgram(hs.shader)
	gl.Uniform4f(hs.vertexColorLocation, 0.0, hs.greenValue, 0.0, 1.0)
	gl.BindVertexArray(hs.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (hs *HelloShaderUniform) Close() {
	gl.DeleteVertexArrays(1, &hs.vao)
	gl.DeleteBuffers(1, &hs.vbo)
	gl.UseProgram(0)
}

func (hs *HelloShaderUniform) GetSubHeader() string {
	return "with a uniform updated by gl.Uniform4f"
}