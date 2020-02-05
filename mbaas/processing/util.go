package processing

func nextPowerOfTwo(orig int) int {
	next := 1
	for next < orig {
		next <<= 1
	}
	return next
}

func numberOfBits(n int) int {
	pow2 := 1
	nbits := 0
	for n > pow2 {
		pow2 *= 2
		nbits++
	}
	return nbits
}
