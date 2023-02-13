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
	chess.Queen:  921,
	chess.Rook:   441,
	chess.Bishop: 346,
	chess.Knight: 333,
	chess.Pawn:   84,
}

var PVM_EG map[chess.PieceType]int = map[chess.PieceType]int{
	chess.King:   king,
	chess.Queen:  886,
	chess.Rook:   478,
	chess.Bishop: 268,
	chess.Knight: 244,
	chess.Pawn:   106,
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
		45, 52, 42, 43, 28, 34, 19, 9,
		-14, -3, 7, 14, 35, 50, 15, -6,
		-27, -6, -8, 13, 16, 4, -3, -25,
		-32, -28, -7, 5, 7, -1, -15, -30,
		-29, -25, -12, -12, -1, -5, 6, -17,
		-34, -23, -27, -18, -14, 10, 13, -22,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-43, -11, -8, -5, 1, -20, -4, -22,
		-31, -22, 19, 7, 5, 13, -8, -11,
		-21, 21, 8, 16, 36, 33, 19, 6,
		-6, 2, 0, 23, 8, 27, 4, 14,
		-3, 10, 12, 8, 16, 10, 19, 1,
		-19, -4, 3, 7, 22, 12, 15, -11,
		-21, -20, -9, 8, 9, 11, -5, 0,
		-19, -13, -20, -14, -2, 3, -11, -8,
	},
	chess.Bishop: {
		-13, 0, -17, -8, -7, -5, -2, -3,
		-21, 0, -16, -10, 4, 1, -6, -41,
		-23, 6, 10, 8, 8, 26, 0, -10,
		-15, -4, 2, 22, 9, 10, -1, -16,
		0, 10, -2, 15, 17, -7, -1, 13,
		-2, 16, 13, 0, 5, 16, 14, 0,
		8, 11, 12, 3, 11, 23, 27, 3,
		-26, 3, -3, -1, 10, -5, -7, -15,
	},
	chess.Rook: {
		3, 1, 0, 7, 7, -1, 0, 0,
		-6, -9, 7, 7, 7, 5, -4, -1,
		-12, 11, 0, 17, -2, 12, 23, -1,
		-17, -9, 4, 0, 3, 15, -1, -2,
		-24, -16, -16, -4, -1, -14, 2, -20,
		-30, -15, -6, -3, 0, 2, 2, -15,
		-25, -6, -6, 5, 8, 6, 8, -46,
		-3, 1, 6, 15, 17, 14, -13, -2,
	},
	chess.Queen: {
		-10, 0, 0, 0, 10, 9, 5, 7,
		-19, -35, -5, 2, -9, 7, 1, 15,
		-10, -7, -4, -9, 15, 29, 24, 22,
		-14, -14, -15, -11, -1, -5, 3, -6,
		-8, -20, -8, -5, -4, -2, 2, -2,
		-13, 5, 2, 1, -1, 8, 4, 2,
		-20, 0, 10, 16, 16, 16, -6, 6,
		-3, -1, 7, 19, 5, -10, -9, -17,
	},
	chess.King: {
		-3, 0, 2, 0, 0, 0, 1, -1,
		1, 4, 0, 7, 4, 2, 3, -2,
		2, 4, 7, 4, 4, 14, 12, 0,
		0, 2, 6, 0, 0, 2, 6, -9,
		-8, 5, 0, -8, -10, -10, -9, -23,
		-3, 5, 1, -8, -12, -12, 8, -24,
		6, 13, 0, -40, -23, -1, 25, 19,
		-28, 29, 17, -53, 2, -25, 34, 15,
	},
}

