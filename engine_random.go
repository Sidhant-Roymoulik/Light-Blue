package main

import (
	"math/rand"

	"github.com/Sidhant-Roymoulik/chess"
)

type e_random struct {
	EngineClass
}

func new_engine_random() e_random {
	return e_random{
		EngineClass{
			name:     "Random",
			upgrades: EngineUpgrades{},
		},
	}
}

func (engine *e_random) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()

	moves := position.ValidMoves()

	best_eval = 0
	best_move = moves[rand.Intn(len(moves))]

	return
}
