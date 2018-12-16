### OpenGL Tutorial in Go
* Based on the excellent [tutorial](https://learnopengl.com/) by Joey de Vries
* uses [go-gl](https://github.com/go-gl) for all OpenGL Go bindings

### Installation

This is setup has only been tested on macOS.
I am currently working with the 4.1 core profile on macOS, which is the latest version that [ships](https://support.apple.com/en-us/HT202823) with it. I didn't want to delve into upgrading that manually.

I figured I would get linux machine at a later time if I needed to use newer features in the API.
The nice thing about go-gl is that you can install multiple profiles and write different programs targeting different version of OpenGL. 

#### go-gl packages

1- **[gl](https://github.com/go-gl/gl)** - OpenGL core  profile v4.1
`go get -u github.com/go-gl/gl/v4.1-core/gl`

2- [**Glow**](https://github.com/go-gl/glow) - Go binding generator for OpenGL 

```shell
go get github.com/go-gl/glow
cd $GOPATH/src/github.com/go-gl/glow
go build
./glow download
./glow generate -api=gl -version=4.1 -profile=core -remext=GL_ARB_cl_event
# the profile is now installed in a gl directory
go install ./gl
```

3- [**GLFW 3.2**](https://github.com/go-gl/glfw) - Go bindings for GLFW 3
`go get -u github.com/go-gl/glfw/v3.2/glfw`

4- [**MathGL**](https://github.com/go-gl/mathgl) - A pure Go 3D math library
`go get github.com/go-gl/mathgl/...`

This package is the equivalent of the GLM library and probably has all the functionality but after some differences.
I didnt dive too deep into it, but I am getting different matrices when running the same sample in C++ with glm vs Go with mgl32.


To test that the installation is working, try the examples from go-gl.

```shell
go get github.com/go-gl/example
cd $GOPATH/src/github.com/go-gl/example
go run gl41core-cube/cube.go
```

#### learnopengl.com tutorial

1- [**assimp**](https://github.com/raedatoui/assimp) - Go wrapper of [Assimp](http://www.assimp.org/)

First, install Assimp on macOS using homebrew `brew install assimp` 

Then install wrapper, `go get github.com/raedatoui/assimp`

2- [**glutils**](https://github.com/raedatoui/glutils)

Some of the utllities developed throughout the tutorials like shader compilation and linking, camera, loading textures, loading models from assimp, other redundant GL commands,etc were packaged together. Initially, these lived within the tutorial repo as the `utils` package and we later moved to a dedicated [repo](https://github.com/raedatoui/glutils) in the hope of being useful for other projects.

`go get github.com/raedatoui/glutils` 

3- [**glfont**](https://github.com/raedatoui/glfont) - A modern opengl text rendering library for golang

`go get github.com/raedatoui/glfont`

I made minor changes to this package where I use the shader functions from the `glutils` package and I explicitly set the profile version in the imports to `4.1` intead of `all-core`

Text rendering sucks and is not intended to look good, but good enough and easy to use for the sake of this tutorial.

### Run

```shell
go get github.com/raedatoui/learn-opengl-golang
cd $GOPATH/src/github.com/raedatoui/learn-opengl-golang
go run tutorial.go
```
and you should see this screen

Use the right and left arrow keys to navigate through the tutorials.

Use the num keys to jump between sections.

![Alt text](/screenshot.png?raw=true "Screenshot")


### Notes

When configuring vertex attribute arrays, the stride is calculated using the size of
a float32 type.
* sizeof(GLfloat) is 4 , float32
* size of uint32 - 4



