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
	"github.com/raedatoui/learn-opengl-golang/sketches/getstarted"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"reflect"

)

// WIDTH is the width of the window
const WIDTH = 800

// HEIGHT is the height of the window
const HEIGHT = 600

var (
	currentSlide  sketches.Slide
	currentSketch sketches.Sketch // polymorphic
	slides        []sketches.Slide
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
	if currentSketch != nil {
		sketches.HandleKeyboardSketch(currentSketch, k, s, a, mk)
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
			currentSketch.Close()
			currentSlide = slides[newIndex]
			currentSketch = getSketch(currentSlide)
			sketches.SetupSlide(currentSlide, window, font)
		}
		switching = false
	}
}

func mouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	sketches.HandleMousePositionSketch(currentSketch, xpos, ypos)
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	sketches.HandleScrollSketch(currentSketch, xoff, yoff)
}

func resizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	sketches.UpdateSlide(currentSlide)
	sketches.DrawSlide(currentSlide)
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

func setupSlides() []sketches.Slide {
	// make a slice of pointers to sketch instances
	return []sketches.Slide{
		&getstarted.HelloCube{},
		&getstarted.HelloWindow{},
		&getstarted.HelloTriangle{},
		&getstarted.HelloSquare{},
		&getstarted.HelloShaders{},
		//},
		//Tutorial{
		//	Name:   "4. Textures",
		//	Color:  utils.RandColor(),
		//	Sketch: &getstarted.HelloTextures{},
		//},
		//Tutorial{
		//	Name:   "5. Transformations",
		//	Color:  utils.RandColor(),
		//	Sketch: &getstarted.HelloTransformations{},
		//},
		//Tutorial{
		//	Name:   "6. Coordinate Systems",
		//	Color:  utils.RandColor(),
		//	Sketch: &getstarted.HelloCoordinates{},
		//},
		//Tutorial{
		//	Name:   "7. Camera (use WSDA and mouse)",
		//	Color:  utils.RandColor(),
		//	Sketch: &getstarted.HelloCamera{},
		//},
	}
}

func getSketch(o interface{}) sketches.Sketch {
	i := reflect.ValueOf(o).Type()
	s := reflect.TypeOf((*sketches.Sketch)(nil)).Elem()
	if i.Implements(s) {
		s, ok := currentSlide.(sketches.Sketch)
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
	f, err := utils.LoadFont("sketches/_assets/fonts/huge_agb_v5.ttf", int32(52), WIDTH, HEIGHT)
	if err != nil {
		log.Fatalf("LoadFont: %v", err)
	}
	font  = f

	slides = setupSlides()
	currentSlide = slides[0]
	currentSketch = getSketch(currentSlide)

	if err := sketches.SetupSlide(currentSlide, window, font); err != nil {
		log.Fatalf("Failed setting up sketch: %v", err)
	}

	// loop
	for !window.ShouldClose() {

		// Update
		sketches.UpdateSlide(currentSlide)

		//Render
		sketches.DrawSlide(currentSlide)

		window.SwapBuffers()
		// Poll Events
		glfw.PollEvents()
	}
	sketches.CloseSlide(currentSlide)
}
