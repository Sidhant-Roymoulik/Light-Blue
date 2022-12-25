package main

import (
	"math"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

type e_minimax_id_cc struct {
	EngineClass
}

func new_engine_minimax_id_cc() e_minimax_id_cc {
	return e_minimax_id_cc{
		EngineClass{
			name:       "Minimax with Iterative Deepening and Concurrency",
			max_ply:    0,
			time_limit: TIME_LIMIT,
			upgrades: EngineUpgrades{
				iterative_deepening: true,
				concurrent:          true,
			},
		},
	}
}

func (engine *e_minimax_id_cc) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()

	engine.start = time.Now()

	for {
		engine.max_ply = engine.max_ply + 1
		if time.Since(engine.start) > engine.time_limit {
			break
		} else {
			best_eval, best_move = engine.minimax_start(position, 0, position.Turn() == chess.White)
		}

		if best_eval >= 100000 {
			break
		}
	}
	print("Depth:", engine.max_ply-1)
	engine.max_ply = 0

	return
}
func (engine *e_minimax_id_cc) minimax_start(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
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
func (engine *e_minimax_id_cc) minimax(position *chess.Position, ply int, turn bool, prev_move *chess.Move, eval_chan chan int, move_chan chan *chess.Move) {
	states++

	if ply >= engine.max_ply || len(position.ValidMoves()) == 0 || time.Since(engine.start) >= engine.time_limit {
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

// func (engine *e_minimax_id_cc) setConfig(max_ply int, time_limit time.Duration) {
// 	engine.max_ply = max_ply
// 	engine.time_limit = time_limit
// }
