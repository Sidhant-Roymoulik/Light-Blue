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
				move_ordering:       true,
				alphabeta:           true,
				iterative_deepening: false,
				q_search:            true,
				concurrent:          false,
			},
		},
	}
}

func (engine *e_minimax_mo_ab_q) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()
	best_eval, best_move = minimax_start_mo_ab_q(position, 0, position.Turn() == chess.White)
	return
}
func minimax_start_mo_ab_q(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
	best_eval = math.MaxInt * -1
	moves := move_ordering_v1(position)
	for _, move := range moves {
		new_eval := minimax_mo_ab_q(position.Update(move), ply+1, !turn, math.MaxInt*-1, math.MaxInt) * -1
		if new_eval > best_eval {
			best_eval = new_eval
			best_move = move
		}
	}
	return best_eval, best_move
}
func minimax_mo_ab_q(position *chess.Position, ply int, turn bool, alpha int, beta int) (best_eval int) {
	if ply > MAX_CONST_DEPTH {
		return q_search(position, ply, turn, alpha, beta)
	}
	states++
	best_eval = math.MaxInt * -1
	moves := move_ordering_v1(position)
	for _, move := range moves {
		new_eval := minimax_mo_ab_q(position.Update(move), ply+1, !turn, -beta, -alpha) * -1
		if new_eval >= best_eval {
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

func q_search(position *chess.Position, ply int, turn bool, alpha int, beta int) (best_eval int) {
	if ply > MAX_CONST_DEPTH*2 {
		return eval_v2(position) * getMultiplier(turn)
	}
	q_states++

	best_eval = eval_v2(position)

	if best_eval >= beta {
		return beta
	}
	if best_eval >= alpha {
		alpha = best_eval
	}

	qmoves := getQMoves(position)
	// print(qmoves)
	for _, qmove := range qmoves {
		new_eval := q_search(position.Update(qmove), ply+1, !turn, -beta, -alpha) * -1
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
