package main

var (
	WindowWidth  int = 1024
	WindowHeight int = 640

	FramesPerSecond  int = 20
	UpdatesPerSecond int = 20

	InitDelay int = 500

	// Board Creation

	BoardRows    float32 = 128
	BoardColumns float32 = 128

	RandomSeed      int64   = 0
	RandomThreshold float64 = 0.15

	CreateEmpty bool = false
	LoadBoard   func() *Frame

	PresentChar string = "x"
	AbsentChar  string = "-"

	SaveOnlyMeta   string = ""
	SaveOnCreation string = ""
	SaveOnKeyPress string = ""
	SaveOnExit     string = ""

	LoadOnlyMeta bool = false
)
