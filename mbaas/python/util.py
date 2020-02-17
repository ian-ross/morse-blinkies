def next_power_of_two(orig):
    next = 1
    while next < orig:
        next <<= 1
    return next

def number_of_bits(n):
    pow2 = 1
    nbits = 0
    while n > pow2:
        pow2 *= 2
        nbits += 1
    return nbits
