package main

import (
	"fmt"
)

var populationSize = 30
var parents = 5

func main() {
	gene := GetGeneFromFile("gene.lua")
	fmt.Println(gene)
	MutateGeneRand(gene, .1, .1)
	fmt.Println(gene)
	WriteGeneToFile(gene, "result.txt")
	fmt.Println(CalcMidFitness(1, 0.5, 0.28, 0.22, 450, 350))
}
