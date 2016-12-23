package sketches

import "github.com/go-gl/glfw/v3.2/glfw"

type LightingColors struct {
	window  *glfw.Window
	keys map[glfw.Key]bool
	lastX, lastY  float64
	firstMouse bool
}