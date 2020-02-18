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

def text_to_bit_sequence(text, seq_type):
    nletters = len([l for l in text if l in MORSE])
    nseqs = nletters if seq_type == 'multi' else 1
    bits = [[] for i in range(nseqs)]
    first = True
    seq = 0
    for let in text:
        if not first:
            for i in range(nseqs):
                bits[i] += [SPACE, SPACE]
        first = False

        if let not in MORSE:
            raise Exception("invalid character '" + let + "'")
        bits[seq] += MORSE[let]
        for i in range(nseqs):
            if i != seq:
                for j in range(len(MORSE[let])):
                    bits[i].append(SPACE)
        if nseqs > 1:
            seq += 1

    for i in range(nseqs):
        bits[i] = bits[i][:len(bits[i])-1]

    return bits

def to_string(seq):
    return ''.join(['1' if b else '0' for b in seq])
