package main

import (
	"math"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

// Upgrades over Version 2.2:
// 		Killer Moves

type engine_version_2_3 struct {
	EngineClass
}

func new_engine_version_2_3() engine_version_2_3 {
	return engine_version_2_3{
		EngineClass{
			name:       "Version 2.3 (Killer Moves)",
			max_ply:    0,
			time_limit: TIME_LIMIT,
			upgrades: EngineUpgrades{
				move_ordering:       true,
				alphabeta:           true,
				iterative_deepening: true,
				q_search:            true,
				delta_pruning:       true,
				transposition_table: true,
				mtd:                 true,
				concurrent:          true,
				killer_moves:        true,
			},
			tt:                TransTable[SearchEntry]{},
			age:               0,
			zobristHistory:    [1024]uint64{},
			zobristHistoryPly: 0,
			prev_guess:        0,
			quit_mtd:          false,
			killer_moves:      [100][2]*chess.Move{},
		},
	}

}

func (engine *engine_version_2_3) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()

	best_eval, best_move = engine.iterative_deepening(position)

	engine.Add_Zobrist_History(Zobrist.GenHash(position.Update(best_move)))

	// print(engine.zobristHistory)
	print("Depth:", engine.max_ply-1)
	engine.max_ply = 0

	return
}

func (engine *engine_version_2_3) iterative_deepening(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.start = time.Now()
	engine.age ^= 1

	for {
		engine.max_ply = engine.max_ply + 1
		engine.quit_mtd = false

		eval_chan := make(chan int, 2)
		move_chan := make(chan *chess.Move, 2)

		go engine.mtd_bi(position, eval_chan, move_chan)
		go engine.mtd_f(position, engine.prev_guess, eval_chan, move_chan)

		new_eval := <-eval_chan
		new_move := <-move_chan

		if engine.time_up() {
			break
		}

		best_eval, best_move = new_eval, new_move
		engine.prev_guess = best_eval

		<-eval_chan
		<-move_chan

		if DEBUG {
			// print("Time:", time.Since(engine.start))
			print("Best Move:", best_move, "Eval:", best_eval, "Depth:", engine.max_ply)
		}

		if best_eval >= CHECKMATE_VALUE/10 {
			break
		}
	}

	return
}

func (engine *engine_version_2_3) mtd_f(position *chess.Position, g int, eval_chan chan int, move_chan chan *chess.Move) {
	mtd_f_iter := 0
	eval := g
	upper := CHECKMATE_VALUE
	lower := -CHECKMATE_VALUE
	var move *chess.Move = nil
	var new_move *chess.Move = nil
	for lower < upper-MTD_EVAL_CUTOFF {
		if engine.time_up() || engine.quit_mtd {
			eval_chan <- eval
			move_chan <- move
			engine.quit_mtd = true
			return
		}
		beta := Max(eval, lower+1)
		eval, new_move = engine.minimax_start(position, position.Turn() == chess.White, beta-1, beta, engine.zobristHistory[:])
		if new_move != nil {
			move = new_move
		}
		if eval < beta {
			upper = eval
		} else {
			lower = eval
		}
		mtd_f_iter++
	}
	if DEBUG {
		print("MTD(f) Iterations:", mtd_f_iter)
	}
	eval_chan <- eval
	move_chan <- move
	engine.quit_mtd = true
}
func (engine *engine_version_2_3) mtd_bi(position *chess.Position, eval_chan chan int, move_chan chan *chess.Move) {
	mtd_bi_iter := 0
	eval := 0
	upper := CHECKMATE_VALUE
	lower := -CHECKMATE_VALUE
	var move *chess.Move = nil
	var new_move *chess.Move = nil
	for lower < upper-MTD_EVAL_CUTOFF {
		if engine.time_up() || engine.quit_mtd {
			eval_chan <- eval
			move_chan <- move
			engine.quit_mtd = true
			return
		}
		beta := (lower + upper + 1) / 2
		eval, new_move = engine.minimax_start(position, position.Turn() == chess.White, beta-1, beta, engine.zobristHistory[:])
		if new_move != nil {
			move = new_move
		}
		if eval < beta {
			upper = eval
		} else {
			lower = eval
		}
		mtd_bi_iter++
	}
	if DEBUG {
		print("MTD(bi) Iterations:", mtd_bi_iter)
	}
	eval_chan <- eval
	move_chan <- move
	engine.quit_mtd = true
}

