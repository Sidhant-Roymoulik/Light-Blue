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
	return chess.NewGame(fen)
}

func getMultiplier(position *chess.Position) int {
	if position.Turn() == chess.White {
		return 1
	} else {
		return -1
	}
}
