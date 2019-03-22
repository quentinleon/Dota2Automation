package main

func CalcMidFitness(winRate float64, kpm float64, dpm float64, apm float64, gpm float64, xpm float64) float64 {
	score := 0.0
	score += winRate * 1000
	score += kpm * 500
	score -= dpm * 900
	score += apm * 400
	score += (gpm / 10)
	score += (xpm / 10) * 0.9
	return score
}

func CalcCarryFitness(winRate float64, kpm float64, dpm float64, apm float64, gpm float64, xpm float64) float64 {
	score := 0.0
	score += winRate * 1000
	score += kpm * 500
	score -= dpm * 1000
	score += apm * 300
	score += (gpm / 10) * 1.2
	score += xpm / 10
	return score
}

func CalcSupportFitness(winRate float64, kpm float64, dpm float64, apm float64, gpm float64, xpm float64) float64 {
	score := 0.0
	score += winRate * 1000
	score += kpm * 200
	score -= dpm * 800
	score += apm * 600
	score += (gpm / 10) * 0.8
	score += xpm / 10
	return score
}
