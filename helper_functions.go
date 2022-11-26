package main

import (
	"fmt"

	"github.com/Sidhant-Roymoulik/chess"
)

func print(str ...any) {
	fmt.Println(str...)
}

func resetCounters() {
	states = 0
	q_states = 0
	hashes = 0
}

func game_from_fen(str string) *chess.Game {
	fen, err := chess.FEN(str)
	if err != nil {
		panic(err)
	}
	return chess.NewGame(fen)
}

func getMultiplier(turn bool) int {
	if turn {
		return 1
	} else {
		return -1
	}
}
