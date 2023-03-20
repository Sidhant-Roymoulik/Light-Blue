package engine

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/Sidhant-Roymoulik/Light-Blue/chess"
)

// -----------------------------------------------------------------------------
//	Search Parameters
// -----------------------------------------------------------------------------

const (
	Window int = 12

	StaticNullMovePruningBaseMargin int = 85
	NMR_Depth_Limit                 int = 2
	FutilityPruningDepthLimit       int = 8
	IID_Depth_Limit                 int = 4
	IID_Depth_Reduction             int = 2
)

var FutilityMargins = [9]int{
	0,
	100, // depth 1
	160, // depth 2
	220, // depth 3
	280, // depth 4
	340, // depth 5
	400, // depth 6
	460, // depth 7
	520, // depth 8
}

var LateMovePruningMargins = [6]int{
	0,
	8,  // depth 1
	12, // depth 2
	16, // depth 3
	20, // depth 4
	24, // depth 5
}

// -----------------------------------------------------------------------------
//	Engine Definition
// -----------------------------------------------------------------------------

func new_light_blue() Engine {
	return Engine{
		EngineClass{
			name:   name,
			author: author,
			upgrades: EngineUpgrades{
				iterative_deepening: true,
			},
		},
		0,
		time.Now(),
		EngineCounters{},
		TimeManager{},
		TransTable[SearchEntry]{},
		0,
		[1024]uint64{},
		0,
		0,
		[MAX_DEPTH][2]*chess.Move{},
		runtime.GOMAXPROCS(0),
	}
}

// -----------------------------------------------------------------------------
//	Search Functions
// -----------------------------------------------------------------------------

func (e *Engine) run(
	position *chess.Position,
) (best_eval int, best_move *chess.Move) {
	e.resetCounters()
	e.resetKillerMoves()

	e.Add_Zobrist_History(Zobrist.GenHash(position))

	pvLine := PVLine{}

	if e.upgrades.iterative_deepening {
		best_eval, best_move = e.iterative_deepening(position, &pvLine)
	} else {
		best_eval = e.aspiration_window(
			position, int(e.timer.MaxDepth), &pvLine,
		)
		best_move = pvLine.getPVMove()
	}

	e.prev_guess = best_eval
	e.Add_Zobrist_History(Zobrist.GenHash(position.Update(best_move)))

	return
}

