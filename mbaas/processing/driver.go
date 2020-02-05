package processing

import (
	"fmt"
	"strings"
)

// MakeBoard is the main driver for constructing a board.
func MakeBoard(text string, padding int) error {
	fmt.Println("Making board for:", text)

	// Convert to Morse mark/space sequence.
	text = strings.ToUpper(text)
	seq, err := textToBitSequence(text)
	if err != nil {
		return err
	}
	fmt.Println("Bit sequence:", seq)

	// Padding: either a fixed padding or to next power of two.
	length := len(seq)
	if padding > 0 {
		length += padding
	} else {
		length = nextPowerOfTwo(length)
	}
	npadding := length - len(seq)
	fmt.Println("Length:", len(seq), "->", length)
	for i := 0; i < npadding; i++ {
		seq = append(seq, space)
	}
	fmt.Println("Padded bit sequence:", seq)

	// Convert to Espresso format.
	esp, err := espresso(seq)
	fmt.Println("Simplified Espresso:", esp)

	return nil
}
