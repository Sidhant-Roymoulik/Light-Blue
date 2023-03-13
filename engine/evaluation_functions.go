package engine

import (
	"encoding/binary"

	"github.com/Sidhant-Roymoulik/Light-Blue/chess"
)

var score_mg = [2]int{}
var score_eg = [2]int{}

// -----------------------------------------------------------------------------
// 		Bonuses + Penalties
// -----------------------------------------------------------------------------

const (
	IsolatedPawnPenatlyMG int = 17
	IsolatedPawnPenatlyEG int = 6

	DoubledPawnPenatlyMG int = 1
	DoubledPawnPenatlyEG int = 16

	KnightOnOutpostBonusMG int = 27
	KnightOnOutpostBonusEG int = 18

	BishopOutPostBonusMG int = 10
	BishopOutPostBonusEG int = 14

	RookOrQueenOnSeventhBonusEG int = 23

	RookOnOpenFileBonusMG int = 23

	BishopPairBonusMG int = 22
	BishopPairBonusEG int = 30

	SemiOpenFileNextToKingPenalty int = 4

	TempoBonusMG int = 14

	DrawishScaleFactor int = 16

	// -------------------------------------------------------------------------
	// 		Tapered Evaluation Values
	// -------------------------------------------------------------------------

	PawnPhase   int = 0
	KnightPhase int = 1
	BishopPhase int = 1
	RookPhase   int = 2
	QueenPhase  int = 4
	TotalPhase  int = PawnPhase*16 +
		KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2
)

// -----------------------------------------------------------------------------
// 		Bitboards
// -----------------------------------------------------------------------------

var DoubledPawnMasks [2][64]chess.Bitboard
var IsolatedPawnMasks [8]chess.Bitboard
var PassedPawnMasks [2][64]chess.Bitboard
var OutpostMasks [2][64]chess.Bitboard

type KingZone struct {
	OuterRing chess.Bitboard
	InnerRing chess.Bitboard
}

// -----------------------------------------------------------------------------
// 		King Safety Stuff
// -----------------------------------------------------------------------------

var KingZones [2]KingZone
var KingZonesMasks [64]KingZone
var KingAttackPoints [2]int
var KingAttackers [2]int

var OuterRingAttackPoints = map[chess.PieceType]int{
	chess.Queen:  1,
	chess.Rook:   1,
	chess.Bishop: 0,
	chess.Knight: 1,
}

var InnerRingAttackPoints = map[chess.PieceType]int{
	chess.Queen:  2,
	chess.Rook:   3,
	chess.Bishop: 4,
	chess.Knight: 3,
}

// -----------------------------------------------------------------------------
// 		Piece Values
// -----------------------------------------------------------------------------

// piece value map
var PVM_MG = map[chess.PieceType]int{
	chess.Queen:  921,
	chess.Rook:   441,
	chess.Bishop: 346,
	chess.Knight: 333,
	chess.Pawn:   84,
}

var PVM_EG = map[chess.PieceType]int{
	chess.Queen:  886,
	chess.Rook:   478,
	chess.Bishop: 268,
	chess.Knight: 244,
	chess.Pawn:   106,
}

var Mobility_MG = map[chess.PieceType]int{
	chess.Queen:  0,
	chess.Rook:   3,
	chess.Bishop: 3,
	chess.Knight: 5,
	chess.Pawn:   0,
}

var Mobility_EG = map[chess.PieceType]int{
	chess.Queen:  6,
	chess.Rook:   2,
	chess.Bishop: 3,
	chess.Knight: 2,
	chess.Pawn:   0,
}

// -----------------------------------------------------------------------------
// 		Piece Square Table Stuff
// -----------------------------------------------------------------------------

var FLIP = [2][64]int{
	{
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	},
	{
		0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 58, 59, 60, 61, 62, 63,
	},
}

var FlipRank = [2][8]chess.Rank{
	{
		chess.Rank1,
		chess.Rank2,
		chess.Rank3,
		chess.Rank4,
		chess.Rank5,
		chess.Rank6,
		chess.Rank7,
		chess.Rank8,
	},
	{
		chess.Rank8,
		chess.Rank7,
		chess.Rank6,
		chess.Rank5,
		chess.Rank4,
		chess.Rank3,
		chess.Rank2,
		chess.Rank1,
	},
}

