package main

import (
	"math"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

// Upgrades over Version 2.0:
// 		Better MTD selection

type engine_version_2_1 struct {
	EngineClass
}

func new_engine_version_2_1() engine_version_2_1 {
	return engine_version_2_1{
		EngineClass{
			name:       "Version 2.1 (MTD Selection)",
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
			},
			tt:                TransTable[SearchEntry]{},
			age:               0,
			zobristHistory:    [1024]uint64{},
			zobristHistoryPly: 0,
			prev_guess:        0,
			use_mtd_f:         false,
		},
	}

}

func (engine *engine_version_2_1) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()

	best_eval, best_move = engine.iterative_deepening(position)

	// print(engine.zobristHistory)
	print("Depth:", engine.max_ply-1)
	engine.max_ply = 0

	return
}

func (engine *engine_version_2_1) iterative_deepening(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.start = time.Now()
	engine.age ^= 1
	engine.use_mtd_f = false

	best_eval = 0
	var new_eval int = 0
	var new_move *chess.Move = nil
	for {
		engine.max_ply = engine.max_ply + 1

		if engine.use_mtd_f {
			new_eval, new_move = engine.mtd_f(position, best_eval)
		} else {
			new_eval, new_move = engine.mtd_bi(position)
		}
		// new_eval, new_move := engine.minimax_start(position, position.Turn() == chess.White, -math.MaxInt, math.MaxInt)

		if engine.time_up() {
			break
		}

		best_eval, best_move = new_eval, new_move
		print("Top Level Move:", best_move, "Eval:", best_eval, "Depth:", engine.max_ply)

		if int(math.Abs(float64(best_eval)-float64(engine.prev_guess))) < MTD_ITER_CUTOFF {
			engine.use_mtd_f = true // Switch from MTD(bi) to MTD(f) if the gap between guesses is low
		}
		engine.prev_guess = best_eval

		if best_eval >= CHECKMATE_VALUE/10 {
			break
		}
	}

	return
}

func (engine *engine_version_2_1) mtd_f(position *chess.Position, g int) (eval int, move *chess.Move) {
	mtd_f_iter := 0
	eval = g
	upper := CHECKMATE_VALUE
	lower := -CHECKMATE_VALUE
	var new_move *chess.Move = nil
	for lower < upper-MTD_EVAL_CUTOFF {
		if engine.time_up() {
			break
		}
		if mtd_f_iter > MTD_ITER_CUTOFF { // If there is an eval jump, use MTD(bi)
			print("MTD(f) Iterations:", mtd_f_iter)
			return engine.mtd_bi(position)
		}
		beta := Max(eval, lower+1)
		eval, new_move = engine.minimax_start(position, position.Turn() == chess.White, beta-1, beta)
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
	print("MTD(f) Iterations:", mtd_f_iter)
	return
}
func (engine *engine_version_2_1) mtd_bi(position *chess.Position) (eval int, move *chess.Move) {
	mtd_bi_iter := 0
	upper := CHECKMATE_VALUE
	lower := -CHECKMATE_VALUE
	var new_move *chess.Move = nil
	for lower < upper-MTD_EVAL_CUTOFF {
		if engine.time_up() {
			break
		}
		beta := (lower + upper + 1) / 2
		eval, new_move = engine.minimax_start(position, position.Turn() == chess.White, beta-1, beta)
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
	print("MTD(bi) Iterations:", mtd_bi_iter)
	return
}

func (engine *engine_version_2_1) minimax_start(position *chess.Position, turn bool, alpha int, beta int) (eval int, move *chess.Move) {
	states++

	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, tt_move = entry.Get(hash, 0, engine.max_ply, -math.MaxInt, math.MaxInt)

	if should_use {
		hash_hits++
		return tt_eval, tt_move
	}

	moves := score_moves_v2(position.ValidMoves(), position.Board())

	var best_eval int = alpha
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag

	for i := 0; i < len(moves); i++ {
		if engine.time_up() {
			break
		}

		move := get_move_v3(moves, i)

		var updated_position = position.Update(move)
		var updated_hash = Zobrist.GenHash(updated_position)

		engine.Add_Zobrist_History(updated_hash)

		new_eval := engine.minimax(position.Update(move), 1, !turn, -beta, -alpha) * -1

		engine.Remove_Zobrist_History()

		// print("Top Level Move:", move, "Eval:", new_eval)
		// if new_eval > best_eval {
		// 	best_eval = new_eval
		// 	best_move = move
		// }

		if new_eval >= beta {
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

func (engine *engine_version_2_1) minimax(position *chess.Position, ply int, turn bool, alpha int, beta int) (eval int) {
	states++

	if engine.time_up() {
		return 0
	}

	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, _ = entry.Get(hash, 0, engine.max_ply-ply, alpha, beta)
	if should_use {
		hash_hits++
		return tt_eval
	}

	if ply > engine.max_ply {
		return engine.q_search(position, ply, turn, alpha, beta)
	}
	if len(position.ValidMoves()) == 0 {
		return eval_v5(position, ply) * getMultiplier(turn)
	}
	if engine.Is_Draw_By_Repetition(hash) {
		return 0
	}

	var best_eval int = alpha
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag

	moves := score_moves_v2(position.ValidMoves(), position.Board())

	for i := 0; i < len(moves); i++ {
		move := get_move_v3(moves, i)

		var updated_position = position.Update(move)
		var updated_hash = Zobrist.GenHash(updated_position)

		engine.Add_Zobrist_History(updated_hash)

		new_eval := engine.minimax(position.Update(move), ply+1, !turn, -beta, -alpha) * -1

		engine.Remove_Zobrist_History()

		if new_eval >= beta {
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

func (engine *engine_version_2_1) q_search(position *chess.Position, ply int, turn bool, alpha int, beta int) (eval int) {
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
