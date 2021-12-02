// Package spcg provides a stateless implementation of a PCG random
// number generator.
//
// The implementation is copied almost verbatim from
// golang.org/x/exp/rand.
package spcg

import "math/bits"

type pcg struct {
	low  uint64
	high uint64
}

const (
	maxUint64 = (1 << 64) - 1

	multiplier = 47026247687942121848144207491837523525
	mulHigh    = multiplier >> 64
	mulLow     = multiplier & maxUint64

	increment = 117397592171526113268558934119004209487
	incHigh   = increment >> 64
	incLow    = increment & maxUint64
)

// Next returns the next random number in the PCG sequence, along with
// the high and low bits for getting the number after the one
// returned.
func Next(high, low uint64) (n, nhigh, nlow uint64) {
	high, low = multiply(high, low)
	high, low = add(high, low)
	return bits.RotateLeft64(high^low, -int(high>>58)), high, low
}

func add(high, low uint64) (uint64, uint64) {
	var carry uint64
	low, carry = bits.Add64(low, incLow, 0)
	high, _ = bits.Add64(high, incHigh, carry)
	return high, low
}

func multiply(high, low uint64) (uint64, uint64) {
	hi, lo := bits.Mul64(low, mulLow)
	hi += high * mulLow
	hi += low * mulHigh
	low = lo
	high = hi
	return high, low
}
