package main

import (
	"math"

	"github.com/Sidhant-Roymoulik/chess"
)

type e_minimax_ab struct {
	EngineClass
}

func new_engine_minimax_ab() e_minimax_ab {
	return e_minimax_ab{
		EngineClass{
			name: "Minimax with Alpha-Beta Pruning",
			upgrades: EngineUpgrades{
				alphabeta: true,
			},
		},
	}
}

func (engine *e_minimax_ab) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()
	best_eval, best_move = minimax_start_ab(position, 0, position.Turn() == chess.White)
	return
}
func minimax_start_ab(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
	best_eval = math.MaxInt * -1
	moves := position.ValidMoves()
	for _, move := range moves {
		new_eval := minimax_ab(position.Update(move), ply+1, !turn, math.MaxInt*-1, math.MaxInt) * -1
		if new_eval > best_eval {
			best_eval = new_eval
			best_move = move
		}
	}
	return best_eval, best_move
}
func minimax_ab(position *chess.Position, ply int, turn bool, alpha int, beta int) (best_eval int) {
	if ply > MAX_CONST_DEPTH {
		return eval_v2(position) * getMultiplier(turn)
	}
	states++
	best_eval = math.MaxInt * -1
	moves := position.ValidMoves()
	for _, move := range moves {
		new_eval := minimax_ab(position.Update(move), ply+1, !turn, -beta, -alpha) * -1
		if new_eval > best_eval {
			best_eval = new_eval
		}
		if beta <= best_eval {
			return beta
		}
		if alpha <= best_eval {
			alpha = best_eval
		}
	}
	return best_eval
}
