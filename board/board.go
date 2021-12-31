package board

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Board represents the chessboard with a key under one of the tiles and a coin on top of each tile,
// either tails or heads up.
type Board struct {
	board        [][]bool // Y by X board. When true, the coin on top shows tails - otherwise heads.
	size         int      // Size of the board as 1 coordinate, the board is size x size large
	bitsWidth    int      // # of bits to represent tile numbers (for display purposes)
	bitstringFmt string   // format string for printing tile numbers or checksums as bits (for display)
	printFmt     string   // format string for printing any %v values to match bitstrings (for display)
	keyY         int      // Y position of the key
	keyX         int      // X position
}

// New returns an empty board with all coins heads-up and the key at (0,0).
func New(sz int) (*Board, error) {
	// The size must be a power of 2 and at least 2.
	if sz < 2 || sz&(sz-1) != 0 {
		return nil, fmt.Errorf("board size %v is not a power of 2", sz)
	}
	b := make([][]bool, sz)
	for y := 0; y < sz; y++ {
		b[y] = make([]bool, sz)
	}
	bitsWidth := int(math.Log2(float64(sz * sz)))
	return &Board{
		board:        b,
		size:         sz,
		bitsWidth:    bitsWidth,
		bitstringFmt: fmt.Sprintf("%%%v.%vb", bitsWidth, bitsWidth),
		printFmt:     fmt.Sprintf(" %%%vv ", bitsWidth),
	}, nil
}

// Clone makes a copy of the board. It can be used to save a copy for comparison with next modified
// versions.
func (b *Board) Clone() *Board {
	board := make([][]bool, b.size)
	for y := 0; y < b.size; y++ {
		board[y] = make([]bool, b.size)
		for x := 0; x < b.size; x++ {
			board[y][x] = b.board[y][x]
		}
	}
	return &Board{
		board:        board,
		size:         b.size,
		bitsWidth:    b.bitsWidth,
		bitstringFmt: b.bitstringFmt,
		printFmt:     b.printFmt,
		keyX:         b.keyX,
		keyY:         b.keyY,
	}
}

// Randomize flips the coins on all tiles to a random value and puts the key under a random tile.
func (b *Board) Randomize() *Board {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	for y := 0; y < b.size; y++ {
		for x := 0; x < b.size; x++ {
			b.board[y][x] = rand.Int()&1 == 1
		}
	}

	b.keyX = rand.Intn(b.size)
	b.keyY = rand.Intn(b.size)
	return b
}

// KeyPosition returns the key position as one value, combining the Y and X coordinate of the key. The returned
// coordinate is just a one-dimensional count from top left. E.g., when the key is on (y=0, x=5) on an 8x8 board
// then the returned value is 5. When the key is on (1,5) then the returned value is 13.
func (b *Board) KeyPosition() int {
	return b.keyY*b.size + b.keyX
}

// String returns the board layout as a printable string.
func (b *Board) String() string {
	layout := ""
	for y := 0; y < b.size; y++ {
		abs := "    "   // Output line holding absolute coords: 0, 1, 2, 3, etc
		bits := "    "  // Output line holding coords as bits: 0000, 0001, 0010, etc
		coins := "    " // Output line showing coin faces (T or -)
		keys := "    "  // Output line showing whether the key is under that tile
		var x int
		for x = 0; x < b.size; x++ {
			abs += fmt.Sprintf(b.printFmt, y*b.size+x)
			bits += " " + b.Bitstring(y*b.size+x) + " "
			coinSymbol := "T"
			if !b.board[y][x] {
				coinSymbol = "-"
			}
			coins += fmt.Sprintf(b.printFmt, coinSymbol)
			keySymbol := "   "
			if x == b.keyX && y == b.keyY {
				keySymbol = "Key"
			}
			keys += fmt.Sprintf(b.printFmt, keySymbol)
		}
		layout += bits + "\n" + abs + "\n" + coins + "\n" + keys + "\n\n"
	}
	return layout
}

// Checksum computes the board checksum. Tiles where the coins are heads-up are not taken into account. The
// returned value is a parity-sum over all tiles with tails-up.
func (b *Board) Checksum() int {
	cs := 0
	first := true
	for y := 0; y < b.size; y++ {
		for x := 0; x < b.size; x++ {
			if !b.board[y][x] {
				continue
			}
			if first {
				cs = y*b.size + x
				first = false
			} else {
				cs ^= y*b.size + x
			}
		}
	}
	return cs
}

// Bistring makes printing easier. The argument is a value that represents the overall position of an (y,x)
// coordinate. The return value is that combined coordinate expressed as bits. The width of the returned
// string is adjusted to the number of bits that are needed. E.g., for 4x4 boards the width is 4 (0000 to 1111);
// for 8x8 boards the width is 6 (000000 to 111111).
func (b *Board) Bitstring(v int) string {
	return fmt.Sprintf(b.bitstringFmt, v)
}

// Flip flips a coin on the given coordinate, which must be given as the count from top-left. E.g., on a 4x4
// board, the tile at y=1, x=3 has to be passed-in as 7.
func (b *Board) Flip(p int) {
	y := p / b.size
	x := p % b.size
	b.board[y][x] = !b.board[y][x]
}
