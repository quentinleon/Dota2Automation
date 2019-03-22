package main

import (
	"math/rand"
	"time"
)

func MutateGeneRand(gene []float64, mutationRate float64, maxMutationPercent float64) {
	src := rand.NewSource(time.Now().UnixNano())
	rando := rand.New(src)
	for i := 0; i < len(gene); i++ {
		rate := rando.Float64()
		if rate < mutationRate {
			amount := (rando.Float64() * (maxMutationPercent * 2)) - maxMutationPercent
			gene[i] *= (1 + amount)
		}
	}
}
