package main

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	board := makeBoard()

	go UpdateLoop(window, board)
	DrawLoop(window, board, program)

}

func UpdateLoop(window *glfw.Window, board *board) {
	for !window.ShouldClose() {
		t := time.Now()

		board.tick()

		time.Sleep(time.Nanosecond*time.Duration(board.UpdateInterval) - time.Since(t))
	}
}

func DrawLoop(window *glfw.Window, board *board, program uint32) {
	for !window.ShouldClose() {
		t := time.Now()

		draw(board, window, program)

		time.Sleep(time.Nanosecond/time.Duration(board.FrameInterval) - time.Since(t))
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func draw(board *board, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	board.draw()

	glfw.PollEvents()
	window.SwapBuffers()
}
