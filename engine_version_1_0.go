package main

import (
	"math"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

// Upgrades over engine_minimax_mo_ab_q_id:
// 		Adds transposition table

type engine_version_1_0 struct {
	EngineClass
}

func new_engine_version_1_0() engine_version_1_0 {
	return engine_version_1_0{
		EngineClass{
			name:       "Version 1.0 (MO, AB, ID, QS, DP, TT)",
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

func (engine *engine_version_1_0) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()

	best_eval, best_move = engine.iterative_deepening(position)

	print("Depth:", engine.max_ply-1)
	engine.max_ply = 0

	return
}

func (engine *engine_version_1_0) iterative_deepening(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.start = time.Now()
	engine.age ^= 1

	for {
		engine.max_ply = engine.max_ply + 1
		new_eval, new_move := engine.minimax_start(position, 0, position.Turn() == chess.White)
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

func (engine *engine_version_1_0) minimax_start(position *chess.Position, ply int, turn bool) (best_eval int, best_move *chess.Move) {
	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, tt_move = entry.Get(hash, 0, engine.max_ply, -math.MaxInt, math.MaxInt)

	if should_use {
		hash_hits++
		return tt_eval, tt_move
	}

	moves := move_ordering_v2(position)

	best_eval = math.MaxInt * -1
	best_move = moves[0]
	for _, move := range moves {
		if time.Since(engine.start) > engine.time_limit {
			break
		}

		new_eval := engine.minimax(position.Update(move), ply+1, !turn, math.MaxInt*-1, -best_eval) * -1
		// print("Top Level Move:", move, "Eval:", new_eval)
		if new_eval > best_eval {
			best_eval = new_eval
			best_move = move
		}
	}

	if !engine.time_up() && best_move != nil { // this is off
		var entry *SearchEntry = engine.tt.Store(hash, ply, engine.age)
		entry.Set(hash, best_eval, best_move, 0, ply, ExactFlag, engine.age)
		hash_writes++
	}

	return best_eval, best_move
}
func (engine *engine_version_1_0) minimax(position *chess.Position, ply int, turn bool, alpha int, beta int) (eval int) {
	if engine.time_up() {
		return 0
	}

	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, _ = entry.Get(hash, 0, engine.max_ply, -math.MaxInt, math.MaxInt)
	if should_use {
		hash_hits++
		return tt_eval
	}

	states++

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

	moves := move_ordering_v2(position)
	for _, move := range moves {
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
		var entry *SearchEntry = engine.tt.Store(hash, ply, engine.age)
		entry.Set(hash, best_eval, best_move, 0, ply, tt_flag, engine.age)

		hash_writes++
	}

	return best_eval
}

func (engine *engine_version_1_0) q_search(position *chess.Position, ply int, turn bool, alpha int, beta int) (best_eval int) {
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

	moves := getQMoves(position)

	if len(moves) == 0 {
		return start_eval
	}

	for _, move := range moves {
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
