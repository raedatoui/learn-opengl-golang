package utils

import "C"

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/raedatoui/assimp"
	"strconv"
	"os"
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

type Mesh struct {
	id       int
	vertices []Vertex
	indices  []uint32
	textures []Texture
	Vao      uint32
	vbo, ebo uint32
}

func NewMesh(v []Vertex, i []uint32, t []Texture) Mesh {
	m := Mesh{
		vertices: v,
		indices:  i,
		textures: t,
	}
	m.setup()
	return m
}

func (m *Mesh) setup() {
	// size of the Vertex struct
	dummy := m.vertices[0]
	structSize := int(unsafe.Sizeof(dummy))
	structSize32 := int32(structSize)

	// Create buffers/arrays
	gl.GenVertexArrays(1, &m.Vao)
	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)

	gl.BindVertexArray(m.Vao)
	// Load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	// A great thing about structs is that their memory layout is sequential for all its items.
	// The effect is that we can simply pass a pointer to the struct and it translates perfectly to a gl.m::vec3/2 array which
	// again translates to 3/2 floats which translates to a byte array.
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertices)* structSize, gl.Ptr(m.vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.indices)*GL_FLOAT32_SIZE, gl.Ptr(m.indices), gl.STATIC_DRAW)

	// Set the vertex attribute pointers
	// Vertex Positions
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, structSize32, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Vertex Normals
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, structSize32, unsafe.Pointer((unsafe.Offsetof(dummy.Normal))))
	gl.EnableVertexAttribArray(1)

	// Vertex Texture Coords
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, structSize32, unsafe.Pointer((unsafe.Offsetof(dummy.TexCoords))))
	gl.EnableVertexAttribArray(2)

	// Vertex Tangent
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, structSize32, unsafe.Pointer(unsafe.Offsetof(dummy.Tangent)))
	// Vertex Bitangent
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointer(4, 3, gl.FLOAT, false, structSize32, unsafe.Pointer(unsafe.Offsetof(dummy.Bitangent)))

	gl.BindVertexArray(0)
}

func (m *Mesh) draw(program uint32) {
	// Bind appropriate textures
	var (
		diffuseNr  uint64
		specularNr uint64
		normalNr   uint64
		heightNr   uint64
		i          uint32
	)
	diffuseNr = 1
	specularNr = 1
	normalNr = 1
	heightNr = 1
	i = 0
	for i = 0; i < uint32(len(m.textures)); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + i) // Active proper texture unit before binding

		// Retrieve texture number (the N in diffuse_textureN)
		ss := ""
		switch m.textures[i].TextureType {
		case "texture_diffuse":
			ss = ss + strconv.FormatUint(diffuseNr, 10) // Transfer GLuint to stream
			diffuseNr++
		case "texture_specular":
			ss = ss + strconv.FormatUint(specularNr, 10) // Transfer GLuint to stream
			specularNr++
		case "texture_normal":
			ss = ss + strconv.FormatUint(normalNr, 10) // Transfer GLuint to stream
			normalNr++
		case "texture_height":
			ss = ss + strconv.FormatUint(heightNr, 10) // Transfer GLuint to stream
			heightNr++
		}

		// Now set the sampler to the correct texture unit
		tu := m.textures[i].TextureType + ss + "\x00"

		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str(tu)), int32(i))
		// And finally bind the texture
		gl.BindTexture(gl.TEXTURE_2D, m.textures[i].Id)
	}

	// Draw mesh
	gl.BindVertexArray(m.Vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)

	// Always good practice to set everything back to defaults once configured.
	for i = 0; i < uint32(len(m.textures)); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + i)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
}

type Vertex struct {
	Position  mgl32.Vec3
	Normal    mgl32.Vec3
	TexCoords mgl32.Vec2
	Tangent   mgl32.Vec3
	Bitangent mgl32.Vec3
}

type Texture struct {
	Id          uint32
	TextureType string
	Path        string
}

type Model struct {
	texturesLoaded  map[string]Texture
	meshes          []Mesh
	director        string
	gammaCorrection bool
	basePath        string
	fileName        string
}

