package main

import (
	"github.com/Sidhant-Roymoulik/chess"
)

// -----------------------------------------------------------------------------
// 		Piece Values
// -----------------------------------------------------------------------------

// piece value map
var PVM map[chess.PieceType]int = map[chess.PieceType]int{
	chess.King:   20000,
	chess.Queen:  921,
	chess.Rook:   441,
	chess.Bishop: 346,
	chess.Knight: 333,
	chess.Pawn:   84,
}

var PVM_EG map[chess.PieceType]int = map[chess.PieceType]int{
	chess.King:   20000,
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

var RookOrQueenOnSeventhBonusMG int = 30
var RookOrQueenOnSeventhBonusEG int = 23

var RookOnOpenFileBonusMG int = 23

var RookOnSemiOpenFileBonusMG int = 10

var IsolatedPawnPenatlyMG int = 17
var IsolatedPawnPenatlyEG int = 6

var DoubledPawnPenatlyMG int = 1
var DoubledPawnPenatlyEG int = 16

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

// Best Evaluation
func eval_pos(position *chess.Position, ply int) int {
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

	var OCC map[chess.Color]map[chess.PieceType]int = map[chess.Color]map[chess.PieceType]int{
		chess.White: {
			chess.King:   0,
			chess.Queen:  0,
			chess.Rook:   0,
			chess.Bishop: 0,
			chess.Knight: 0,
			chess.Pawn:   0,
		},
		chess.Black: {
			chess.King:   0,
			chess.Queen:  0,
			chess.Rook:   0,
			chess.Bishop: 0,
			chess.Knight: 0,
			chess.Pawn:   0,
		},
	}

	var P_FILE map[chess.Color]map[chess.File]int = map[chess.Color]map[chess.File]int{
		chess.White: {
			chess.FileA: 0,
			chess.FileB: 0,
			chess.FileC: 0,
			chess.FileD: 0,
			chess.FileE: 0,
			chess.FileF: 0,
			chess.FileG: 0,
			chess.FileH: 0,
		},
		chess.Black: {
			chess.FileA: 0,
			chess.FileB: 0,
			chess.FileC: 0,
			chess.FileD: 0,
			chess.FileE: 0,
			chess.FileF: 0,
			chess.FileG: 0,
			chess.FileH: 0,
		},
	}

	squares := position.Board().SquareMap()
	for square, piece := range squares {
		OCC[piece.Color()][piece.Type()]++

		if piece.Type() == chess.Pawn {
			P_FILE[piece.Color()][square.File()]++
		}
	}

	var white_knights int = OCC[chess.White][chess.Knight]
	var white_bishops int = OCC[chess.White][chess.Bishop]
	var black_knights int = OCC[chess.Black][chess.Knight]
	var black_bishops int = OCC[chess.Black][chess.Bishop]

	var pawns int = OCC[chess.White][chess.Pawn] + OCC[chess.Black][chess.Pawn]
	var knights int = white_knights + black_knights
	var bishops int = white_bishops + black_bishops
	var rooks int = OCC[chess.White][chess.Rook] + OCC[chess.Black][chess.Rook]
	var queens int = OCC[chess.White][chess.Queen] + OCC[chess.Black][chess.Queen]

	var majors int = queens + rooks
	var minors int = bishops + knights

	// Draw by Insufficient Material
	if pawns+majors+minors == 0 {
		return 0
	} else if majors+pawns == 0 {
		if minors == 1 {
			return 0
		} else if minors == 2 {
			if white_knights == 1 && black_knights == 1 {
				return 0
			} else if white_bishops == 1 && black_bishops == 1 {
				return 0
			}
		}
	}

	var turn chess.Color = position.Turn()
	var other_turn chess.Color = turn.Other()

	var score_mg map[chess.Color]int = map[chess.Color]int{
		chess.White: 0,
		chess.Black: 0,
	}
	var score_eg map[chess.Color]int = map[chess.Color]int{
		chess.White: 0,
		chess.Black: 0,
	}

	if OCC[chess.White][chess.Bishop] == 2 {
		score_mg[chess.White] += BishopPairBonusMG
		score_eg[chess.White] += BishopPairBonusEG
	}
	if OCC[chess.Black][chess.Bishop] == 2 {
		score_mg[chess.Black] += BishopPairBonusMG
		score_eg[chess.Black] += BishopPairBonusEG
	}

	valid_moves := position.ValidMoves()
	other_valid_moves := position.NullMove().ValidMoves()

	score_mg[turn] += len(valid_moves)
	score_eg[turn] += len(valid_moves)
	score_mg[other_turn] += len(other_valid_moves)
	score_eg[other_turn] += len(other_valid_moves)

	for square, piece := range squares {
		var square_file chess.File = square.File()
		var piece_color chess.Color = piece.Color()
		var piece_type chess.PieceType = piece.Type()

		if piece_color == chess.White {
			score_mg[chess.White] += PVM[piece_type] +
				PST_MG[piece_type][FLIP[square]]
			score_eg[chess.White] += PVM_EG[piece_type] +
				PST_EG[piece_type][FLIP[square]]

			if square.Rank() == chess.Rank7 &&
				(piece_type == chess.Queen || piece_type == chess.Rook) {
				score_mg[chess.White] += RookOrQueenOnSeventhBonusMG
				score_eg[chess.White] += RookOrQueenOnSeventhBonusEG
			}
		} else {
			score_mg[chess.Black] += PVM[piece_type] +
				PST_MG[piece_type][square]
			score_eg[chess.Black] += PVM_EG[piece_type] +
				PST_EG[piece_type][square]

			if square.Rank() == chess.Rank2 &&
				(piece_type == chess.Queen || piece_type == chess.Rook) {
				score_mg[chess.Black] += RookOrQueenOnSeventhBonusMG
				score_eg[chess.Black] += RookOrQueenOnSeventhBonusEG
			}
		}

		if piece_type == chess.Pawn {
			if P_FILE[piece_color][square_file] > 1 {
				score_mg[piece_color] -= DoubledPawnPenatlyMG
				score_eg[piece_color] -= DoubledPawnPenatlyEG
			}

			var isIsolated bool = true
			if square_file > chess.FileA &&
				P_FILE[piece_color][square_file-1] > 0 {
				isIsolated = false
			}
			if square_file < chess.FileH &&
				P_FILE[piece_color][square_file+1] > 0 {
				isIsolated = false
			}
			if isIsolated {
				score_mg[piece_color] -= IsolatedPawnPenatlyMG
				score_eg[piece_color] -= IsolatedPawnPenatlyEG
			}
		}

		if piece_type == chess.Rook {
			if P_FILE[piece_color][square_file] == 0 {
				if P_FILE[piece_color.Other()][square_file] == 0 {
					score_mg[piece_color] += RookOnOpenFileBonusMG
				} else {
					score_mg[piece_color] += RookOnSemiOpenFileBonusMG
				}
			}

		}
	}

	// Tapered Evaluation
	var delta_mg int = score_mg[turn] - score_mg[other_turn]
	var delta_eg int = score_eg[turn] - score_eg[other_turn]

	var phase int = TotalPhase
	phase -= pawns * PawnPhase
	phase -= knights * KnightPhase
	phase -= bishops * BishopPhase
	phase -= rooks * RookPhase
	phase -= queens * QueenPhase
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase

	var delta int = ((delta_mg * (256 - phase)) + (delta_eg * phase)) / 256

	return delta
}
