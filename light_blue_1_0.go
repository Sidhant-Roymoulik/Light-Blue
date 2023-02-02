package main

import (
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

// Upgrades over Version 4.1:
// 		Refactored

type light_blue_1_0 struct {
	EngineClass
}

func new_light_blue_1_0() light_blue_1_0 {
	return light_blue_1_0{
		EngineClass{
			name:       "Light Blue 1.0",
			author:     "Sidhant Roymoulik",
			max_ply:    0,
			max_q_ply:  0,
			time_limit: TIME_LIMIT,
			counters: EngineCounters{
				nodes_searched:   0,
				q_nodes_searched: 0,
				hashes_used:      0,
			},
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
			killer_moves:      [100][2]*chess.Move{},
			threads:           runtime.NumCPU(),
		},
	}
}

func (engine *light_blue_1_0) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.resetCounters()
	engine.resetKillerMoves()

	engine.Add_Zobrist_History(Zobrist.GenHash(position))

	if engine.upgrades.lazy_smp {
		best_eval, best_move = engine.lazy_smp(position)
	} else if engine.upgrades.iterative_deepening {
		best_eval, best_move = engine.iterative_deepening(position)
	} else {
		best_eval, best_move = engine.minimax_start(position, -math.MaxInt, math.MaxInt, engine.max_ply, engine.zobristHistory[:])
	}

	engine.Add_Zobrist_History(Zobrist.GenHash(position.Update(best_move)))

	// print(engine.zobristHistory)

	return
}

func (engine *light_blue_1_0) lazy_smp(position *chess.Position) (best_eval int, best_move *chess.Move) {
	result_chan := make(chan *Result, 200)

	best_depth := 0

	// Adapted from CounterGo (github.com/ChizhovVadim/CounterGo)
	if engine.threads == 1 {
		if engine.upgrades.iterative_deepening {
			engine.iterative_deepening_cc(position, result_chan)
		} else {
			engine.minimax_start_cc(position, -math.MaxInt, math.MaxInt, engine.max_ply, engine.zobristHistory[:], result_chan)
		}
	} else {
		var wg = &sync.WaitGroup{}

		for i := 0; i < engine.threads; i++ {
			wg.Add(1)
			go func(i int) {
				if engine.upgrades.iterative_deepening {
					engine.iterative_deepening_cc(position, result_chan)
				} else {
					engine.minimax_start_cc(position, -math.MaxInt, math.MaxInt, engine.max_ply, engine.zobristHistory[:], result_chan)
				}
				wg.Done()
			}(i)
		}

		wg.Wait()
		close(result_chan)
	}

	result := <-result_chan
	for result != nil {
		if result.depth > best_depth {
			best_eval, best_move, best_depth = result.eval, result.move, result.depth
			if DEBUG {
				print("Lazy Move:", result.move, "Eval:", result.eval, "Depth:", result.depth)
			}
		}
		result = <-result_chan
	}
	engine.max_ply = best_depth

	return
}

func (engine *light_blue_1_0) iterative_deepening_cc(position *chess.Position, result_chan chan *Result) {
	engine.start = time.Now()
	engine.age ^= 1
	max_depth := 0

	for {
		max_depth += 1

		new_eval, new_move := engine.minimax_start(position, -math.MaxInt, math.MaxInt, max_depth, engine.zobristHistory[:])

		if engine.time_up() {
			engine.max_q_ply -= 2
			break
		}

		result_chan <- &Result{new_eval, new_move, max_depth}

		if new_eval >= CHECKMATE_VALUE/10 {
			break
		}

		engine.resetKillerMoves()
	}
}

func (engine *light_blue_1_0) iterative_deepening(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.start = time.Now()
	engine.age ^= 1
	max_depth := 0

	for {
		max_depth += 1

		new_eval, new_move := engine.minimax_start(position, -math.MaxInt, math.MaxInt, max_depth, engine.zobristHistory[:])

		if engine.time_up() {
			engine.max_q_ply -= 2
			return best_eval, best_move
		}

		best_eval, best_move = new_eval, new_move
		engine.max_ply = max_depth

		if DEBUG {
			print("Time:", time.Since(engine.start))
			print("Best Move:", best_move, "Eval:", best_eval, "Depth:", max_depth)
		}

		if best_eval >= CHECKMATE_VALUE/10 {
			return best_eval, best_move
		}

		engine.resetKillerMoves()
	}
}

func (engine *light_blue_1_0) minimax_start_cc(position *chess.Position, alpha int, beta int, max_depth int, hash_history []uint64, result_chan chan *Result) {
	engine.counters.nodes_searched++

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

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

	result_chan <- &Result{best_eval, best_move, max_depth}
}

func (engine *light_blue_1_0) minimax_start(position *chess.Position, alpha int, beta int, max_depth int, hash_history []uint64) (eval int, move *chess.Move) {
	engine.counters.nodes_searched++

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

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

func (engine *light_blue_1_0) pv_search(position *chess.Position, ply int, alpha int, beta int, max_depth int, hash_history []uint64) (eval int) {
	engine.counters.nodes_searched++

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
		engine.counters.hashes_used++
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

func (engine *light_blue_1_0) zw_search(position *chess.Position, ply int, beta int, max_depth int, hash_history []uint64) (eval int) {
	engine.counters.nodes_searched++

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

func (engine *light_blue_1_0) q_search(position *chess.Position, ply int, alpha int, beta int, max_depth int) (eval int) {
	engine.counters.q_nodes_searched++
	engine.max_q_ply = Max(engine.max_q_ply, ply)

	start_eval := eval_v5(position, ply) * getMultiplier(position.Turn() == chess.White)

	// Delta Pruning
	if start_eval >= beta {
		return beta
	}
	if start_eval >= alpha {
		alpha = start_eval
	}

	if ply >= max_depth*2 {
		return start_eval
	}

	// Sort Moves
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

func (engine *light_blue_1_0) Is_Draw_By_Repetition_Local(hash uint64, hash_history []uint64) bool {
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
