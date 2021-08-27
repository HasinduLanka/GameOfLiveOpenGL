package main

import (
	"math/rand"
	"time"
)

type cell struct {
	drawable uint32
}

func makeBoard() *board {
	rand.Seed(time.Now().UnixNano())
	intcolumns := int(columns)
	introws := int(rows)

	FrameBufferSize := int(UpdatesPerSecond/FramesPerSecond) + 4

	cells := make([][]*cell, introws)
	B := board{
		Rows:           introws,
		Columns:        intcolumns,
		Cells:          cells,
		Current:        &Frame{Alive: make([][]bool, introws)},
		UpdateInterval: int64(time.Second / time.Duration(UpdatesPerSecond)),
		FrameInterval:  int64(time.Second / time.Duration(FramesPerSecond)),
		DrawableFrames: make(chan *Frame, FrameBufferSize),
		DirtyFrames:    make(chan *Frame, FrameBufferSize),
	}

	for i := 0; i < FrameBufferSize; i++ {
		frm := Frame{
			Alive: make([][]bool, introws),
		}

		B.DirtyFrames <- &frm
	}

	for x := 0; x < introws; x++ {

		cells[x] = make([]*cell, intcolumns)
		B.Current.Alive[x] = make([]bool, intcolumns)

		for y := 0; y < intcolumns; y++ {
			c := newCell(x, y)

			alive := rand.Float64() < threshold
			B.Current.Alive[x][y] = alive

			cells[x][y] = c
		}

		for i := 0; i < FrameBufferSize; i++ {
			frm := <-B.DirtyFrames

			frm.Alive[x] = make([]bool, intcolumns)

			B.DirtyFrames <- frm
		}

	}

	return &B
}

func newCell(x, y int) *cell {
	points := make([]float32, len(square))
	copy(points, square)

	sizex := 2.0 / rows
	sizey := 2.0 / columns
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
