### OpenGL Tutorial in Go
* Based on the excellent [tutorial](https://learnopengl.com/) by Joey de Vries
* uses [go-gl](https://github.com/go-gl) for all OpenGL Go bindings

### Installation

This is setup has only been tested on macOS.
I am currently working with the 4.1 core profile on macOS. This is the latest version that [ships](https://support.apple.com/en-us/HT202823) with macOS and I didn't want to delve upgrading that manually.

I figured I would get linux machine at a later if I need to use newer features in the API.
The nice thing about go-gl is that you can install multiple profiles and write different programs targeting different version of OpenGL. 

#### go-gl packages

1- **[gl](https://github.com/go-gl/gl)** - OpenGL core  profile v4.1
`go get -u github.com/go-gl/v4.1-core/gl`

2- [**Glow**](https://github.com/go-gl/glow) - Go binding generator for OpenGL 

```bash
go get github.com/go-gl/glow
cd $GOPATH/src/github.com/go-gl/glow
go build
./glow download
./glow generate -api=gl -version=4.1 -profile=core -remext=GL_ARB_cl_event
go install ./gl-core/3.3/gl
```

3- [**GLFW 3.2**](https://github.com/go-gl/glfw) - Go bindings for GLFW 3
`go get -u github.com/go-gl/glfw/v3.2/glfw`

4- [**MathGL**](https://github.com/go-gl/mathgl) - A pure Go 3D math library
`go get github.com/go-gl/mathgl`

To test that the installation is working, try the examples from go-gl.

`go get github.com/go-gl/examples` 

Run the `gl41core-cube` example by executing `go run cube.go`

#### learnopengl.com tutorial

1- [**glutils**](https://github.com/raedatoui/glutils)

Some of the utllities developed throughout the tutorials like shader compilation and linking, camera, loading textures, loading models from assimp, other redundant GL commands,etc were packaged together. Initially, these lived within the tutorial repo as the `utils` package and we later moved to a dedicated [repo](https://github.com/raedatoui/glutils) in the hope of being useful for other projects.

`go get github.com/raedatoui/glutils` 

I had to fork 2 libraries and update them to get everything working.

2- [**glfont**](https://github.com/raedatoui/glfont) - A modern opengl text rendering library for golang

`go get github.com/raedatoui/glfont`

I made minor changes to this package where I use the shader functions from the `glutils` package and I explicitly set the profile version in the imports to `4.1` intead of `all-core`

Text rendering sucks and is not intended to look good, but good enough and easy to use for the sake of this tutorial.

3- [**assimp**](https://github.com/raedatoui/assimp) - Go wrapper of [Assimp](http://www.assimp.org/)

First, install Assimp on macOS using homebrew `brew install assimp` 

Then install wrapper, `go get github.com/raedatoui/assimp`

I fixed some minor bugs and changed the cgo directives for linking assimp. Intead of using `LDFLAGS` and other windows specific flags, I used the `pkg-config` flag.

### Run

`go run tutorial.go` and you should see this screen

Use the right and left arrow keys to navigate through the tutorials.

Use the num keys to jump between sections.

![Alt text](/screenshot.png?raw=true "Screenshot")


### Notes

When configuring vertex attribute arrays, the stride is calculated using the size of
a float32 type.
* sizeof(GLfloat) is 4 , float32
* size of uint32 - 4



