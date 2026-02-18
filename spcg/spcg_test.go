package spcg_test

import (
	"math/rand/v2"
	"testing"

	"github.com/DeedleFake/sirdsc/spcg"
)

func TestNext(t *testing.T) {
	s1 := rand.Uint64()
	s2 := rand.Uint64()

	pcg := rand.NewPCG(s1, s2)

	for i := range 100000 {
		var sn uint64
		sn, s1, s2 = spcg.Next(s1, s2)
		n := pcg.Uint64()
		if sn != n {
			t.Fatalf("%v != %v, i = %v", sn, n, i)
		}
	}
}
