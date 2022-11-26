package main

// Debug
const DEBUG bool = true

const CHESS_START_POSITION string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Parameters
const TIME_LIMIT int = 5 // Time in sec, set to 0 for no time limit
const MAX_CONST_DEPTH int = 3
const CHECKMATE_VALUE int = 100000

// Strategies
const MOVE_ORDERING bool = false
const ALPHA_BETA_PRUNING bool = false
const QUIESCENCE_SEARCH bool = false
const ITERATIVE_DEEPENING bool = false
const TRANSPOSITION_TABLE bool = false

// Counters
var states int = 0
var q_states int = 0
var hashes int = 0
