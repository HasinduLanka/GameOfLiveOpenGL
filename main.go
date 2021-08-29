package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {

	//args := strings.Split("-i save1.board", " ")
	// args := strings.Split("-i test.board --SaveOnCreation save1.board --SaveOnlyMeta save1.meta.board --SaveOnExit save1.exit.board", " ")

	args := os.Args
	CheckAndLoadFromFile(args)

	if RunArgs(args) {
		return
	}

	fitWindow()

	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	if len(SaveOnlyMeta) != 0 {
		realLoadOnlyMeta := LoadOnlyMeta
		LoadOnlyMeta = true
		SaveMataToFile(SaveOnlyMeta)
		LoadOnlyMeta = realLoadOnlyMeta
	}

	board := makeBoard()

	if len(SaveOnCreation) != 0 {
		SaveToFile(SaveOnCreation, board, board.Current)
	}

	drawFrame(window, board, program)

	time.Sleep(time.Millisecond * time.Duration(InitDelay))

	go UpdateLoop(window, board)
	DrawLoop(window, board, program)

	OnExit(board)
}

func OnExit(board *board) {
	if len(SaveOnExit) != 0 {
		SaveToFile(SaveOnExit, board, board.Current)
	}

}

// Returns true if the program needs to exit
func RunArgs(args []string) bool {

	SkipNext := false

	for i := 0; i < len(args); i++ {
		if SkipNext {
			SkipNext = false
			continue
		}

		switch args[i] {
		case "-h", "--help":
			PrintHelp()
			return true

		case "-e", "--empty":
			CreateEmpty = true

		case "-m", "--LoadOnlyMeta":
			LoadOnlyMeta = true

		case "-r", "--rows":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					BoardRows = float32(val)
					SkipNext = true
				}
			}
		case "-c", "--columns":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					BoardColumns = float32(val)
					SkipNext = true
				}
			}

		case "-s", "--seed":
			if i+1 < len(args) {
				if val, err := strconv.ParseInt(args[i+1], 10, 64); err == nil {
					RandomSeed = int64(val)
					SkipNext = true
				}
			}

		case "-thr", "--threshold":
			if i+1 < len(args) {
				if val, err := strconv.ParseFloat(args[i+1], 64); err == nil {
					RandomThreshold = val
					SkipNext = true
				}
			}

		case "-ww", "--width":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					WindowWidth = val
					SkipNext = true
				}
			}
		case "-wh", "--height":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					WindowHeight = val
					SkipNext = true
				}
			}
		case "-u", "--ups":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					UpdatesPerSecond = val
					SkipNext = true
				}
			}
		case "-f", "--fps":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					FramesPerSecond = val
					SkipNext = true
				}
			}
		case "-d", "--delay":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					InitDelay = val
					SkipNext = true
				}
			}

		case "-pc", "--PresentChar":
			if i+1 < len(args) {
				if val := args[i+1]; len([]byte(val)) == 1 {
					PresentChar = val
					SkipNext = true
				}
			}
		case "-ac", "--AbsentChar":
			if i+1 < len(args) {
				if val := args[i+1]; len([]byte(val)) == 1 {
					AbsentChar = val
					if AbsentChar == "s" {
						AbsentChar = " "
					}
					SkipNext = true
				}
			}

		case "-om", "--SaveOnlyMeta":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					SaveOnlyMeta = val
					SkipNext = true
				}
			}

		case "-os", "--SaveOnCreation":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					SaveOnCreation = val
					SkipNext = true
				}
			}
		case "-oe", "--SaveOnExit":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					SaveOnExit = val
					SkipNext = true
				}
			}
		case "-o", "--SaveOnKeyPress":
			if i+1 < len(args) {
				if val := args[i+1]; len(val) != 0 {
					SaveOnKeyPress = val
					SkipNext = true
				}
			}

		}
	}

	return false
}

func CheckAndLoadFromFile(args []string) bool {

	for i := 0; i < len(args); i++ {

		switch args[i] {
		case "-i", "--input-file":
			if i+1 < len(args) {
				filename := args[i+1]

				return LoadFromFile(filename)

			}
		}
	}

	return false
}

func fitWindow() {
	CellsAR := BoardRows / BoardColumns
	WindowAR := float32(WindowWidth) / float32(WindowHeight)

	if CellsAR > WindowAR {
		WindowHeight = int(float32(WindowWidth) / CellsAR)
	} else {
		WindowWidth = int(float32(WindowHeight) * CellsAR)
	}

}

func PrintHelp() {
	print("Usage:")
}

func UpdateLoop(window *glfw.Window, board *board) {
	for !window.ShouldClose() {
		t := time.Now()

		board.tick()

		time.Sleep(time.Nanosecond*time.Duration(board.UpdateInterval) - time.Since(t))
	}
}

func drawFrame(window *glfw.Window, board *board, program uint32) {
	t := time.Now()

	draw(board, window, program)

	time.Sleep(time.Nanosecond/time.Duration(board.FrameInterval) - time.Since(t))
}

func DrawLoop(window *glfw.Window, board *board, program uint32) {
	for !window.ShouldClose() {
		drawFrame(window, board, program)
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
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