func NewModel(b, f string, g bool) (Model, error) {

	m := Model{
		basePath:        b,
		fileName:        f,
		gammaCorrection: g,
	}
	m.texturesLoaded = make(map[string]Texture)
	err := m.loadModel()
	return m, err
}

func (m *Model) Draw(shader uint32) {
	for i := 0; i < len(m.meshes); i++ {
		m.meshes[i].draw(shader)
	}
}

// Loads a model with supported ASSIMP extensions from file and stores the resulting meshes in the meshes vector.
func (m *Model) loadModel() error {
	// Read file via ASSIMP
	path := m.basePath + m.fileName
	scene := assimp.ImportFile(path, uint(
		assimp.Process_Triangulate|assimp.Process_FlipUVs))

	// Check for errors
	if scene.Flags() & assimp.SceneFlags_Incomplete != 0 { // if is Not Zero
		fmt.Println("ERROR::ASSIMP:: %s\n", scene.Flags())
		return errors.New("shit failed")
	}

	// Process ASSIMP's root node recursively
	m.processNode(scene.RootNode(), scene)

	return nil
}

func (m *Model) processNode(n *assimp.Node, s *assimp.Scene) {
	// Process each mesh located at the current node
	for i := 0; i < n.NumMeshes(); i++ {
		// The node object only contains indices to index the actual objects in the scene.
		// The scene contains all the data, node is just to keep stuff organized (like relations between nodes).
		mesh := s.Meshes()[n.Meshes()[i]]
		ms := m.processMesh(mesh, s)
		ms.id = i
		m.meshes = append(m.meshes, ms)
	}

	// After we've processed all of the meshes (if any) we then recursively process each of the children nodes
	c := n.Children()
	for j := 0; j < len(c); j++ {
		m.processNode(c[j], s)
	}
}

func (m *Model) processMeshVertices(mesh *assimp.Mesh) []Vertex {
	// Walk through each of the mesh's vertices
	vertices := []Vertex{}
	//p, _ := os.Create("data/positions"+ mesh.Name()+ ".txt")
	//n, _ := os.Create("data/normals"+ mesh.Name()+ ".txt")
	//t, _ := os.Create("data/texcoords"+ mesh.Name()+ ".txt")
	//
	//defer p.Close()
	//defer n.Close()
	//defer t.Close()
	//p.WriteString(mesh.Name() + "\n")
	//n.WriteString(mesh.Name() + "\n")
	//t.WriteString(mesh.Name() + "\n")

	positions := mesh.Vertices()

	normals := mesh.Normals()
	useNormals := len(normals) > 0

	tex :=  mesh.TextureCoords(0)
	useTex := true
	if tex == nil {
		useTex = false
	}

	tangents := mesh.Tangents()
	useTangents := len(tangents) > 0

	bitangents := mesh.Bitangents()
	useBitTangents := len(bitangents) > 0

	for i := 0; i < mesh.NumVertices(); i++ {
		// We declare a placeholder vector since assimp uses its own vector class that
		// doesn't directly convert to glm's vec3 class so we transfer the data to this placeholder glm::vec3 first.
		vertex := Vertex{}

		// Positions
		vertex.Position = mgl32.Vec3{positions[i].X(), positions[i].Y(), positions[i].Z()}

		// Normals
		if useNormals {
			vertex.Normal = mgl32.Vec3{normals[i].X(), normals[i].Y(), normals[i].Z()}
			//n.WriteString(fmt.Sprintf("[%f, %f, %f]\n", tmp[i].X(), tmp[i].Y(), tmp[i].Z()))
		}

		// Texture Coordinates
		if useTex {
			// Does the mesh contain texture coordinates?
			// A vertex can contain up to 8 different texture coordinates. We thus make the assumption that we won't
			// use models where a vertex can have multiple texture coordinates so we always take the first set (0).
			vertex.TexCoords = mgl32.Vec2{tex[i].X(),  tex[i].Y()}
		} else {
			vertex.TexCoords = mgl32.Vec2{0.0, 0.0}
		}

		// Tangent
		if useTangents {
			vertex.Tangent = mgl32.Vec3{tangents[i].X(), tangents[i].Y(), tangents[i].Z()}
		}

		// Bitangent
		if useBitTangents {
			vertex.Bitangent = mgl32.Vec3{bitangents[i].X(), bitangents[i].Y(), bitangents[i].Z()}
		}

		vertices = append(vertices, vertex)
	}

	return vertices
}

