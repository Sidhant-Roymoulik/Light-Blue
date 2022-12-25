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
			name:     "Minimax",
			upgrades: EngineUpgrades{},
		},
	}
}

func (engine *e_minimax) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()
	best_eval, best_move = engine.minimax_start(position, 0, position.Turn() == chess.White)
	return
}
func (engine *e_minimax) minimax_start(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
	moves := position.ValidMoves()

	best_eval = math.MaxInt * -1
	for _, move := range moves {
		new_eval := engine.minimax(position.Update(move), ply+1, !turn) * -1

		// print("Top Level Move:", move, "Eval:", new_eval)

		if new_eval > best_eval {
			best_eval = new_eval
			best_move = move
		}
	}
	return best_eval, best_move
}
func (engine *e_minimax) minimax(position *chess.Position, ply int, turn bool) (best_eval int) {
	states++

	if ply > MAX_CONST_DEPTH {
		return eval_v2(position) * getMultiplier(turn)
	}

	moves := position.ValidMoves()

	best_eval = math.MaxInt * -1
	for _, move := range moves {
		new_eval := engine.minimax(position.Update(move), ply+1, !turn) * -1

		if new_eval > best_eval {
			best_eval = new_eval
		}
	}
	return best_eval
}
