import  sys
import socket
import os

if (len(sys.argv) < 2):
    print("Usage: python " + sys.argv[0] + " <num games> [bot script folder path]")
    quit()

PORT = 1042
IP_ADDR = str(socket.gethostbyname(socket.gethostname()))
MOUNT = "/Volumes/AOEU/dota/:/dota/"
WORKERS = [
    "dota2@10.10.125.52"
]

NUM_GAMES = int(sys.argv[1])
if (len(sys.argv) > 2): BOT_SCRIPTS = sys.argv[2]
    

print ("Running " + str(NUM_GAMES) + " DotA games...")

def printStatus():
    print("Requested Games: " + str(NUM_GAMES) + "  Remaining Games: " + str(remainingGames) + "  Running Games: " + str(runningGames) + "  Completed Games: " + str(completedGames))

#open socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind((socket.gethostname(), PORT))
s.listen(5)

#TODO
#put new bot scripts onto nfs and set up executable for bot games

#spin up remote workers
print ""
for worker in WORKERS:
    command = "ssh " + worker + " './Dota2Automation/old-client/client.pl " + IP_ADDR + " " + str(PORT) + " " + MOUNT + " &> /dev/null' &"
    if NUM_GAMES > 0:
        os.system(command)
        print ("Started worker " + worker)
print ""

#accept connections until there are no more remaining or running games
#TODO make server handle multiple client connections at once
runningGames = 0
completedGames = 0
remainingGames = NUM_GAMES
while remainingGames > 0 or runningGames > 0:
    (client, address) = s.accept()
    #print("Connection from: " + str(address))
    data = client.recv(1024)
    #print("Recieved: " + data.rstrip())
    if "new" in data:
        if remainingGames > 0:
            client.sendall("yes")
            runningGames += 1
            remainingGames -= 1
        else:
            client.sendall("no")
    elif "done" in data:
        completedGames += 1
        if remainingGames > 0:
            client.sendall("yes")
            remainingGames -= 1
        else:
            client.sendall("no")
            runningGames -= 1
    printStatus()
    client.close()

s.close()

print ("")
print ("Done!")



