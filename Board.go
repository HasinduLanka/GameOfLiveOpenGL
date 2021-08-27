package main

import "github.com/go-gl/gl/v4.1-core/gl"

type board struct {
	FrameInterval  int64
	UpdateInterval int64

	Rows    int
	Columns int
	Cells   [][]*cell
	Current *Frame

	DrawableFrames chan *Frame
	DirtyFrames    chan *Frame
}

type Frame struct {
	Alive [][]bool
}

// Board tick determines the state of the board for the next tick of the game.
func (board *board) tick() {

	// Get a dirty frame to work with. We do not create a new frame, to avoid unnecessary memory allocations.
	NextFrame := board.GetNextDirtyFrame()

	// Get the current frame at the momment of the tick. board.Current could change during the excution of this function.
	Frm := board.Current

	Routes := 0
	RoutesCh := make(chan interface{})

	for X := 0; X < board.Rows; X++ {

		// Start a goroutine for each row.
		// Utilize CPU cores for parallelism.
		Routes++
		go func(x int) {

			for y := 0; y < board.Columns; y++ {
				liveCount := board.liveNeighbors(Frm, x, y)

				isAlive := Frm.Alive[x][y]
				if isAlive {

					if liveCount < 2 {
						// 1. Any live cell with fewer than two live neighbours dies, as if caused by underpopulation.
						NextFrame.Alive[x][y] = false

					} else if liveCount > 3 {
						// 2. Any live cell with more than three live neighbours dies, as if by overpopulation.
						NextFrame.Alive[x][y] = false

					} else {
						// 3. Any live cell with two or three live neighbours lives on to the next generation.
						NextFrame.Alive[x][y] = true
					}

				} else {
					// 4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
					if liveCount == 3 {
						NextFrame.Alive[x][y] = true
					} else {
						NextFrame.Alive[x][y] = isAlive
					}
				}
			}
			RoutesCh <- nil

		}(X)
	}

	for r := 0; r < Routes; r++ {
		<-RoutesCh
	}

	board.Current = NextFrame

	// Send the next frame to the drawable frames. We do not change it any time soon.
	board.DrawableFrames <- NextFrame

}

func (board *board) draw() {

	// If there are more than 1 drawable frames, put the rest in the dirty frames.
	for len(board.DrawableFrames) > 1 {
		board.DirtyFrames <- <-board.DrawableFrames
	}

	// Get a drawable frame to work with.
	frm := <-board.DrawableFrames

	// Draw the frame.
	for x := 0; x < board.Rows; x++ {
		for y := 0; y < board.Columns; y++ {
			if frm.Alive[x][y] {
				drawable := board.Cells[x][y].drawable
				gl.BindVertexArray(drawable)
				gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
			}
		}
	}

	// Put the frame in the dirty frames.
	board.DirtyFrames <- frm

}

// liveNeighbors returns the number of live neighbors for a cell.
func (board *board) liveNeighbors(frame *Frame, cx int, cy int) int {
	var liveCount int
	add := func(x, y int) {
		// If we're at an edge, check the other side of the board.
		if x == board.Rows {
			x = 0
		} else if x == -1 {
			x = board.Rows - 1
		}
		if y == board.Columns {
			y = 0
		} else if y == -1 {
			y = board.Columns - 1
		}

		if frame.Alive[x][y] {
			liveCount++
		}
	}

	add(cx-1, cy)   // To the left
	add(cx+1, cy)   // To the right
	add(cx, cy+1)   // up
	add(cx, cy-1)   // down
	add(cx-1, cy+1) // top-left
	add(cx+1, cy+1) // top-right
	add(cx-1, cy-1) // bottom-left
	add(cx+1, cy-1) // bottom-right

	return liveCount
}

func (board *board) GetNextDirtyFrame() *Frame {

	// If there are no dirty frames, get from drawable frames.
	if len(board.DirtyFrames) == 0 {

		for len(board.DrawableFrames) > 2 {
			board.DirtyFrames <- <-board.DrawableFrames
		}

		return <-board.DrawableFrames
	}

	return <-board.DirtyFrames
}
