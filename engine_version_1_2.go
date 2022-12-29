package main

import (
	"math"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

// Upgrades over Version 1.1:
// 		Faster move picking

type engine_version_1_2 struct {
	EngineClass
}

func new_engine_version_1_2() engine_version_1_2 {
	return engine_version_1_2{
		EngineClass{
			name:       "Version 1.2 (MO, AB, ID, QS, DP, TT)",
			max_ply:    0,
			time_limit: TIME_LIMIT,
			upgrades: EngineUpgrades{
				move_ordering:       true,
				alphabeta:           true,
				iterative_deepening: true,
				q_search:            true,
				delta_pruning:       true,
				transposition_table: true,
			},
			tt:                TransTable[SearchEntry]{},
			age:               0,
			zobristHistory:    [1024]uint64{},
			zobristHistoryPly: 0,
		},
	}

}

func (engine *engine_version_1_2) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()

	best_eval, best_move = engine.iterative_deepening(position)

	// print(engine.zobristHistory)
	print("Depth:", engine.max_ply-1)
	engine.max_ply = 0

	return
}

func (engine *engine_version_1_2) iterative_deepening(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.start = time.Now()
	engine.age ^= 1

	for {
		engine.max_ply = engine.max_ply + 1
		new_eval, new_move := engine.minimax_start(position, position.Turn() == chess.White)
		if time.Since(engine.start) > engine.time_limit {
			break
		} else {
			best_eval, best_move = new_eval, new_move
		}

		if best_eval >= 100000 {
			break
		}
	}

	return
}

func (engine *engine_version_1_2) minimax_start(position *chess.Position, turn bool) (best_eval int, best_move *chess.Move) {
	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, tt_move = entry.Get(hash, 0, engine.max_ply, -math.MaxInt, math.MaxInt)

	if should_use {
		hash_hits++
		return tt_eval, tt_move
	}

	moves := score_moves(position.ValidMoves(), position.Board())

	best_eval = math.MaxInt * -1
	best_move = moves[0].move
	for i := 0; i < len(moves); i++ {
		if time.Since(engine.start) > engine.time_limit {
			break
		}
		move := get_move_v2(moves, i)

		var updated_position = position.Update(move)
		var updated_hash = Zobrist.GenHash(updated_position)

		engine.Add_Zobrist_History(updated_hash)

		new_eval := engine.minimax(position.Update(move), 1, !turn, math.MaxInt*-1, -best_eval) * -1

		engine.Remove_Zobrist_History()

		// print("Top Level Move:", move, "Eval:", new_eval)
		if new_eval > best_eval {
			best_eval = new_eval
			best_move = move
		}
	}

	if !engine.time_up() && best_move != nil { // this is off
		var entry *SearchEntry = engine.tt.Store(hash, engine.max_ply, engine.age)
		entry.Set(hash, best_eval, best_move, 0, engine.max_ply, ExactFlag, engine.age)
		hash_writes++
	}

	return best_eval, best_move
}
func (engine *engine_version_1_2) minimax(position *chess.Position, ply int, turn bool, alpha int, beta int) (eval int) {
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

	moves := score_moves(position.ValidMoves(), position.Board())

	for i := 0; i < len(moves); i++ {
		move := get_move_v2(moves, i)

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

func (engine *engine_version_1_2) q_search(position *chess.Position, ply int, turn bool, alpha int, beta int) (best_eval int) {
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

	moves := score_moves(get_q_moves(position), position.Board())

	if len(moves) == 0 {
		return start_eval
	}

	for i := 0; i < len(moves); i++ {
		move := get_move_v2(moves, i)

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
