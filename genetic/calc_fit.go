package main

import (
	"fmt"
)

func CalcFitness(winRate float64, kpm float64, dpm float64, apm float64, gpm float64, xpm float64) float64 {
	score := 0.0
	score += winRate * 1000
	score += kpm * 500
	score -= dpm * 800
	score += apm * 200
	score += gpm / 10
	score += xpm / 10
	/*
		for i, _ := range towerSituation {
			if strings.Contains(towerSituation[i], "goodguys") {
				if strings.Contains(towerSituation[i], "rax") {
					score -= 200
				}
			}  else if strings.Contains(towerSituation[i], "badguys") {
				if strings.Contains(towerSituation[i], "rax") {

				}
			}
		}
	*/
	return score
}

func mein() {
	//towerSituation := []string{"goodguys", "badguys"}
	fmt.Println(CalcFitness(1, 0.5, 0.28, 0.22, 450, 350))
}
