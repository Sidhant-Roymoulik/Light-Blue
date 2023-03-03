package main

import (
	"fmt"
	"math"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

// Upgrades over Version 4.1:
// 		Refactored

type light_blue struct {
	EngineClass
}

func new_light_blue() light_blue {
	return light_blue{
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
				move_ordering:                true,
				alphabeta:                    true,
				iterative_deepening:          true,
				q_search:                     true,
				delta_pruning:                true,
				transposition_table:          true,
				killer_moves:                 true,
				pvs:                          true,
				aspiration_window:            true,
				internal_iterative_deepening: true,
			},
			timer:             TimeManager{},
			tt:                TransTable[SearchEntry]{},
			age:               0,
			zobristHistory:    [1024]uint64{},
			zobristHistoryPly: 0,
			killer_moves:      [MAX_DEPTH][2]*chess.Move{},
		},
	}
}

func (engine *light_blue) run(position *chess.Position) (best_eval int, best_move *chess.Move) {
	engine.resetCounters()
	engine.resetKillerMoves()

	engine.Add_Zobrist_History(Zobrist.GenHash(position))

	pvLine := PVLine{}

	if engine.upgrades.iterative_deepening {
		best_eval, best_move = engine.iterative_deepening(position, &pvLine)
	} else {
		best_eval = engine.aspiration_window(
			position, int(engine.timer.MaxDepth), &pvLine,
		)
		best_move = pvLine.getPVMove()
	}

	engine.prev_guess = best_eval
	engine.Add_Zobrist_History(Zobrist.GenHash(position.Update(best_move)))

	// print(engine.zobristHistory)
	// print(engine.killer_moves)

	return
}

func (engine *light_blue) iterative_deepening(
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
		total_time := time.Since(engine.start).Milliseconds() + 1

		fmt.Printf(
			"info depth %d score %s nodes %d nps %d time %d pv %s\n",
			engine.max_ply,
			getMateOrCPScore(best_eval),
			engine.getTotalNodesSearched(),
			int64(engine.getTotalNodesSearched()*1000)/total_time,
			total_time,
			pvLine,
		)

		if best_eval >= MATE_CUTOFF {
			break
		}
	}

	return best_eval, best_move
}

func (engine *light_blue) aspiration_window(
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
			print("Aspiration tight fail low")
		}
		alpha = engine.prev_guess - WINDOW_VALUE
		eval = engine.minimax_start(position, alpha, beta, max_depth, pvLine)
	} else if eval >= beta {
		if DEBUG {
			print("Aspiration tight fail high")
		}
		beta = engine.prev_guess + WINDOW_VALUE
		eval = engine.minimax_start(position, alpha, beta, max_depth, pvLine)
	}

	if eval <= alpha || eval >= beta {
		if DEBUG {
			print("Aspiration loose fail")
		}
		eval = engine.minimax_start(
			position, -math.MaxInt, math.MaxInt, max_depth, pvLine,
		)
	}

	return eval
}

func (engine *light_blue) minimax_start(
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

		// Bug with move gen where engine tries to castle when in check
		if position.InCheck() &&
			(move.HasTag(chess.KingSideCastle) ||
				move.HasTag(chess.QueenSideCastle)) {
			continue
		}

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
			&childPVLine,
			true,
		)

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

