package sections

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"strings"
)

type TitleSlide struct {
	BaseSlide
	Name  string
	lines []string
}

func (s *TitleSlide) Setup(w *glfw.Window, f *utils.Font) error {
	s.Window = w
	s.Font = f
	s.Color = utils.ColorA{R: 236.0 / 255.0, G: 0, B: 140.0 / 255.0, A: 1.0}

	if strings.Contains(s.Name, "\n") {
		s.lines = strings.Split(s.Name, "\n")
	} else {
		s.lines = []string{s.Name}
	}

	return nil
}

func (s *TitleSlide) Update() {

}

func (s *TitleSlide) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(s.Color.R, s.Color.G, s.Color.B, s.Color.A)

	s.Font.SetColor(1.0, 1.0, 1.0, 1.0)
	for i := 0; i < len(s.lines); i++ {
		s.Font.Printf(30, 200+60*float32(i), 1, s.lines[i])

	}
}

func (s *TitleSlide) Close() {

}
