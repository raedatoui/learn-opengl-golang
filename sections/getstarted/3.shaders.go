package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
	"math"
)

type HelloShaders struct {
	sections.BaseSketch
	shader *glutils.Shader
	va     glutils.VertexArray
}

func (hs *HelloShaders) createShader(v, f string) error {
	var err error
	hs.shader, err = glutils.NewShader(v, f, "")

	if err != nil {
		return err
	}
	return nil
}

func (hs *HelloShaders) createBuffers() {
	var vertices = []float32{
		// Positions      // Colors
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // Bottom Left
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // Top
	}
	attr := make(map[uint32]int32)
	attr[hs.shader.Attributes["position"]] = 3
	attr[hs.shader.Attributes["color"]] = 3
	hs.va = glutils.VertexArray{
		Data:       vertices,
		Stride:     6,
		Normalized: false,
		DrawMode:   gl.STATIC_DRAW,
		Attributes: attr,
	}

	hs.va.Setup()
}

func (hs *HelloShaders) InitGL() error {
	hs.Name = "3a. Shaders"

	if err := hs.createShader(
		"_assets/getting_started/3.shaders/basic.vs",
		"_assets/getting_started/3.shaders/basic.frag"); err != nil {
		return err
	}

	hs.createBuffers()

	return nil
}

func (hs *HelloShaders) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hs.Color32.R, hs.Color32.G, hs.Color32.B, hs.Color32.A)

	// Draw the triangle
	gl.UseProgram(hs.shader.Program)
	gl.BindVertexArray(hs.va.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (hs *HelloShaders) Close() {
	hs.shader.Delete()
	hs.va.Delete()
}

type ShaderEx1 struct {
	HelloShaders
	timeValue           float64
	greenValue          float32
	vertexColorLocation int32
}

func (hs *ShaderEx1) InitGL() error {
	hs.Name = "3b. Shaders Ex1"
	if err := hs.createShader(
		"_assets/getting_started/3.shaders/basic.vs",
		"_assets/getting_started/3.shaders/uniform.frag"); err != nil {
		return err
	}

	hs.createBuffers()
	return nil
}

func (hs *ShaderEx1) Update() {
	hs.timeValue = glfw.GetTime()
	hs.greenValue = float32(math.Sin(hs.timeValue)/2) + 0.5
}

func (hs *ShaderEx1) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hs.Color32.R, hs.Color32.G, hs.Color32.B, hs.Color32.A)

	// Draw the triangle
	gl.UseProgram(hs.shader.Program)
	gl.Uniform4f(hs.shader.Uniforms["ourColor"], 0.0, hs.greenValue, 0.0, 1.0)
	gl.BindVertexArray(hs.va.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (hs *ShaderEx1) GetSubHeader() string {
	return "with a uniform updated by gl.Uniform4f"
}

type ShaderEx2 struct {
	HelloShaders
}

func (hs *ShaderEx2) InitGL() error {
	hs.Name = "3c. Shaders Ex2"
	if err := hs.createShader(
		"_assets/getting_started/3.shaders/reverse.vs",
		"_assets/getting_started/3.shaders/basic.frag"); err != nil {
		return err
	}

	hs.createBuffers()

	return nil
}

type ShaderEx3 struct {
	HelloShaders
}

func (hs *ShaderEx3) InitGL() error {
	hs.Name = "3b. Shaders Ex3"
	if err := hs.createShader(
		"_assets/getting_started/3.shaders/offset.vs",
		"_assets/getting_started/3.shaders/basic.frag"); err != nil {
		return err
	}

	hs.createBuffers()

	gl.UseProgram(hs.shader.Program)
	gl.Uniform1f(hs.shader.Uniforms["xOffset"], 0.5)

	return nil
}

type ShaderEx4 struct {
	HelloShaders
}

func (hs *ShaderEx4) InitGL() error {
	hs.Name = "3c. Shaders Ex4"
	if err := hs.createShader(
		"_assets/getting_started/3.shaders/ex4.vs",
		"_assets/getting_started/3.shaders/ex4.frag"); err != nil {
		return err
	}

	hs.createBuffers()

	return nil
}
