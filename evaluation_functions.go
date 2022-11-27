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

func eval_move_v1(move *chess.Move, board *chess.Board) int {
	if move.HasTag(chess.MoveTag(chess.Checkmate)) {
		return CHECKMATE_VALUE
	}

	delta := 0

	if move.HasTag(chess.MoveTag(chess.Capture)) {
		delta += PVM[chess.PieceType(board.Piece(move.S2()))]
		delta -= PVM[chess.PieceType(board.Piece(move.S1()))]
	}
	if move.HasTag(chess.MoveTag(chess.Check)) {
		delta += 25
	}
	if move.HasTag(chess.MoveTag(chess.KingSideCastle)) {
		delta += 50
	}
	if move.HasTag(chess.MoveTag(chess.QueenSideCastle)) {
		delta += 40
	}

	return delta
}

// Sums value of pieces and checks for game-over states

func eval_v1(position *chess.Position) int {
	squares := position.Board().SquareMap()
	var delta int = 0
	for _, piece := range squares {
		var turn int = getMultiplier(position)
		delta += PVM[piece.Type()] * turn
	}

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

	return delta
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
		100, 100, 100, 100, 105, 100, 100, 100,
		78, 83, 86, 73, 102, 82, 85, 90,
		7, 29, 21, 44, 40, 31, 44, 7,
		-17, 16, -2, 15, 14, 0, 15, -13,
		-26, 3, 10, 9, 6, 1, 0, -23,
		-22, 9, 5, -11, -10, -2, 3, -19,
		-31, 8, -7, -37, -36, -14, 3, -31,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-66, -53, -75, -75, -10, -55, -58, -70,
		-3, -6, 100, -36, 4, 62, -4, -14,
		10, 67, 1, 74, 73, 27, 62, -2,
		24, 24, 45, 37, 33, 41, 25, 17,
		-1, 5, 31, 21, 22, 35, 2, 0,
		-18, 10, 13, 22, 18, 15, 11, -14,
		-23, -15, 2, 0, 2, 0, -23, -20,
		-74, -23, -26, -24, -19, -35, -22, -69,
	},
	chess.Bishop: {
		-59, -78, -82, -76, -23, -107, -37, -50,
		-11, 20, 35, -42, -39, 31, 2, -22,
		-9, 39, -32, 41, 52, -10, 28, -14,
		25, 17, 20, 34, 26, 25, 15, 10,
		13, 10, 17, 23, 17, 16, 0, 7,
		14, 25, 24, 15, 8, 25, 20, 15,
		19, 20, 11, 6, 7, 6, 20, 16,
		-7, 2, -15, -12, -14, -15, -10, -10,
	},
	chess.Rook: {
		35, 29, 33, 4, 37, 33, 56, 50,
		55, 29, 56, 67, 55, 62, 34, 60,
		19, 35, 28, 33, 45, 27, 25, 15,
		0, 5, 16, 13, 18, -4, -9, -6,
		-28, -35, -16, -21, -13, -29, -46, -30,
		-42, -28, -42, -25, -25, -35, -26, -46,
		-53, -38, -31, -26, -29, -43, -44, -53,
		-30, -24, -18, 5, -2, -18, -31, -32,
	},
	chess.Queen: {
		6, 1, -8, -104, 69, 24, 88, 26,
		14, 32, 60, -10, 20, 76, 57, 24,
		-2, 43, 32, 60, 72, 63, 43, 2,
		1, -16, 22, 17, 25, 20, -13, -6,
		-14, -15, -2, -5, -1, -10, -20, -22,
		-30, -6, -13, -11, -16, -11, -16, -27,
		-36, -18, 0, -19, -15, -15, -21, -38,
		-39, -30, -31, -13, -31, -36, -34, -42,
	},
	chess.King: {
		4, 54, 47, -99, -99, 60, 83, -62,
		-32, 10, 55, 56, 56, 55, 10, 3,
		-62, 12, -57, 44, -67, 28, 37, -31,
		-55, 50, 11, -4, -19, 13, 0, -49,
		-55, -43, -52, -28, -51, -47, -8, -50,
		-47, -42, -43, -79, -64, -32, -29, -32,
		-4, 3, -14, -50, -57, -18, 13, 4,
		17, 30, -3, -14, 6, -1, 40, 18,
	},
}

// Also uses PST
func eval_v2(position *chess.Position) int {
	squares := position.Board().SquareMap()
	var delta int = 0
	for square, piece := range squares {
		if piece.Color() == chess.White {
			delta += PVM[piece.Type()]
			delta += PST_MG[piece.Type()][square]
		} else {
			delta -= PVM[piece.Type()]
			delta -= PST_MG[piece.Type()][FLIP[square]]
		}
	}

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

	return delta
}
