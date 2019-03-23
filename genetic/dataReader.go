package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type lanePosition int

const (
	mid     lanePosition = 0
	carry   lanePosition = 1
	support lanePosition = 2
)

type HeroEvaluation struct {
	heroName	string
	lane    	lanePosition
	winRate 	float64
	kpm     	float64
	dpm     	float64
	apm     	float64
	gpm     	float64
	xpm     	float64
}

type rawGameData struct {
	GameDuration       float64     `json:"gameDuration"`
	RunTime            float64     `json:"runTime"`
	DestroyedBuildings []string    `json:"destroyedBuildings"`
	Winner             string      `json:"winner"`
	Goodguys           rawTeamData `json:"goodguys"`
}

type rawTeamData struct {
	Bane        rawHeroData `json:"npc_dota_hero_bane(3)"`
	ChaosKnight rawHeroData `json:"npc_dota_hero_chaos_knight(81)"`
	Medusa      rawHeroData `json:"npc_dota_hero_medusa(94)"`
	Lich        rawHeroData `json:"npc_dota_hero_lich(31)"`
	OgreMagi    rawHeroData `json:"npc_dota_hero_ogre_magi(84)"`
}

type rawHeroData struct {
	Kill       float64 `json:"kill"`
	Death      float64 `json:"death"`
	Assist     float64 `json:"assist"`
	XpPerMin   float64 `json:"xpPerMin"`
	GoldPerMin float64 `json:"goldPerMin"`
}

func addFields(hero *HeroEvaluation, data *rawHeroData, totalMinutes float64) {
	hero.kpm += data.Kill / totalMinutes
	hero.dpm += data.Death / totalMinutes
	hero.apm += data.Assist / totalMinutes
	hero.gpm += data.XpPerMin
	hero.xpm += data.GoldPerMin
}

func averageFields(hero *HeroEvaluation, nGames float64) {
	hero.kpm /= nGames
	hero.dpm /= nGames
	hero.apm /= nGames
	hero.gpm /= nGames
	hero.xpm /= nGames
}

func ReadDotaFiles(fileDir string) [5]HeroEvaluation {
	var files []string

	err := filepath.Walk(fileDir,
		func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
	if err != nil {
		fmt.Println(" file walking failed")
	}
	var rawJSON rawGameData
	var teamEval [5]HeroEvaluation

	//ogremagi
	//lich
	//medusa
	//chaosknight
	//bane
	teamEval[0].lane = carry
	teamEval[1].lane = support
	teamEval[2].lane = mid
	teamEval[3].lane = carry
	teamEval[4].lane = support
	teamEval[0].heroName = "ogre_magi"
	teamEval[1].heroName = "lich"
	teamEval[2].heroName = "medusa"
	teamEval[3].heroName = "chaos_knight"
	teamEval[4].heroName = "bane"
	
	var winCount int32
	var gameCount int32

	for _, file := range files {
		rawData, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(file + " read failed")
		} else {
			err := json.Unmarshal([]byte(rawData), &rawJSON)
			if err != nil {
				fmt.Println(file + " is not a valid JSON file.")
			} else {
				gameCount++
				if rawJson.Winner == "goodguys" {
					winCount++
				}
				totalMinutes := rawJson.GameDuration / 60.0
				addFields(&teamEval.bane, &rawJSON.Goodguys.Bane, totalMinutes)
				addFields(&teamEval.chaosKnight, &rawJSON.Goodguys.ChaosKnight, totalMinutes)
				addFields(&teamEval.medusa, &rawJSON.Goodguys.Medusa, totalMinutes)
				addFields(&teamEval.lich, &rawJSON.Goodguys.Lich, totalMinutes)
				addFields(&teamEval.ogreMagi, &rawJSON.Goodguys.OgreMagi, totalMinutes)
			}
		}
	}
	if gameCount == 0 {
		fmt.Println("no games were parsed")
	} else {
		winRate := float64(winCount) / float64(gameCount)
		teamEval.bane.winRate = winRate
		teamEval.chaosKnight.winRate = winRate
		teamEval.medusa.winRate = winRate
		teamEval.lich.winRate = winRate
		teamEval.ogreMagi.winRate = winRate
		averageFields(&teamEval.bane, float64(gameCount))
		averageFields(&teamEval.chaosKnight, float64(gameCount))
		averageFields(&teamEval.medusa, float64(gameCount))
		averageFields(&teamEval.lich, float64(gameCount))
		averageFields(&teamEval.ogreMagi, float64(gameCount))
	}
	return teamEval
}


type geneFile struct {
	fileName string
	fitness	float64
}

type top5genes struct {
	minimum       int32
	gene          [5]geneFile
	numberOfGenes int32
}

type topTeamGenes struct {
	bane			top5genes
	chaos_knight	top5genes
	medusa			top5genes
	lich			top5genes
	ogre_magi	top5genes
}

func assignHeroData(gene *geneFile, result *TeamEvaluation, heroName string, geneDir string) {
	gene.fileName = geneDir + "genefolder/gene_" + heroName + ".lua"
	gene.fitness = CalcFitness(result[heroName])

	if gene.fitness > top5genes[0].minimum
	{
		if top5genes[0].numberOfGenes < 5 {
			top5genes[0].gene[numberOfGenes] = gene

			top5genes[0].numberOfGenes.numberOfGenes++
			if top5genes[0].numberOfGenes == 5 {
				top5genes[0].minimum = findMinimum(top5genes[0])
			}
		} else {
			removeLast(top5genes[0])
			binarysearch(gene)
			top5genes[0].minimum = top5genes[0].gene[4].fitness
		}
	}

}

func findTop5() {
	// survivedGenes [5]top5genes
	var geneResults topTeamGenes

	//go through all game_data. folders for one bot and repeat
	//call 
	//save them as a string of multiple top 5 genes
	potentialParents := [5]string{"bane", "lich", "ogre_magi", "medusa", "chaos_knight"}//[5]string{"gene_bane.lua", "gene_lich.lua", "gene_ogre_magi.lua", "gene_medusa.lua", "gene_chaos_knight.lua"}
	fileDir, err := ioutil.ReadDir("./gene_pool")
	if err != nil {
		log.Fatal(err)
	}

	geneDataNum := 0
	var gene geneFile
	for bots := 0; bots < len(potentialParents); bots++ {
		for i,geneDir := range fileDir {
			result := ReadDotaFiles(geneDir + "gamedata")
			assignHeroData(&gene, &result, potentialParents[i], geneDir)
			geneDataNum++
		}
	}
}

//if its above the minimum, knoock it out, find the minimum