var PST_MG = map[chess.PieceType][]int{
	chess.Pawn: {
		0, 0, 0, 0, 0, 0, 0, 0,
		45, 52, 42, 43, 28, 34, 19, 9,
		-14, -3, 7, 14, 35, 50, 15, -6,
		-27, -6, -8, 13, 16, 4, -3, -25,
		-32, -28, -7, 5, 7, -1, -15, -30,
		-29, -25, -12, -12, -1, -5, 6, -17,
		-34, -23, -27, -18, -14, 10, 13, -22,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-43, -11, -8, -5, 1, -20, -4, -22,
		-31, -22, 19, 7, 5, 13, -8, -11,
		-21, 21, 8, 16, 36, 33, 19, 6,
		-6, 2, 0, 23, 8, 27, 4, 14,
		-3, 10, 12, 8, 16, 10, 19, 1,
		-19, -4, 3, 7, 22, 12, 15, -11,
		-21, -20, -9, 8, 9, 11, -5, 0,
		-19, -13, -20, -14, -2, 3, -11, -8,
	},
	chess.Bishop: {
		-13, 0, -17, -8, -7, -5, -2, -3,
		-21, 0, -16, -10, 4, 1, -6, -41,
		-23, 6, 10, 8, 8, 26, 0, -10,
		-15, -4, 2, 22, 9, 10, -1, -16,
		0, 10, -2, 15, 17, -7, -1, 13,
		-2, 16, 13, 0, 5, 16, 14, 0,
		8, 11, 12, 3, 11, 23, 27, 3,
		-26, 3, -3, -1, 10, -5, -7, -15,
	},
	chess.Rook: {
		3, 1, 0, 7, 7, -1, 0, 0,
		-6, -9, 7, 7, 7, 5, -4, -1,
		-12, 11, 0, 17, -2, 12, 23, -1,
		-17, -9, 4, 0, 3, 15, -1, -2,
		-24, -16, -16, -4, -1, -14, 2, -20,
		-30, -15, -6, -3, 0, 2, 2, -15,
		-25, -6, -6, 5, 8, 6, 8, -46,
		-3, 1, 6, 15, 17, 14, -13, -2,
	},
	chess.Queen: {
		-10, 0, 0, 0, 10, 9, 5, 7,
		-19, -35, -5, 2, -9, 7, 1, 15,
		-10, -7, -4, -9, 15, 29, 24, 22,
		-14, -14, -15, -11, -1, -5, 3, -6,
		-8, -20, -8, -5, -4, -2, 2, -2,
		-13, 5, 2, 1, -1, 8, 4, 2,
		-20, 0, 10, 16, 16, 16, -6, 6,
		-3, -1, 7, 19, 5, -10, -9, -17,
	},
	chess.King: {
		-3, 0, 2, 0, 0, 0, 1, -1,
		1, 4, 0, 7, 4, 2, 3, -2,
		2, 4, 7, 4, 4, 14, 12, 0,
		0, 2, 6, 0, 0, 2, 6, -9,
		-8, 5, 0, -8, -10, -10, -9, -23,
		-3, 5, 1, -8, -12, -12, 8, -24,
		6, 13, 0, -40, -23, -1, 25, 19,
		-28, 29, 17, -53, 2, -25, 34, 15,
	},
}

