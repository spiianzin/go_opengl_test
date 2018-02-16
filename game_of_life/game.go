package main

import (
    "log"
    "runtime"
    "strings"
    "fmt"
	"time"

    "github.com/go-gl/gl/v4.1-core/gl" 
    "github.com/go-gl/glfw/v3.2/glfw"
)

const (
	rows = 10
	columns = 10

    width  = 500
    height = 500
    vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 1, 1, 1);
    }
` + "\x00"
)

var (
    square = []float32 {
        -0.5, 0.5, 0,
        -0.5, -0.5, 0, 
        0.5, -0.5, 0,
		
		-0.5, 0.5, 0,
        0.5, -0.5, 0,
		0.5, 0.5, 0,		
    }
)

type cell struct {
	drawable uint32
	
	x int
	y int
	
	alive int
	aliveNextGen int
}

func makeCells() [][]*cell {
	cells := make([][]*cell, rows, columns)
	for x := 0; x < rows; x++ {
		for y := 0; y < columns; y++ {
			c := newCell(x, y)
			cells[x] = append(cells[x], c)
		}
	}
	
	return cells
}

func countAliveNeighboars(c *cell, cells [][]*cell) int {
	var aliveCount int
	var add = func (x int, y int) {
		if (x < 0) {
			x = rows - 1
		}
		if (x >= rows) {
			x = 0
		}
		if (y < 0) {
			y = columns - 1
		}
		
		if (y >= columns) {
			y = 0
		}
		
		if (cells[x][y].alive == 1) {
			aliveCount = aliveCount + cells[x][y].alive
		}
	}
	
	add(c.x - 1, c.y - 1)
	add(c.x - 1, c.y)
	add(c.x - 1, c.y + 1)
	add(c.x, c.y + 1)
	add(c.x, c.y - 1)
	add(c.x + 1, c.y - 1)
	add(c.x + 1, c.y)
	add(c.x + 1, c.y + 1)
	
	return aliveCount
}

func newCell(x int, y int) *cell {
	points := make([]float32, len(square), len(square))
    copy(points, square)

    for i := 0; i < len(points); i++ {
        var position float32
        var size float32
        switch i % 3 {
        case 0:
                size = 1.0 / float32(columns)
                position = float32(x) * size
        case 1:
                size = 1.0 / float32(rows)
                position = float32(y) * size
        default:
                continue
        }

        if points[i] < 0 {
                points[i] = (position * 2) - 1
        } else {
                points[i] = ((position + size) * 2) - 1
        }
    }

    return &cell{
        drawable: makeVao(points),

        x: x,
        y: y,
    }
}

func draw(cells [][]*cell, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)
	
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			
			var c = cells[i][j]
			if (c.alive == 1) {
				c.draw()
			}
			
			var alive = countAliveNeighboars(c, cells)
			
			c.aliveNextGen = 0
			
			if (c.alive == 1 && (alive == 2 || alive == 3)) {
				c.aliveNextGen = 1
			}
			if (c.alive == 0 && alive == 3) {
				c.aliveNextGen = 1
			}
		}
	}
	
	glfw.PollEvents()
	window.SwapBuffers()
	
	time.Sleep(500 * time.Millisecond)
	
	var aliveCells int
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			cells[i][j].swapGeneration()
			aliveCells = aliveCells + cells[i][j].alive
		}
	}
	
	if (aliveCells == 0) {
		log.Println("End of the game")
		time.Sleep(5 * time.Second)
		window.SetShouldClose(true)
	}
	
	//Will close the window when it will reach end of the loop
	// 
}

func main() {
    runtime.LockOSThread()

    window := initGlfw()
    defer glfw.Terminate()
    
    program := initOpenGL()

    cells := makeCells()

    cells[2][2].alive = 1
	cells[3][2].alive = 1
	cells[4][2].alive = 1
	cells[4][3].alive = 1
	cells[4][4].alive = 1
	cells[3][5].alive = 1

    for !window.ShouldClose() {
        draw(cells, window, program)
    }
}

func (c *cell) draw() {
    gl.BindVertexArray(c.drawable)
    gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square) / 3))
}

func (c *cell) swapGeneration() {
	c.alive = c.aliveNextGen
}


func makeVao(points []float32) uint32 {
    var vbo1 uint32
    gl.GenBuffers(1, &vbo1)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo1)
    gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
    
    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    gl.EnableVertexAttribArray(0)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo1)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
    
    return vao

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

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
    if err := glfw.Init(); err != nil {
            panic(err)
    }
    
    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 4) 
    glfw.WindowHint(glfw.ContextVersionMinor, 1)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
    if err != nil {
            panic(err)
    }
    window.MakeContextCurrent()

    return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
    if err := gl.Init(); err != nil {
            panic(err)
    }
    version := gl.GoStr(gl.GetString(gl.VERSION))
    log.Println("OpenGL version", version)

    vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil {
        panic(err)
    }
    fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil {
        panic(err)
    }

    prog := gl.CreateProgram()
    gl.AttachShader(prog, vertexShader)
    gl.AttachShader(prog, fragmentShader) 
    gl.LinkProgram(prog)
    
    return prog
}
