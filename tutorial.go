package main

import (
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"fmt"

	"github.com/raedatoui/learn-opengl-golang/sketches"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"github.com/raedatoui/learn-opengl-golang/sketches/getstarted"
)

// WIDTH is the width of the window
const WIDTH = 800
// HEIGHT is the height of the window
const HEIGHT = 600

// Tutorial is a simple struct for sketches
type Tutorial struct {
	Name   string
	Color  utils.ColorA
	Sketch sketches.Sketch
}

var (
	theSketch     sketches.Sketch
	tutorialIndex = 0
	tutorials     []Tutorial
	switching     bool
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func init() {
	dir, err := utils.ImportPathToDir("github.com/raedatoui/learn-opengl-golang")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	if err := os.Chdir(dir); err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

func keyCallBack(w *glfw.Window, k glfw.Key, scancode int, a glfw.Action, mk glfw.ModifierKey) {
	theSketch.HandleKeyboard(k, scancode, a, mk)
	if !switching && a== glfw.Press {
		switching = true
		newIndex := tutorialIndex
		if a == glfw.Press && scancode == 124 {
			newIndex = tutorialIndex + 1
			if newIndex > len(tutorials)-1 {
				newIndex = len(tutorials) - 1
			}
		}
		if a == glfw.Press && scancode == 123 {
			newIndex = tutorialIndex - 1
			if newIndex < 0 {
				newIndex = 0
			}
		}
		if a == glfw.Press && newIndex != tutorialIndex {
			tutorialIndex = newIndex
			theSketch.Close()
			tut := &tutorials[newIndex]
			tut.Color = utils.RandColor()
			theSketch = tut.Sketch
			theSketch.Setup()
		}
		switching = false
	}
}

func mouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	theSketch.HandleMousePosition(xpos, ypos)
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	theSketch.HandleScroll(xoff, xoff)
}

func resizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	theSketch.Update()
	theSketch.Draw()
}

func setup() (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "learnopengl.com in Golang", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	// Initialize Glow - this is the equivalent of glew
	if err := gl.Init(); err != nil {
		return nil, err
	}

	// Resize Callback
	window.SetFramebufferSizeCallback(resizeCallback)

	//Keyboard Callback
	window.SetKeyCallback(keyCallBack)
	window.SetCursorPosCallback(mouseCallback)
	window.SetScrollCallback(scrollCallback)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	glsl := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	fmt.Println("OpenGL version", version, glsl)

	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	return window, nil
}

func initTutorials(w *glfw.Window) []Tutorial {
	// make a slice of pointers to sketch instances
	return []Tutorial{
		Tutorial{
			Name:   "0. Test Cube From github.com/go-gl/examples",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloCube{Window: w},
		},
		Tutorial{
			Name:   "1. Hello Window",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloWindow{Window: w},
		},
		Tutorial{
			Name:   "2. Hello Triangles",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloTriangle{Window: w},
		},
		Tutorial{
			Name:   "2a. Hello Cube",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloSquare{Window: w},
		},
		Tutorial{
			Name:   "3. Shaders",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloShaders{Window: w},
		},
		Tutorial{
			Name:   "4. Textures",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloTextures{Window: w},
		},
		Tutorial{
			Name:   "5. Transformations",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloTransformations{Window: w},
		},
		Tutorial{
			Name:   "6. Coordinate Systems",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloCoordinates{Window: w},
		},
		Tutorial{
			Name:   "7. Camera (use WSDA and mouse)",
			Color:  utils.RandColor(),
			Sketch: &getstarted.HelloCamera{Window: w},
		},
	}
}

func main() {
	// init GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// create window
	window, err := setup()
	if err != nil {
		log.Fatalf("cant create window %v", err)
	}

	tutorials = initTutorials(window)
	tutorial := tutorials[tutorialIndex]
	theSketch = tutorial.Sketch
	if err := theSketch.Setup(); err != nil {
		log.Fatalf("Failed setting up sketch: %v", err)
	}

	//load font (fontfile, font scale, window width, window height
	font, err := utils.LoadFont("sketches/_assets/fonts/huge_agb_v5.ttf", int32(52), WIDTH, HEIGHT)
	if err != nil {
		log.Fatalf("LoadFont: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// loop
	for !window.ShouldClose() {
		t := tutorials[tutorialIndex]

		// Update
		theSketch.Update()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(t.Color.R, t.Color.G, t.Color.B, t.Color.A)

		//Render
		theSketch.Draw()

		font.SetColor(0.0, 0.0, 0.0, 1.0)
		font.Printf(30, 30, 0.5, t.Name)

		window.SwapBuffers()
		// Poll Events
		glfw.PollEvents()
	}
	theSketch.Close()
}