var PST_EG = map[chess.PieceType][]int{
	chess.Pawn: {
		0, 0, 0, 0, 0, 0, 0, 0,
		77, 74, 63, 53, 59, 60, 72, 77,
		17, 11, 11, 11, 11, -6, 14, 8,
		-3, -14, -18, -31, -29, -25, -20, -18,
		-12, -14, -24, -31, -29, -28, -27, -28,
		-22, -20, -25, -20, -21, -24, -34, -34,
		-16, -22, -11, -19, -13, -23, -32, -34,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-36, -16, -7, -14, -4, -20, -20, -29,
		-17, 2, -7, 14, 2, -7, -9, -19,
		-13, -7, 14, 12, 4, 6, 0, -13,
		-5, 8, 24, 18, 22, 15, 11, -4,
		-3, 4, 20, 30, 22, 25, 15, -2,
		-7, 1, 3, 19, 10, -2, -4, -4,
		-10, -2, -1, 0, 6, -8, -3, -13,
		-12, -28, -8, 1, -5, -12, -27, -12,
	},
	chess.Bishop: {
		-9, -5, -9, -5, -2, -4, -5, -8,
		0, 2, 8, -7, 1, 0, -2, -8,
		8, 0, 0, 1, 0, 1, 5, 6,
		0, 7, 7, 8, 3, 5, 2, 6,
		-1, 0, 12, 8, 0, 6, 0, -5,
		0, 0, 3, 6, 8, -1, 0, -1,
		-6, -12, -7, 0, 0, -8, -9, -13,
		-11, 0, -6, 0, -3, -4, -5, -9,
	},
	chess.Rook: {
		8, 9, 11, 13, 13, 12, 13, 9,
		3, 5, 1, 0, -1, 0, 6, 2,
		9, 5, 7, 2, 2, 1, 0, 0,
		3, 3, 6, 0, 0, 0, 0, 4,
		5, 4, 9, 0, -3, -2, -6, -2,
		0, 0, -6, -5, -9, -14, -7, -12,
		-2, -5, -1, -7, -9, -11, -13, -1,
		-7, -3, 0, -8, -13, -12, -4, -24,
	},
	chess.Queen: {
		-12, 4, 8, 4, 10, 9, 3, 6,
		-17, -7, -1, 7, 3, 6, 1, 0,
		-5, -1, -4, 12, 14, 20, 12, 14,
		-2, 2, 2, 9, 13, 7, 18, 22,
		-9, 3, 1, 15, 5, 10, 12, 10,
		-6, -20, 0, -15, 0, -1, 10, 7,
		-6, -14, -31, -27, -19, -12, -11, -4,
		-12, -22, -19, -30, -8, -13, -6, -15,
	},
	chess.King: {
		-15, -11, -11, -6, -2, 3, 4, -9,
		-9, 14, 11, 13, 13, 28, 19, 1,
		-1, 18, 19, 15, 16, 35, 34, 4,
		-12, 14, 21, 25, 19, 25, 18, -5,
		-23, -6, 14, 21, 20, 18, 5, -16,
		-21, -6, 5, 13, 15, 9, -2, -12,
		-27, -10, 2, 9, 9, 1, -12, -26,
		-43, -34, -20, -5, -26, -9, -35, -55,
	},
}

var PassedPawn_MG = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	45, 52, 42, 43, 28, 34, 19, 9,
	48, 43, 43, 30, 24, 31, 12, 2,
	28, 17, 13, 10, 10, 19, 6, 1,
	14, 0, -9, -7, -13, -7, 9, 16,
	5, 3, -3, -14, -3, 10, 13, 19,
	8, 9, 2, -8, -3, 8, 16, 9,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawn_EG = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	77, 74, 63, 53, 59, 60, 72, 77,
	91, 83, 66, 40, 30, 61, 67, 84,
	55, 52, 42, 35, 30, 34, 56, 52,
	29, 26, 21, 18, 17, 19, 34, 30,
	8, 6, 5, 1, 1, -1, 14, 7,
	2, 3, -4, 0, -2, -1, 7, 6,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// -----------------------------------------------------------------------------
// 		Position Evaluation Function
// -----------------------------------------------------------------------------

