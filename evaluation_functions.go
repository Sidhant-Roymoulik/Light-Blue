package main

import "github.com/Sidhant-Roymoulik/chess"

// piece weights
const pawn int = 100
const knight int = 320
const bishop int = 330
const rook int = 500
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
		0, 0, 0, 0, 0, 0, 0, 0,
		15, 15, 15, 20, 20, 15, 15, 15,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 10, 10, 10, 0, 0,
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
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 20, 20, 0, 0, 0,
		0, 0, 0, 20, 20, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, -10, -10, 0, 0, 0,
		0, 0, 20, -10, -10, 0, 20, 0,
	},
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
