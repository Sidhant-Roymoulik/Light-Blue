package chess

import "math/bits"

// Reverse returns a Bitboard where the bit order is reversed.
func (b Bitboard) Reverse() Bitboard {
	return Bitboard(bits.Reverse64(uint64(b)))
}

// Occupied returns true if the square's Bitboard position is 1.
func (b Bitboard) Occupied(sq Square) bool {
	return (bits.RotateLeft64(uint64(b), int(sq)+1) & 1) == 1
}
