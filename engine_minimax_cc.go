package main

import (
	"math"

	"github.com/Sidhant-Roymoulik/chess"
)

type e_minimax_cc struct {
	EngineClass
}

func new_engine_minimax_cc() e_minimax_cc {
	return e_minimax_cc{
		EngineClass{
			name: "Minimax with Concurrency",
			upgrades: EngineUpgrades{
				move_ordering:       false,
				alphabeta:           false,
				iterative_deepening: false,
				q_search:            false,
				concurrent:          true,
			},
		},
	}
}

func (engine *e_minimax_cc) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()
	best_eval, best_move = engine.minimax_start(position, 0, position.Turn() == chess.White)
	return
}
func (engine *e_minimax_cc) minimax_start(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
	moves := position.ValidMoves()
	eval_chan_local := make(chan int, len(moves))
	move_chan_local := make(chan *chess.Move, len(moves))

	for _, move := range moves {
		go engine.minimax(position.Update(move), ply+1, !turn, move, eval_chan_local, move_chan_local)
	}

	best_eval = math.MaxInt * -1
	best_move = moves[0]
	for i := 0; i < len(moves); i++ {
		new_eval := -1 * <-eval_chan_local
		new_move := <-move_chan_local
		// print("Top Level Move:", new_move, "Eval:", new_eval)
		if new_eval > best_eval {
			best_eval = new_eval
			best_move = new_move
		}
	}

	return best_eval, best_move
}
func (engine *e_minimax_cc) minimax(position *chess.Position, ply int, turn bool, prev_move *chess.Move, eval_chan chan int, move_chan chan *chess.Move) {
	states++

	if ply >= MAX_CONST_DEPTH || len(position.ValidMoves()) == 0 {
		eval_chan <- eval_v3(position, ply) * getMultiplier(turn)
		if ply == 1 {
			move_chan <- prev_move
		}
		return
	}

	moves := position.ValidMoves()
	eval_chan_local := make(chan int, len(moves))

	for _, move := range moves {
		go engine.minimax(position.Update(move), ply+1, !turn, move, eval_chan_local, move_chan)
	}

	best_eval := math.MaxInt * -1
	for i := 0; i < len(moves); i++ {
		new_eval := -1 * <-eval_chan_local
		if new_eval > best_eval {
			best_eval = new_eval
		}
	}
	eval_chan <- best_eval
	if ply == 1 {
		move_chan <- prev_move
	}
}