var PST_EG = map[chess.PieceType][]int{
	chess.Pawn: {
		0, 0, 0, 0, 0, 0, 0, 0,
		77, 74, 63, 53, 59, 60, 72, 77,
		17, 11, 11, 11, 11, -6, 14, 8,
		-3, -14, -18, -31, -29, -25, -20, -18,
		-12, -14, -24, -31, -29, -28, -27, -28,
		-22, -20, -25, -20, -21, -24, -34, -34,
		-16, -22, -11, -19, -13, -23, -32, -34,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-36, -16, -7, -14, -4, -20, -20, -29,
		-17, 2, -7, 14, 2, -7, -9, -19,
		-13, -7, 14, 12, 4, 6, 0, -13,
		-5, 8, 24, 18, 22, 15, 11, -4,
		-3, 4, 20, 30, 22, 25, 15, -2,
		-7, 1, 3, 19, 10, -2, -4, -4,
		-10, -2, -1, 0, 6, -8, -3, -13,
		-12, -28, -8, 1, -5, -12, -27, -12,
	},
	chess.Bishop: {
		-9, -5, -9, -5, -2, -4, -5, -8,
		0, 2, 8, -7, 1, 0, -2, -8,
		8, 0, 0, 1, 0, 1, 5, 6,
		0, 7, 7, 8, 3, 5, 2, 6,
		-1, 0, 12, 8, 0, 6, 0, -5,
		0, 0, 3, 6, 8, -1, 0, -1,
		-6, -12, -7, 0, 0, -8, -9, -13,
		-11, 0, -6, 0, -3, -4, -5, -9,
	},
	chess.Rook: {
		8, 9, 11, 13, 13, 12, 13, 9,
		3, 5, 1, 0, -1, 0, 6, 2,
		9, 5, 7, 2, 2, 1, 0, 0,
		3, 3, 6, 0, 0, 0, 0, 4,
		5, 4, 9, 0, -3, -2, -6, -2,
		0, 0, -6, -5, -9, -14, -7, -12,
		-2, -5, -1, -7, -9, -11, -13, -1,
		-7, -3, 0, -8, -13, -12, -4, -24,
	},
	chess.Queen: {
		-12, 4, 8, 4, 10, 9, 3, 6,
		-17, -7, -1, 7, 3, 6, 1, 0,
		-5, -1, -4, 12, 14, 20, 12, 14,
		-2, 2, 2, 9, 13, 7, 18, 22,
		-9, 3, 1, 15, 5, 10, 12, 10,
		-6, -20, 0, -15, 0, -1, 10, 7,
		-6, -14, -31, -27, -19, -12, -11, -4,
		-12, -22, -19, -30, -8, -13, -6, -15,
	},
	chess.King: {
		-15, -11, -11, -6, -2, 3, 4, -9,
		-9, 14, 11, 13, 13, 28, 19, 1,
		-1, 18, 19, 15, 16, 35, 34, 4,
		-12, 14, 21, 25, 19, 25, 18, -5,
		-23, -6, 14, 21, 20, 18, 5, -16,
		-21, -6, 5, 13, 15, 9, -2, -12,
		-27, -10, 2, 9, 9, 1, -12, -26,
		-43, -34, -20, -5, -26, -9, -35, -55,
	},
}

// -----------------------------------------------------------------------------
// 		Bonuses + Penalties
// -----------------------------------------------------------------------------

var BishopPairBonusMG int = 22
var BishopPairBonusEG int = 30

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

	var OCC_WHITE map[chess.PieceType]int = map[chess.PieceType]int{
		chess.King:   0,
		chess.Queen:  0,
		chess.Rook:   0,
		chess.Bishop: 0,
		chess.Knight: 0,
		chess.Pawn:   0,
	}
	var OCC_BLACK map[chess.PieceType]int = map[chess.PieceType]int{
		chess.King:   0,
		chess.Queen:  0,
		chess.Rook:   0,
		chess.Bishop: 0,
		chess.Knight: 0,
		chess.Pawn:   0,
	}

	var delta_mg int = 0
	var delta_eg int = 0

	squares := position.Board().SquareMap()
	for square, piece := range squares {
		if piece.Color() == chess.White {
			OCC_WHITE[piece.Type()]++
			delta_mg += PVM[piece.Type()] + PST_MG[piece.Type()][FLIP[square]]
			delta_eg += PVM_EG[piece.Type()] + PST_EG[piece.Type()][FLIP[square]]
		} else {
			OCC_BLACK[piece.Type()]++
			delta_mg -= PVM[piece.Type()] + PST_MG[piece.Type()][square]
			delta_eg -= PVM_EG[piece.Type()] + PST_EG[piece.Type()][square]
		}
	}

	if OCC_WHITE[chess.Bishop] == 2 {
		delta_mg += BishopPairBonusMG
		delta_eg += BishopPairBonusMG
	}
	if OCC_BLACK[chess.Bishop] == 2 {
		delta_mg -= BishopPairBonusMG
		delta_eg -= BishopPairBonusMG
	}

	// Tapered Evaluation
	var phase int = TotalPhase
	phase -= (OCC_WHITE[chess.Pawn] + OCC_BLACK[chess.Pawn]) * PawnPhase
	phase -= (OCC_WHITE[chess.Knight] + OCC_BLACK[chess.Knight]) * KnightPhase
	phase -= (OCC_WHITE[chess.Bishop] + OCC_BLACK[chess.Bishop]) * BishopPhase
	phase -= (OCC_WHITE[chess.Rook] + OCC_BLACK[chess.Rook]) * RookPhase
	phase -= (OCC_WHITE[chess.Queen] + OCC_BLACK[chess.Queen]) * QueenPhase
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
