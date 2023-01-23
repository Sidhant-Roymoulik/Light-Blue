package main

import (
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

// Upgrades over Version 4.0:
// 		Lazy SMP

type engine_version_4_1 struct {
	EngineClass
}

func new_engine_version_4_1() engine_version_4_1 {
	return engine_version_4_1{
		EngineClass{
			name:       "Version 4.1 (PVS/NWS + Lazy SMP)",
			max_ply:    0,
			time_limit: TIME_LIMIT,
			upgrades: EngineUpgrades{
				move_ordering:       true,
				alphabeta:           true,
				iterative_deepening: true,
				q_search:            true,
				delta_pruning:       true,
				transposition_table: true,
				killer_moves:        true,
				pvs:                 true,
				lazy_smp:            true,
			},
			tt:                TransTable[SearchEntry]{},
			age:               0,
			zobristHistory:    [1024]uint64{},
			zobristHistoryPly: 0,
			prev_guess:        0,
			quit_mtd:          false,
			killer_moves:      [100][2]*chess.Move{},
			threads:           runtime.NumCPU(),
		},
	}

}

func (engine *engine_version_4_1) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	resetCounters()

	engine.Add_Zobrist_History(Zobrist.GenHash(position))

	best_eval, best_move = engine.lazy_smp(position)

	engine.Add_Zobrist_History(Zobrist.GenHash(position.Update(best_move)))

	// print(engine.zobristHistory)
	print("Depth:", engine.max_ply)
	engine.max_ply = 0

	return
}

func (engine *engine_version_4_1) lazy_smp(position *chess.Position) (best_eval int, best_move *chess.Move) {
	result_chan := make(chan *Result, 200)

	best_depth := 0

	// Adapted from CounterGo (github.com/ChizhovVadim/CounterGo)
	if engine.threads == 1 {
		go engine.iterative_deepening(position, result_chan)
	} else {
		var wg = &sync.WaitGroup{}

		for i := 0; i < engine.threads; i++ {
			wg.Add(1)
			go func(i int) {
				engine.iterative_deepening(position, result_chan)
				wg.Done()
			}(i)
		}

		wg.Wait()
		close(result_chan)
	}

	result := <-result_chan
	for result != nil {
		if result.depth > best_depth {
			best_eval, best_move = result.eval, result.move
			best_depth = result.depth
			// if DEBUG {
			// 	print("Lazy Move:", result.move, "Eval:", result.eval, "Depth:", result.depth)
			// }
		}
		result = <-result_chan
	}
	engine.max_ply = best_depth

	return
}

func (engine *engine_version_4_1) iterative_deepening(position *chess.Position, result_chan chan *Result) {
	engine.start = time.Now()
	engine.age ^= 1
	max_depth := 0

	for {
		max_depth += 1

		new_eval, new_move := engine.minimax_start(position, -math.MaxInt, math.MaxInt, max_depth, engine.zobristHistory[:])

		if engine.time_up() {
			break
		}

		result_chan <- &Result{new_eval, new_move, max_depth}
		engine.prev_guess = new_eval

		if DEBUG {
			print("Time:", time.Since(engine.start))
			print("Best Move:", new_move, "Eval:", new_eval, "Depth:", max_depth)
		}

		if new_eval >= CHECKMATE_VALUE/10 {
			break
		}

		for i := 0; i < len(engine.killer_moves); i++ {
			engine.killer_moves[i][0] = nil
			engine.killer_moves[i][1] = nil
		}
	}
}

