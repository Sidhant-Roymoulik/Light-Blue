package chess

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

// A Board represents a chess board and its relationship between squares and pieces.
type Board struct {
	BBWhiteKing   Bitboard
	BBWhiteQueen  Bitboard
	BBWhiteRook   Bitboard
	BBWhiteBishop Bitboard
	BBWhiteKnight Bitboard
	BBWhitePawn   Bitboard
	BBBlackKing   Bitboard
	BBBlackQueen  Bitboard
	BBBlackRook   Bitboard
	BBBlackBishop Bitboard
	BBBlackKnight Bitboard
	BBBlackPawn   Bitboard
	WhiteSqs      Bitboard
	BlackSqs      Bitboard
	EmptySqs      Bitboard
	WhiteKingSq   Square
	BlackKingSq   Square
}

// NewBoard returns a board from a square to piece mapping.
func NewBoard(m map[Square]Piece) *Board {
	b := &Board{}
	for _, p1 := range allPieces {
		bm := map[Square]bool{}
		for sq, p2 := range m {
			if p1 == p2 {
				bm[sq] = true
			}
		}
		BB := newBitboard(bm)
		b.setBBForPiece(p1, BB)
	}
	b.calcConvienceBBs(nil)
	return b
}

// SquareMap returns a mapping of squares to pieces.  A square is only added to the map if it is occupied.
func (b *Board) SquareMap() map[Square]Piece {
	m := map[Square]Piece{}
	for sq := 0; sq < numOfSquaresInBoard; sq++ {
		p := b.Piece(Square(sq))
		if p != NoPiece {
			m[Square(sq)] = p
		}
	}
	return m
}

// Rotate rotates the board 90 degrees clockwise.
func (b *Board) Rotate() *Board {
	return b.Flip(UpDown).Transpose()
}

// FlipDirection is the direction for the Board.Flip method
type FlipDirection int

const (
	// UpDown flips the board's rank values
	UpDown FlipDirection = iota
	// LeftRight flips the board's file values
	LeftRight
)

// Flip flips the board over the vertical or hoizontal
// center line.
func (b *Board) Flip(fd FlipDirection) *Board {
	m := map[Square]Piece{}
	for sq := 0; sq < numOfSquaresInBoard; sq++ {
		var mv Square
		switch fd {
		case UpDown:
			file := Square(sq).File()
			rank := Rank(7 - Square(sq).Rank())
			mv = NewSquare(file, rank)
		case LeftRight:
			file := File(7 - Square(sq).File())
			rank := Square(sq).Rank()
			mv = NewSquare(file, rank)
		}
		m[mv] = b.Piece(Square(sq))
	}
	return NewBoard(m)
}

// Transpose flips the board over the A8 to H1 diagonal.
func (b *Board) Transpose() *Board {
	m := map[Square]Piece{}
	for sq := 0; sq < numOfSquaresInBoard; sq++ {
		file := File(7 - Square(sq).Rank())
		rank := Rank(7 - Square(sq).File())
		mv := NewSquare(file, rank)
		m[mv] = b.Piece(Square(sq))
	}
	return NewBoard(m)
}

// Draw returns visual representation of the board useful for debugging.
func (b *Board) Draw() string {
	s := "\n A B C D E F G H\n"
	for r := 7; r >= 0; r-- {
		s += Rank(r).String()
		for f := 0; f < numOfSquaresInRow; f++ {
			p := b.Piece(NewSquare(File(f), Rank(r)))
			if p == NoPiece {
				s += "-"
			} else {
				s += p.String()
			}
			s += " "
		}
		s += "\n"
	}
	return s
}

// String implements the fmt.Stringer interface and returns
// a string in the FEN board format: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR
func (b *Board) String() string {
	fen := ""
	for r := 7; r >= 0; r-- {
		for f := 0; f < numOfSquaresInRow; f++ {
			sq := NewSquare(File(f), Rank(r))
			p := b.Piece(sq)
			if p != NoPiece {
				fen += p.getFENChar()
			} else {
				fen += "1"
			}
		}
		if r != 0 {
			fen += "/"
		}
	}
	for i := 8; i > 1; i-- {
		repeatStr := strings.Repeat("1", i)
		countStr := strconv.Itoa(i)
		fen = strings.Replace(fen, repeatStr, countStr, -1)
	}
	return fen
}

// Piece returns the piece for the given square.
func (b *Board) Piece(sq Square) Piece {
	for _, p := range allPieces {
		BB := b.BBForPiece(p)
		if BB.Occupied(sq) {
			return p
		}
	}
	return NoPiece
}

// MarshalText implements the encoding.TextMarshaler interface and returns
// a string in the FEN board format: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR
func (b *Board) MarshalText() (text []byte, err error) {
	return []byte(b.String()), nil
}

