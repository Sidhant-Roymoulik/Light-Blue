package main

import (
	"github.com/Sidhant-Roymoulik/chess"
)

// -----------------------------------------------------------------------------
// 		Piece Values
// -----------------------------------------------------------------------------

// piece weights
const pawn int = 100
const knight int = 320
const bishop int = 330
const rook int = 525
const queen int = 900
const king int = 20000

// piece value map
var PVM map[chess.PieceType]int = map[chess.PieceType]int{
	chess.King:   king,
	chess.Queen:  queen,
	chess.Rook:   rook,
	chess.Bishop: bishop,
	chess.Knight: knight,
	chess.Pawn:   pawn,
}

// -----------------------------------------------------------------------------
// 		Piece Square Table Stuff
// -----------------------------------------------------------------------------

var FLIP = []int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

var PST_MG = map[chess.PieceType][]int{
	chess.Pawn: {
		0, 0, 0, 0, 0, 0, 0, 0,
		60, 60, 60, 60, 70, 60, 60, 60,
		40, 40, 40, 50, 60, 40, 40, 40,
		20, 20, 20, 40, 50, 20, 20, 20,
		5, 5, 15, 30, 40, 10, 5, 5,
		5, 5, 10, 20, 30, 5, 5, 5,
		5, 5, 5, -30, -30, 5, 5, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, -5, -5, -5, -5, -5, -5, -10,
		-10, -5, 15, 15, 15, 15, -5, -10,
		-10, -5, 15, 15, 15, 15, -5, -10,
		-10, -5, 15, 15, 15, 15, -5, -10,
		-10, -5, 10, 15, 15, 15, -5, -10,
		-10, -5, -5, -5, -5, -5, -5, -10,
		-20, 0, -10, -10, -10, -10, 0, -20,
	},
	chess.Bishop: {
		-20, 0, 0, 0, 0, 0, 0, -20,
		-15, 0, 0, 0, 0, 0, 0, -15,
		-10, 0, 0, 5, 5, 0, 0, -10,
		-10, 10, 10, 30, 30, 10, 10, -10,
		5, 5, 10, 25, 25, 10, 5, 5,
		5, 5, 5, 10, 10, 5, 5, 5,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	},
	chess.Rook: {
		0, 0, 5, 10, 10, 5, 0, 0,
		25, 25, 25, 25, 25, 25, 25, 25,
		0, 0, 5, 10, 10, 5, 0, 0,
		0, 0, 5, 10, 10, 5, 0, 0,
		0, 0, 5, 10, 10, 5, 0, 0,
		0, 0, 5, 10, 10, 5, 0, 0,
		0, 0, 5, 10, 10, 5, 0, 0,
		0, 0, 5, 10, 10, 5, 0, 0,
	},
	chess.Queen: {
		-30, -20, -10, -10, -10, -10, -20, -30,
		-20, -10, -5, -5, -5, -5, -10, -20,
		-10, -5, 10, 10, 10, 10, -5, -10,
		-10, -5, 10, 20, 20, 10, -5, -10,
		-10, -5, 10, 20, 20, 10, -5, -10,
		-10, -5, -5, -5, -5, -5, -5, -10,
		-20, -10, -5, -5, -5, -5, -10, -20,
		-30, -20, -10, -10, -10, -10, -20, -30,
	},
	chess.King: {
		-70, -70, -70, -70, -70, -70, -70, -70,
		-70, -70, -70, -70, -70, -70, -70, -70,
		-70, -70, -70, -70, -70, -70, -70, -70,
		-70, -70, -70, -70, -70, -70, -70, -70,
		-70, -70, -70, -70, -70, -70, -70, -70,
		-10, -10, -10, -10, -10, -10, -10, -10,
		0, 0, 5, -10, -10, 0, 5, 0,
		0, 0, 20, -10, -10, 0, 20, 0,
	},
}

var PST_EG = map[chess.PieceType][]int{
	chess.Pawn:   PST_MG[chess.Pawn],
	chess.Knight: PST_MG[chess.Knight],
	chess.Bishop: PST_MG[chess.Bishop],
	chess.Rook:   PST_MG[chess.Rook],
	chess.Queen:  PST_MG[chess.Queen],
	chess.King: {
		-50, -10, 0, 0, 0, 0, -10, -50,
		-10, 0, 10, 10, 10, 10, 0, -10,
		0, 10, 15, 15, 15, 15, 10, 0,
		0, 10, 15, 20, 20, 15, 10, 0,
		0, 10, 15, 20, 20, 15, 10, 0,
		0, 10, 15, 15, 15, 15, 10, 0,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-50, -10, 0, 0, 0, 0, -10, -50,
	},
}

// -----------------------------------------------------------------------------
// 		Tapered Evaluation Values
// -----------------------------------------------------------------------------

var PawnPhase int = 0
var KnightPhase int = 1
var BishopPhase int = 1
var RookPhase int = 2
var QueenPhase int = 4
var TotalPhase int = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

// -----------------------------------------------------------------------------
// 		Position Evaluation Function
// -----------------------------------------------------------------------------

