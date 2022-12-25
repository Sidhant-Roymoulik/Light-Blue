package main

import "time"

// Debug
const DEBUG bool = true

const CHESS_START_POSITION string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Parameters
const TIME_LIMIT time.Duration = 5 * 1000000000 // Time in nanosec, set to 0 for no time limit
const MAX_CONST_DEPTH int = 2
const CHECKMATE_VALUE int = 1000000

// Counters
var start time.Time
var states int = 0
var q_states int = 0
var hashes int = 0
