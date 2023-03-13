package engine

import "time"

// -----------------------------------------------------------------------------
//	Identifying Constants
// -----------------------------------------------------------------------------

const (
	name   string = "Light Blue 0"
	author string = "Sidhant Roymoulik"
)

// -----------------------------------------------------------------------------
//	Debug
// -----------------------------------------------------------------------------

const DEBUG bool = false

// -----------------------------------------------------------------------------
//	Openings
// -----------------------------------------------------------------------------

var CHESS_OPENINGS map[string][]string = map[string][]string{
	"Start Position":   {},
	"Sicilian Defense": {"e4", "c5"},
	"Italian Game":     {"e4", "e5", "Nf3", "Nc6", "Bc4"},
}

// -----------------------------------------------------------------------------
//	Parameters
// -----------------------------------------------------------------------------

const (
	TIME_LIMIT      time.Duration = 2 * time.Second // Time in sec
	CHECKMATE_VALUE int           = 1000000
	MATE_CUTOFF     int           = CHECKMATE_VALUE / 2
	MAX_DEPTH       int           = 100
	TIMER_CHECK     uint64        = (1 << 3) - 1

	WINDOW_VALUE_TIGHT int = 25
	WINDOW_VALUE       int = 100

	StaticNullMovePruningBaseMargin int = 85
	NMR_Depth_Limit                 int = 2
	FutilityPruningDepthLimit       int = 8
	IID_Depth_Limit                 int = 4
	IID_Depth_Reduction             int = 2
)

var FutilityMargins = [9]int{
	0,
	100, // depth 1
	160, // depth 2
	220, // depth 3
	280, // depth 4
	340, // depth 5
	400, // depth 6
	460, // depth 7
	520, // depth 8
}
