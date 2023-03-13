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
		21, 38, 56, 58, 41, 30, 14, 5,
		-8, 11, 35, 46, 51, 45, 42, 19,
		-29, -16, -16, -4, 9, 3, 16, -5,
		-39, -24, -24, -18, -12, -10, 2, -16,
		-44, -29, -35, -36, -20, -13, 6, -17,
		-44, -32, -35, -46, -25, -13, 4, -26,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-40, 0, 7, 5, -2, -3, -1, -4,
		-4, 12, 10, 23, 22, 18, 0, 12,
		-1, 18, 36, 49, 58, 38, 35, 8,
		-3, -3, 2, 21, 6, 27, 4, 21,
		-21, 3, -2, -6, -1, 4, 13, -12,
		-30, -17, -14, 1, 0, -7, 2, -16,
		-39, -32, -16, -17, -15, -16, -31, -28,
		-37, -30, -43, -40, -28, -27, -24, -25,
	},
	chess.Bishop: {
		-8, -6, 2, -10, -3, -6, -3, -8,
		-3, -6, -3, 0, 5, 10, -15, -10,
		-4, 20, 9, 28, 25, 25, 29, 15,
		-7, -11, 6, 22, 8, 24, -5, -4,
		-6, -3, -6, 3, 2, -1, -1, 0,
		-12, 3, -5, -7, -7, -12, -2, -6,
		-5, -12, 4, -18, -16, -9, -6, -19,
		-20, -6, -22, -33, -28, -23, -21, -18,
	},
	chess.Rook: {
		15, 15, 12, 18, 7, 6, 6, 9,
		3, -3, 6, 25, 23, 17, 5, 16,
		-5, 13, 6, 19, 33, 35, 22, 8,
		-2, -7, -11, 11, 0, -5, 1, 7,
		-34, -29, -29, -15, -34, -21, -4, -28,
		-52, -28, -27, -29, -36, -30, -20, -20,
		-50, -38, -28, -29, -35, -37, -22, -41,
		-43, -37, -29, -26, -26, -40, -18, -38,
	},
	chess.Queen: {
		-12, -1, 4, 7, 2, 3, 5, 2,
		-11, -16, 4, -5, 6, 10, -3, 6,
		-6, 0, 0, 13, 25, 30, 26, 31,
		-8, 0, 1, -2, 2, 10, -2, 7,
		-15, -11, -11, -14, -13, 2, 0, 7,
		-8, -11, -8, -16, -3, -7, 0, -1,
		-27, -14, -6, -3, -3, 1, -8, -9,
		-8, -21, -18, -12, -8, -33, -2, 6,
	},
	chess.King: {
		0, 0, 0, 0, 0, 0, 0, 0,
		1, 3, 2, 2, 1, 1, 4, 0,
		0, 8, 6, 4, 3, 6, 5, 1,
		1, 4, 6, 5, 6, 5, 7, -1,
		-1, 4, 7, 5, 5, 4, -1, -8,
		-7, 4, 11, -11, -12, -6, -16, -22,
		7, 7, -7, -31, -24, -24, 5, -2,
		4, 43, 12, -54, -1, -35, 21, 22,
	},
}

var PST_EG = map[chess.PieceType][]int{
	chess.Pawn: {
		0, 0, 0, 0, 0, 0, 0, 0,
		58, 60, 53, 48, 48, 22, 21, 22,
		46, 31, 14, 14, 10, 13, 38, 14,
		14, 4, -5, -14, -22, -19, -6, -8,
		-5, -3, -28, -35, -26, -23, -18, -25,
		-16, -5, -27, -17, -19, -25, -18, -31,
		-6, 0, -12, -7, 6, -11, -14, -29,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	chess.Knight: {
		-21, -1, 4, 4, -4, -2, -1, -3,
		-11, 6, 6, 21, 15, 6, -1, 0,
		1, 7, 25, 20, 27, 13, 18, 6,
		13, 21, 31, 41, 37, 27, 23, -3,
		-17, 4, 28, 31, 37, 28, 7, -15,
		-36, -5, 9, 22, 16, -3, -21, -15,
		-30, -15, -21, -16, -10, -8, -18, -18,
		-21, -65, -29, -18, -37, -38, -47, -13,
	},
	chess.Bishop: {
		-4, -8, 9, -2, -1, -4, -7, -4,
		-8, 0, 2, 9, 0, 3, 0, -14,
		4, 8, 10, 2, 17, 15, 13, 18,
		-5, 20, 11, 23, 24, 7, 8, 7,
		-8, 14, 25, 10, 14, 10, 0, -18,
		-14, 9, 9, 17, 16, 4, -7, -1,
		-18, -23, -10, 1, 7, -14, -7, -16,
		-23, -18, -37, -11, -17, -8, -8, -17,
	},
	chess.Rook: {
		7, 16, 21, 19, 20, 5, 11, 11,
		9, 23, 25, 30, 29, 10, 8, 11,
		21, 15, 19, 19, 14, 13, 18, 1,
		1, 8, 18, 13, 3, 6, -3, 1,
		-2, -9, 7, -3, 1, -5, -6, -12,
		-10, -13, -14, -14, -11, -23, -13, -18,
		-17, -21, -15, -11, -26, -26, -29, -13,
		-17, -13, -6, -2, -14, -10, -28, -28,
	},
	chess.Queen: {
		0, 5, 4, 3, 4, 5, 5, 3,
		-10, 5, 11, 6, 4, 7, 5, -3,
		-10, 6, 13, 14, 15, 27, 11, 11,
		-6, -7, 11, 20, 17, 18, 10, 1,
		-11, 0, 3, 25, 4, -8, 0, -4,
		-10, -13, 0, 2, -7, 6, -9, -1,
		-12, -12, -23, -19, -32, -33, -15, -10,
		-11, -8, -29, -28, -29, -26, -10, 1,
	},
	chess.King: {
		0, 1, 0, 3, 4, 2, 2, -2,
		1, 12, 9, 6, 5, 8, 14, -2,
		1, 25, 28, 15, 15, 24, 21, 4,
		7, 19, 29, 28, 34, 22, 23, 2,
		2, 19, 22, 24, 21, 16, 6, -15,
		-7, -3, 6, 2, 5, -3, -10, -24,
		-5, -19, -8, -7, -5, -8, -22, -30,
		-16, -27, -18, -27, -61, -22, -39, -70,
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
