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

func getMultiplier(turn bool) int {
	if turn {
		return 1
	} else {
		return -1
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
