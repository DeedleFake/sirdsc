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
func Next(high, low uint64) (r, rhigh, rlow uint64) {
	pcg := pcg{high, low}
	pcg.multiply()
	pcg.add()
	return bits.RotateLeft64(pcg.high^pcg.low, -int(pcg.high>>58)), pcg.high, pcg.low
}

func (pcg *pcg) add() {
	var carry uint64
	pcg.low, carry = bits.Add64(pcg.low, incLow, 0)
	pcg.high, _ = bits.Add64(pcg.high, incHigh, carry)
}

func (pcg *pcg) multiply() {
	hi, lo := bits.Mul64(pcg.low, mulLow)
	hi += pcg.high * mulLow
	hi += pcg.low * mulHigh
	pcg.low = lo
	pcg.high = hi
}
