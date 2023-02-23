package main

import (
	"fmt"

	"github.com/Sidhant-Roymoulik/chess"
)

func print(str ...any) {
	fmt.Println(str...)
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

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func getMateOrCPScore(score int) string {
	if score > MATE_CUTOFF {
		pliesToMate := CHECKMATE_VALUE - score
		mateInN := (pliesToMate / 2) + (pliesToMate % 2)
		return fmt.Sprintf("mate %d", mateInN)
	}

	if score < -MATE_CUTOFF {
		pliesToMate := -CHECKMATE_VALUE - score
		mateInN := (pliesToMate / 2) + (pliesToMate % 2)
		return fmt.Sprintf("mate %d", mateInN)
	}

	return fmt.Sprintf("cp %d", score)
}
