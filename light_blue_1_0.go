package main

import (
	"fmt"
	"math"
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
			time_limit: TIME_LIMIT,
			counters: EngineCounters{
				nodes_searched:   0,
				q_nodes_searched: 0,
				hashes_used:      0,
				hashes_written:   0,
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
			timer:             TimeManager{},
			tt:                TransTable[SearchEntry]{},
			age:               0,
			zobristHistory:    [1024]uint64{},
			zobristHistoryPly: 0,
			killer_moves:      [100][2]*chess.Move{},
		},
	}
}

func (engine *light_blue_1_0) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.resetCounters()
	engine.resetKillerMoves()

	engine.Add_Zobrist_History(Zobrist.GenHash(position))

	pvLine := PVLine{}

	if engine.upgrades.iterative_deepening {
		best_eval, best_move = engine.iterative_deepening(position, &pvLine)
	} else {
		best_eval = engine.aspiration_window(position, engine.max_ply, &pvLine)
		best_move = pvLine.getPVMove()
	}

	engine.prev_guess = best_eval
	engine.Add_Zobrist_History(Zobrist.GenHash(position.Update(best_move)))

	// print(engine.zobristHistory)
	// print(engine.killer_moves)

	return
}

func (engine *light_blue_1_0) iterative_deepening(
	position *chess.Position, pvLine *PVLine,
) (best_eval int, best_move *chess.Move) {
	engine.start = time.Now()
	engine.age ^= 1
	engine.timer.Start()

	for depth := 1; depth <= MAX_DEPTH &&
		depth <= int(engine.timer.MaxDepth) &&
		engine.timer.MaxNodeCount > 0; depth++ {

		pvLine.clear()

		new_eval := engine.aspiration_window(position, depth, pvLine)

		if engine.timer.IsStopped() {
			break
		}

		best_eval = new_eval
		engine.prev_guess = best_eval
		engine.max_ply = depth

		best_move = pvLine.getPVMove()
		total_nodes := engine.counters.nodes_searched +
			engine.counters.q_nodes_searched
		total_time := time.Since(engine.start).Milliseconds() + 1

		fmt.Printf(
			"info depth %d score %s nodes %d nps %d time %d pv %s\n",
			engine.max_ply,
			getMateOrCPScore(best_eval),
			total_nodes,
			int64(total_nodes*1000)/total_time,
			total_time,
			pvLine,
		)

		if best_eval >= MATE_CUTOFF {
			break
		}
	}

	return best_eval, best_move
}

func (engine *light_blue_1_0) aspiration_window(
	position *chess.Position, max_depth int, pvLine *PVLine,
) (eval int) {

	if max_depth == 1 {
		eval = engine.minimax_start(
			position, -math.MaxInt, math.MaxInt, max_depth, pvLine,
		)
		return eval
	}

	var alpha int = engine.prev_guess - WINDOW_VALUE_TIGHT
	var beta int = engine.prev_guess + WINDOW_VALUE_TIGHT

	eval = engine.minimax_start(position, alpha, beta, max_depth, pvLine)

	if eval <= alpha {
		if DEBUG {
			// print("Aspiration tight fail low")
		}
		alpha = engine.prev_guess - WINDOW_VALUE
		eval = engine.minimax_start(position, alpha, beta, max_depth, pvLine)
	} else if eval >= beta {
		if DEBUG {
			// print("Aspiration tight fail high")
		}
		beta = engine.prev_guess + WINDOW_VALUE
		eval = engine.minimax_start(position, alpha, beta, max_depth, pvLine)
	}

	if eval <= alpha || eval >= beta {
		if DEBUG {
			// print("Aspiration loose fail")
		}
		eval = engine.minimax_start(
			position, -math.MaxInt, math.MaxInt, max_depth, pvLine,
		)
	}

	return eval
}

func (engine *light_blue_1_0) minimax_start(
	position *chess.Position,
	alpha int,
	beta int,
	max_depth int,
	pvLine *PVLine,
) (eval int) {
	engine.counters.nodes_searched++

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

	// Sort Moves
	moves := score_moves(
		position.ValidMoves(),
		position.Board(),
		engine.killer_moves[0],
		pvLine.getPVMove(),
	)

	// Initialize variables
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag
	var childPVLine PVLine = PVLine{}

	// Loop through moves
	for i := 0; i < len(moves); i++ {
		// Check for search over
		if engine.timer.IsStopped() {
			break
		}

		// Pick move
		get_move(moves, i)
		move := moves[i].move

		// Generate new position
		var new_position *chess.Position = position.Update(move)

		// Add to move history
		engine.Add_Zobrist_History(Zobrist.GenHash(new_position))

		// Principal-Variation Search
		new_eval := -engine.pv_search(
			new_position,
			1,
			-beta,
			-alpha,
			max_depth,
			&childPVLine)

		// Clear move from history
		engine.Remove_Zobrist_History()

		if new_eval > alpha {
			best_move = move

			if new_eval >= beta { // Fail-hard beta-cutoff
				// Add killer move
				engine.addKillerMove(move, 0)

				alpha = beta
				tt_flag = BetaFlag

				break
			}

			alpha = new_eval
			tt_flag = ExactFlag
			pvLine.update(move, childPVLine)
		}

		childPVLine.clear()
	}

	// Save position to transposition table
	if !engine.timer.IsStopped() {
		var entry *SearchEntry = engine.tt.Store(hash, max_depth, engine.age)
		entry.Set(hash, alpha, best_move, 0, max_depth, tt_flag, engine.age)

		engine.counters.hashes_written++
	}

	return alpha
}

