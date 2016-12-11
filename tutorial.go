package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
	"github.com/raedatoui/learn-opengl/sketches"
)

const WIDTH = 800
const HEIGHT = 600


func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func keyCallBack(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	theSketch.HandleKeyboard(key, scancode, action, mods)
}

func resizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	//theSketch.draw()
	w.SwapBuffers()
}

func setup() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "Test Tutorial", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()


	// Initialize Glow - this is the equivalent of glew
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	return window
}

var theSketch sketches.Sketch

func main() {
	// init GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// create window
	window := setup()
	theSketch = sketches.SimpleKetch{Window: window}
	// Resize Callback
	window.SetKeyCallback(keyCallBack)
	window.SetFramebufferSizeCallback(resizeCallback)

	// loop
	for !window.ShouldClose() {
		// Maintenance
		glfw.PollEvents()

		theSketch.Update()

		//Render
		theSketch.Draw()
		window.SwapBuffers()
	}
	glfw.Terminate()
}
