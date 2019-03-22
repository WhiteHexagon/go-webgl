package main

import (
	"fmt"

	"github.com/whitehexagon/go-webgl"
)

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Go/WASM main()")

	canvas := webgl.CreateCanvasAt("app", "canvas42")
	gl, err := webgl.GetContext(canvas)
	if err != nil {
		fmt.Println("failed:", err)
		return
	}

	vBuffer, iBuffer, icount := createBuffers(gl)
	fmt.Println("createBuffers", vBuffer, iBuffer, icount)

	prog := setupShaders(gl)
	fmt.Println("prog", prog)

	//// Associating shaders to buffer objects ////

	// Bind vertex buffer object
	gl.BindBuffer(webgl.ARRAY_BUFFER, vBuffer)

	// Bind index buffer object
	gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, iBuffer)

	// Get the attribute location
	coord := gl.GetAttribLocation(prog, "coordinates")

	// Point an attribute to the currently bound VBO
	gl.VertexAttribPointer(coord, 3, webgl.FLOAT, false, 0, 0)

	// Enable the attribute
	gl.EnableVertexAttribArray(coord)

	//// Drawing the triangle ////

	// Clear the canvas
	gl.ClearColor(0.5, 0.5, 0.5, 0.9)
	gl.Clear(webgl.COLOR_BUFFER_BIT)

	// Enable the depth test
	gl.Enable(webgl.DEPTH_TEST)

	// Set the view port
	gl.Viewport(0, 0, canvas.Width, canvas.Height)

	// Draw the triangle
	gl.DrawElements(webgl.TRIANGLES, uint(icount), webgl.UNSIGNED_SHORT, uint(0))

	fmt.Println("done")
	<-c
}

func createBuffers(gl *webgl.Context) (*webgl.Buffer, *webgl.Buffer, int) {
	//// VERTEX BUFFER ////
	var vertices = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	// Create buffer
	vBuffer := gl.CreateBuffer(webgl.ARRAY_BUFFER)
	// Bind to buffer
	gl.BindBuffer(webgl.ARRAY_BUFFER, vBuffer)
	// Pass data to buffer
	gl.BufferData(webgl.ARRAY_BUFFER, vertices, webgl.STATIC_DRAW)
	// Unbind buffer
	gl.BindBuffer(webgl.ARRAY_BUFFER, nil)

	// INDEX BUFFER ////
	var indices = []uint32{
		2, 1, 0,
	}
	// Create buffer
	iBuffer := gl.CreateBuffer(webgl.ELEMENT_ARRAY_BUFFER)
	// Bind to buffer
	gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, iBuffer)
	// Pass data to buffer
	gl.BufferData(webgl.ELEMENT_ARRAY_BUFFER, indices, webgl.STATIC_DRAW)
	// Unbind buffer
	gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, nil)
	return vBuffer, iBuffer, len(indices)
}

func setupShaders(gl *webgl.Context) *webgl.Program {
	// Vertex shader source code
	vertCode := `
	attribute vec3 coordinates;
	void main(void) {
		gl_Position = vec4(coordinates, 1.0);
	}`

	// Create a vertex shader object
	vShader := gl.CreateShader(webgl.VERTEX_SHADER)

	// Attach vertex shader source code
	gl.ShaderSource(vShader, vertCode)

	// Compile the vertex shader
	gl.CompileShader(vShader)

	//fragment shader source code
	fragCode := `
	void main(void) {
		gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
	}`

	// Create fragment shader object
	fShader := gl.CreateShader(webgl.FRAGMENT_SHADER)

	// Attach fragment shader source code
	gl.ShaderSource(fShader, fragCode)

	// Compile the fragmentt shader
	gl.CompileShader(fShader)

	// Create a shader program object to store
	// the combined shader program
	prog := gl.CreateProgram()

	// Attach a vertex shader
	gl.AttachShader(prog, vShader)

	// Attach a fragment shader
	gl.AttachShader(prog, fShader)

	// Link both the programs
	gl.LinkProgram(prog)

	// Use the combined shader program object
	gl.UseProgram(prog)

	return prog
}