// Best Evaluation
func eval_pos(position *chess.Position, ply int) int {
	data, err := position.Board().MarshalBinary()
	if err != nil {
		print(err)
		return 0
	}

	var pieces = [2]map[chess.PieceType]chess.Bitboard{
		{
			chess.King:   chess.Bitboard(binary.BigEndian.Uint64(data[:8])),
			chess.Queen:  chess.Bitboard(binary.BigEndian.Uint64(data[8:16])),
			chess.Rook:   chess.Bitboard(binary.BigEndian.Uint64(data[16:24])),
			chess.Bishop: chess.Bitboard(binary.BigEndian.Uint64(data[24:32])),
			chess.Knight: chess.Bitboard(binary.BigEndian.Uint64(data[32:40])),
			chess.Pawn:   chess.Bitboard(binary.BigEndian.Uint64(data[40:48])),
		},
		{
			chess.King:   chess.Bitboard(binary.BigEndian.Uint64(data[48:56])),
			chess.Queen:  chess.Bitboard(binary.BigEndian.Uint64(data[56:64])),
			chess.Rook:   chess.Bitboard(binary.BigEndian.Uint64(data[64:72])),
			chess.Bishop: chess.Bitboard(binary.BigEndian.Uint64(data[72:80])),
			chess.Knight: chess.Bitboard(binary.BigEndian.Uint64(data[80:88])),
			chess.Pawn:   chess.Bitboard(binary.BigEndian.Uint64(data[88:96])),
		},
	}

	// Draw by Insufficient Material
	if is_draw(pieces) {
		return 0
	}

	var sides = [2]chess.Bitboard{
		chess.White: 0,
		chess.Black: 0,
	}

	for i := 0; i < 12; i++ {
		if i < 6 {
			sides[chess.White] |= chess.Bitboard(
				binary.BigEndian.Uint64(data[i*8 : i*8+8]),
			)
		} else {
			sides[chess.Black] |= chess.Bitboard(
				binary.BigEndian.Uint64(data[i*8 : i*8+8]),
			)
		}
	}

	score_mg = [2]int{0, 0}
	score_eg = [2]int{0, 0}

	turn := position.Turn()

	squares := position.Board().SquareMap()
	all_bb := sides[chess.White] | sides[chess.Black]

	for all_bb != 0 {
		square := all_bb.PopBit()
		piece := squares[chess.Square(square)]
		color := piece.Color()
		if color == chess.NoColor {
			print(square)
		}

		score_mg[color] += PVM_MG[piece.Type()]
		score_mg[color] += PST_MG[piece.Type()][FLIP[color][square]]

		score_eg[color] += PVM_EG[piece.Type()]
		score_eg[color] += PST_EG[piece.Type()][FLIP[color][square]]

		switch piece.Type() {
		case chess.Pawn:
			ally := pieces[color][chess.Pawn]
			enemy := pieces[color^1][chess.Pawn]

			// Isolated Pawns
			if IsolatedPawnMasks[FileOf(square)]&ally != 0 {
				score_mg[color] -= IsolatedPawnPenatlyMG
				score_eg[color] -= IsolatedPawnPenatlyEG
			}

			// Doubled Pawns
			if DoubledPawnMasks[color][square]&ally != 0 {
				score_mg[color] -= DoubledPawnPenatlyMG
				score_eg[color] -= DoubledPawnPenatlyEG
			} else {
				// Check for Passed Pawn only if not Doubled
				if PassedPawnMasks[color][square]&enemy == 0 {
					score_mg[color] += PassedPawn_MG[FLIP[color][square]]
					score_eg[color] += PassedPawn_EG[FLIP[color][square]]
				}
			}

		case chess.Knight:
			ally := pieces[color][chess.Pawn]
			enemy := pieces[color^1][chess.Pawn]

			// Check for Outposts
			if OutpostMasks[color][square]&enemy == 0 &&
				PawnAttacks[color][square]&ally != 0 &&
				FlipRank[color][RankOf(square)] >= chess.Rank5 {
				score_mg[color] += KnightOnOutpostBonusMG
				score_eg[color] += KnightOnOutpostBonusEG
			}

			moves := chess.BBKnightMoves[square] & ^sides[color]

			// Mobility Bonus
			safe_moves := moves

			for enemy != 0 {
				square = enemy.PopBit()
				safe_moves &= ^PawnAttacks[color^1][square]
			}

			mobility := safe_moves.CountBits()
			score_mg[color] += (mobility - 4) * Mobility_MG[chess.Knight]
			score_eg[color] += (mobility - 4) * Mobility_EG[chess.Knight]

			// King Attacks
			outer_ring_attacks := moves & KingZones[color^1].OuterRing
			inner_ring_attacks := moves & KingZones[color^1].InnerRing

			if outer_ring_attacks > 0 || inner_ring_attacks > 0 {
				KingAttackers[color]++
				KingAttackPoints[color] += outer_ring_attacks.CountBits() *
					OuterRingAttackPoints[chess.Knight]
				KingAttackPoints[color] += inner_ring_attacks.CountBits() *
					InnerRingAttackPoints[chess.Knight]
			}

		case chess.Bishop:
			ally := pieces[color][chess.Pawn]
			enemy := pieces[color^1][chess.Pawn]

			// Check for Outposts
			if OutpostMasks[color][square]&enemy == 0 &&
				PawnAttacks[color][square]&ally != 0 &&
				FlipRank[color][RankOf(square)] >= chess.Rank5 {
				score_mg[color] += BishopOutPostBonusMG
				score_eg[color] += BishopOutPostBonusEG
			}

			// Mobility Bonus
			full_bb := sides[color] | sides[color^1]
			moves := chess.DiaAttack(full_bb, chess.Square(square)) & ^sides[color]

			mobility := moves.CountBits()
			score_mg[color] += (mobility - 7) * Mobility_MG[chess.Bishop]
			score_eg[color] += (mobility - 7) * Mobility_EG[chess.Bishop]

			// King Attacks
			outer_ring_attacks := moves & KingZones[color^1].OuterRing
			inner_ring_attacks := moves & KingZones[color^1].InnerRing

			if outer_ring_attacks > 0 || inner_ring_attacks > 0 {
				KingAttackers[color]++
				KingAttackPoints[color] += outer_ring_attacks.CountBits() *
					OuterRingAttackPoints[chess.Bishop]
				KingAttackPoints[color] += inner_ring_attacks.CountBits() *
					InnerRingAttackPoints[chess.Bishop]
			}

		case chess.Rook:
			// Seventh Rank Bonus
			enemy_king := pieces[color^1][chess.King].Msb()
			if FlipRank[color][RankOf(square)] == chess.Rank7 &&
				FlipRank[color][RankOf(enemy_king)] >= chess.Rank7 {
				score_eg[color] += RookOrQueenOnSeventhBonusEG
			}

			// Open File Bonus
			pawns := pieces[color][chess.Pawn] | pieces[color^1][chess.Pawn]
			if MaskFile[FileOf(square)]&pawns == 0 {
				score_mg[color] += RookOnOpenFileBonusMG
			}

			// Mobility Bonus
			full_bb := sides[color] | sides[color^1]
			moves := chess.HvAttack(full_bb, chess.Square(square)) & ^sides[color]

			mobility := moves.CountBits()
			score_mg[color] += (mobility - 7) * Mobility_MG[chess.Rook]
			score_eg[color] += (mobility - 7) * Mobility_EG[chess.Rook]

			// King Attacks
			outer_ring_attacks := moves & KingZones[color^1].OuterRing
			inner_ring_attacks := moves & KingZones[color^1].InnerRing

			if outer_ring_attacks > 0 || inner_ring_attacks > 0 {
				KingAttackers[color]++
				KingAttackPoints[color] += outer_ring_attacks.CountBits() *
					OuterRingAttackPoints[chess.Rook]
				KingAttackPoints[color] += inner_ring_attacks.CountBits() *
					InnerRingAttackPoints[chess.Rook]
			}

		case chess.Queen:
			// Seventh Rank Bonus
			enemy_king := pieces[color^1][chess.King].Msb()
			if FlipRank[color][RankOf(square)] == chess.Rank7 &&
				FlipRank[color][RankOf(enemy_king)] >= chess.Rank7 {
				score_eg[color] += RookOrQueenOnSeventhBonusEG
			}

			// Mobility Bonus
			full_bb := sides[color] | sides[color^1]
			moves := (chess.DiaAttack(full_bb, chess.Square(square)) |
				chess.HvAttack(full_bb, chess.Square(square))) & ^sides[color]

			mobility := moves.CountBits()
			score_mg[color] += (mobility - 14) * Mobility_MG[chess.Queen]
			score_eg[color] += (mobility - 14) * Mobility_EG[chess.Queen]

			// King Attacks
			outer_ring_attacks := moves & KingZones[color^1].OuterRing
			inner_ring_attacks := moves & KingZones[color^1].InnerRing

			if outer_ring_attacks > 0 || inner_ring_attacks > 0 {
				KingAttackers[color]++
				KingAttackPoints[color] += outer_ring_attacks.CountBits() *
					OuterRingAttackPoints[chess.Queen]
				KingAttackPoints[color] += inner_ring_attacks.CountBits() *
					InnerRingAttackPoints[chess.Queen]
			}
		}

	}

	if pieces[chess.White][chess.Bishop].CountBits() == 2 {
		score_mg[chess.White] += BishopPairBonusMG
		score_eg[chess.White] += BishopPairBonusEG
	}
	if pieces[chess.Black][chess.Bishop].CountBits() == 2 {
		score_mg[chess.Black] += BishopPairBonusMG
		score_eg[chess.Black] += BishopPairBonusEG
	}

	evalKing(
		&pieces, chess.White, pieces[chess.White][chess.King].Msb(),
	)
	evalKing(
		&pieces, chess.Black, pieces[chess.Black][chess.King].Msb(),
	)

	score_mg[turn] += TempoBonusMG

	// Tapered Evaluation
	eval_mg := score_mg[turn] - score_mg[turn^1]
	eval_eg := score_eg[turn] - score_eg[turn^1]

	phase := TotalPhase
	phase -= (pieces[chess.White][chess.Pawn].CountBits() +
		pieces[chess.Black][chess.Pawn].CountBits()) * PawnPhase
	phase -= (pieces[chess.White][chess.Knight].CountBits() +
		pieces[chess.Black][chess.Knight].CountBits()) * KnightPhase
	phase -= (pieces[chess.White][chess.Bishop].CountBits() +
		pieces[chess.Black][chess.Bishop].CountBits()) * BishopPhase
	phase -= (pieces[chess.White][chess.Rook].CountBits() +
		pieces[chess.Black][chess.Rook].CountBits()) * RookPhase
	phase -= (pieces[chess.White][chess.Queen].CountBits() +
		pieces[chess.Black][chess.Queen].CountBits()) * QueenPhase
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase

	eval := ((eval_mg * (256 - phase)) + (eval_eg * phase)) / 256

	// Check if position is likely a draw
	if is_drawish(pieces) {
		eval /= DrawishScaleFactor
	}

	return eval
}

