import numpy
import json
import sys
import os
import numpy as np
import matplotlib.pyplot as plot

if len (sys.argv) < 2:
	print ("Usage: python " + argv[0] + " <data folder>")
	exit
dataDir = sys.argv[1]

winCount = 0.0
totalGames = 0

longestGame = 0
winTimes = []
lossTimes = []

maxKills = 0

heroes = {}

class Hero:
	def __init__(self):
		self.kills = []
		self.assists = []
		self.deaths = []
		self.lastHits = []
		self.denies = []
		self.levels = []
		self.goldPerMinute = []
		self.xpPerMinue = []

	kills = []
	assists = []
	deaths = []
	lastHits = []
	denies = []
	levels = []
	goldPerMinute = []
	xpPerMinue = []

for fileName in os.listdir(dataDir):
	try:
		f = open(dataDir+fileName)
		data = json.load(f)
	except:
		print (fileName + " is not a valid JSON file.")
		continue
	
	for hero in data["goodguys"]:
		if totalGames == 0:
			#print (hero)
			heroes[hero] = Hero()
		if data["goodguys"][hero]["kill"] > maxKills:
			maxKills = data["goodguys"][hero]["kill"]
		heroes[hero].kills.append(data["goodguys"][hero]["kill"])
		heroes[hero].assists.append(data["goodguys"][hero]["assist"])
		heroes[hero].deaths.append(data["goodguys"][hero]["death"])
		heroes[hero].lastHits.append(data["goodguys"][hero]["lastHit"])
		heroes[hero].denies.append(data["goodguys"][hero]["deny"])
		heroes[hero].levels.append(data["goodguys"][hero]["level"])
		heroes[hero].goldPerMinute.append(data["goodguys"][hero]["goldPerMin"])
		heroes[hero].xpPerMinue.append(data["goodguys"][hero]["xpPerMin"])
		

	totalGames += 1
	timeInMin = int(data["gameDuration"]) / 60
	if (timeInMin > longestGame):
		longestGame = timeInMin
	if data["winner"] == "goodguys":
		winCount += 1
		winTimes.append(timeInMin)
	else:
		lossTimes.append(timeInMin)

print("")

if totalGames == 0:
	print("There is no gamedata in " + dataDir)
	sys.exit()

print ("Total Games: " + str(totalGames))
print ("Win Percentage: " + str(round((winCount/totalGames) * 100, 2)) + "%")

print
bins = np.linspace(0, maxKills, maxKills)

avgKillTotal = 0
avgDeathTotal = 0
avgAssistTotal = 0

for key in heroes:
	hero = heroes[key]
	print(key)
	avgKillTotal += np.mean(hero.kills)
	avgDeathTotal += np.mean(hero.deaths)
	avgAssistTotal += np.mean(hero.assists)
	print ("\tK:" + str(round(np.mean(hero.kills), 1)) + " " + str(np.median(hero.kills)) + "\n\tD:" + str(round(np.mean(hero.deaths), 1)) + " " + str(np.median(hero.deaths)) + "\n\tA:" + str(round(np.mean(hero.assists), 1)) + " " + str(np.median(hero.assists)) + "\n")
	#plot.hist(hero.deaths, bins, label="Kills", color="blue", alpha=1)
	#plot.show()

print ("Total Avg. Kills: " + str(round(avgKillTotal, 2)))
print ("Total Avg. Deaths: " + str(round(avgDeathTotal, 2)))
print ("Total Avg. Assists: " + str(round(avgAssistTotal, 2)))

bins = np.linspace(0, longestGame, longestGame/2)
plot.xlabel = "Game Duration"
plot.ylabel = "Num Wins/Losses"
if len(winTimes) > 0:
	plot.hist(winTimes, bins, label="Wins", color="green", alpha=1, ec="black")
if len(lossTimes) > 0:
	plot.hist(lossTimes, bins, label="Losses", color="red", alpha=0.7, ec="black")
plot.legend(loc='upper left')
plot.show()


