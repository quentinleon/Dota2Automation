package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var Roaster = [5]string{"ogre_magi", "jakiro", "medusa", "skeleton_king", "bane"}

var populationSize = 30

var root = "/Users/dota1/Downloads/Dota2Automation-master-2/genetic/data"
var initialGenes = root + "/InitialGenes"

func PrintUsage(progName string) {
	fmt.Println("Usage: " + progName + " <num games per individual> [bot script folder path] [generation#]")
}

func main() {
	args := os.Args
	if len(args) < 4 {
		PrintUsage(args[0])
		return
	}
	gamesPerIndividual, err := strconv.Atoi(args[1])
	if err != nil || gamesPerIndividual <= 0 {
		log.Fatal("number of games incorrect")
	}
	botScriptPath := args[2]
	currentGeneration, err := strconv.Atoi(args[3])
	if err != nil || currentGeneration <= 0 {
		log.Fatal("generation must start from 1")
	}

	go StartServer()

	var genePool [5][][]float64
	var scriptPool [][5][]float64
	for true {
		previousGenerationPath := root + "/G" + strconv.Itoa(currentGeneration-1)
		currentGenerationPath := root + "/G" + strconv.Itoa(currentGeneration)
		if currentGeneration == 1 {
			for i := 0; i < 5; i++ {
				genePool[i] = make([][]float64, 2)
				genePool[i][0] = GetGeneFromFile(initialGenes + "/gene_" + Roaster[i] + ".lua")
				genePool[i][1] = make([]float64, len(genePool[i][0]))
				copy(genePool[i][1], genePool[i][0])
				genePool[i] = MixIn(genePool[i], populationSize)
			}
			scriptPool = CombineGenes(genePool)
		} else {
			top5 := FindTop5(previousGenerationPath)
			convertedTop5 := ConvertGeneData(top5)
			WriteBestGenes(previousGenerationPath, convertedTop5)
			for i := 0; i < 5; i++ {
				genePool[i] = MixIn(convertedTop5[i], populationSize)
			}
			scriptPool = CombineGenes(genePool)
		}
		os.Mkdir(currentGenerationPath, 0777)
		for i := range scriptPool {
			currentIndividualPath := currentGenerationPath + "/I" + strconv.Itoa(i)
			os.Mkdir(currentIndividualPath, 0777)
			os.Mkdir(currentIndividualPath+"/genes", 0777)
			for j := 0; j < 5; j++ {
				WriteGeneToFile(scriptPool[i][j], currentIndividualPath+"/genes/gene_"+Roaster[j]+".lua")
			}
		}
		for i := range scriptPool {
			currentIndividualPath := currentGenerationPath + "/I" + strconv.Itoa(i)
			RunGames(gamesPerIndividual, currentIndividualPath, botScriptPath)
		}
		currentGeneration++
	}
}
