package webgl

import (
	"errors"
	"syscall/js"
)

type Buffer struct {
	js.Value
}

type Shader struct {
	js.Value
}

type Program struct {
	js.Value
}

type Context struct {
	js.Value
}

type Canvas struct {
	canvas js.Value
	id     string
	Width  uint
	Height uint
}

func CreateCanvasAt(divID string, withID string) *Canvas {
	doc := js.Global().Get("document")
	app := doc.Call("getElementById", divID)
	cw := app.Get("clientWidth")
	ch := app.Get("clientHeight")

	canvas := doc.Call("createElement", "canvas")
	canvas.Set("id", withID)
	canvas.Set("width", cw)
	canvas.Set("height", ch)
	app.Call("appendChild", canvas)
	return &Canvas{canvas, divID, uint(cw.Int()), uint(ch.Int())}
}

func GetContext(canvas *Canvas) (*Context, error) {
	context := canvas.canvas.Call("getContext", "webgl")
	if context == js.Undefined() {
		return nil, errors.New("WebGL not supported?")
	}
	return &Context{context}, nil
}

func (gl *Context) ClearColor(r float64, g float64, b float64, a float64) {
	var args [4]js.Value
	args[0] = js.ValueOf(r)
	args[1] = js.ValueOf(g)
	args[2] = js.ValueOf(a)
	args[3] = js.ValueOf(r)
	gl.Call("clearColor", r, g, b, a)
}

type ClearBufferMask uint

const (
	DEPTH_BUFFER_BIT   ClearBufferMask = 0x00000100
	STENCIL_BUFFER_BIT ClearBufferMask = 0x00000400
	COLOR_BUFFER_BIT   ClearBufferMask = 0x00004000
)

func (gl *Context) Clear(mask ClearBufferMask) {
	gl.Call("clear", uint(mask))
}

type BufferTarget uint

const (
	ARRAY_BUFFER                 BufferTarget = 0x8892
	ELEMENT_ARRAY_BUFFER         BufferTarget = 0x8893
	ARRAY_BUFFER_BINDING         BufferTarget = 0x8894
	ELEMENT_ARRAY_BUFFER_BINDING BufferTarget = 0x8895
)

func (gl *Context) CreateBuffer(target BufferTarget) *Buffer {
	b := gl.Call("createBuffer", uint(target))
	return &Buffer{b}
}

type Usage uint

const (
	STREAM_DRAW  Usage = 0x88E0
	STATIC_DRAW  Usage = 0x88E4
	DYNAMIC_DRAW Usage = 0x88E8
)

func (gl *Context) BufferData(target BufferTarget, data interface{}, usage Usage) {
	gl.Call("bufferData", uint(target), js.TypedArrayOf(data), uint(usage)) //leaks
}

func (gl *Context) BindBuffer(target BufferTarget, buffer *Buffer) {
	if buffer == nil {
		gl.Call("bindBuffer", uint(target), nil)
	} else {
		gl.Call("bindBuffer", uint(target), buffer)
	}
}

type Shaders uint

const (
	FRAGMENT_SHADER                  Shaders = 0x8B30
	VERTEX_SHADER                    Shaders = 0x8B31
	MAX_VERTEX_ATTRIBS               Shaders = 0x8869
	MAX_VERTEX_UNIFORM_VECTORS       Shaders = 0x8DFB
	MAX_VARYING_VECTORS              Shaders = 0x8DFC
	MAX_COMBINED_TEXTURE_IMAGE_UNITS Shaders = 0x8B4D
	MAX_VERTEX_TEXTURE_IMAGE_UNITS   Shaders = 0x8B4C
	MAX_TEXTURE_IMAGE_UNITS          Shaders = 0x8872
	MAX_FRAGMENT_UNIFORM_VECTORS     Shaders = 0x8DFD
	SHADER_TYPE                      Shaders = 0x8B4F
	DELETE_STATUS                    Shaders = 0x8B80
	LINK_STATUS                      Shaders = 0x8B82
	VALIDATE_STATUS                  Shaders = 0x8B83
	ATTACHED_SHADERS                 Shaders = 0x8B85
	ACTIVE_UNIFORMS                  Shaders = 0x8B86
	ACTIVE_ATTRIBUTES                Shaders = 0x8B89
	SHADING_LANGUAGE_VERSION         Shaders = 0x8B8C
	CURRENT_PROGRAM                  Shaders = 0x8B8D
)

func (gl *Context) CreateShader(sType Shaders) *Shader {
	s := gl.Call("createShader", uint(sType))
	return &Shader{s}
}

func (gl *Context) ShaderSource(shader *Shader, source string) {
	gl.Call("shaderSource", shader, source)
}

func (gl *Context) CompileShader(shader *Shader) {
	gl.Call("compileShader", shader)
}

func (gl *Context) CreateProgram() *Program {
	p := gl.Call("createProgram")
	return &Program{p}
}

func (gl *Context) AttachShader(prog *Program, shader *Shader) {
	gl.Call("attachShader", prog, shader)
}

func (gl *Context) LinkProgram(prog *Program) {
	gl.Call("linkProgram", prog)
}

func (gl *Context) UseProgram(prog *Program) {
	gl.Call("useProgram", prog)
}

type Location struct {
	js.Value
}

func (gl *Context) GetAttribLocation(prog *Program, attr string) *Location {
	loc := gl.Call("getAttribLocation", prog, attr)
	return &Location{loc}
}

type DataType uint

const (
	BYTE           DataType = 0x1400
	UNSIGNED_BYTE  DataType = 0x1401
	SHORT          DataType = 0x1402
	UNSIGNED_SHORT DataType = 0x1403
	INT            DataType = 0x1404
	UNSIGNED_INT   DataType = 0x1405
	FLOAT          DataType = 0x1406
)

func (gl *Context) VertexAttribPointer(loc *Location, size uint, dtype DataType, normalized bool, stride uint, offset uint) {
	gl.Call("vertexAttribPointer", loc, size, uint(dtype), normalized, stride, offset)
}

func (gl *Context) EnableVertexAttribArray(loc *Location) {
	gl.Call("enableVertexAttribArray", loc)
}

type EnableCap uint

/* EnableCap */
/* TEXTURE_2D */
const (
	CULL_FACE                EnableCap = 0x0B44
	BLEND                    EnableCap = 0x0BE2
	DITHER                   EnableCap = 0x0BD0
	STENCIL_TEST             EnableCap = 0x0B90
	DEPTH_TEST               EnableCap = 0x0B71
	SCISSOR_TEST             EnableCap = 0x0C11
	POLYGON_OFFSET_FILL      EnableCap = 0x8037
	SAMPLE_ALPHA_TO_COVERAGE EnableCap = 0x809E
	SAMPLE_COVERAGE          EnableCap = 0x80A0
)

func (gl *Context) Enable(cap EnableCap) {
	gl.Call("enable", uint(cap))
}

func (gl *Context) Viewport(x uint, y uint, width uint, height uint) {
	gl.Call("viewport", x, y, width, height)
}

type BeginMode uint

const (
	POINTS         BeginMode = 0x0000
	LINES          BeginMode = 0x0001
	LINE_LOOP      BeginMode = 0x0002
	LINE_STRIP     BeginMode = 0x0003
	TRIANGLES      BeginMode = 0x0004
	TRIANGLE_STRIP BeginMode = 0x0005
	TRIANGLE_FAN   BeginMode = 0x0006
)

func (gl *Context) DrawElements(mode BeginMode, count uint, dtype DataType, offset uint) {
	gl.Call("drawElements", uint(mode), count, uint(dtype), offset)
}
