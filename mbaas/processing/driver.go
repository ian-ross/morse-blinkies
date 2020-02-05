package processing

import (
	"errors"
	"fmt"
	"strings"
)

// MakeBoard is the main driver for constructing a board.
func MakeBoard(text string) error {
	text = strings.ToUpper(text)
	fmt.Println("Making board for:", text)

	bits, err := textToBitSequence(text)
	if err != nil {
		return err
	}
	fmt.Println("Bits:", bits)

	return nil
}

type symbol uint

const (
	dot  symbol = 0
	dash symbol = 1
)

func letter(symbols ...symbol) []bool {
	bits := []bool{}
	for _, s := range symbols {
		switch s {
		case dot:
			bits = append(bits, true)
		case dash:
			bits = append(bits, true)
			bits = append(bits, true)
			bits = append(bits, true)
		}
		// Inter-symbol space: every letter also has one of these at the
		// end!
		bits = append(bits, false)
	}
	return bits
}

var morse = map[rune][]bool{
	'A': letter(dot, dash),
	'B': letter(dash, dot, dot, dot),
	'C': letter(dash, dot, dash, dot),
	'D': letter(dash, dot, dot),
	'E': letter(dot),
	'F': letter(dot, dot, dash, dot),
	'G': letter(dash, dash, dot),
	'H': letter(dot, dot, dot, dot),
	'I': letter(dot, dot),
	'J': letter(dot, dash, dash, dash),
	'K': letter(dash, dot, dash),
	'L': letter(dot, dash, dot, dot),
	'M': letter(dash, dash),
	'N': letter(dash, dot),
	'O': letter(dash, dash, dash),
	'P': letter(dot, dash, dash, dot),
	'Q': letter(dash, dash, dot, dash),
	'R': letter(dot, dash, dot),
	'S': letter(dot, dot, dot),
	'T': letter(dash),
	'U': letter(dot, dot, dash),
	'V': letter(dot, dot, dot, dash),
	'W': letter(dot, dash, dash),
	'X': letter(dash, dot, dot, dash),
	'Y': letter(dash, dot, dash, dash),
	'Z': letter(dash, dash, dot, dot),
	' ': []bool{false, false}, // TODO: MAKE SURE THIS IS RIGHT
}

func textToBitSequence(text string) ([]bool, error) {
	bits := []bool{}
	for i, let := range text {
		if i != 0 {
			bits = append(bits, false)
			bits = append(bits, false)
		}
		seq, ok := morse[let]
		if !ok {
			return nil, errors.New("invalid character '" + string(let) + "'")
		}
		// TODO: DO THIS BETTER
		for _, b := range seq {
			bits = append(bits, b)
		}
	}

	return bits[:len(bits)-1], nil
}
