package main

import (
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/raedatoui/learn-opengl/sketches"
	"github.com/raedatoui/learn-opengl/utils"
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
	render(w)
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

	// Resize Callback
	window.SetFramebufferSizeCallback(resizeCallback)

	//Keyboard Callback
	window.SetKeyCallback(keyCallBack)

	//width, height := window.GetFramebufferSize()
	//gl.Viewport(0, 0, int32(width), int32(height))
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

	//Use a pointer to the sketch in order to call mutating functions
	theSketch = &sketches.HelloCube{Window: window}
	theSketch.Setup()

	// loop
	for !window.ShouldClose() {
		// Update
		theSketch.Update()

		//Render
		render(window)

		// Poll Events
		glfw.PollEvents()
	}
	theSketch.Close()
	glfw.Terminate()
}

func init() {
	dir, err := utils.ImportPathToDir("github.com/raedatoui/learn-opengl")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

func render(window *glfw.Window) {
	theSketch.Draw()
	window.SwapBuffers()
}