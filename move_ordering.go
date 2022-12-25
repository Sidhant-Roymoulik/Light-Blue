package main

import (
	"sort"

	"github.com/Sidhant-Roymoulik/chess"
)

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

func move_ordering_v1(position *chess.Position) []*chess.Move {
	moves := position.ValidMoves()
	sort.SliceStable(moves, func(i, j int) bool {
		return eval_move_v1(moves[i], position.Board()) < eval_move_v1(moves[j], position.Board())
	})
	return moves
}

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

func move_ordering_v2(position *chess.Position) []*chess.Move {
	moves := position.ValidMoves()
	sort.SliceStable(moves, func(i, j int) bool {
		return eval_move_v2(moves[i], position.Board()) < eval_move_v2(moves[j], position.Board())
	})
	return moves
}

func MVV_LVA(move *chess.Move, board *chess.Board) int {
	return mvv_lva[board.Piece(move.S2()).Type()][board.Piece(move.S1()).Type()]
}
