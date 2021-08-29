package main

import (
	"strconv"
	"strings"
)

type LineProvider struct {
	ch chan string
}

// Returns true if successful
func LoadFromFile(filename string) bool {

	LP, LPErr := NewLineProvider_FromFile(filename)

	if PrintError(LPErr) {
		return false
	}

	metaline := <-LP.ch

	rawmetas := strings.Split(metaline, " ")
	metas := make([]string, len(rawmetas))

	for i, rawmeta := range rawmetas {
		metas[i] = strings.Trim(rawmeta, " ")
	}

	if len(metas) < 2 || metas[0] != "GameOfLife" {
		Print("File load error : Not enough metadata")
		Print("File load error : File should begin with GameOfLife, followed by a PresentChar char")

		return false
	}

	PresentChar := metas[1]

	RunArgs(metas)

	if LoadOnlyMeta {
		LoadBoard = nil
		LoadOnlyMeta = false

		go func() {
			for range LP.ch {
				//Discard All and close the channel
			}
		}()

	} else {
		IgnoredLine := <-LP.ch
		Print(IgnoredLine)

		LoadBoard = func() *Frame {

			intBoardColumns := int(BoardColumns)
			intBoardRows := int(BoardRows)
			F := Frame{Alive: make([][]bool, intBoardRows)}
			for r := 0; r < int(BoardRows); r++ {
				F.Alive[r] = make([]bool, intBoardColumns)
			}

			col := intBoardColumns
			for line := range LP.ch {
				col--

				// Recorrect little mistakes that can heppen during manual editing of board files
				lenline := len(line)
				if lenline > intBoardRows {
					lenline = intBoardRows
				}

				for row, chr := range line[:lenline] {
					F.Alive[row][col] = string(chr) == PresentChar
				}
			}

			return &F
		}
	}

	return true
}

func GlobalsToMetaString() string {
	meta := strings.Builder{}

	meta.WriteString(" --rows " + strconv.Itoa(int(BoardRows)))
	meta.WriteString(" --columns " + strconv.Itoa(int(BoardColumns)))

	meta.WriteString(" --width " + strconv.Itoa(int(WindowWidth)))
	meta.WriteString(" --height " + strconv.Itoa(int(WindowHeight)))

	meta.WriteString(" --ups " + strconv.Itoa(int(UpdatesPerSecond)))
	meta.WriteString(" --fps " + strconv.Itoa(int(FramesPerSecond)))

	meta.WriteString(" --delay " + strconv.Itoa(InitDelay))

	meta.WriteString(" --seed " + strconv.Itoa(int(RandomSeed)))
	meta.WriteString(" --threshold " + strconv.FormatFloat(RandomThreshold, 'f', -1, 64))

	meta.WriteString(" --PresentChar " + PresentChar)

	AbsentCharSave := AbsentChar
	if AbsentChar == " " {
		AbsentCharSave = "s"
	}
	meta.WriteString(" --AbsentChar " + AbsentCharSave)

	if LoadOnlyMeta {
		meta.WriteString(" --LoadOnlyMeta ")
	}

	if len(SaveOnExit) != 0 {
		meta.WriteString(" --SaveOnExit " + SaveOnExit)
	}

	return meta.String()
}

func SaveMataToFile(filename string) {
	DeleteFiles(filename)
	AppendLine(filename, "GameOfLife "+PresentChar+GlobalsToMetaString()+"\n")
}

func SaveToFile(filename string, board *board, frm *Frame) {

	SaveMataToFile(filename)

	ch := make(chan []byte, 64)
	go AppendChan(filename, ch)

	PresendByte := ([]byte(PresentChar))[0]
	AbsentByte := ([]byte(AbsentChar))[0]
	NewLineByte := ([]byte("\n"))[0]

	for y := board.Columns - 1; y > -1; y-- {
		row := make([]byte, board.Rows+1)
		row[board.Rows] = NewLineByte
		for x := 0; x < board.Rows; x++ {
			if frm.Alive[x][y] {
				row[x] = PresendByte
			} else {
				row[x] = AbsentByte
			}
		}

		ch <- row

	}
}
