package main

import (
	"math/rand"
	"time"
)

type cell struct {
	drawable uint32
}

func makeBoard() *board {

	if RandomSeed == 0 {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(RandomSeed)
	}

	intcolumns := int(BoardColumns)
	introws := int(BoardRows)

	FrameBufferSize := 14

	cells := make([][]*cell, introws)
	B := board{
		Rows:           introws,
		Columns:        intcolumns,
		Cells:          cells,
		Current:        &Frame{},
		UpdateInterval: int64(time.Second / time.Duration(UpdatesPerSecond)),
		FrameInterval:  int64(time.Second / time.Duration(FramesPerSecond)),
		DrawableFrames: make(chan *Frame, FrameBufferSize+2),
		DirtyFrames:    make(chan *Frame, FrameBufferSize+2),
	}

	for i := 0; i < FrameBufferSize; i++ {
		frm := Frame{
			Alive: make([][]bool, introws),
		}

		B.DirtyFrames <- &frm
	}

	for x := 0; x < introws; x++ {

		cells[x] = make([]*cell, intcolumns)

		for y := 0; y < intcolumns; y++ {
			c := newCell(x, y)

			cells[x][y] = c
		}

		for i := 0; i < FrameBufferSize; i++ {
			frm := <-B.DirtyFrames

			frm.Alive[x] = make([]bool, intcolumns)

			B.DirtyFrames <- frm
		}

	}

	if LoadBoard != nil {
		B.Current.Alive = LoadBoard().Alive

	} else if CreateEmpty {
		B.Current.Alive = make([][]bool, introws)

		for x := 0; x < introws; x++ {
			B.Current.Alive[x] = make([]bool, intcolumns)
		}

	} else {
		B.Current.Alive = make([][]bool, introws)

		for x := 0; x < introws; x++ {
			B.Current.Alive[x] = make([]bool, intcolumns)

			for y := 0; y < intcolumns; y++ {

				alive := rand.Float64() < RandomThreshold
				B.Current.Alive[x][y] = alive
			}
		}
	}

	B.DrawableFrames <- B.Current
	return &B
}

func newCell(x, y int) *cell {
	points := make([]float32, len(square))
	copy(points, square)

	sizex := 2.0 / BoardRows
	sizey := 2.0 / BoardColumns
	positionx := (float32(x) * sizex) - 1
	positiony := (float32(y) * sizey) - 1

	for ix := 0; ix < len(points); ix += 3 {
		points[ix] = points[ix]*sizex + positionx
		points[ix+1] = points[ix+1]*sizey + positiony
	}

	return &cell{
		drawable: makeVao(points),
	}
}
