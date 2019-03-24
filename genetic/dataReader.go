package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type lanePosition int

const (
	mid     lanePosition = 0
	carry   lanePosition = 1
	support lanePosition = 2
)

//HeroEvaluation : hero game data for evalutaion
type HeroEvaluation struct {
	heroName string
	lane     lanePosition
	winRate  float64
	kpm      float64
	dpm      float64
	apm      float64
	gpm      float64
	xpm      float64
}

type rawGameData struct {
	GameDuration       float64     `json:"gameDuration"`
	RunTime            float64     `json:"runTime"`
	DestroyedBuildings []string    `json:"destroyedBuildings"`
	Winner             string      `json:"winner"`
	Goodguys           rawTeamData `json:"goodguys"`
}

type rawTeamData struct {
	Bane         rawHeroData `json:"npc_dota_hero_bane(3)"`
	SkeletonKing rawHeroData `json:"npc_dota_hero_skeleton_king(42)"`
	Medusa       rawHeroData `json:"npc_dota_hero_medusa(94)"`
	Jakiro       rawHeroData `json:"npc_dota_jakiro(64)"`
	OgreMagi     rawHeroData `json:"npc_dota_hero_ogre_magi(84)"`
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

func ReadDotaFiles(fileDir string) ([5]HeroEvaluation, int) {
	var files []string

	err := filepath.Walk(fileDir,
		func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
	if err != nil {
		fmt.Println(" file walking failed")
	}

	var teamEval [5]HeroEvaluation

	teamEval[0].lane = carry
	teamEval[1].lane = support
	teamEval[2].lane = mid
	teamEval[3].lane = carry
	teamEval[4].lane = support
	teamEval[0].heroName = Roaster[0]
	teamEval[1].heroName = Roaster[1]
	teamEval[2].heroName = Roaster[2]
	teamEval[3].heroName = Roaster[3]
	teamEval[4].heroName = Roaster[4]

	var winCount int32
	var gameCount int32
	for _, file := range files {
		if file != fileDir {
			rawData, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Println(file + ": read failed")
			} else {
				var rawJSON rawGameData
				err := json.Unmarshal([]byte(rawData), &rawJSON)
				if err != nil {
					fmt.Println(file + " is not a valid JSON file.")
				} else {
					gameCount++
					if rawJSON.Winner == "goodguys" {
						winCount++
					}
					totalMinutes := rawJSON.GameDuration / 60.0
					addFields(&teamEval[0], &rawJSON.Goodguys.OgreMagi, totalMinutes)
					addFields(&teamEval[1], &rawJSON.Goodguys.Jakiro, totalMinutes)
					addFields(&teamEval[2], &rawJSON.Goodguys.Medusa, totalMinutes)
					addFields(&teamEval[3], &rawJSON.Goodguys.SkeletonKing, totalMinutes)
					addFields(&teamEval[4], &rawJSON.Goodguys.Bane, totalMinutes)
				}
			}
		}
	}
	if gameCount == 0 {
		return teamEval, 0
	}
	winRate := float64(winCount) / float64(gameCount)
	teamEval[0].winRate = winRate
	teamEval[1].winRate = winRate
	teamEval[2].winRate = winRate
	teamEval[3].winRate = winRate
	teamEval[4].winRate = winRate
	averageFields(&teamEval[0], float64(gameCount))
	averageFields(&teamEval[1], float64(gameCount))
	averageFields(&teamEval[2], float64(gameCount))
	averageFields(&teamEval[3], float64(gameCount))
	averageFields(&teamEval[4], float64(gameCount))
	return teamEval, 1
}

type geneFile struct {
	fileName string
	fitness  float64
}

type Top5genes struct {
	gene          [5]geneFile
	numberOfGenes int32
}

func (ranking Top5genes) Len() int {
	return len(ranking.gene)
}
func (ranking Top5genes) Swap(i, j int) {
	temp := ranking.gene[i]
	ranking.gene[i] = ranking.gene[j]
	ranking.gene[j] = temp
}

//it's actually greater.
func (ranking Top5genes) Less(i, j int) bool {
	return ranking.gene[i].fitness > ranking.gene[j].fitness
}

func assignHeroData(hero *HeroEvaluation, ranking *Top5genes, heroName string, geneDir string) {
	var gene geneFile
	gene.fileName = geneDir + "/genes/gene_" + heroName + ".lua"
	gene.fitness = CalcFitness(hero)

	if ranking.numberOfGenes < 5 {
		ranking.gene[ranking.numberOfGenes] = gene
		ranking.numberOfGenes++
	} else if gene.fitness > ranking.gene[4].fitness {
		ranking.gene[4] = gene
		sort.Sort(ranking)
	}
}

func FindTop5(path string) [5]Top5genes {
	//ogremagi
	//lich
	//medusa
	//chaosknight
	//bane
	var geneResults [5]Top5genes

	fileDir, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	geneDataNum := 0
	for _, geneDir := range fileDir {
		if !geneDir.IsDir() {
			fmt.Println(path + geneDir.Name() + " is not a directory")
		} else {
			result, err := ReadDotaFiles(path + "/" + geneDir.Name() + "/gamedata")
			if err != 1 {
				fmt.Println(path + "/" + geneDir.Name() + ": no games were parsed")
			} else {
				for i := range Roaster {
					assignHeroData(&result[i], &geneResults[i], Roaster[i], path+"/"+geneDir.Name())
				}
				geneDataNum++
			}
		}
	}
	return geneResults
}

//if its above the minimum, knoock it out, find the minimum