func evalKing(
	pieces *[2]map[chess.PieceType]chess.Bitboard,
	color chess.Color,
	sq uint8,
) {

	enemyPoints := KingAttackPoints[color^1]

	// Evaluate semi-open files adjacent to the enemy king
	kingFile := MaskFile[FileOf(sq)]
	ally := (*pieces)[color][chess.Pawn]

	leftFile := ((kingFile & ClearFile[FileA]) << 1)
	rightFile := ((kingFile & ClearFile[FileH]) >> 1)

	if kingFile&ally == 0 {
		enemyPoints += SemiOpenFileNextToKingPenalty
	}

	if leftFile != 0 && leftFile&ally == 0 {
		enemyPoints += SemiOpenFileNextToKingPenalty
	}

	if rightFile != 0 && rightFile&ally == 0 {
		enemyPoints += SemiOpenFileNextToKingPenalty
	}

	// Take all the king saftey points collected for the enemy,
	// and see what kind of penatly we should get.
	penatly := (enemyPoints * enemyPoints) / 4
	if KingAttackers[color^1] >= 2 && (*pieces)[color^1][chess.Queen] != 0 {
		score_mg[color] -= penatly
	}
}

func is_draw(pieces [2]map[chess.PieceType]chess.Bitboard) bool {
	white_knights := pieces[chess.White][chess.Knight].CountBits()
	white_bishops := pieces[chess.White][chess.Bishop].CountBits()

	black_knights := pieces[chess.Black][chess.Knight].CountBits()
	black_bishops := pieces[chess.Black][chess.Bishop].CountBits()

	pawns := pieces[chess.White][chess.Pawn].CountBits() +
		pieces[chess.Black][chess.Pawn].CountBits()
	knights := white_knights + black_knights
	bishops := white_bishops + black_bishops
	rooks := pieces[chess.White][chess.Rook].CountBits() +
		pieces[chess.Black][chess.Rook].CountBits()
	queens := pieces[chess.White][chess.Queen].CountBits() +
		pieces[chess.Black][chess.Queen].CountBits()

	minors := knights + bishops
	majors := rooks + queens

	if pawns+majors+minors == 0 {
		return true
	} else if majors+pawns == 0 {
		if minors == 1 {
			return true
		} else if minors == 2 {
			if white_knights == 1 && black_knights == 1 {
				return true
			} else if white_bishops == 1 && black_bishops == 1 {
				white_bishop_square := pieces[chess.White][chess.Bishop].Msb()
				black_bishop_square := pieces[chess.Black][chess.Bishop].Msb()

				return isSqDark(white_bishop_square) ==
					isSqDark(black_bishop_square)
			}
		}
	}

	return false
}

