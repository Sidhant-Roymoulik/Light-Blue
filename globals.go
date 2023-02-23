package main

import "time"

// -----------------------------------------------------------------------------
//	Debug
// -----------------------------------------------------------------------------

const DEBUG bool = true

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

const TIME_LIMIT time.Duration = 2 * time.Second // Time in sec
const CHECKMATE_VALUE int = 1000000
const MATE_CUTOFF int = CHECKMATE_VALUE / 2
const MAX_DEPTH int = 100

const WINDOW_VALUE_TIGHT int = 25
const WINDOW_VALUE int = 100
const IID_Depth_Limit int = 4

const R int = 2

const MAX_CONST_DEPTH int = 2
const MTD_ITER_CUTOFF int = 14
const MTD_EVAL_CUTOFF int = 0

// -----------------------------------------------------------------------------
//	Counters
// -----------------------------------------------------------------------------

var start time.Time
var states int = 0
var q_states int = 0
var hash_hits int = 0
var hash_reads int = 0
var hash_writes int = 0
var hash_collisions int = 0
