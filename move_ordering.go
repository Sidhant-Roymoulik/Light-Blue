package main

import (
	"sort"

	"github.com/Sidhant-Roymoulik/chess"
)

//	--------------------------------------------------------------------------------------
// 		MVV-LVA Stuff
//	--------------------------------------------------------------------------------------

// x axis is attacker: no piece, king, queen, rook, bishop, knight, pawn
var mvv_lva = [7][7]int{
	{0, 0, 0, 0, 0, 0, 0},       // victim no piece
	{0, 0, 0, 0, 0, 0, 0},       // victim king
	{0, 50, 51, 52, 53, 54, 55}, // victim queen
	{0, 40, 41, 42, 43, 44, 45}, // victim rook
	{0, 30, 31, 32, 33, 34, 35}, // victim bishop
	{0, 20, 21, 22, 23, 24, 25}, // victim knight
	{0, 10, 11, 12, 13, 14, 15}, // victim pawn
}

func MVV_LVA(move *chess.Move, board *chess.Board) int {
	return mvv_lva[board.Piece(move.S2()).Type()][board.Piece(move.S1()).Type()]
}

//	--------------------------------------------------------------------------------------
// 		Move Evaluation Functions
//	--------------------------------------------------------------------------------------

func eval_move_v1(move *chess.Move, board *chess.Board) int {
	if move.HasTag(chess.MoveTag(chess.Checkmate)) {
		return CHECKMATE_VALUE
	}

	delta := 0

	if move.HasTag(chess.MoveTag(chess.Capture)) {
		delta += PVM[chess.PieceType(board.Piece(move.S2()))]
	}
	if move.HasTag(chess.MoveTag(chess.Check)) {
		delta += 25
	}
	if move.HasTag(chess.MoveTag(chess.KingSideCastle)) {
		delta += 50
	}
	if move.HasTag(chess.MoveTag(chess.QueenSideCastle)) {
		delta += 40
	}

	return delta
}

// Adding MVV_LVA
func eval_move_v2(move *chess.Move, board *chess.Board) int {
	if move.HasTag(chess.MoveTag(chess.Checkmate)) {
		return CHECKMATE_VALUE
	}

	if move.HasTag(chess.MoveTag(chess.KingSideCastle)) {
		return 50
	}
	if move.HasTag(chess.MoveTag(chess.QueenSideCastle)) {
		return 40
	}

	return MVV_LVA(move, board)
}

//	--------------------------------------------------------------------------------------
// 		Q-Move Stuff
//	--------------------------------------------------------------------------------------

func is_q_move(move *chess.Move) bool {
	if move.HasTag(chess.Capture) {
		return true
	}
	if move.HasTag(chess.EnPassant) {
		return true
	}
	// if move.HasTag(chess.MoveTag(chess.Checkmate)) {
	// 	return true
	// }
	// if move.HasTag(chess.Check) {
	// 	return true
	// }
	if move.Promo() != chess.NoPieceType {
		return true
	}
	return false
}

func get_q_moves(position *chess.Position) []*chess.Move {
	moves := position.ValidMoves()
	n := 0
	for _, move := range moves {
		if is_q_move(move) {
			moves[n] = move
			n++
		}
	}
	return moves[:n]
}

//	--------------------------------------------------------------------------------------
// 		Very Inefficient Move Ordering
//	--------------------------------------------------------------------------------------

func move_ordering_v1(position *chess.Position) []*chess.Move {
	moves := position.ValidMoves()
	sort.SliceStable(moves, func(i, j int) bool {
		return eval_move_v1(moves[i], position.Board()) < eval_move_v1(moves[j], position.Board())
	})
	return moves
}

func move_ordering_v2(position *chess.Position) []*chess.Move {
	moves := position.ValidMoves()
	sort.SliceStable(moves, func(i, j int) bool {
		return eval_move_v2(moves[i], position.Board()) < eval_move_v2(moves[j], position.Board())
	})
	return moves
}

//	--------------------------------------------------------------------------------------
// 		Less Efficient Move Picking
//	--------------------------------------------------------------------------------------

func get_move_v1(position *chess.Position, moves []*chess.Move, start int) *chess.Move {
	best_eval := eval_move_v2(moves[start], position.Board())
	for i := start + 1; i < len(moves); i++ {
		new_eval := eval_move_v2(moves[i], position.Board())
		if new_eval > best_eval {
			moves[start], moves[i] = moves[i], moves[start]
			best_eval = new_eval
		}
	}
	return moves[start]
}

//	--------------------------------------------------------------------------------------
// 		Slightly Less Efficient Move Picking
//	--------------------------------------------------------------------------------------

type scored_move struct {
	move *chess.Move
	eval int
}

func score_moves(moves []*chess.Move, board *chess.Board) []scored_move {
	scores := make([]scored_move, len(moves))
	for i := 0; i < len(moves); i++ {
		scores[i] = scored_move{moves[i], eval_move_v2(moves[i], board)}
	}
	return scores
}

func get_move_v2(moves []scored_move, start int) *chess.Move {
	best_eval := moves[start].eval
	for i := start + 1; i < len(moves); i++ {
		if moves[i].eval > best_eval {
			moves[start], moves[i] = moves[i], moves[start]
			best_eval = moves[i].eval
		}
	}
	return moves[start].move
}

//	--------------------------------------------------------------------------------------
// 		Efficient Move Picking
//	--------------------------------------------------------------------------------------

func score_moves_v2(moves []*chess.Move, board *chess.Board) []scored_move {
	scores := make([]scored_move, len(moves))
	for i := 0; i < len(moves); i++ {
		scores[i] = scored_move{moves[i], eval_move_v2(moves[i], board)}
		if scores[i].eval > scores[0].eval { // Use first guaranteed iteration to sort first move
			scores[i], scores[0] = scores[0], scores[i]
		}
	}
	return scores
}

func get_move_v3(moves []scored_move, start int) *chess.Move {
	if start == 0 { //	First move is already sorted
		return moves[0].move
	}
	best_eval := moves[start].eval
	for i := start + 1; i < len(moves); i++ {
		if moves[i].eval > best_eval {
			moves[start], moves[i] = moves[i], moves[start]
			best_eval = moves[i].eval
		}
	}
	return moves[start].move
}