// Iterative Deepening
func (e *Engine) iterative_deepening(
	position *chess.Position, pvLine *PVLine,
) (best_eval int, best_move *chess.Move) {
	e.start = time.Now()
	e.age ^= 1
	e.timer.Start()

	for depth := 1; depth <= MAX_DEPTH &&
		depth <= int(e.timer.MaxDepth) &&
		e.timer.MaxNodeCount > 0; depth++ {

		pvLine.clear()

		new_eval := e.aspiration_window(position, depth, pvLine)

		if e.timer.IsStopped() {
			break
		}

		best_eval = new_eval
		e.prev_guess = best_eval
		e.max_ply = depth

		best_move = pvLine.getPVMove()
		total_nodes := e.counters.nodes_searched + e.counters.q_nodes_searched
		total_time := time.Since(e.start).Milliseconds() + 1

		fmt.Printf(
			"info depth %d score %s nodes %d nps %d time %d pv %s\n",
			e.max_ply,
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

// Aspiration window adapted from https://github.com/Aryan1508/Bit-Genie
func (e *Engine) aspiration_window(
	position *chess.Position, max_depth int, pvLine *PVLine,
) (eval int) {

	alpha := -CHECKMATE_VALUE
	beta := CHECKMATE_VALUE
	delta := Window

	if max_depth > 1 {
		w := Window - Min(6, max_depth/3)
		alpha = Max(e.prev_guess-w, -CHECKMATE_VALUE)
		beta = Min(e.prev_guess+w, CHECKMATE_VALUE)
	}

	for {
		eval = e.pv_search(position, 0, max_depth, alpha, beta, pvLine, true)

		if eval <= alpha {
			beta = (alpha + beta) / 2
			alpha = Max(alpha-delta, -CHECKMATE_VALUE)
		} else if eval >= beta {
			beta = Min(beta+delta, CHECKMATE_VALUE)
		} else {
			break
		}

		delta *= 3
		delta /= 2
	}

	return eval
}

// Principal-Variation Search
func (e *Engine) pv_search(
	position *chess.Position,
	ply int,
	max_depth int,
	alpha int,
	beta int,
	pvLine *PVLine,
	do_null bool,
) (eval int) {
	e.counters.nodes_searched++

	if ply >= MAX_DEPTH {
		return eval_pos(position)
	}

	// Check if search is over
	if e.counters.nodes_searched+e.counters.q_nodes_searched >=
		e.timer.MaxNodeCount {
		e.timer.ForceStop()
	}
	if (e.counters.nodes_searched+e.counters.q_nodes_searched)&
		TIMER_CHECK == 0 {
		e.timer.CheckIfTimeIsUp()
	}

	if e.timer.IsStopped() {
		return 0
	}

	// Generate hash for position
	hash := Zobrist.GenHash(position)

	// Initialize variables
	childPVLine := PVLine{}
	isPVNode := beta-alpha != 1
	isRoot := ply == 0
	inCheck := position.InCheck()
	canFutilityPrune := false
	var tt_move *chess.Move = nil

	// Check Extension
	if inCheck {
		e.counters.check_extensions++
		max_depth++
	}

	depth := max_depth - ply

	// Start Q-Search
	if depth <= 0 {
		e.counters.nodes_searched--
		return e.q_search(position, max_depth, alpha, beta)
	}

	// Check for draw by repetition
	// Check for draw by 50-move rule but let mate-in-1 trump it
	possibleMateInOne := inCheck && depth == 1
	if !isRoot &&
		((position.HalfMoveClock() >= 100 && !possibleMateInOne) ||
			e.Is_Draw_By_Repetition(hash)) {
		entry := e.tt.Probe(hash)
		_, _, tt_move := entry.Get(
			hash,
			ply,
			depth,
			-math.MaxInt,
			math.MaxInt,
		)
		pvLine.update(tt_move, childPVLine)
		return 0
	}

	// Check for usable entry in transposition table
	entry := e.tt.Probe(hash)
	tt_eval, should_use, tt_move := entry.Get(
		hash, ply, depth, alpha, beta,
	)
	if !isRoot && should_use {
		e.counters.hashes_used++
		return tt_eval
	}

	if !inCheck && !isPVNode {
		// Static Eval Calculation for Pruning
		static_eval := eval_pos(position)

		// Static Move Pruning
		if abs(beta) < MATE_CUTOFF {
			eval_margin := StaticNullMovePruningBaseMargin * depth
			if static_eval-eval_margin >= beta {
				e.counters.smp_pruned++
				return static_eval - eval_margin
			}
		}

		// Null Move Pruning
		if do_null && depth >= NMR_Depth_Limit {
			R := 3 + depth/6
			eval := -e.pv_search(
				position.NullMove(),
				ply+1,
				Max(max_depth-R, ply+1),
				-beta,
				-beta+1,
				&childPVLine,
				false,
			)
			childPVLine.clear()
			if eval >= beta && abs(eval) < MATE_CUTOFF {
				e.counters.nmp_pruned++
				return beta
			}
		}

		// Razoring
		if depth <= 2 {
			if static_eval+FutilityMargins[depth]*3 < beta {
				eval := e.q_search(position, ply, alpha, beta)
				if eval < beta {
					e.counters.razor_pruned++
					return eval
				}
			}
		}

		// Futility Pruning
		if depth <= FutilityPruningDepthLimit &&
			alpha < MATE_CUTOFF &&
			beta < MATE_CUTOFF {
			canFutilityPrune = static_eval+FutilityMargins[depth] <= alpha
		}
	}

	// Internal Iterative Deepening
	// if depth > IID_Depth_Limit &&
	// 	(isPVNode || entry.GetFlag() == BetaFlag) &&
	// 	tt_move == nil {
	// 	e.pv_search(
	// 		position,
	// 		ply+1,
	// 		max_depth-IID_Depth_Reduction,
	// 		-beta,
	// 		-alpha,
	// 		&childPVLine,
	// 		do_null,
	// 	)
	// 	if len(childPVLine.Moves) > 0 {
	// 		e.counters.iid_move_found++
	// 		tt_move = childPVLine.getPVMove()
	// 		childPVLine.clear()
	// 	}
	// }

	// Sort Moves
	moves := score_moves(
		position.ValidMoves(),
		position.Board(),
		e.killer_moves[ply],
		tt_move,
	)

	// Initialize variables
	var best_move *chess.Move = nil
	var tt_flag = AlphaFlag

	// Loop through moves
	for i := 0; i < len(moves); i++ {
		// Pick move
		get_move(moves, i)
		move := moves[i].move

		// Late Move Pruning
		if !isPVNode && !inCheck && depth <= 5 &&
			i >= LateMovePruningMargins[depth] {
			if !move.HasTag(chess.Check) &&
				move.Promo() == chess.NoPieceType {
				e.counters.lmp_pruned++
				continue
			}
		}

		// Futility Pruning
		if canFutilityPrune &&
			i > 0 &&
			!move.HasTag(chess.Check) &&
			!move.HasTag(chess.Capture) &&
			!move.HasTag(chess.EnPassant) &&
			move.Promo() == chess.NoPieceType {
			e.counters.futility_pruned++
			continue
		}

		// Generate new position
		new_position := position.Update(move)

		// Add to move history
		e.Add_Zobrist_History(Zobrist.GenHash(new_position))

		new_eval := 0

		if i == 0 {
			// Principal-Variation Search
			new_eval = -e.pv_search(
				new_position,
				ply+1,
				max_depth,
				-beta,
				-alpha,
				&childPVLine,
				do_null,
			)
		} else {
			// Null-Window Search
			new_eval = -e.pv_search(
				new_position,
				ply+1,
				max_depth,
				-(alpha + 1),
				-alpha,
				&childPVLine,
				true,
			)
			if new_eval > alpha && new_eval < beta {
				// Principal-Variation Search
				new_eval = -e.pv_search(
					new_position,
					ply+1,
					max_depth,
					-beta,
					-new_eval,
					&childPVLine,
					do_null,
				)
			}
		}

		// Clear move from history
		e.Remove_Zobrist_History()

		if new_eval > alpha {
			best_move = move

			if new_eval >= beta { // Fail-hard beta-cutoff
				// Add killer move
				e.addKillerMove(move, ply)

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

	// If there are no moves return either checkmate or draw
	if len(moves) == 0 {
		if inCheck {
			return -CHECKMATE_VALUE + ply
		}
		return 0
	}

	// Save position to transposition table
	if !e.timer.IsStopped() {
		entry := e.tt.Store(
			hash, depth, e.age,
		)
		entry.Set(
			hash, alpha, best_move, ply, depth, tt_flag, e.age,
		)
	}

	return alpha
}

// Quiescence Search
func (e *Engine) q_search(
	position *chess.Position,
	depth int,
	alpha int,
	beta int,
) int {
	e.counters.q_nodes_searched++

	// Check if search is over
	if e.counters.nodes_searched+e.counters.q_nodes_searched >=
		e.timer.MaxNodeCount {
		e.timer.ForceStop()
	}
	if (e.counters.nodes_searched+e.counters.q_nodes_searched)&
		TIMER_CHECK == 0 {
		e.timer.CheckIfTimeIsUp()
	}

	if e.timer.IsStopped() {
		return 0
	}

	if depth <= 0 {
		return eval_pos(position)
	}

	eval := eval_pos(position)

	// Delta Pruning
	if eval >= beta {
		return beta
	}
	alpha = Max(alpha, eval)

	// Sort Moves
	moves := score_moves(
		get_q_moves(position),
		position.Board(),
		[2]*chess.Move{nil, nil},
		nil,
	)

	for i := 0; i < len(moves); i++ {
		get_move(moves, i)
		move := moves[i].move

		new_eval := -e.q_search(
			position.Update(move), depth-1, -beta, -alpha,
		)

		alpha = Max(alpha, new_eval)

		if alpha >= beta {
			return beta
		}
	}

	return alpha
}
