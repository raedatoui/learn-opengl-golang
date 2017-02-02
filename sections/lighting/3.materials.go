package lighting

import (
	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"github.com/go-gl/mathgl/mgl32"
)

type Materials struct {
	sections.BaseSketch
	lightingShader, lampShader  uint32
	vbo, containerVAO, lightVAO uint32
	lastX                       float64
	lastY                       float64
	firstMouse                  bool
	deltaTime, lastFrame        float64
	camera                      utils.Camera
	lightPos                    mgl32.Vec3
	w, a, s, d                  bool
}
//
//func (m *Materials) InitGL() {
//	m.lightingShader = utils.Shader()
//			"materials.vs", "materials.frag");
//    Shader lampShader("lamp.vs", "lamp.frag");
//}