package chess

import (
	"strconv"
)

// Bitboard is a board representation encoded in an unsigned 64-bit integer.  The
// 64 board positions begin with A1 as the most significant bit and H8 as the least.
type Bitboard uint64

func newBitboard(m map[Square]bool) Bitboard {
	s := ""
	for sq := 0; sq < numOfSquaresInBoard; sq++ {
		if m[Square(sq)] {
			s += "1"
		} else {
			s += "0"
		}
	}
	bb, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		panic(err)
	}
	return Bitboard(bb)
}

func (b Bitboard) Mapping() map[Square]bool {
	m := map[Square]bool{}
	for sq := 0; sq < numOfSquaresInBoard; sq++ {
		if b&BBForSquare(Square(sq)) > 0 {
			m[Square(sq)] = true
		}
	}
	return m
}

// String returns a 64 character string of 1s and 0s starting with the most significant bit.
// func (b Bitboard) String() string {
// 	s := strconv.FormatUint(uint64(b), 2)
// 	return strings.Repeat("0", numOfSquaresInBoard-len(s)) + s
// }

// Draw returns visual representation of the Bitboard useful for debugging.
func (b Bitboard) Draw() string {
	s := "\n A B C D E F G H\n"
	for r := 7; r >= 0; r-- {
		s += Rank(r).String()
		for f := 0; f < numOfSquaresInRow; f++ {
			sq := NewSquare(File(f), Rank(r))
			if b.Occupied(sq) {
				s += "1"
			} else {
				s += "0"
			}
			s += " "
		}
		s += "\n"
	}
	return s
}
