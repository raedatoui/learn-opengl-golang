package sections

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"errors"
)

// Slide is the most basic slide. it has to setup, update, draw and cloae
type Slide interface {
	Init(a ...interface{}) error
	InitGL() error
	Update()
	Draw()
	Close()
	GetName() string
}

// Sketch is an interactive Slide and process user interactions
type Sketch interface {
	Slide
	HandleKeyboard(k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey)
	HandleMousePosition(xpos, ypos float64)
	HandleScroll(xoff, yoff float64)
}

// BaseSlide is the base implementation of Slide with the min required fields
type BaseSlide struct {
	Slide
	Name   string
	Color  utils.ColorA
}

func (s *BaseSlide) GetName() string {
	return s.Name
}


type BaseSketch struct {
	Sketch
	BaseSlide
}

func (b *BaseSketch) Init(a ...interface{}) error {
	c, ok := a[0].(utils.ColorA)
	if  ok == false {
		return errors.New("first argument isnt a color")
	}
	b.Color = c
	return nil
}
func (s *BaseSketch) GetName() string {
	return s.Name
}