func (engine *engine_version_2_3) minimax_start(position *chess.Position, turn bool, alpha int, beta int, hash_history []uint64) (eval int, move *chess.Move) {
	states++

	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, tt_move = entry.Get(hash, 0, engine.max_ply, -math.MaxInt, math.MaxInt)

	if should_use {
		hash_hits++
		return tt_eval, tt_move
	}

	moves := score_moves_v3(position.ValidMoves(), position.Board(), engine.killer_moves[0])

	var best_eval int = alpha
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag

	for i := 0; i < len(moves); i++ {
		if engine.time_up() {
			break
		}

		move := get_move_v3(moves, i)

		new_eval := engine.minimax(position.Update(move), 1, !turn, -beta, -alpha, append(hash_history, hash)) * -1

		// print("Top Level Move:", move, "Eval:", new_eval)
		// if new_eval > best_eval {
		// 	best_eval = new_eval
		// 	best_move = move
		// }

		if new_eval >= beta {

			if !move.HasTag(chess.Capture) && move != engine.killer_moves[0][0] {
				engine.killer_moves[0][1] = engine.killer_moves[0][0]
				engine.killer_moves[0][0] = move
			}

			best_eval = beta

			best_move = move
			tt_flag = BetaFlag
			break
		}
		if new_eval > alpha {
			best_eval = new_eval

			alpha = new_eval
			best_move = move
			tt_flag = ExactFlag
		}
	}

	if !engine.time_up() && best_move != nil {
		var entry *SearchEntry = engine.tt.Store(hash, engine.max_ply, engine.age)
		entry.Set(hash, best_eval, best_move, 0, engine.max_ply, tt_flag, engine.age)
		hash_writes++
	}

	return best_eval, best_move
}

func (engine *engine_version_2_3) minimax(position *chess.Position, ply int, turn bool, alpha int, beta int, hash_history []uint64) (eval int) {
	states++

	if engine.time_up() {
		return 0
	}

	var hash uint64 = Zobrist.GenHash(position)

	if engine.Is_Draw_By_Repetition_Local(hash, hash_history) {
		return 0
	}

	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, _ = entry.Get(hash, 0, engine.max_ply-ply, alpha, beta)
	if should_use {
		hash_hits++
		return tt_eval
	}

	if ply > engine.max_ply {
		return engine.q_search(position, ply, turn, alpha, beta)
	}

	moves := score_moves_v3(position.ValidMoves(), position.Board(), engine.killer_moves[ply])

	if len(moves) == 0 {
		return eval_v5(position, ply) * getMultiplier(turn)
	}

	var best_eval int = alpha
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag

	for i := 0; i < len(moves); i++ {
		move := get_move_v3(moves, i)

		new_eval := engine.minimax(position.Update(move), ply+1, !turn, -beta, -alpha, append(hash_history, hash)) * -1

		if new_eval >= beta {

			if !move.HasTag(chess.Capture) && move != engine.killer_moves[ply][0] {
				engine.killer_moves[ply][1] = engine.killer_moves[ply][0]
				engine.killer_moves[ply][0] = move
			}

			best_eval = beta

			best_move = move
			tt_flag = BetaFlag
			break
		}
		if new_eval > alpha {
			best_eval = new_eval

			alpha = new_eval
			best_move = move
			tt_flag = ExactFlag
		}
	}

	if !engine.time_up() {
		var entry *SearchEntry = engine.tt.Store(hash, engine.max_ply-ply, engine.age)
		entry.Set(hash, best_eval, best_move, 0, engine.max_ply-ply, tt_flag, engine.age)

		hash_writes++
	}

	return best_eval
}

func (engine *engine_version_2_3) q_search(position *chess.Position, ply int, turn bool, alpha int, beta int) (eval int) {
	q_states++

	start_eval := eval_v5(position, ply) * getMultiplier(turn)

	if start_eval >= beta {
		return beta
	}
	if start_eval >= alpha {
		alpha = start_eval
	}

	if ply > engine.max_ply*2 {
		return start_eval
	}

	moves := score_moves_v2(get_q_moves(position), position.Board())

	if len(moves) == 0 {
		return start_eval
	}

	for i := 0; i < len(moves); i++ {
		move := get_move_v3(moves, i)

		new_eval := engine.q_search(position.Update(move), ply+1, !turn, -beta, -alpha) * -1

		if new_eval >= beta {
			return beta
		}
		if new_eval > alpha {
			alpha = new_eval
		}
	}
	return alpha
}

func (engine *engine_version_2_3) Is_Draw_By_Repetition_Local(hash uint64, hash_history []uint64) bool {
	for i := 0; i < len(hash_history); i++ {
		if hash_history[i] == 0 {
			return false
		}
		if hash_history[i] == hash {
			return true
		}
	}
	return false
}
