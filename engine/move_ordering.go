package engine

import "github.com/Sidhant-Roymoulik/Light-Blue/chess"

// -----------------------------------------------------------------------------
// 		Constants
// -----------------------------------------------------------------------------

const (
	MvvLvaOffset          int = 10000 - 256
	PVMoveScore           int = 65
	FirstKillerMoveScore  int = 10
	SecondKillerMoveScore int = 20
)

// -----------------------------------------------------------------------------
// 		MVV-LVA Stuff
// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
// 		Q-Move Stuff
// -----------------------------------------------------------------------------

func is_q_move(move *chess.Move) bool {
	return move.HasTag(chess.Capture) ||
		move.HasTag(chess.EnPassant) ||
		move.HasTag(chess.Check) ||
		move.Promo() != chess.NoPieceType
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

// -----------------------------------------------------------------------------
// 		Best Move Picking
// -----------------------------------------------------------------------------

type scored_move struct {
	move *chess.Move
	eval int
}

func get_move(moves []scored_move, start int) {
	best_index := start
	best_eval := moves[best_index].eval
	for i := best_index; i < len(moves); i++ {
		if moves[i].eval > best_eval {
			best_index = i
			best_eval = moves[i].eval
		}
	}
	temp := moves[start]
	moves[start] = moves[best_index]
	moves[best_index] = temp
}

// Adding Killer Moves
func score_moves(moves []*chess.Move, board *chess.Board, killer_moves [2]*chess.Move, pv_move *chess.Move) []scored_move {
	scores := make([]scored_move, len(moves))
	for i := 0; i < len(moves); i++ {
		if pv_move != nil &&
			moves[i].S1() == pv_move.S1() &&
			moves[i].S2() == pv_move.S2() {
			scores[i] = scored_move{
				moves[i], MvvLvaOffset + PVMoveScore,
			}
		} else if moves[i].HasTag(chess.Capture) {
			scores[i] = scored_move{
				moves[i], MvvLvaOffset + MVV_LVA(moves[i], board),
			}
		} else if moves[i] == killer_moves[0] {
			scores[i] = scored_move{
				moves[i], MvvLvaOffset - FirstKillerMoveScore,
			}
		} else if moves[i] == killer_moves[1] {
			scores[i] = scored_move{
				moves[i], MvvLvaOffset - SecondKillerMoveScore,
			}
		} else {
			scores[i] = scored_move{moves[i], 0}
		}
	}
	return scores
}
