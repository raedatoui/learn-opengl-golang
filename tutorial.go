package main

import (
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"fmt"

	"reflect"

	"github.com/raedatoui/learn-opengl-golang/sections"
	"github.com/raedatoui/learn-opengl-golang/sections/getstarted"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

// WIDTH is the width of the window
const WIDTH = 800

// HEIGHT is the height of the window
const HEIGHT = 600

var (
	currentSlide  sections.Slide
	currentSketch sections.Sketch // polymorphic
	slides        []sections.Slide
	slideIndex    = 0
	switching     bool
	window        *glfw.Window
	font          *utils.Font
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

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if k == glfw.KeyEscape && a== glfw.Press {
		window.SetShouldClose(true)
	}

	if currentSketch != nil {
		currentSketch.HandleKeyboard(k, s, a, mk)
	}

	if !switching && a == glfw.Press {
		switching = true
		newIndex := slideIndex
		if a == glfw.Press && s == 124 {
			newIndex = slideIndex + 1
			if newIndex > len(slides)-1 {
				newIndex = len(slides) - 1
			}
		}
		if a == glfw.Press && s == 123 {
			newIndex = slideIndex - 1
			if newIndex < 0 {
				newIndex = 0
			}
		}
		if a == glfw.Press && newIndex != slideIndex {
			slideIndex = newIndex
			currentSlide.Close()
			currentSlide = slides[newIndex]
			currentSketch = getSketch(currentSlide)
			currentSlide.Setup(window, font)
		}
		switching = false
	}
}

func mouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	if currentSketch != nil {
		currentSketch.HandleMousePosition(xpos, ypos)
	}
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	if currentSketch != nil {
		currentSketch.HandleScroll(xoff, yoff)
	}
}

func resizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	currentSlide.Update()
	currentSlide.Draw()
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

func setupSlides() []sections.Slide {
	// make a slice of pointers to sketch instances
	return []sections.Slide{
		&sections.TitleSlide{Name: "Test Installation of gl,\nglwf, glow"},
		new(getstarted.HelloCube),
		&sections.TitleSlide{Name: "Section 1: Getting Started"},
		new(getstarted.HelloWindow),
		new(getstarted.HelloTriangle),
		new(getstarted.HelloSquare),
		new(getstarted.HelloShaders),
		new(getstarted.HelloTextures),
		new(getstarted.HelloTransformations),
		new(getstarted.HelloCoordinates),
		new(getstarted.HelloCamera),
		&sections.TitleSlide{Name: "Section 2: Lighting"},
	}
}

func getSketch(o interface{}) sections.Sketch {
	i := reflect.ValueOf(o).Type()
	s := reflect.TypeOf((*sections.Sketch)(nil)).Elem()
	if i.Implements(s) {
		s, ok := currentSlide.(sections.Sketch)
		if !ok {
			log.Fatalf("cant convert Slide to Sketch. Bravo!")
		}
		return s
	}
	return nil
}

func main() {
	// init GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// create window
	w, err := setup()
	if err != nil {
		log.Fatalf("cant create window %v", err)
	}
	window = w

	//load font (fontfile, font scale, window width, window height
	f, err := utils.LoadFont("_assets/fonts/huge_agb_v5.ttf", int32(52), WIDTH, HEIGHT)
	if err != nil {
		log.Fatalf("LoadFont: %v", err)
	}
	font = f

	slides = setupSlides()
	currentSlide = slides[0]
	currentSketch = getSketch(currentSlide)

	if err := currentSlide.Setup(window, font); err != nil {
		log.Fatalf("Failed setting up sketch: %v", err)
	}

	// loop
	for !window.ShouldClose() {

		// Update
		currentSlide.Update()

		//Render
		currentSlide.Draw()

		window.SwapBuffers()
		// Poll Events
		glfw.PollEvents()
	}
	currentSlide.Close()
}
