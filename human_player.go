package main

import (
	"fmt"

	"github.com/Sidhant-Roymoulik/chess"
)

type e_human struct {
	EngineClass
}

func new_engine_human() e_human {
	return e_human{
		EngineClass{
			name: "Human",
			// upgrades: EngineUpgrades{
			// 	move_ordering:       false,
			// 	alphabeta:           false,
			// 	iterative_deepening: false,
			// 	q_search:            false,
			// 	concurrent:          false,
			// },
		},
	}
}

func (engine *e_human) run(position *chess.Position) (best_eval int, move *chess.Move) {
	best_eval = 0

	print("Enter move: ")
	var input string
	fmt.Scanln(&input)
	move, err := chess.AlgebraicNotation{}.Decode(position, input)
	if err != nil {
		print("Invalid move.")
		print("Did you mean?", valid_move_strings(position))
		return engine.run(position)
	}
	return
}

func valid_move_strings(position *chess.Position) []string {
	moves := position.ValidMoves()
	move_strings := make([]string, len(moves))
	for i, move := range moves {
		move_strings[i] = chess.AlgebraicNotation{}.Encode(position, move)
	}
	return move_strings
}
