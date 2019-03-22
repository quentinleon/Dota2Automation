package main

import (
	"fmt"
)

var gamesPerIndividual = 250
var populationSize = 30
var parents = 5

func main() {
	//gene := GetGeneFromFile("gene.lua")
	//fmt.Println(gene)
	//MutateGeneRand(gene, .1, .1)
	//fmt.Println(gene)
	//WriteGeneToFile(gene, "result.txt")
	//fmt.Println(CalcFitness(1, 0.5, 0.28, 0.22, 450, 350))
	teamEval := ReadFiles("D:\\Dota2AI\\Dota2Automation\\game1553219414")
	fmt.Println(CalcFitness(teamEval.bane))
	fmt.Println(CalcFitness(teamEval.chaosKnight))
	fmt.Println(CalcFitness(teamEval.juggernaut))
	fmt.Println(CalcFitness(teamEval.lich))
	fmt.Println(CalcFitness(teamEval.ogreMagi))

	/*
		a := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		b := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

		c := SpliceBreed(a, b)

		fmt.Println()
		fmt.Println(a)
		fmt.Println(b)
		fmt.Println(c)

		c = ShuffleBreed(a, b)

		fmt.Println()
		fmt.Println(a)
		fmt.Println(b)
		fmt.Println(c)

		c = AverageBreed(a, b)

		fmt.Println()
		fmt.Println(a)
		fmt.Println(b)
		fmt.Println(c)

		fmt.Println()
		gene := GetGeneFromFile("gene.lua")
		fmt.Println(gene)
		MutateGeneRand(gene, .1, .1)
		fmt.Println(gene)
		WriteGeneToFile(gene, "result.txt")

		fmt.Println()
		fmt.Println(CalcMidFitness(1, 0.5, 0.28, 0.22, 450, 350))
	*/
}
