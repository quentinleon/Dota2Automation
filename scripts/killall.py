import os

firstW = 2
lastW = 15

for i in range(firstW, lastW + 1, 1):
    worker = "dota" + str(i) + "@10.10.10." + str(i)
    command = "ssh " + worker + " '/usr/local/bin/docker kill $(/usr/local/bin/docker ps -q) &> /dev/null'"
    os.system(command)