func (engine *light_blue_1_0) pv_search(
	position *chess.Position,
	ply int,
	alpha int,
	beta int,
	max_depth int,
	pvLine *PVLine,
) (eval int) {
	engine.counters.nodes_searched++

	if ply >= MAX_DEPTH {
		return eval_pos(position, ply)
	}

	if engine.getTotalNodesSearched() >= engine.timer.MaxNodeCount {
		engine.timer.ForceStop()
	}
	if (engine.getTotalNodesSearched() & 2047) == 0 {
		engine.timer.CheckIfTimeIsUp()
	}

	if engine.timer.IsStopped() {
		return 0
	}

	// Generate hash for position
	var hash uint64 = Zobrist.GenHash(position)

	// Initialize variables
	var childPVLine PVLine = PVLine{}
	isPVNode := beta-alpha != 1

	// Check for draw by repetition
	if engine.Is_Draw_By_Repetition(hash) {
		return 0
	}

	// Start Q-Search
	if ply >= max_depth {
		engine.counters.nodes_searched--
		return engine.q_search(position, ply, alpha, beta, max_depth)
	}

	// Check for usable entry in transposition table
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, tt_move = entry.Get(
		hash, ply, max_depth-ply, alpha, beta,
	)
	if should_use {
		engine.counters.hashes_used++
		return tt_eval
	}

	// Internal Iterative Deepening
	if max_depth-ply > IID_Depth_Limit &&
		(isPVNode || entry.GetFlag() == BetaFlag) &&
		tt_move == nil {
		engine.pv_search(
			position,
			ply+1,
			-beta,
			-alpha,
			max_depth-IID_Depth_Limit,
			&childPVLine)
		if len(childPVLine.Moves) > 0 {
			tt_move = childPVLine.getPVMove()
			childPVLine.clear()
		}
	}

	// Sort Moves
	moves := score_moves(
		position.ValidMoves(),
		position.Board(),
		engine.killer_moves[ply],
		tt_move,
	)

	// If there are no moves, return the eval
	if len(moves) == 0 {
		return eval_pos(position, ply)
	}

	// Initialize variables
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag
	var bSearchPv bool = true

	// Loop through moves
	for i := 0; i < len(moves); i++ {
		// Pick move
		get_move(moves, i)
		move := moves[i].move

		new_eval := 0

		// Generate new position
		var new_position *chess.Position = position.Update(move)

		// Add to move history
		engine.Add_Zobrist_History(Zobrist.GenHash(new_position))

		if bSearchPv {
			// Principal-Variation Search
			new_eval = -engine.pv_search(
				new_position,
				ply+1,
				-beta,
				-alpha,
				max_depth,
				&childPVLine,
			)
		} else {
			// Zero-Window Search
			new_eval = -engine.pv_search(
				new_position,
				ply+1,
				-alpha-1,
				-alpha,
				max_depth,
				&childPVLine,
			)
			if new_eval > alpha && new_eval < beta {
				// Principal-Variation Search
				new_eval = -engine.pv_search(
					new_position,
					ply+1,
					-beta,
					-alpha,
					max_depth,
					&childPVLine,
				)
			}
		}

		// Clear move from history
		engine.Remove_Zobrist_History()

		if new_eval > alpha {
			best_move = move

			if new_eval >= beta { // Fail-hard beta-cutoff
				// Add killer move
				engine.addKillerMove(move, ply)

				alpha = beta
				tt_flag = BetaFlag

				break
			}

			alpha = new_eval
			tt_flag = ExactFlag
			bSearchPv = false
			pvLine.update(move, childPVLine)
		}

		childPVLine.clear()
	}

	// Save position to transposition table
	if !engine.timer.IsStopped() {
		var entry *SearchEntry = engine.tt.Store(
			hash, max_depth-ply, engine.age,
		)
		entry.Set(
			hash, alpha, best_move, ply, max_depth-ply, tt_flag, engine.age,
		)

		engine.counters.hashes_written++
	}

	return alpha
}

func (engine *light_blue_1_0) q_search(
	position *chess.Position,
	ply int,
	alpha int,
	beta int,
	max_depth int,
) (eval int) {
	engine.counters.q_nodes_searched++

	start_eval := eval_pos(position, ply)

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

	if engine.getTotalNodesSearched() >= engine.timer.MaxNodeCount {
		engine.timer.ForceStop()
	}
	if (engine.getTotalNodesSearched() & 2047) == 0 {
		engine.timer.CheckIfTimeIsUp()
	}

	if engine.timer.IsStopped() {
		return 0
	}

	// Sort Moves
	// moves := score_moves_v2(get_q_moves(position), position.Board())
	moves := score_moves(
		get_q_moves(position),
		position.Board(),
		[2]*chess.Move{nil, nil},
		nil,
	)

	if len(moves) == 0 {
		return start_eval
	}

	for i := 0; i < len(moves); i++ {
		get_move(moves, i)
		move := moves[i].move

		new_eval := -engine.q_search(
			position.Update(move), ply+1, -beta, -alpha, max_depth,
		)

		if new_eval > alpha {
			if new_eval >= beta {
				return beta
			}

			alpha = new_eval
		}
	}

	return alpha
}