func (engine *engine_version_4_1) minimax_start(position *chess.Position, alpha int, beta int, max_depth int, hash_history []uint64) (eval int, move *chess.Move) {
	states++

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

	// Check for usable entry in transposition table
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, tt_move = entry.Get(hash, 0, max_depth, -math.MaxInt, math.MaxInt)
	if should_use {
		hash_hits++
		return tt_eval, tt_move
	}

	// Sort Moves
	moves := score_moves_v3(position.ValidMoves(), position.Board(), engine.killer_moves[0])

	// Initialize variables
	var best_eval int = alpha
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag

	// Loop through moves
	for i := 0; i < len(moves); i++ {
		// Check for time up
		if engine.time_up() {
			break
		}

		// Pick move
		move := get_move_v3(moves, i)

		// Principal-Variation Search
		new_eval := engine.pv_search(position.Update(move), 1, -beta, -alpha, max_depth, append(hash_history, hash)) * -1

		if new_eval >= beta { // Fail-hard beta-cutoff
			// Add killer move
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

	// Save position to transposition table
	if !engine.time_up() && best_move != nil {
		var entry *SearchEntry = engine.tt.Store(hash, max_depth, engine.age)
		entry.Set(hash, best_eval, best_move, 0, max_depth, tt_flag, engine.age)
		hash_writes++
	}

	return best_eval, best_move
}

func (engine *engine_version_4_1) pv_search(position *chess.Position, ply int, alpha int, beta int, max_depth int, hash_history []uint64) (eval int) {
	states++

	// Check for time up
	if engine.time_up() {
		return 0
	}

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

	// Check for draw by repetition
	if engine.Is_Draw_By_Repetition_Local(hash, hash_history) {
		return 0
	}

	// Check for usable entry in transposition table
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, _ = entry.Get(hash, 0, max_depth-ply, alpha, beta)
	if should_use {
		hash_hits++
		return tt_eval
	}

	// Start Q-Search
	if ply > max_depth {
		return engine.q_search(position, ply, alpha, beta, max_depth)
	}

	// Sort Moves
	moves := score_moves_v3(position.ValidMoves(), position.Board(), engine.killer_moves[ply])

	// If there are no moves, return the eval
	if len(moves) == 0 {
		return eval_v5(position, ply) * getMultiplier(position.Turn() == chess.White)
	}

	// Initialize variables
	var best_eval int = alpha
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag
	var bSearchPv bool = true

	// Loop through moves
	for i := 0; i < len(moves); i++ {
		// Pick move
		move := get_move_v3(moves, i)

		new_eval := 0

		if bSearchPv {
			// Principal-Variation Search
			new_eval = -engine.pv_search(position.Update(move), ply+1, -beta, -alpha, max_depth, append(hash_history, hash))
		} else {
			// Zero-Window Search
			new_eval = -engine.zw_search(position.Update(move), ply+1, -alpha, max_depth, append(hash_history, hash))
			if new_eval > alpha && new_eval < beta {
				// Principal-Variation Search
				new_eval = -engine.pv_search(position.Update(move), ply+1, -beta, -alpha, max_depth, append(hash_history, hash))
			}
		}

		if new_eval >= beta { // Fail-hard beta-cutoff
			// Add killer move
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
			bSearchPv = false
		}
	}

	// Save position to transposition table
	if !engine.time_up() {
		var entry *SearchEntry = engine.tt.Store(hash, max_depth-ply, engine.age)
		entry.Set(hash, best_eval, best_move, 0, max_depth-ply, tt_flag, engine.age)

		hash_writes++
	}

	return best_eval
}

func (engine *engine_version_4_1) zw_search(position *chess.Position, ply int, beta int, max_depth int, hash_history []uint64) (eval int) {
	states++

	// Check for time up
	if engine.time_up() {
		return 0
	}

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

	// Check for draw by repetition
	if engine.Is_Draw_By_Repetition_Local(hash, hash_history) {
		return 0
	}

	alpha := beta - 1

	// Check for usable entry in transposition table
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, _ = entry.Get(hash, 0, max_depth-ply, alpha, beta)
	if should_use {
		hash_hits++
		return tt_eval
	}

	// Start Q-Search
	if ply > max_depth {
		return engine.q_search(position, ply, alpha, beta, max_depth)
	}

	// Sort Moves
	moves := score_moves_v3(position.ValidMoves(), position.Board(), engine.killer_moves[ply])

	// If there are no moves, return the eval
	if len(moves) == 0 {
		return eval_v5(position, ply) * getMultiplier(position.Turn() == chess.White)
	}

	// Initialize variables
	var best_eval int = alpha
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag

	// Loop through moves
	for i := 0; i < len(moves); i++ {
		// Pick move
		move := get_move_v3(moves, i)

		// Zero-Window Search
		new_eval := -engine.zw_search(position.Update(move), ply+1, 1-beta, max_depth, append(hash_history, hash))

		if new_eval >= beta { // Fail-hard beta-cutoff
			// Add killer move
			if !move.HasTag(chess.Capture) && move != engine.killer_moves[ply][0] {
				engine.killer_moves[ply][1] = engine.killer_moves[ply][0]
				engine.killer_moves[ply][0] = move
			}

			best_eval = beta

			best_move = move
			tt_flag = BetaFlag
			break
		}
		best_eval = alpha // Fail-hard, return alpha
	}

	// Save position to transposition table
	if !engine.time_up() {
		var entry *SearchEntry = engine.tt.Store(hash, max_depth-ply, engine.age)
		entry.Set(hash, best_eval, best_move, 0, max_depth-ply, tt_flag, engine.age)

		hash_writes++
	}

	return best_eval
}

func (engine *engine_version_4_1) q_search(position *chess.Position, ply int, alpha int, beta int, max_depth int) (eval int) {
	q_states++

	start_eval := eval_v5(position, ply) * getMultiplier(position.Turn() == chess.White)

	if start_eval >= beta {
		return beta
	}
	if start_eval >= alpha {
		alpha = start_eval
	}

	if ply > max_depth*2 {
		return start_eval
	}

	moves := score_moves_v2(get_q_moves(position), position.Board())

	if len(moves) == 0 {
		return start_eval
	}

	for i := 0; i < len(moves); i++ {
		move := get_move_v3(moves, i)

		new_eval := engine.q_search(position.Update(move), ply+1, -beta, -alpha, max_depth) * -1

		if new_eval >= beta {
			return beta
		}
		if new_eval > alpha {
			alpha = new_eval
		}
	}
	return alpha
}

func (engine *engine_version_4_1) Is_Draw_By_Repetition_Local(hash uint64, hash_history []uint64) bool {
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
