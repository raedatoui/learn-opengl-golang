package sections

import (
	"errors"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/raedatoui/glutils"
	"strings"
	"github.com/raedatoui/glfont"
)

type TitleSlide struct {
	BaseSlide
	font  *glfont.Font
	lines []string
}

func (s *TitleSlide) Init(a ...interface{}) error {
	f, ok := a[0].(*glfont.Font)
	if ok == false {
		return errors.New("first argument isnt a font")
	}
	s.font = f

	c, ok := a[1].(glutils.Color)
	if ok == false {
		return errors.New("second argument isnt a ColorA")
	}
	s.Color = c
	s.Color32 = c.To32()
	s.ColorHex = glutils.Rgb2Hex(c)

	n, ok := a[2].(string)
	if ok == false {
		return errors.New("third argument isnt a string")
	}
	s.Name = n

	if strings.Contains(s.Name, "\n") {
		s.lines = strings.Split(s.Name, "\n")
	} else {
		s.lines = []string{s.Name}
	}

	return nil
}

func (s *TitleSlide) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(s.Color32.R, s.Color32.G, s.Color32.B, s.Color32.A)

	s.font.SetColor(1.0, 1.0, 1.0, 1.0)
	for i := 0; i < len(s.lines); i++ {
		s.font.Printf(30, 100+60*float32(i), 0.85, s.lines[i])
	}
}

func (b *TitleSlide) DrawText() bool {
	return false
}
