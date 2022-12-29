package main

import "time"

// --------------------------------------------------------------------------------------
//	Debug
// --------------------------------------------------------------------------------------

const DEBUG bool = true

// --------------------------------------------------------------------------------------
//	Openings
// --------------------------------------------------------------------------------------

var CHESS_FENs map[string]string = map[string]string{
	"Start Position": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	"Italian Game":   "r1bqk1nr/pppp1ppp/2n5/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1",
}

// --------------------------------------------------------------------------------------
//	Parameters
// --------------------------------------------------------------------------------------

const TIME_LIMIT time.Duration = 5 * 1000000000 // Time in nanosec
const MAX_CONST_DEPTH int = 2
const CHECKMATE_VALUE int = 1000000

// --------------------------------------------------------------------------------------
//	Counters
// --------------------------------------------------------------------------------------

var start time.Time
var states int = 0
var q_states int = 0
var hash_hits int = 0
var hash_reads int = 0
var hash_writes int = 0
var hash_collisions int = 0