// Checks for Checkmate, Stalemate, and Total Piece Delta
func eval_v1(position *chess.Position) int {
	// faster than doing two comparisons
	if position.Status() != chess.NoMethod {
		if position.Status() == chess.Stalemate {
			return 0
		}
		if position.Status() == chess.Checkmate {
			if position.Turn() == chess.White {
				return -CHECKMATE_VALUE
			} else {
				return CHECKMATE_VALUE
			}
		}
	}

	squares := position.Board().SquareMap()
	var delta int = 0
	for _, piece := range squares {
		if piece.Color() == chess.Black {
			delta -= PVM[piece.Type()]
		} else {
			delta += PVM[piece.Type()]
		}
	}
	return delta
}

// Also uses PST
func eval_v2(position *chess.Position) int {
	// faster than doing two comparisons
	if position.Status() != chess.NoMethod {
		if position.Status() == chess.Stalemate {
			return 0
		}
		if position.Status() == chess.Checkmate {
			if position.Turn() == chess.White {
				return -CHECKMATE_VALUE
			} else {
				return CHECKMATE_VALUE
			}
		}
	}

	squares := position.Board().SquareMap()
	var delta int = 0
	for square, piece := range squares {
		if piece.Color() == chess.Black {
			delta -= PVM[piece.Type()]
			delta -= PST_MG[piece.Type()][square]
		} else {
			delta += PVM[piece.Type()]
			delta += PST_MG[piece.Type()][FLIP[square]]
		}
	}
	return delta
}

// Checkmate ply penalty
func eval_v3(position *chess.Position, ply int) int {
	// faster than doing two comparisons
	if position.Status() != chess.NoMethod {
		if position.Status() == chess.Stalemate {
			return 0
		}
		if position.Status() == chess.Checkmate {
			if position.Turn() == chess.White {
				return -CHECKMATE_VALUE + ply
			} else {
				return CHECKMATE_VALUE - ply
			}
		}
	}

	squares := position.Board().SquareMap()
	var delta int = 0
	for square, piece := range squares {
		if piece.Color() == chess.White {
			delta += PVM[piece.Type()]
			delta += PST_MG[piece.Type()][FLIP[square]]
		} else {
			delta -= PVM[piece.Type()]
			delta -= PST_MG[piece.Type()][square]
		}
	}

	return delta
}

// Includes Piece Mobility
func eval_v4(position *chess.Position, ply int) int {
	// faster than doing two comparisons
	if position.Status() != chess.NoMethod {
		if position.Status() == chess.Checkmate {
			if position.Turn() == chess.White {
				return -CHECKMATE_VALUE + ply
			} else {
				return CHECKMATE_VALUE - ply
			}
		}
		return 0
	}

	squares := position.Board().SquareMap()
	var delta int = 0
	for square, piece := range squares {
		if piece.Color() == chess.White {
			delta += PVM[piece.Type()]
			delta += PST_MG[piece.Type()][FLIP[square]]
		} else {
			delta -= PVM[piece.Type()]
			delta -= PST_MG[piece.Type()][square]
		}
	}

	if position.Turn() == chess.White {
		delta += len(position.ValidMoves())
		delta -= len(position.NullMove().ValidMoves())
	} else {
		delta += len(position.NullMove().ValidMoves())
		delta -= len(position.ValidMoves())
	}

	return delta
}

// Tapered Evaluation
func eval_v5(position *chess.Position, ply int) int {
	// faster than doing two comparisons
	if position.Status() != chess.NoMethod {
		if position.Status() == chess.Checkmate {
			if position.Turn() == chess.White {
				return -CHECKMATE_VALUE + ply
			} else {
				return CHECKMATE_VALUE - ply
			}
		}
		return 0
	}

	var POCC map[chess.PieceType]int = map[chess.PieceType]int{
		chess.King:   0,
		chess.Queen:  0,
		chess.Rook:   0,
		chess.Bishop: 0,
		chess.Knight: 0,
		chess.Pawn:   0,
	}

	squares := position.Board().SquareMap()
	var delta_mg int = 0
	var delta_eg int = 0
	for square, piece := range squares {
		POCC[piece.Type()] += 1
		if piece.Color() == chess.White {
			delta_mg += PVM[piece.Type()] + PST_MG[piece.Type()][FLIP[square]]
			delta_eg += PVM[piece.Type()] + PST_EG[piece.Type()][FLIP[square]]
		} else {
			delta_mg -= PVM[piece.Type()] + PST_MG[piece.Type()][square]
			delta_eg -= PVM[piece.Type()] + PST_EG[piece.Type()][square]
		}
	}

	var phase int = TotalPhase
	phase -= POCC[chess.Pawn] * PawnPhase
	phase -= POCC[chess.Knight] * KnightPhase
	phase -= POCC[chess.Bishop] * BishopPhase
	phase -= POCC[chess.Rook] * RookPhase
	phase -= POCC[chess.Queen] * QueenPhase
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase

	var delta int = ((delta_mg * (256 - phase)) + (delta_eg * phase)) / 256

	if position.Turn() == chess.White {
		delta += len(position.ValidMoves())
		delta -= len(position.NullMove().ValidMoves())
	} else {
		delta += len(position.NullMove().ValidMoves())
		delta -= len(position.ValidMoves())
	}

	return delta
}