func (m *Model) processMeshIndices(mesh *assimp.Mesh) []uint32 {
	indices := []uint32{}
	ind, _ := os.Create("data/indices"+ mesh.Name()+ ".txt")
	// Now wak through each of the mesh's faces (a face is a mesh its triangle) and retrieve the corresponding vertex indices.
	for i := 0; i < mesh.NumFaces(); i++ {
		face := mesh.Faces()[i]
		t := face.CopyIndices()
		ind.WriteString(fmt.Sprintf("[%d, %d, %d]\n", t[0], t[1], t[2]))
		// Retrieve all indices of the face and store them in the indices vector
		indices = append(indices,t ...)
	}
	ind.WriteString("\n--------")
	return indices
}

func (m *Model) processMeshTextures(mesh *assimp.Mesh, s *assimp.Scene) []Texture {
	textures := []Texture{}
	// Process materials
	if mesh.MaterialIndex() >= 0 {
		material := s.Materials()[mesh.MaterialIndex()]

		// We assume a convention for sampler names in the shaders. Each diffuse texture should be named
		// as 'texture_diffuseN' where N is a sequential number ranging from 1 to MAX_SAMPLER_NUMBER.
		// Same applies to other texture as the following list summarizes:
		// Diffuse: texture_diffuseN
		// Specular: texture_specularN
		// Normal: texture_normalN

		// 1. Diffuse maps
		diffuseMaps := m.loadMaterialTextures(material, assimp.TextureMapping_Diffuse, "texture_diffuse")
		textures = append(textures, diffuseMaps...)
		// 2. Specular maps
		specularMaps := m.loadMaterialTextures(material, assimp.TextureMapping_Specular, "texture_specular")
		textures = append(textures, specularMaps...)
		// 3. Normal maps
		normalMaps := m.loadMaterialTextures(material, assimp.TextureMapping_Height, "texture_normal")
		textures = append(textures, normalMaps...)
		// 4. Height maps
		heightMaps := m.loadMaterialTextures(material, assimp.TextureMapping_Ambient, "texture_height")
		textures = append(textures, heightMaps...)
	}
	return textures
}

func (ml *Model) processMesh(m *assimp.Mesh, s *assimp.Scene) Mesh {
	// Return a mesh object created from the extracted mesh data
	return NewMesh(
		ml.processMeshVertices(m),
		ml.processMeshIndices(m),
		ml.processMeshTextures(m, s))
}

func (m *Model) loadMaterialTextures(ms *assimp.Material, tm assimp.TextureMapping, tt string) []Texture {
	textureType := assimp.TextureType(tm)
	textureCount := ms.GetMaterialTextureCount(textureType)
	result := []Texture{}

	for i := 0; i < textureCount; i++ {
		file, _, _, _, _, _, _, _ := ms.GetMaterialTexture(textureType, 0)
		filename := m.basePath + file
		if val, ok := m.texturesLoaded[filename]; ok {
			result = append(result, val)
		} else {
			texId := m.textureFromFile(filename)
			texture := Texture{Id: texId, TextureType: tt, Path: file}
			result = append(result, texture)
			m.texturesLoaded[filename] = texture
		}
	}
	return result
}

func (ml *Model) textureFromFile(f string) uint32 {
	//Generate texture ID and load texture data
	var textureID uint32

	gl.GenTextures(1, &textureID)
	rgba, err := ImageToPixelData(f)
	if err != nil {
		panic(err)
	}
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	// Assign texture to ID
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// Parameters

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return textureID
}
