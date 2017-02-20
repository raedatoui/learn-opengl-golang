package lighting

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Materials struct {
	BasicSpecular
}

func (m *Materials) GetHeader() string {
	return "2. Basic Specular Lighting"
}

func (m *Materials) InitGL() error {
	m.initCamera()
	if err := m.initShaders(
		"_assets/lighting/3.materials/materials.vs",
		"_assets/lighting/3.materials/materials.frag",
		"_assets/lighting/3.materials/lamp.vs",
		"_assets/lighting/3.materials/lamp.frag",
	); err != nil {
		return err
	}
	m.initContainers(m.getVertices())
	return nil
}

func (m *Materials) setLightingUniforms() {
	gl.Uniform3f(m.lightingShader.Uniforms["light.position"], m.lightPos.X(), m.lightPos.Y(), m.lightPos.Z())
	gl.Uniform3f(m.lightingShader.Uniforms["viewPos"], m.camera.Position.X(), m.camera.Position.Y(), m.camera.Position.Z())

	// Set lights properties
	lightColor := mgl32.Vec3{
		float32(math.Sin(glfw.GetTime() * 2.0)),
		float32(math.Sin(glfw.GetTime() * 0.7)),
		float32(math.Sin(glfw.GetTime() * 1.3)),
	}

	// Decrease the influence
	diffuseColor := mgl32.Vec3{
		lightColor.X() * 0.5,
		lightColor.Y() * 0.5,
		lightColor.Z() * 0.5,
	}
	// Low influence
	ambientColor := mgl32.Vec3{
		diffuseColor.X() * 0.2,
		diffuseColor.Y() * 0.2,
		diffuseColor.Z() * 0.2,
	}

	gl.Uniform3f(m.lightingShader.Uniforms["light.ambient"], ambientColor.X(), ambientColor.Y(), ambientColor.Z())
	gl.Uniform3f(m.lightingShader.Uniforms["light.diffuse"], diffuseColor.X(), diffuseColor.Y(), diffuseColor.Z())
	gl.Uniform3f(m.lightingShader.Uniforms["light.specular"], 1.0, 1.0, 1.0)
	// Set material properties
	gl.Uniform3f(m.lightingShader.Uniforms["material.ambient"], 1.0, 0.5, 0.31)
	gl.Uniform3f(m.lightingShader.Uniforms["material.diffuse"], 1.0, 0.5, 0.31)
	gl.Uniform3f(m.lightingShader.Uniforms["material.specular"], 0.5, 0.5, 0.5) // Specular doesn't have full effect on this object's material
	gl.Uniform1f(m.lightingShader.Uniforms["material.shininess"], 32.0)
}

func (m *Materials) Draw() {
	m.clear()
	gl.UseProgram(m.lightingShader.Program)
	m.setLightingUniforms()
	v, p := m.getCameraTransforms()
	m.transformShader(m.lightingShader, v, p)
	m.drawContainer()

	// Also draw the lamp object, again binding the appropriate shader
	gl.UseProgram(m.lampShader.Program)
	// Set matrices
	m.transformShader(m.lampShader, v, p)
	m.drawLamp()
}
