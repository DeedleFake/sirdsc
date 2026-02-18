// Package spcg provides a stateless implementation of a PCG random
// number generator.
//
// The implementation is copied almost verbatim from math/rand/v2.
package spcg

import "math/bits"

// Next returns the next random number in the PCG sequence, along with
// the high and low bits for getting the number after the one
// returned.
func Next(high, low uint64) (n, nhigh, nlow uint64) {
	const (
		mulHi = 2549297995355413924
		mulLo = 4865540595714422341
		incHi = 6364136223846793005
		incLo = 1442695040888963407
	)

	nhigh, nlow = bits.Mul64(low, mulLo)
	nhigh += high*mulLo + low*mulHi
	nlow, c := bits.Add64(nlow, incLo, 0)
	nhigh, _ = bits.Add64(nhigh, incHi, c)
	return next(nhigh, nlow), nhigh, nlow
}

func next(high, low uint64) uint64 {
	const cheapMul = 0xda942042e4dd58b5
	high ^= high >> 32
	high *= cheapMul
	high ^= high >> 48
	high *= (low | 1)
	return high
}
