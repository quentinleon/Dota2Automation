package main

import (
	"os"
	"strconv"
)

var Roaster = [5]string{"ogre_magi", "lich", "medusa", "chaos_knight", "bane"}

var gamesPerIndividual = 250
var populationSize = 30
var parents = 5

func main() {
	var gen Generation
	var genePool [5][][]float64
	var scriptPool [][5][]float64
	gen.path = "D:\\Dota2AI\\Dota2AutomationFork\\genetic\\data\\G1-Test"
	top5 := FindTop5(gen)
	convertedTop5 := ConvertGeneData(top5)
	WriteBestGenes(gen, convertedTop5)

	for i := 0; i < 5; i++ {
		genePool[i] = MixIn(convertedTop5[i], populationSize)
	}
	scriptPool = CombineGenes(genePool)

	testing := gen.path + "\\script_pool"
	os.Mkdir(testing, 0777)
	for i := range scriptPool {
		dir := testing + "\\" + strconv.FormatInt(int64(i), 10)
		os.Mkdir(dir, 0777)
		for j := 0; j < 5; j++ {
			WriteGeneToFile(scriptPool[i][j], dir+"\\gene_"+Roaster[j]+".lua")
		}
	}
}
