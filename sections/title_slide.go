package sections

import (
	"errors"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"strings"
)

type TitleSlide struct {
	BaseSlide
	font  *utils.Font
	Name  string
	lines []string
}

func (s *TitleSlide) Init(a ...interface{}) error {
	f, ok := a[0].(*utils.Font)
	if ok == false {
		return errors.New("first argument isnt a font")
	}
	s.font = f

	c, ok := a[1].(utils.Color)
	if ok == false {
		return errors.New("second argument isnt a ColorA")
	}
	s.Color = c
	s.Color32 = c.To32()


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
		s.font.Printf(30, 200+60*float32(i), 1, s.lines[i])

	}
}
