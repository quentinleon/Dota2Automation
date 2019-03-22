package main

import (
	"log"
	"math/rand"
	"time"
)

func SpliceBreed(a, b []float64) []float64 {
	if len(a) != len(b) {
		log.Fatal("Genes are not the same length.")
	}

	var gene []float64
	src := rand.NewSource(time.Now().UnixNano())
	rando := rand.New(src)
	parentA := (rando.Float64() >= .5)
	sliceLen := 0
	for i := 0; i < len(a); i++ {
		if sliceLen == 0 {
			parentA = !parentA
			sliceLen = int(rando.Float64() * float64(len(a)) / 2)
		}

		if parentA {
			gene = append(gene, a[i])
		} else {
			gene = append(gene, b[i])
		}
		sliceLen--
	}
	return gene
}

func ShuffleBreed(a, b []float64) []float64 {
	if len(a) != len(b) {
		log.Fatal("Genes are not the same length.")
	}

	var gene []float64
	src := rand.NewSource(time.Now().UnixNano())
	rando := rand.New(src)
	for i := 0; i < len(a); i++ {
		if rando.Float64() >= .5 {
			gene = append(gene, a[i])
		} else {
			gene = append(gene, b[i])
		}
	}
	return gene
}

func AverageBreed(a, b []float64) []float64 {
	if len(a) != len(b) {
		log.Fatal("Genes are not the same length.")
	}

	var gene []float64
	for i := 0; i < len(a); i++ {
		gene = append(gene, (a[i]+b[i])/2)
	}
	return gene
}