// UnmarshalText implements the encoding.TextUnarshaler interface and takes
// a string in the FEN board format: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR
func (b *Board) UnmarshalText(text []byte) error {
	cp, err := fenBoard(string(text))
	if err != nil {
		return err
	}
	*b = *cp
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface and returns
// the Bitboard representations as a array of bytes.  Bitboads are encoded
// in the following order: WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight
// WhitePawn, BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn
func (b *Board) MarshalBinary() (data []byte, err error) {
	BBs := []Bitboard{b.BBWhiteKing, b.BBWhiteQueen, b.BBWhiteRook, b.BBWhiteBishop, b.BBWhiteKnight, b.BBWhitePawn,
		b.BBBlackKing, b.BBBlackQueen, b.BBBlackRook, b.BBBlackBishop, b.BBBlackKnight, b.BBBlackPawn}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, BBs)
	return buf.Bytes(), err
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface and parses
// the Bitboard representations as a array of bytes.  Bitboads are decoded
// in the following order: WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight
// WhitePawn, BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn
func (b *Board) UnmarshalBinary(data []byte) error {
	if len(data) != 96 {
		return errors.New("chess: invalid number of bytes for board unmarshal binary")
	}
	b.BBWhiteKing = Bitboard(binary.BigEndian.Uint64(data[:8]))
	b.BBWhiteQueen = Bitboard(binary.BigEndian.Uint64(data[8:16]))
	b.BBWhiteRook = Bitboard(binary.BigEndian.Uint64(data[16:24]))
	b.BBWhiteBishop = Bitboard(binary.BigEndian.Uint64(data[24:32]))
	b.BBWhiteKnight = Bitboard(binary.BigEndian.Uint64(data[32:40]))
	b.BBWhitePawn = Bitboard(binary.BigEndian.Uint64(data[40:48]))
	b.BBBlackKing = Bitboard(binary.BigEndian.Uint64(data[48:56]))
	b.BBBlackQueen = Bitboard(binary.BigEndian.Uint64(data[56:64]))
	b.BBBlackRook = Bitboard(binary.BigEndian.Uint64(data[64:72]))
	b.BBBlackBishop = Bitboard(binary.BigEndian.Uint64(data[72:80]))
	b.BBBlackKnight = Bitboard(binary.BigEndian.Uint64(data[80:88]))
	b.BBBlackPawn = Bitboard(binary.BigEndian.Uint64(data[88:96]))
	b.calcConvienceBBs(nil)
	return nil
}

func (b *Board) update(m *Move) {
	p1 := b.Piece(m.s1)
	s1BB := BBForSquare(m.s1)
	s2BB := BBForSquare(m.s2)

	// move s1 piece to s2
	for _, p := range allPieces {
		BB := b.BBForPiece(p)
		// remove what was at s2
		b.setBBForPiece(p, BB & ^s2BB)
		// move what was at s1 to s2
		if BB.Occupied(m.s1) {
			BB = b.BBForPiece(p)
			b.setBBForPiece(p, (BB & ^s1BB)|s2BB)
		}
	}
	// check promotion
	if m.promo != NoPieceType {
		newPiece := NewPiece(m.promo, p1.Color())
		// remove pawn
		BBPawn := b.BBForPiece(p1)
		b.setBBForPiece(p1, BBPawn & ^s2BB)
		// add promo piece
		BBPromo := b.BBForPiece(newPiece)
		b.setBBForPiece(newPiece, BBPromo|s2BB)
	}
	// remove captured en passant piece
	if m.HasTag(EnPassant) {
		if p1.Color() == White {
			b.BBBlackPawn = ^(BBForSquare(m.s2) << 8) & b.BBBlackPawn
		} else {
			b.BBWhitePawn = ^(BBForSquare(m.s2) >> 8) & b.BBWhitePawn
		}
	}
	// move rook for castle
	if p1.Color() == White && m.HasTag(KingSideCastle) {
		b.BBWhiteRook = (b.BBWhiteRook & ^BBForSquare(H1) | BBForSquare(F1))
	} else if p1.Color() == White && m.HasTag(QueenSideCastle) {
		b.BBWhiteRook = (b.BBWhiteRook & ^BBForSquare(A1)) | BBForSquare(D1)
	} else if p1.Color() == Black && m.HasTag(KingSideCastle) {
		b.BBBlackRook = (b.BBBlackRook & ^BBForSquare(H8) | BBForSquare(F8))
	} else if p1.Color() == Black && m.HasTag(QueenSideCastle) {
		b.BBBlackRook = (b.BBBlackRook & ^BBForSquare(A8)) | BBForSquare(D8)
	}
	b.calcConvienceBBs(m)
}

func (b *Board) calcConvienceBBs(m *Move) {
	WhiteSqs := b.BBWhiteKing | b.BBWhiteQueen | b.BBWhiteRook | b.BBWhiteBishop | b.BBWhiteKnight | b.BBWhitePawn
	BlackSqs := b.BBBlackKing | b.BBBlackQueen | b.BBBlackRook | b.BBBlackBishop | b.BBBlackKnight | b.BBBlackPawn
	EmptySqs := ^(WhiteSqs | BlackSqs)
	b.WhiteSqs = WhiteSqs
	b.BlackSqs = BlackSqs
	b.EmptySqs = EmptySqs
	if m == nil {
		b.WhiteKingSq = NoSquare
		b.BlackKingSq = NoSquare

		for sq := 0; sq < numOfSquaresInBoard; sq++ {
			sqr := Square(sq)
			if b.BBWhiteKing.Occupied(sqr) {
				b.WhiteKingSq = sqr
			} else if b.BBBlackKing.Occupied(sqr) {
				b.BlackKingSq = sqr
			}
		}
	} else if m.s1 == b.WhiteKingSq {
		b.WhiteKingSq = m.s2
	} else if m.s1 == b.BlackKingSq {
		b.BlackKingSq = m.s2
	}
}

func (b *Board) copy() *Board {
	return &Board{
		WhiteSqs:      b.WhiteSqs,
		BlackSqs:      b.BlackSqs,
		EmptySqs:      b.EmptySqs,
		WhiteKingSq:   b.WhiteKingSq,
		BlackKingSq:   b.BlackKingSq,
		BBWhiteKing:   b.BBWhiteKing,
		BBWhiteQueen:  b.BBWhiteQueen,
		BBWhiteRook:   b.BBWhiteRook,
		BBWhiteBishop: b.BBWhiteBishop,
		BBWhiteKnight: b.BBWhiteKnight,
		BBWhitePawn:   b.BBWhitePawn,
		BBBlackKing:   b.BBBlackKing,
		BBBlackQueen:  b.BBBlackQueen,
		BBBlackRook:   b.BBBlackRook,
		BBBlackBishop: b.BBBlackBishop,
		BBBlackKnight: b.BBBlackKnight,
		BBBlackPawn:   b.BBBlackPawn,
	}
}

func (b *Board) isOccupied(sq Square) bool {
	return !b.EmptySqs.Occupied(sq)
}

func (b *Board) hasSufficientMaterial() bool {
	// queen, rook, or pawn exist
	if (b.BBWhiteQueen | b.BBWhiteRook | b.BBWhitePawn |
		b.BBBlackQueen | b.BBBlackRook | b.BBBlackPawn) > 0 {
		return true
	}
	// if king is missing then it is a test
	if b.BBWhiteKing == 0 || b.BBBlackKing == 0 {
		return true
	}
	count := map[PieceType]int{}
	pieceMap := b.SquareMap()
	for _, p := range pieceMap {
		count[p.Type()]++
	}
	// 	king versus king
	if count[Bishop] == 0 && count[Knight] == 0 {
		return false
	}
	// king and bishop versus king
	if count[Bishop] == 1 && count[Knight] == 0 {
		return false
	}
	// king and knight versus king
	if count[Bishop] == 0 && count[Knight] == 1 {
		return false
	}
	// king and bishop(s) versus king and bishop(s) with the bishops on the same colour.
	if count[Knight] == 0 {
		whiteCount := 0
		blackCount := 0
		for sq, p := range pieceMap {
			if p.Type() == Bishop {
				switch sq.color() {
				case White:
					whiteCount++
				case Black:
					blackCount++
				}
			}
		}
		if whiteCount == 0 || blackCount == 0 {
			return false
		}
	}
	return true
}

func (b *Board) BBForPiece(p Piece) Bitboard {
	switch p {
	case WhiteKing:
		return b.BBWhiteKing
	case WhiteQueen:
		return b.BBWhiteQueen
	case WhiteRook:
		return b.BBWhiteRook
	case WhiteBishop:
		return b.BBWhiteBishop
	case WhiteKnight:
		return b.BBWhiteKnight
	case WhitePawn:
		return b.BBWhitePawn
	case BlackKing:
		return b.BBBlackKing
	case BlackQueen:
		return b.BBBlackQueen
	case BlackRook:
		return b.BBBlackRook
	case BlackBishop:
		return b.BBBlackBishop
	case BlackKnight:
		return b.BBBlackKnight
	case BlackPawn:
		return b.BBBlackPawn
	}
	return Bitboard(0)
}

func (b *Board) setBBForPiece(p Piece, BB Bitboard) {
	switch p {
	case WhiteKing:
		b.BBWhiteKing = BB
	case WhiteQueen:
		b.BBWhiteQueen = BB
	case WhiteRook:
		b.BBWhiteRook = BB
	case WhiteBishop:
		b.BBWhiteBishop = BB
	case WhiteKnight:
		b.BBWhiteKnight = BB
	case WhitePawn:
		b.BBWhitePawn = BB
	case BlackKing:
		b.BBBlackKing = BB
	case BlackQueen:
		b.BBBlackQueen = BB
	case BlackRook:
		b.BBBlackRook = BB
	case BlackBishop:
		b.BBBlackBishop = BB
	case BlackKnight:
		b.BBBlackKnight = BB
	case BlackPawn:
		b.BBBlackPawn = BB
	default:
		panic("invalid piece")
	}
}
