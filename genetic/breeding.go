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
	for i := range a {
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

func refillGenes(dest *[][]float64, src *[][]float64) {
	if len(*src) == 0 {
		log.Fatal("source genes are empty")
	}

	*dest = make([][]float64, len(*src))
	for i := range *dest {
		(*dest)[i] = make([]float64, len((*src)[i]))
		copy((*dest)[i], (*src)[i])
	}
}

func MixIn(genes [][]float64, numberOfGenes int) [][]float64 {
	if len(genes) < 2 {
		log.Fatal("Less than two genes in pool")
	}

	src := rand.NewSource(time.Now().UnixNano())
	rando := rand.New(src)
	var nonReproducedGenes [][]float64
	var result [][]float64
	for i := 0; i < numberOfGenes; i++ {
		if len(nonReproducedGenes) < 2 {
			refillGenes(&nonReproducedGenes, &genes)
			rando.Shuffle(len(nonReproducedGenes), func(i, j int) {
				var temp float64
				for k := range nonReproducedGenes[i] {
					temp = nonReproducedGenes[i][k]
					nonReproducedGenes[i][k] = nonReproducedGenes[j][k]
					nonReproducedGenes[j][k] = temp
				}
			})
		}
		gene := SpliceBreed(nonReproducedGenes[len(nonReproducedGenes)-1], nonReproducedGenes[len(nonReproducedGenes)-2])
		MutateGeneRand(gene, 0.1, 0.1)
		result = append(result, gene)
		nonReproducedGenes = nonReproducedGenes[:len(nonReproducedGenes)-2]
	}
	return result
}

func CombineGenes(genes [5][][]float64) [][5][]float64 {
	if len(genes[0]) != len(genes[1]) || len(genes[1]) != len(genes[2]) ||
		len(genes[2]) != len(genes[3]) || len(genes[3]) != len(genes[4]) {
		log.Fatal("number of gene pool do not match with other heroes")
	}
	result := make([][5][]float64, len(genes[0]))
	for i := range genes[0] {
		for j := 0; j < 5; j++ {
			result[i][j] = make([]float64, len(genes[j][i]))
			for k := range result[i][j] {
				result[i][j][k] = genes[j][i][k]
			}
		}
	}
	return result
}
