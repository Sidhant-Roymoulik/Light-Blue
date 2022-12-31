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
	hash_hits = 0
	hash_writes = 0
}

func game_from_fen(str string) *chess.Game {
	fen, err := chess.FEN(str)
	if err != nil {
		panic(err)
	}
	return chess.NewGame(fen, chess.UseNotation(chess.AlgebraicNotation{}))
}

func game_from_opening(opening string) *chess.Game {
	fen, err := chess.FEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		panic(err)
	}
	game := chess.NewGame(fen, chess.UseNotation(chess.AlgebraicNotation{}))
	for _, move := range CHESS_OPENINGS[opening] {
		move, err := chess.AlgebraicNotation{}.Decode(game.Position(), move)
		if err != nil {
			panic(err)
		}
		game.Move(move)
	}
	return game
}

func getMultiplier(turn bool) int {
	if turn {
		return 1
	} else {
		return -1
	}
}
