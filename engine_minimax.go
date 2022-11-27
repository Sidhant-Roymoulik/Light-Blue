package main

import (
	"math"

	"github.com/Sidhant-Roymoulik/chess"
)

type e_minimax struct {
	EngineClass
}

func new_engine_minimax() e_minimax {
	return e_minimax{
		EngineClass{
			name: "Minimax",
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

func (engine *e_minimax) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()
	best_eval, best_move = minimax_start(position, 0, position.Turn() == chess.White)
	return
}
func minimax_start(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
	best_eval = math.MaxInt * -1
	moves := position.ValidMoves()
	for _, move := range moves {
		new_eval := minimax(position.Update(move), ply+1, !turn) * -1
		if new_eval > best_eval {
			best_eval = new_eval
			best_move = move
		}
	}
	return best_eval, best_move
}
func minimax(position *chess.Position, ply int, turn bool) (best_eval int) {
	if ply > MAX_CONST_DEPTH {
		return eval_v1(position) * getMultiplier(position)
	}
	states++
	best_eval = math.MaxInt * -1
	moves := position.ValidMoves()
	for _, move := range moves {
		new_eval := minimax(position.Update(move), ply+1, !turn) * -1
		if new_eval > best_eval {
			best_eval = new_eval
		}
	}
	return best_eval
}
