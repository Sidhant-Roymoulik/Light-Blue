package main

import (
	"github.com/Sidhant-Roymoulik/chess"
)

// -----------------------------------------------------------------------------
// 		Bonuses + Penalties
// -----------------------------------------------------------------------------

const (
	BishopPairBonusMG int = 22
	BishopPairBonusEG int = 30

	RookOrQueenOnSeventhBonusEG int = 23

	RookOnOpenFileBonusMG int = 23

	RookOnSemiOpenFileBonusMG int = 10

	IsolatedPawnPenatlyMG int = 17
	IsolatedPawnPenatlyEG int = 6

	DoubledPawnPenatlyMG int = 1
	DoubledPawnPenatlyEG int = 16

	// -------------------------------------------------------------------------
	// 		Tapered Evaluation Values
	// -------------------------------------------------------------------------

	PawnPhase   int = 0
	KnightPhase int = 1
	BishopPhase int = 1
	RookPhase   int = 2
	QueenPhase  int = 4
	TotalPhase  int = PawnPhase*16 +
		KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2
)

var SEVENTH_RANK = map[chess.Color]chess.Rank{
	chess.White: chess.Rank7,
	chess.Black: chess.Rank2,
}

// -----------------------------------------------------------------------------
// 		Piece Values
// -----------------------------------------------------------------------------

// piece value map
var PVM_MG = map[chess.PieceType]int{
	chess.King:   20000,
	chess.Queen:  921,
	chess.Rook:   441,
	chess.Bishop: 346,
	chess.Knight: 333,
	chess.Pawn:   84,
}

