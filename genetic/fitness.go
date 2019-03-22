package main

func CalcFitness(hero heroEvaluation) float64 {
	score := 0.0
	score += hero.winRate * 1000
	if hero.lane == mid {
		score += hero.kpm * 500
		score -= hero.dpm * 900
		score += hero.apm * 400
		score += (hero.gpm / 10)
		score += (hero.xpm / 10) * 0.9
	} else if hero.lane == carry {
		score += hero.kpm * 500
		score -= hero.dpm * 1000
		score += hero.apm * 300
		score += (hero.gpm / 10) * 1.2
		score += hero.xpm / 10
	} else if hero.lane == support {
		score += hero.kpm * 200
		score -= hero.dpm * 800
		score += hero.apm * 600
		score += (hero.gpm / 10) * 0.8
		score += hero.xpm / 10
	}

	return score
}
