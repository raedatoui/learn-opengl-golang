package utils

import "C"

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/assimp"
	"strconv"
	"unsafe"
)

type Model struct {
	texturesLoaded  []assimp.Texture
	meshes          []assimp.Mesh
	director        string
	gammaCorrection bool
}

func NewModel(path string, gamma bool) (*Model, error) {
	m := &Model{gammaCorrection: gamma}
	err := m.loadModel(path)
	return m, err
}

type Mesh struct {
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
	dummy := (*Vertex)(nil)
	structSize := C.size_t(unsafe.Sizeof(dummy))

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
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertices)*structSize, gl.Ptr(m.vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.vertices)*GL_FLOAT32_SIZE, gl.Ptr(m.indices), gl.STATIC_DRAW)

	// Set the vertex attribute pointers
	// Vertex Positions
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, structSize, gl.PtrOffset(0))
	// Vertex Normals
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, structSize, unsafe.Pointer(unsafe.Offsetof(dummy.Normal)))
	// Vertex Texture Coords
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, structSize, unsafe.Pointer(unsafe.Offsetof(dummy.TexCoords)))
	// Vertex Tangent
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, structSize, unsafe.Pointer(unsafe.Offsetof(dummy.Tangent)))
	// Vertex Bitangent
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointer(4, 3, gl.FLOAT, false, structSize, unsafe.Pointer(unsafe.Offsetof(dummy.Bitangent)))

	gl.BindVertexArray(0)
}

func (m *Mesh) Draw(program uint32) {
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

// Loads a model with supported ASSIMP extensions from file and stores the resulting meshes in the meshes vector.
func (m *Model) loadModel(path string) error {
	// Read file via ASSIMP
	scene := assimp.ImportFile(path, uint(
		assimp.Process_Triangulate|assimp.Process_FlipUVs|assimp.Process_CalcTangentSpace))

	// Check for errors
	if scene != nil || scene.Flags()&assimp.SceneFlags_Incomplete || scene.RootNode() == nil { // if is Not Zero
		fmt.Println("ERROR::ASSIMP:: %s\n", scene.Flags())
		return error("shit failed")
	}
	// Retrieve the directory path of the filepath
	//m.directory = path.substr(0, path.find_last_of('/'))

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
		m.meshes = append(m.meshes, mesh)
		m.meshes = append(m.meshes, m.processMesh(mesh, s))
		//m.meshes.push_back(m.processMesh(mesh, scene))
	}
	// After we've processed all of the meshes (if any) we then recursively process each of the children nodes
	for j := 0; j < n.NumChildren(); j++ {
		m.processNode(n.Children[j])
	}
}

func (m *Model) processMeshVertices(mesh *assimp.Mesh) []Vertex {
	// Walk through each of the mesh's vertices
	vertices := []Vertex{}
	for i := 0; i < mesh.NumVertices(); i++ {
		// We declare a placeholder vector since assimp uses its own vector class that
		// doesn't directly convert to glm's vec3 class so we transfer the data to this placeholder glm::vec3 first.
		var vertex Vertex
		vector := mgl32.Vec3{}
		var tmp []assimp.Vector3

		// Positions
		tmp = mesh.Vertices()
		vector[0] = tmp[i].X()
		vector[1] = tmp[i].Y()
		vector[2] = tmp[i].Z()
		vertex.Position = vector

		// Normals
		vector = mgl32.Vec3{}
		tmp = mesh.Normals()
		vector[0] = tmp[i].X()
		vector[1] = tmp[i].Y()
		vector[2] = tmp[i].Z()
		vertex.Normal = vector

		// Texture Coordinates
		if mesh.TextureCoords(0) { // Does the mesh contain texture coordinates?
			vec := mgl32.Vec2{}
			// A vertex can contain up to 8 different texture coordinates. We thus make the assumption that we won't
			// use models where a vertex can have multiple texture coordinates so we always take the first set (0).
			tex := mesh.TextureCoords(0)
			vec[0] = tex[i].X()
			vec[1] = tex[i].Y()
			vertex.TexCoords = vec
		} else {
			vertex.TexCoords = mgl32.Vec2{0.0, 0.0}
		}

		// Tangent
		vector = mgl32.Vec3{}
		tmp = mesh.Tangents()
		vector[0] = tmp[i].X()
		vector[1] = tmp[i].Y()
		vector[2] = tmp[i].Z()
		vertex.Tangent = vector

		// Bitangent
		vector = mgl32.Vec3{}
		tmp = mesh.Bitangents()
		vector[0] = tmp[i].X()
		vector[1] = tmp[i].Y()
		vector[2] = tmp[i].Z()
		vertex.Bitangent = vector

		vertices = append(vertices, vertex)
	}
	return vertices
}

func (m *Model) processMeshIndices(mesh *assimp.Mesh) []uint32 {
	indices := []uint32{}

	// Now wak through each of the mesh's faces (a face is a mesh its triangle) and retrieve the corresponding vertex indices.
	for i := 0; i < mesh.NumFaces(); i++ {
		face := mesh.Faces()[i]
		// Retrieve all indices of the face and store them in the indices vector
		indices = append(indices, face.CopyIndices()...)
	}
	return indices
}

func (m *Model) processMesh(mesh *assimp.Mesh, s *assimp.Scene) Mesh {
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
		diffuseMaps := loadMaterialTextures(material, assimp.TextureMapping_Diffuse, "texture_diffuse")
		textures = append(textures, diffuseMaps...)
		// 2. Specular maps
		specularMaps := loadMaterialTextures(material, assimp.TextureMapping_Specular, "texture_specular")
		textures = append(textures, specularMaps...)
		// 3. Normal maps
		normalMaps := loadMaterialTextures(material, assimp.TextureMapping_Height, "texture_normal")
		textures = append(textures, normalMaps...)
		// 4. Height maps
		heightMaps := loadMaterialTextures(material, assimp.TextureMapping_Ambient, "texture_height")
		textures = append(textures, heightMaps...)

	}

	// Return a mesh object created from the extracted mesh data
	return Mesh{
		vertices: m.processMeshVertices(mesh),
		indices:  m.processMeshIndices(mesh),
		textures: textures,
	}
}

func loadMaterialTextures(m *assimp.Material, tm assimp.TextureMapping, tt string) []Texture {
	//mat := cScene.Materials()[cScene.Meshes()[0].MaterialIndex()]
	textureType := assimp.TextureType(tm)
	textureCount := m.GetMaterialTextureCount(textureType)
	result := []Texture{}

	for i := 0; textureCount; i++ {
		file, _, _, _, _, _, _, _ := m.GetMaterialTexture(tt, 0)
		texId := textureFromFile(file)
		texture := Texture{Id: texId, TextureType: tt, Path: file}
		result = append(result, texture)
	}
	return result
}

func textureFromFile(f string) uint32 {
	//Generate texture ID and load texture data
	var textureID uint32

	filename := "_assets/objects/cyborg/" + f

	gl.GenTextures(1, &textureID)
	var width, height int32
	rgba, err := ImageToPixelData(filename)
	if err != nil {
		panic(err)
	}
	// Assign texture to ID
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0, gl.RGB, gl.UNSIGNED_BYTE, rgba)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// Parameters

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return textureID
}