func is_drawish(pieces [2]map[chess.PieceType]chess.Bitboard) bool {

	white_pawns := pieces[chess.White][chess.Pawn].CountBits()
	white_knights := pieces[chess.White][chess.Knight].CountBits()
	white_bishops := pieces[chess.White][chess.Bishop].CountBits()
	white_rooks := pieces[chess.White][chess.Rook].CountBits()
	white_queens := pieces[chess.White][chess.Queen].CountBits()

	black_pawns := pieces[chess.Black][chess.Pawn].CountBits()
	black_knights := pieces[chess.Black][chess.Knight].CountBits()
	black_bishops := pieces[chess.Black][chess.Bishop].CountBits()
	black_rooks := pieces[chess.Black][chess.Rook].CountBits()
	black_queens := pieces[chess.Black][chess.Queen].CountBits()

	pawns := white_pawns + black_pawns
	knights := white_knights + black_knights
	bishops := white_bishops + black_bishops
	rooks := white_rooks + black_rooks
	queens := white_queens + black_queens

	white_minors := white_knights + white_bishops
	black_minors := black_knights + black_bishops

	majors := queens + rooks
	minors := bishops + knights

	all := majors + minors

	if pawns == 0 {
		if all == 2 {
			// KQ v KQ
			if white_queens == 1 && black_queens == 1 {
				return true
			}
			// KR v KR
			if white_rooks == 1 && black_rooks == 1 {
				return true
			}
			// KN v KN
			// KN v KB
			// KB v KB
			if white_minors == 1 && black_minors == 1 {
				return true
			}
			// KNN v K
			if white_knights == 2 || black_knights == 2 {
				return true
			}
		} else if all == 3 {
			// KQ v KRR
			if (white_queens == 1 && black_rooks == 2) ||
				(black_queens == 1 && white_rooks == 2) {
				return true
			}

			// KQ v KBB
			if (white_queens == 1 && black_bishops == 2) ||
				(black_queens == 1 && white_bishops == 2) {
				return true
			}

			// KQ v KNN
			if (white_queens == 1 && black_knights == 2) ||
				(black_queens == 1 && white_knights == 2) {
				return true
			}

			// KNN v KN
			// KNN v KB
			if (white_knights == 2 && black_minors == 1) ||
				(black_knights == 2 && white_minors == 1) {
				return true
			}

		} else if all == 4 {
			// KRR v KRB
			// KRR v KRN
			if (white_rooks == 2 && black_rooks == 1 && black_minors == 1) ||
				(black_rooks == 2 && white_rooks == 1 && white_minors == 1) {
				return true
			}
		}
	}

	return false
}

