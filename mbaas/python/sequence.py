MARK = True
SPACE = False

DOT = 0
DASH = 1


def letter(*symbols):
    bits = []
    for s in symbols:
        if s == DOT:
            bits += [MARK]
        else:
            bits += [MARK, MARK, MARK]

        # Inter-symbol space: every letter also has one of these at
        # the end!
        bits.append(SPACE)

    return bits

MORSE = {
    'A': letter(DOT, DASH),
    'B': letter(DASH, DOT, DOT, DOT),
    'C': letter(DASH, DOT, DASH, DOT),
    'D': letter(DASH, DOT, DOT),
    'E': letter(DOT),
    'F': letter(DOT, DOT, DASH, DOT),
    'G': letter(DASH, DASH, DOT),
    'H': letter(DOT, DOT, DOT, DOT),
    'I': letter(DOT, DOT),
    'J': letter(DOT, DASH, DASH, DASH),
    'K': letter(DASH, DOT, DASH),
    'L': letter(DOT, DASH, DOT, DOT),
    'M': letter(DASH, DASH),
    'N': letter(DASH, DOT),
    'O': letter(DASH, DASH, DASH),
    'P': letter(DOT, DASH, DASH, DOT),
    'Q': letter(DASH, DASH, DOT, DASH),
    'R': letter(DOT, DASH, DOT),
    'S': letter(DOT, DOT, DOT),
    'T': letter(DASH),
    'U': letter(DOT, DOT, DASH),
    'V': letter(DOT, DOT, DOT, DASH),
    'W': letter(DOT, DASH, DASH),
    'X': letter(DASH, DOT, DOT, DASH),
    'Y': letter(DASH, DOT, DASH, DASH),
    'Z': letter(DASH, DASH, DOT, DOT),
    ' ': [SPACE, SPACE],
}

def text_to_bit_sequence(text):
    bits = []
    first = True
    for let in text:
        if not first:
            bits += [SPACE, SPACE]
        first = False

        if let not in MORSE:
            raise Exception("invalid character '" + let + "'")
        bits += MORSE[let]

    return bits[:len(bits)-1]

def to_string(seq):
    return ''.join(['1' if b else '0' for b in seq])