func (engine *light_blue) pv_search(
	position *chess.Position,
	ply int,
	alpha int,
	beta int,
	max_depth int,
	pvLine *PVLine,
	do_null bool,
) (eval int) {
	engine.counters.nodes_searched++

	if ply >= MAX_DEPTH {
		return eval_pos(position, ply)
	}

	// Check if search is over
	if engine.getTotalNodesSearched() >= engine.timer.MaxNodeCount {
		engine.timer.ForceStop()
	}
	if (engine.getTotalNodesSearched() & 1023) == 0 {
		engine.timer.CheckIfTimeIsUp()
	}

	if engine.timer.IsStopped() {
		return 0
	}

	// Generate hash for position
	hash := Zobrist.GenHash(position)

	// Initialize variables
	childPVLine := PVLine{}
	isPVNode := beta-alpha != 1
	inCheck := position.InCheck()
	canFutilityPrune := false

	// Check Extension
	if inCheck {
		max_depth++
	}

	depth := max_depth - ply

	// Start Q-Search
	if ply >= max_depth {
		engine.counters.nodes_searched--
		return engine.q_search(position, 0, alpha, beta, max_depth)
	}

	// Check for draw by repetition
	// Check for draw by 50-move rule but let mate-in-1 trump it
	possibleMateInOne := inCheck && depth == 1
	if (position.HalfMoveClock() >= 100 && !possibleMateInOne) ||
		engine.Is_Draw_By_Repetition(hash) {
		return 0
	}

	// Check for usable entry in transposition table
	entry := engine.tt.Probe(hash)
	tt_eval, should_use, tt_move := entry.Get(
		hash, ply, depth, alpha, beta,
	)
	if should_use {
		engine.counters.hashes_used++
		return tt_eval
	}

	// Static Move Pruning
	if !inCheck && !isPVNode && abs(beta) < MATE_CUTOFF {
		static_eval := eval_pos(position, ply)
		eval_margin := StaticNullMovePruningBaseMargin * depth
		if static_eval-eval_margin >= beta {
			return static_eval - eval_margin
		}
	}

	// Null Move Pruning
	if do_null && !inCheck && !isPVNode && depth >= NMR_Depth_Limit {
		R := 3 + depth/6
		eval := -engine.pv_search(
			position.NullMove(),
			ply+1,
			-beta,
			-beta+1,
			max_depth-R,
			&childPVLine,
			false,
		)

		childPVLine.clear()

		if eval >= beta && abs(eval) < MATE_CUTOFF {
			return beta
		}
	}

	// Razoring
	if depth <= 2 && !inCheck && !isPVNode {
		static_eval := eval_pos(position, ply)
		if static_eval+FutilityMargins[depth]*3 < alpha {
			eval := engine.q_search(position, 0, alpha, beta, ply)
			if eval < alpha {
				return alpha
			}
		}
	}

	// Futility Pruning
	if depth <= FutilityPruningDepthLimit &&
		!inCheck &&
		!isPVNode &&
		alpha < MATE_CUTOFF &&
		beta < MATE_CUTOFF {
		static_eval := eval_pos(position, ply)
		eval_margin := FutilityMargins[depth]

		canFutilityPrune = static_eval+eval_margin <= alpha
	}

	// Internal Iterative Deepening
	if depth > IID_Depth_Limit &&
		(isPVNode || entry.GetFlag() == BetaFlag) &&
		tt_move == nil {
		engine.pv_search(
			position,
			ply+1,
			-beta,
			-alpha,
			max_depth-IID_Depth_Reduction,
			&childPVLine,
			do_null,
		)
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

	// If there are no moves return either checkmate or draw
	if len(moves) == 0 {
		if inCheck {
			return -CHECKMATE_VALUE + ply
		}
		return 0
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

		// Bug with move gen where engine tries to castle when in check
		if inCheck &&
			(move.HasTag(chess.KingSideCastle) ||
				move.HasTag(chess.QueenSideCastle)) {
			continue
		}

		// Futility Pruning
		if canFutilityPrune &&
			i > 0 &&
			!is_q_move(move) {
			continue
		}

		// Generate new position
		new_position := position.Update(move)

		// Add to move history
		engine.Add_Zobrist_History(Zobrist.GenHash(new_position))

		new_eval := 0

		if bSearchPv {
			// Principal-Variation Search
			new_eval = -engine.pv_search(
				new_position,
				ply+1,
				-beta,
				-alpha,
				max_depth,
				&childPVLine,
				do_null,
			)
		} else {
			// Null-Window Search
			new_eval = -engine.pv_search(
				new_position,
				ply+1,
				-alpha-1,
				-alpha,
				max_depth,
				&childPVLine,
				do_null,
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
					do_null,
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
		entry := engine.tt.Store(
			hash, depth, engine.age,
		)
		entry.Set(
			hash, alpha, best_move, ply, depth, tt_flag, engine.age,
		)

		engine.counters.hashes_written++
	}

	return alpha
}

func (engine *light_blue) q_search(
	position *chess.Position,
	ply int,
	alpha int,
	beta int,
	max_depth int,
) int {
	engine.counters.q_nodes_searched++

	eval := eval_pos(position, ply+max_depth)
	inCheck := ply <= 2 && position.InCheck()

	// Delta Pruning
	if !inCheck && eval >= beta {
		return beta
	}
	if eval >= alpha {
		alpha = eval
	}

	if ply >= max_depth {
		return eval
	}

	// Check if search is over
	if engine.getTotalNodesSearched() >= engine.timer.MaxNodeCount {
		engine.timer.ForceStop()
	}
	if (engine.getTotalNodesSearched() & 1023) == 0 {
		engine.timer.CheckIfTimeIsUp()
	}

	if engine.timer.IsStopped() {
		return 0
	}

	// Sort Moves
	var moves []scored_move = nil

	if inCheck {
		moves = score_moves(
			position.ValidMoves(),
			position.Board(),
			[2]*chess.Move{nil, nil},
			nil,
		)
	} else {
		moves = score_moves(
			get_q_moves(position),
			position.Board(),
			[2]*chess.Move{nil, nil},
			nil,
		)
	}

	for i := 0; i < len(moves); i++ {
		get_move(moves, i)
		move := moves[i].move

		new_eval := -engine.q_search(
			position.Update(move), ply+1, -beta, -alpha, max_depth,
		)

		if new_eval > eval {
			eval = new_eval
		}

		if new_eval >= beta {
			break
		}

		if new_eval > alpha {
			alpha = new_eval
		}
	}

	return eval
}
