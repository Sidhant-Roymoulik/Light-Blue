package main

import (
	"sort"

	"github.com/Sidhant-Roymoulik/chess"
)

func move_ordering_v1(position *chess.Position) []*chess.Move {
	moves := position.ValidMoves()
	sort.SliceStable(moves, func(i, j int) bool {
		return eval_move_v1(moves[i], position.Board()) < eval_move_v1(moves[j], position.Board())
	})
	return moves
}
