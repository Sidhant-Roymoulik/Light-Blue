package main

import (
	"math"

	"github.com/Sidhant-Roymoulik/chess"
)

type e_minimax_mo_ab_q struct {
	EngineClass
}

func new_engine_minimax_mo_ab_q() e_minimax_mo_ab_q {
	return e_minimax_mo_ab_q{
		EngineClass{
			name: "Minimax with Move Ordering, Alpha-Beta Pruning, and Quiesence Search",
			upgrades: EngineUpgrades{
				move_ordering: true,
				alphabeta:     true,
				q_search:      true,
			},
		},
	}
}

func (engine *e_minimax_mo_ab_q) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()
	best_eval, best_move = engine.minimax_start(position, 0, position.Turn() == chess.White)
	return
}
func (engine *e_minimax_mo_ab_q) minimax_start(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
	moves := move_ordering_v2(position)

	best_eval = math.MaxInt * -1
	best_move = moves[0]
	for _, move := range moves {
		new_eval := engine.minimax(position.Update(move), ply+1, !turn, math.MaxInt*-1, -best_eval) * -1
		// print("Top Level Move:", move, "Eval:", new_eval)
		if new_eval > best_eval {
			best_eval = new_eval
			best_move = move
		}
	}
	return best_eval, best_move
}
func (engine *e_minimax_mo_ab_q) minimax(position *chess.Position, ply int, turn bool, alpha int, beta int) (best_eval int) {
	states++

	if ply > MAX_CONST_DEPTH {
		return engine.q_search(position, ply, turn, alpha, beta)
	}
	if len(position.ValidMoves()) == 0 {
		return eval_v4(position, ply) * getMultiplier(turn)
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

func (engine *e_minimax_mo_ab_q) q_search(position *chess.Position, ply int, turn bool, alpha int, beta int) (best_eval int) {
	q_states++

	start_eval := eval_v4(position, ply) * getMultiplier(turn)

	if start_eval >= beta {
		return beta
	}
	if start_eval >= alpha {
		alpha = start_eval
	}

	if ply > MAX_CONST_DEPTH*2 {
		return start_eval
	}

	moves := get_q_moves(position)

	if len(moves) == 0 {
		return start_eval
	}

	best_eval = math.MaxInt * -1
	for _, move := range moves {
		new_eval := engine.q_search(position.Update(move), ply+1, !turn, -beta, -alpha) * -1

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
