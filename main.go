package main

import (
	"fmt"
	"log"

	"aicp/board"
)

const (
	// Size in one dimension, 8 for a normal chess board. Allowed: 2, 4, 8, 16, etc. (powers of 2)
	// Beyond 16 you'll need a really wide screen for output.
	boardSize = 8

	// # of puzzles to try in a batch run
	nRuns = 1000000
)

func main() {
	// Single random board.
	if err := run(true); err != nil {
		log.Fatalln(err)
	}

	// Let's do a bunch of puzzles and try to catch some edge case that fails.
	fmt.Println("Trying", nRuns, "puzzles...")
	for i := 0; i < nRuns; i++ {
		if err := run(false); err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("No errors were seen. The solution probably works :-)")
}

// run solves one puzzle. The argument is whether to print on stdout the initial board, checksum, what to flip,
// and the flipped board. That may be illustrative for a single run.
func run(beVerbose bool) error {
	org, err := board.New(boardSize)
	if err != nil {
		return err
	}
	org.Randomize()
	diff := org.Checksum() ^ org.KeyPosition()
	flipped := org.Clone()
	flipped.Flip(diff)
	if flipped.Checksum() != org.KeyPosition() {
		return fmt.Errorf(
			"Guess failed: flipped board checksum=%v, but key is at %v\n"+
				"Original board:\n%v\n"+
				"Flipped board:\n%v\n",
			flipped.Bitstring(flipped.Checksum()), flipped.Bitstring(flipped.KeyPosition()),
			org,
			flipped)
	}

	if !beVerbose {
		return nil
	}

	fmt.Printf("The warden has prepared this board:\n%v", org)
	fmt.Printf("Alice computes the checksum : %v\n", org.Bitstring(org.Checksum()))
	fmt.Printf("The key is at position      : %v (abs:%v)\n", org.Bitstring(org.KeyPosition()), org.KeyPosition())
	fmt.Printf("The difference              : %v\n", org.Bitstring(diff))
	fmt.Printf("Alice asks the difference-coin (abs:%v) to be flipped and leaves the room.\n", diff)

	fmt.Printf("Bob enters the room and sees the modified board (though not where the key is):\n%v", flipped)
	fmt.Printf("Bob computes the checksum: %v (abs:%v)\n", flipped.Bitstring(flipped.Checksum()), flipped.Checksum())
	fmt.Println("That's the tile Bob wants to get uncovered. Whee! Found the key!\n")

	return nil
}