var PVM_EG = map[chess.PieceType]int{
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

var FLIP = map[chess.Color][]int{
	chess.White: {
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	},
	chess.Black: {
		0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 58, 59, 60, 61, 62, 63,
	},
}

var PST_MG = map[chess.PieceType][]int{
	chess.Pawn: {
		0, 0, 0, 0, 0, 0, 0, 0,
		49, 49, 50, 50, 51, 49, 49, 49,
		11, 11, 19, 31, 29, 21, 10, 9,
		5, 4, 9, 25, 25, 11, 4, 6,
		1, 0, 1, 19, 19, -1, -1, -1,
		4, -6, -10, 0, 0, -10, -4, 4,
		5, 11, 11, -19, -19, 11, 11, 4,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-51, -40, -30, -29, -29, -30, -40, -50,
		-40, -19, 0, -1, 0, -1, -20, -40,
		-29, 1, 11, 14, 15, 9, 0, -29,
		-30, 6, 16, 19, 19, 15, 6, -31,
		-31, 0, 15, 19, 19, 16, -1, -29,
		-29, 4, 9, 16, 14, 11, 6, -29,
		-39, -19, 1, 4, 6, 0, -19, -39,
		-49, -39, -30, -29, -31, -29, -39, -51,
	},
	chess.Bishop: {
		-19, -9, -9, -9, -10, -9, -11, -20,
		-10, 1, 1, 1, -1, -1, 1, -9,
		-11, 1, 4, 11, 9, 5, -1, -10,
		-9, 4, 4, 11, 9, 4, 6, -11,
		-9, 1, 11, 11, 9, 9, 1, -11,
		-9, 9, 11, 9, 10, 11, 11, -10,
		-11, 4, 0, 1, 1, 1, 6, -9,
		-19, -9, -11, -9, -11, -10, -10, -19,
	},
	chess.Rook: {
		0, 0, 1, -1, 1, 1, -1, -1,
		6, 9, 11, 10, 11, 11, 11, 5,
		-4, 0, 1, 1, -1, 1, -1, -4,
		-6, 0, 1, -1, 1, 1, 1, -4,
		-4, 1, 1, -1, 1, -1, -1, -6,
		-5, 1, 0, -1, 1, 0, -1, -5,
		-6, 0, -1, 0, 1, -1, -1, -6,
		0, 1, 1, 4, 6, 1, -1, 0,
	},
	chess.Queen: {
		-21, -10, -9, -6, -5, -10, -9, -19,
		-11, 1, 1, -1, 1, -1, 0, -11,
		-11, -1, 4, 5, 6, 5, 1, -9,
		-4, 0, 5, 4, 5, 6, 0, -4,
		-1, 1, 6, 4, 4, 6, -1, -4,
		-11, 4, 4, 6, 6, 6, 1, -11,
		-9, 1, 6, -1, 1, -1, 1, -9,
		-19, -9, -11, -4, -6, -11, -10, -19,
	},
	chess.King: {
		-30, -40, -40, -50, -50, -40, -39, -30,
		-30, -39, -39, -50, -51, -39, -40, -30,
		-30, -39, -39, -50, -49, -39, -41, -31,
		-31, -41, -41, -49, -49, -40, -40, -30,
		-20, -31, -30, -40, -39, -29, -30, -20,
		-10, -21, -20, -19, -20, -20, -20, -10,
		19, 21, 0, 0, 1, 1, 19, 19,
		19, 31, 9, 1, 0, 11, 30, 19,
	},
}

var PST_EG = map[chess.PieceType][]int{
	chess.Pawn: {
		0, 0, 0, 0, 0, 0, 0, 0,
		-1, 50, 49, 49, 50, 50, 49, 1,
		-1, 40, 41, 40, 40, 39, 39, 1,
		-1, -1, -1, -1, 1, -1, 0, 1,
		0, -1, 1, -1, 0, -1, -1, 0,
		-1, 1, -1, 1, 0, 1, 0, 0,
		0, 0, 1, 1, -1, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-51, -39, -29, -29, -30, -31, -41, -50,
		-39, -19, 1, -1, 1, -1, -19, -41,
		-30, 0, 11, 14, 15, 10, -1, -30,
		-30, 4, 15, 21, 20, 14, 6, -29,
		-31, -1, 16, 19, 19, 16, 0, -29,
		-29, 4, 9, 16, 14, 11, 4, -29,
		-40, -20, 0, 4, 5, -1, -21, -40,
		-50, -40, -31, -31, -31, -30, -40, -50,
	},
	chess.Bishop: {
		-19, -9, -9, -10, -11, -9, -11, -19,
		-11, 0, -1, -1, 0, -1, 1, -9,
		-10, 0, 6, 10, 9, 6, 1, -9,
		-10, 4, 5, 9, 9, 6, 6, -11,
		-10, 1, 10, 11, 9, 10, 1, -9,
		-9, 9, 11, 10, 9, 9, 10, -11,
		-10, 5, 0, 1, 0, 1, 4, -10,
		-20, -9, -11, -9, -10, -9, -11, -21,
	},
	chess.Rook: {
		-1, 1, 1, -1, -1, 0, 1, 0,
		-9, 0, 1, 0, 1, -1, 1, -11,
		-10, 0, 0, 0, -1, 1, 1, -9,
		-10, 1, 0, -1, 0, -1, 1, -9,
		-11, 0, 0, -1, 1, 0, -1, -9,
		-9, 0, 0, -1, 1, 0, -1, -10,
		-11, 0, 1, 0, 0, -1, -1, -10,
		1, 1, 0, 0, 0, 0, -1, 0,
	},
	chess.Queen: {
		-21, -11, -9, -4, -5, -11, -9, -19,
		-11, 1, 1, 1, 1, 1, 1, -10,
		-11, -1, 6, 4, 6, 5, 1, -10,
		-4, 1, 4, 4, 5, 5, 0, -4,
		-1, 0, 5, 5, 5, 5, 0, -4,
		-11, 4, 4, 5, 5, 6, 1, -11,
		-10, 1, 6, 0, 0, -1, 1, -10,
		-19, -9, -11, -5, -5, -9, -10, -19,
	},
	chess.King: {
		-50, -41, -30, -19, -19, -30, -39, -51,
		-30, -20, -11, 0, 1, -9, -19, -29,
		-31, -11, 21, 30, 30, 19, -11, -30,
		-31, -10, 30, 41, 39, 31, -10, -30,
		-29, -11, 30, 39, 39, 29, -9, -29,
		-30, -10, 21, 31, 31, 20, -11, -31,
		-30, -29, 1, 0, 1, 0, -30, -30,
		-51, -29, -29, -30, -30, -31, -31, -51,
	},
}

// -----------------------------------------------------------------------------
// 		Position Evaluation Function
// -----------------------------------------------------------------------------

// Best Evaluation
func eval_pos(position *chess.Position, ply int) int {
	// faster than doing two comparisons
	if position.Status() != chess.NoMethod {
		if position.Status() == chess.Checkmate {
			return -CHECKMATE_VALUE + ply
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
	if majors+minors+pawns == 0 {
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

	valid_moves := position.ValidMoves()
	other_valid_moves := position.NullMove().ValidMoves()

	score_mg[turn] += len(valid_moves)
	score_eg[turn] += len(valid_moves)
	score_mg[other_turn] += len(other_valid_moves)
	score_eg[other_turn] += len(other_valid_moves)

	if OCC[chess.White][chess.Bishop] == 2 {
		score_mg[chess.White] += BishopPairBonusMG
		score_eg[chess.White] += BishopPairBonusEG
	}
	if OCC[chess.Black][chess.Bishop] == 2 {
		score_mg[chess.Black] += BishopPairBonusMG
		score_eg[chess.Black] += BishopPairBonusEG
	}

	for square, piece := range squares {
		var square_file chess.File = square.File()
		var piece_color chess.Color = piece.Color()
		var piece_type chess.PieceType = piece.Type()

		score_mg[piece_color] += PVM_MG[piece_type]
		score_mg[piece_color] += PST_MG[piece_type][FLIP[piece_color][square]]
		score_eg[piece_color] += PVM_EG[piece_type]
		score_eg[piece_color] += PST_EG[piece_type][FLIP[piece_color][square]]

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

			if square.Rank() == SEVENTH_RANK[piece_color] {
				score_eg[piece_color] += RookOrQueenOnSeventhBonusEG
			}
		}

		if piece_type == chess.Queen {
			if square.Rank() == SEVENTH_RANK[piece_color] {
				score_eg[piece_color] += RookOrQueenOnSeventhBonusEG
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
