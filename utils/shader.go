package utils

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

func BasicProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func Shader(vertFile, fragFile, geomFile string) (uint32, error) {
	vertexShaderSource, err := readFile(vertFile)
	if err != nil {
		return 0, err
	}

	fragmentShaderSource, err := readFile(fragFile)
	if err != nil {
		return 0, err
	}

	var geometryShaderSource []byte
	if geomFile != "" {
		geometryShaderSource, err = readFile(geomFile)
		if err != nil {
			return 0, err
		}
	}

	vertexShader, err := compileShader(string(vertexShaderSource)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(string(fragmentShaderSource)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	var geometryShader uint32
	if geomFile != "" {
		geometryShader, err = compileShader(string(geometryShaderSource)+"\x00", gl.GEOMETRY_SHADER)
		if err != nil {
			return 0, err
		}
	}

	program, err := createProgram(vertexShader, fragmentShader, geometryShader)
	if err != nil {
		return 0, err
	}

	gl.DetachShader(program, vertexShader)
	gl.DetachShader(program, fragmentShader)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	if geomFile != "" {
		gl.DetachShader(program, geometryShader)
		gl.DeleteShader(geometryShader)
	}

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func createProgram(vertexShader, fragmentShader, geometryShader uint32) (uint32, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	if geometryShader != 0 {
		gl.AttachShader(program, geometryShader)
	}

	gl.LinkProgram(program)
	// check for program linking errors
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return program, nil
}
