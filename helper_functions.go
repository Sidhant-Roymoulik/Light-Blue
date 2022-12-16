package main

import (
	"fmt"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

func print(str ...any) {
	fmt.Println(str...)
}

func resetCounters() {
	start = time.Now()
	states = 0
	q_states = 0
	hashes = 0
}

func game_from_fen(str string) *chess.Game {
	fen, err := chess.FEN(str)
	if err != nil {
		panic(err)
	}
	return chess.NewGame(fen, chess.UseNotation(chess.AlgebraicNotation{}))
}

func getMultiplier(turn bool) int {
	if turn {
		return 1
	} else {
		return -1
	}
}

func isQMove(move *chess.Move) bool {
	if move.HasTag(chess.MoveTag(chess.Checkmate)) {
		return true
	}
	if move.HasTag(chess.MoveTag(chess.Capture)) {
		return true
	}
	if move.HasTag(chess.MoveTag(chess.Check)) {
		return true
	}
	return false
}

func getQMoves(position *chess.Position) []*chess.Move {
	moves := move_ordering_v1(position)
	n := 0
	for _, move := range moves {
		if isQMove(move) {
			moves[n] = move
			n++
		}
	}
	return moves[:n]
}
