package processing

import "errors"

type unit bool

const (
	mark  unit = true
	space unit = false
)

type sequence []unit

type symbol uint

const (
	dot  symbol = 0
	dash symbol = 1
)

func letter(symbols ...symbol) sequence {
	bits := sequence{}
	for _, s := range symbols {
		switch s {
		case dot:
			bits = append(bits, mark)
		case dash:
			bits = append(bits, sequence{mark, mark, mark}...)
		}
		// Inter-symbol space: every letter also has one of these at the
		// end!
		bits = append(bits, space)
	}
	return bits
}

var morse = map[rune]sequence{
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
	' ': sequence{space, space},
}

func textToBitSequence(text string) (sequence, error) {
	bits := sequence{}
	for i, let := range text {
		if i != 0 {
			bits = append(bits, sequence{space, space}...)
		}
		seq, ok := morse[let]
		if !ok {
			return nil, errors.New("invalid character '" + string(let) + "'")
		}
		bits = append(bits, seq...)
	}

	return bits[:len(bits)-1], nil
}

func (u unit) String() string {
	if u == mark {
		return "1"
	}
	return "0"
}

func (seq sequence) String() string {
	s := ""
	for _, u := range seq {
		s += u.String()
	}
	return s
}
