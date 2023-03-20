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
	TIMER_CHECK     uint64        = (1 << 10) - 1
)
