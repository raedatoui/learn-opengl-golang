package sections

import (
	"errors"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/utils"
)

// Slide is the most basic slide. it has to setup, update, draw and close
type Slide interface {
	Init(a ...interface{}) error
	InitGL() error
	Update()
	Draw()
	Close()
	GetName() string
	GetColor() utils.Color
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
	Name  string
	Color utils.Color
	Color32 utils.Color32
}

func (s *BaseSlide) GetName() string {
	return s.Name
}

func (s *BaseSlide) GeColor() utils.Color {
	return s.Color
}

type BaseSketch struct {
	Sketch
	BaseSlide
}

func (b *BaseSketch) Init(a ...interface{}) error {
	c, ok := a[0].(utils.Color)
	if ok == false {
		return errors.New("first argument isnt a color")
	}
	b.Color = c
	b.Color32 = c.To32()
	return nil
}

func (b *BaseSketch) GetName() string {
	return b.Name
}

func (b *BaseSketch) GeColor() utils.Color {
	return b.Color
}
