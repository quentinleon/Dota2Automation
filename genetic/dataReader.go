package main

import (
	"encoding/json"
	"path/filepath"
	"os"
	"io/ioutil"
	"fmt"
)

type lanePosition int

const (
	mid  	lanePosition = 0
	carry   lanePosition = 1
	support lanePosition = 2
 )

type heroEvaluation struct {
	lane lanePosition
	winRate float64
	kpm float64
	dpm float64
	apm float64
	gpm float64
	xpm float64
}

type TeamEvaluation struct {
	bane heroEvaluation
	chaosKnight heroEvaluation
	juggernaut heroEvaluation
	lich heroEvaluation
	ogreMagi heroEvaluation
}

type rawGameData struct {
	GameDuration float64 `json:"gameDuration"`
	RunTime float64 `json:"runTime"`
	DestroyedBuildings []string `json:"destroyedBuildings"`
	Winner string `json:"winner"`
	Goodguys rawTeamData `json:"goodguys"`
}

type rawTeamData struct {
	Bane rawHeroData `json:"npc_dota_hero_bane(3)"`
	ChaosKnight rawHeroData `json:"npc_dota_hero_chaos_knight(81)"`
	Juggernaut rawHeroData `json:"npc_dota_hero_juggernaut(8)"`
	Lich rawHeroData `json:"npc_dota_hero_lich(31)"`
	OgreMagi rawHeroData `json:"npc_dota_hero_ogre_magi(84)"`
}

type rawHeroData struct {
	Kill float64 `json:"kill"`
	Death float64 `json:"death"`
	Assist float64 `json:"assist"`
	XpPerMin float64 `json:"xpPerMin"`
	GoldPerMin float64 `json:"goldPerMin"`
}

func addFields(hero *heroEvaluation, data *rawHeroData, totalMinutes float64) {
	hero.kpm += data.Kill / totalMinutes
	hero.dpm += data.Death / totalMinutes
	hero.apm += data.Assist / totalMinutes
	hero.gpm += data.XpPerMin
	hero.xpm += data.GoldPerMin
}

func averageFields(hero *heroEvaluation, nGames float64) {
	hero.kpm /= nGames
	hero.dpm /= nGames
	hero.apm /= nGames
	hero.gpm /= nGames
	hero.xpm /= nGames
}

func ReadFiles(fileDir string) TeamEvaluation {
	var files []string

	err := filepath.Walk(fileDir, 
		func(path string, info os.FileInfo, err error) error {
       		files = append(files, path)
        	return nil
		})
    if err != nil {
		fmt.Println(" file walking failed")
	}
	var rawJson rawGameData
	var teamEval TeamEvaluation

	teamEval.bane.lane = support
	teamEval.chaosKnight.lane = carry
	teamEval.juggernaut.lane = carry
	teamEval.lich.lane = support
	teamEval.ogreMagi.lane = mid

	var winCount int32
	var gameCount int32

	for _, file := range files {
		rawData, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(file + " read failed")
		} else {
			err := json.Unmarshal([]byte(rawData), &rawJson)
			if err != nil {
				fmt.Println(file + " is not a valid JSON file.")
			} else {
				gameCount++
				if rawJson.Winner == "goodguys" {
					winCount++
				}
				totalMinutes := rawJson.GameDuration / 60.0
				addFields(&teamEval.bane, &rawJson.Goodguys.Bane, totalMinutes)
				addFields(&teamEval.chaosKnight, &rawJson.Goodguys.ChaosKnight, totalMinutes)
				addFields(&teamEval.juggernaut, &rawJson.Goodguys.Juggernaut, totalMinutes)
				addFields(&teamEval.lich, &rawJson.Goodguys.Lich, totalMinutes)
				addFields(&teamEval.ogreMagi, &rawJson.Goodguys.OgreMagi, totalMinutes)
			}
		}
	}
	if (gameCount == 0) {
		fmt.Println("no games parsed")
	} else {
		winRate := float64(winCount) / float64(gameCount)
		teamEval.bane.winRate = winRate
		teamEval.chaosKnight.winRate = winRate
		teamEval.juggernaut.winRate = winRate
		teamEval.lich.winRate = winRate
		teamEval.ogreMagi.winRate = winRate
		averageFields(&teamEval.bane, float64(gameCount))
		averageFields(&teamEval.chaosKnight, float64(gameCount))
		averageFields(&teamEval.juggernaut, float64(gameCount))
		averageFields(&teamEval.lich, float64(gameCount))
		averageFields(&teamEval.ogreMagi, float64(gameCount))
	}
	return teamEval
}
