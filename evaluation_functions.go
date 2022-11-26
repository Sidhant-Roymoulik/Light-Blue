package main

import "github.com/Sidhant-Roymoulik/chess"

// piece weights
const pawn int = 100
const knight int = 320
const bishop int = 330
const rook int = 500
const queen int = 900
const king int = 20000

// piece map
var piece_map map[chess.PieceType]int = map[chess.PieceType]int{
	1: king,
	2: queen,
	3: rook,
	4: bishop,
	5: knight,
	6: pawn,
}

// Sums value of pieces and checks for game-over states

func eval_v1(position *chess.Position) int {
	squares := position.Board().SquareMap()
	var delta int = 0
	for _, piece := range squares {
		var turn int = getMultiplier(piece.Color() == chess.White)
		delta += piece_map[piece.Type()] * turn
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
