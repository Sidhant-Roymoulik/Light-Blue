package main

// movegen.go implements the move generator for the engine.

const (
	// These masks help determine whether or not the squares between
	// the king and it's rooks are clear for castling
	F1_G1, B1_C1_D1 = 0x600000000000000, 0x7000000000000000
	F8_G8, B8_C8_D8 = 0x6, 0x70
)

// Generate rook moves.
func GenRookMoves(sq uint8, blockers Bitboard) Bitboard {
	magic := &RookMagics[sq]
	blockers &= magic.BlockerMask
	return RookMoves[sq][(uint64(blockers)*magic.MagicNo)>>magic.Shift]
}

// Generate rook moves.
func GenBishopMoves(sq uint8, blockers Bitboard) Bitboard {
	magic := &BishopMagics[sq]
	blockers &= magic.BlockerMask
	return BishopMoves[sq][(uint64(blockers)*magic.MagicNo)>>magic.Shift]
}
