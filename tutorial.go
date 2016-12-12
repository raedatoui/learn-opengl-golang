package main

import (
	"fmt"
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
	fmt.Println("init1")
	runtime.LockOSThread()
}

func keyCallBack(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	theSketch.HandleKeyboard(key, scancode, action, mods)
}

func resizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	//render(w)
}

func setup() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
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
	//window.SetKeyCallback(keyCallBack)
	//window.SetFramebufferSizeCallback(resizeCallback)

	//width, height := window.GetFramebufferSize()
	//gl.Viewport(0, 0, int32(width), int32(height))
	return window
}

var theSketch sketches.HelloCube

func main() {
	// init GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// create window
	window := setup()
	theSketch = sketches.HelloCube{Window: window}
	theSketch.Setup()
	fmt.Printf("Prog: %d, Vao: %d, Vbo: %d, Tex: %d, MU: %d, Model:\n%v\n",
			theSketch.Program, theSketch.Vao, theSketch.Vbo, theSketch.Texture, theSketch.ModelUniform, theSketch.Model)
	// loop
	for !window.ShouldClose() {
		theSketch.Update()

		//Render
		theSketch.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
	glfw.Terminate()
}

func init() {
	fmt.Println("init2")
	dir, err := utils.ImportPathToDir("github.com/raedatoui/learn-opengl")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	fmt.Println(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}