func InitEvalBitboards() {
	for file := FileA; file <= FileH; file++ {
		fileBB := MaskFile[file]
		mask := (fileBB & ClearFile[FileA]) << 1
		mask |= (fileBB & ClearFile[FileH]) >> 1
		IsolatedPawnMasks[file] = mask
	}

	for sq := 0; sq < 64; sq++ {
		// Create king zones.
		sqBB := chess.SquareBB[sq]
		zone := ((sqBB & ClearFile[FileH]) >> 1) | ((sqBB & (ClearFile[FileG] & ClearFile[FileH])) >> 2)
		zone |= ((sqBB & ClearFile[FileA]) << 1) | ((sqBB & (ClearFile[FileB] & ClearFile[FileA])) << 2)
		zone |= sqBB

		zone |= ((zone >> 8) | (zone >> 16))
		zone |= ((zone << 8) | (zone << 16))

		KingZonesMasks[sq] = KingZone{OuterRing: zone & ^(KingMoves[sq] | sqBB), InnerRing: KingMoves[sq] | sqBB}

		file := FileOf(uint8(sq))
		fileBB := MaskFile[file]
		rank := int(RankOf(uint8(sq)))

		// Create doubled pawns masks.
		mask := fileBB
		for r := 0; r <= rank; r++ {
			mask &= ClearRank[r]
		}
		DoubledPawnMasks[chess.White][sq] = mask

		mask = fileBB
		for r := 7; r >= rank; r-- {
			mask &= ClearRank[r]
		}
		DoubledPawnMasks[chess.Black][sq] = mask

		// Passed pawn masks and outpost masks.
		frontSpanMask := fileBB
		frontSpanMask |= (fileBB & ClearFile[FileA]) << 1
		frontSpanMask |= (fileBB & ClearFile[FileH]) >> 1

		whiteFrontSpan := frontSpanMask
		for r := 0; r <= rank; r++ {
			whiteFrontSpan &= ClearRank[r]
		}

		PassedPawnMasks[chess.White][sq] = whiteFrontSpan
		OutpostMasks[chess.White][sq] = whiteFrontSpan & ^fileBB

		blackFrontSpan := frontSpanMask
		for r := 7; r >= rank; r-- {
			blackFrontSpan &= ClearRank[r]
		}

		PassedPawnMasks[chess.Black][sq] = blackFrontSpan
		OutpostMasks[chess.Black][sq] = blackFrontSpan & ^fileBB
	}
}
