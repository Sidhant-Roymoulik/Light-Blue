package main

import (
	"github.com/Sidhant-Roymoulik/chess"
)

type e_capture struct {
	EngineClass
}

func new_engine_capture() e_capture {
	return e_capture{
		EngineClass{
			name: "Capture",
			upgrades: EngineUpgrades{
				move_ordering:       false,
				alphabeta:           false,
				iterative_deepening: false,
				q_search:            false,
				concurrent:          false,
			},
		},
	}
}

func (engine *e_capture) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()

	moves := position.ValidMoves()
	var captures []*chess.Move
	for i := 0; i < len(moves); i++ {
		if moves[i].HasTag(chess.Capture) {
			captures = append(captures, moves[i])
		}
	}

	best_eval = 0
	if len(captures) > 0 {
		best_move = captures[0]
	} else {
		best_move = moves[0]
	}

	return
}
