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
	"fmt"
)

const WIDTH = 800
const HEIGHT = 600

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func keyCallBack(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Println(scancode, action)
	newIndex := sketchIndex
	if action == glfw.Press && scancode == 124 {
		newIndex = sketchIndex + 1
		if newIndex > len(theSketches) - 1 {
			newIndex = len(theSketches) - 1
		}
	}
	if action == glfw.Press && scancode == 123 {
		newIndex = sketchIndex - 1
		if newIndex < 0 {
			newIndex = 0
		}
	}
	if action == glfw.Press && newIndex != sketchIndex {
		sketchIndex = newIndex
		theSketch.Close()
		theSketch = theSketches[newIndex]
		theSketch.Setup()
	}
	theSketch.HandleKeyboard(key, scancode, action, mods)
}

func resizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	theSketch.Update()
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

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	return window
}

var theSketch sketches.Sketch
var sketchIndex = 0
var theSketches []sketches.Sketch

func main() {
	// init GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// create window
	window := setup()

	// make a slice of pointers to sketch instances
	theSketches = []sketches.Sketch{
		&sketches.HelloWindow{Window: window},
		&sketches.HelloCube{Window: window},
		&sketches.HelloTriangle{Window: window},
		&sketches.HelloSquare{Window: window},
	}

	//Use a pointer to the sketch in order to call mutating functions
	theSketch = theSketches[sketchIndex]
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