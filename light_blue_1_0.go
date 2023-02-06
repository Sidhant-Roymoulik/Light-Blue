package main

import (
	"math"
	"runtime"
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
			name:       "Light Blue 1",
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
				aspiration_window:   true,
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

	if engine.upgrades.iterative_deepening {
		best_eval, best_move = engine.iterative_deepening(position)
	} else {
		best_eval, best_move = engine.aspiration_window(position, engine.max_ply)
	}

	engine.prev_guess = best_eval
	engine.Add_Zobrist_History(Zobrist.GenHash(position.Update(best_move)))

	// print(engine.zobristHistory)

	return
}

func (engine *light_blue_1_0) iterative_deepening(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.start = time.Now()
	engine.age ^= 1
	max_depth := 0

	for {
		max_depth += 1

		new_eval, new_move := engine.aspiration_window(position, max_depth)

		if engine.time_up() {
			engine.max_q_ply -= 2
			return best_eval, best_move
		}

		best_eval, best_move = new_eval, new_move
		engine.prev_guess = best_eval
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

func (engine *light_blue_1_0) aspiration_window(position *chess.Position, max_depth int) (eval int, move *chess.Move) {

	if max_depth == 1 {
		eval, move = engine.minimax_start(position, -math.MaxInt, math.MaxInt, max_depth)
		return eval, move
	}

	var alpha int = engine.prev_guess - WINDOW_VALUE_TIGHT
	var beta int = engine.prev_guess + WINDOW_VALUE_TIGHT

	eval, move = engine.minimax_start(position, alpha, beta, max_depth)

	if eval <= alpha {
		if DEBUG {
			print("Aspiration tight no work :(")
		}
		alpha = engine.prev_guess - WINDOW_VALUE
		eval, move = engine.minimax_start(position, alpha, beta, max_depth)
	} else if eval >= beta {
		if DEBUG {
			print("Aspiration tight no work :(")
		}
		beta = engine.prev_guess + WINDOW_VALUE
		eval, move = engine.minimax_start(position, alpha, beta, max_depth)
	}

	if eval <= alpha || eval >= beta {
		if DEBUG {
			print("Aspiration no work :(")
		}
		eval, move = engine.minimax_start(position, -math.MaxInt, math.MaxInt, max_depth)
	}

	return eval, move
}

func (engine *light_blue_1_0) minimax_start(position *chess.Position, alpha int, beta int, max_depth int) (eval int, move *chess.Move) {
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
		// Check for search over
		if engine.check_search_over(max_depth) {
			break
		}

		// Pick move
		move := get_move_v3(moves, i)

		// Generate new position
		var new_position *chess.Position = position.Update(move)

		// Generate hash for new position
		var new_hash uint64 = Zobrist.GenHash(new_position)

		// Add to move history
		engine.Add_Zobrist_History(new_hash)

		// Principal-Variation Search
		new_eval := -engine.pv_search(new_position, 1, -beta, -alpha, max_depth)

		// Clear move from history
		engine.Remove_Zobrist_History()

		if new_eval > alpha {
			best_eval = new_eval
			best_move = move

			if new_eval >= beta { // Fail-hard beta-cutoff
				// Add killer move
				engine.addKillerMove(move, 0)

				best_eval = beta
				best_move = move
				tt_flag = BetaFlag
				break
			}

			alpha = new_eval
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

func (engine *light_blue_1_0) pv_search(position *chess.Position, ply int, alpha int, beta int, max_depth int) (eval int) {
	engine.counters.nodes_searched++

	// Check for search over
	if engine.check_search_over(max_depth) {
		return 0
	}

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

	// Check for draw by repetition
	if engine.Is_Draw_By_Repetition(hash) {
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

		// Generate new position
		var new_position *chess.Position = position.Update(move)

		// Generate hash for new position
		var new_hash uint64 = Zobrist.GenHash(new_position)

		// Add to move history
		engine.Add_Zobrist_History(new_hash)

		if bSearchPv {
			// Principal-Variation Search
			new_eval = -engine.pv_search(new_position, ply+1, -beta, -alpha, max_depth)
		} else {
			// Zero-Window Search
			new_eval = -engine.zw_search(new_position, ply+1, -alpha, max_depth)
			if new_eval > alpha && new_eval < beta {
				// Principal-Variation Search
				new_eval = -engine.pv_search(new_position, ply+1, -beta, -alpha, max_depth)
			}
		}

		// Clear move from history
		engine.Remove_Zobrist_History()

		if new_eval > alpha {
			best_eval = new_eval
			best_move = move

			if new_eval >= beta { // Fail-hard beta-cutoff
				// Add killer move
				engine.addKillerMove(move, 0)

				best_eval = beta
				best_move = move
				tt_flag = BetaFlag
				break
			}

			alpha = new_eval
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

func (engine *light_blue_1_0) zw_search(position *chess.Position, ply int, beta int, max_depth int) (eval int) {
	engine.counters.nodes_searched++

	// Check for search over
	if engine.check_search_over(max_depth) {
		return 0
	}

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

	// Check for draw by repetition
	if engine.Is_Draw_By_Repetition(hash) {
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

		// Generate new position
		var new_position *chess.Position = position.Update(move)

		// Generate hash for new position
		var new_hash uint64 = Zobrist.GenHash(new_position)

		// Add to move history
		engine.Add_Zobrist_History(new_hash)

		// Zero-Window Search
		new_eval := -engine.zw_search(new_position, ply+1, 1-beta, max_depth)

		// Clear move from history
		engine.Remove_Zobrist_History()

		if new_eval >= beta { // Fail-hard beta-cutoff
			// Add killer move
			engine.addKillerMove(move, 0)

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

	// Check for search over
	if engine.check_search_over(max_depth) {
		return 0
	}

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

		// Generate new position
		var new_position *chess.Position = position.Update(move)

		new_eval := -engine.q_search(new_position, ply+1, -beta, -alpha, max_depth)

		if new_eval >= beta {
			return beta
		}
		if new_eval > alpha {
			alpha = new_eval
		}
	}
	return alpha
}

func (engine *light_blue_1_0) check_search_over(max_depth int) bool {
	return (engine.time_up())
}
