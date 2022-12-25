package main

import (
	"math"

	"github.com/Sidhant-Roymoulik/chess"
)

type e_minimax_mo_ab struct {
	EngineClass
}

func new_engine_minimax_mo_ab() e_minimax_mo_ab {
	return e_minimax_mo_ab{
		EngineClass{
			name: "Minimax with Move Ordering and Alpha-Beta Pruning",
			upgrades: EngineUpgrades{
				move_ordering:       true,
				alphabeta:           true,
				iterative_deepening: false,
				q_search:            false,
				concurrent:          false,
			},
		},
	}
}

func (engine *e_minimax_mo_ab) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()
	best_eval, best_move = engine.minimax_start(position, 0, position.Turn() == chess.White)
	return
}
func (engine *e_minimax_mo_ab) minimax_start(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
	moves := position.ValidMoves()

	best_eval = math.MaxInt * -1
	best_move = moves[0]
	for _, move := range moves {
		new_eval := engine.minimax(position.Update(move), ply+1, !turn, math.MaxInt*-1, math.MaxInt) * -1
		// print("Top Level Move:", move, "Eval:", new_eval)
		if new_eval > best_eval {
			best_eval = new_eval
			best_move = move
		}
	}
	return best_eval, best_move
}
func (engine *e_minimax_mo_ab) minimax(position *chess.Position, ply int, turn bool, alpha int, beta int) (best_eval int) {
	states++

	if ply > MAX_CONST_DEPTH || len(position.ValidMoves()) == 0 {
		return eval_v3(position, ply) * getMultiplier(turn)
	}

	moves := move_ordering_v2(position)

	best_eval = math.MaxInt * -1
	for _, move := range moves {
		new_eval := engine.minimax(position.Update(move), ply+1, !turn, -beta, -alpha) * -1

		if new_eval > best_eval {
			best_eval = new_eval
		}

		if best_eval >= beta {
			return beta
		}
		if best_eval >= alpha {
			alpha = best_eval
		}
	}
	return alpha
}
