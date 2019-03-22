package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const IP_ADDR = "10.10.10.1"
const PORT = 1042

const FIRST_WORKER = 2
const LAST_WORKER = 15

var workerTasks [LAST_WORKER + 1]int
var workerCompleted [LAST_WORKER + 1]int

var debug = false
var dataPath = ""
var botsPath = ""
var individualPath = ""
var requestedGames = 0
var remainingGames = 0
var runningGames = 0
var completedGames = 0
var finished sync.WaitGroup
var handling sync.WaitGroup

func RunGames(numGames int, individualGenesPath string, botScriptPath string) {
	var err error
	requestedGames = numGames
	individualPath = individualGenesPath
	botsPath = botScriptPath
	remainingGames = requestedGames
	PrepFiles()

	//make data directory
	if requestedGames > 0 {
		dataPath := individualPath + "/gamedata"
		os.Mkdir(dataPath, 0777)
	}

	//start server
	fmt.Printf("Running %d DotA games...\n\n", requestedGames)
	finished.Add(1)
	go StartServer()

	//spin up all clients
	for i := FIRST_WORKER; i <= LAST_WORKER; i++ {
		worker := GetWorkerIP(i)
		fmt.Printf("Starting worker %s\n", worker)
		cmd := exec.Command("ssh", worker, "./Dota2Automation/client/worker", IP_ADDR, strconv.Itoa(PORT), "&> /dev/null")
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println()

	//wait for last worker to finish
	finished.Wait()
	fmt.Println("\nDone!")
}

func StartServer() {
	http.HandleFunc("/new", HandleNew)
	http.HandleFunc("/done", HandleDone)
	http.HandleFunc("/bots", HandleBots)
	http.HandleFunc("/", HandleUnknown)
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
}

func HandleNew(w http.ResponseWriter, r *http.Request) {
	handling.Wait()
	handling.Add(1)
	if remainingGames > 0 {
		remainingGames--
		runningGames++
		workerTasks[IPToWorkerID(r.RemoteAddr)]++
		fmt.Fprintf(w, "yes")
		fmt.Printf("Worker %s starting new game.\n", r.RemoteAddr)
		printStatus()
	} else {
		fmt.Fprintf(w, "no")
	}
	handling.Done()
}

func HandleDone(w http.ResponseWriter, r *http.Request) {
	handling.Wait()
	handling.Add(1)
	runningGames--
	workerTasks[IPToWorkerID(r.RemoteAddr)]--
	completedGames++
	workerCompleted[IPToWorkerID(r.RemoteAddr)]++
	fmt.Printf("Worker %s finished game.\n", r.RemoteAddr)
	//get json file
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	filename := "game" + strconv.Itoa(completedGames) + ".json"
	err = ioutil.WriteFile(dataPath+"/"+filename, body, 0644)
	if err != nil {
		log.Fatal(err)
	}
	printStatus()
	if remainingGames == 0 && runningGames == 0 {
		finished.Done()
	}
	handling.Done()
}

func HandleBots(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving bot files to %s\n", r.RemoteAddr)
	if _, err := os.Stat("bots.tar.gz"); err == nil {
		http.ServeFile(w, r, "bots.tar.gz")
	} else {
		fmt.Fprintf(w, "no")
	}
}

func HandleUnknown(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Unknown request. GET: /new or POST: /done")
}

func PrepFiles() {
	//clear tmp
	cmd := exec.Command("rm", "-r", "/tmp/bots")
	err := cmd.Run()

	//copy to tmp
	cmd = exec.Command("cp", "-r", botsPath, "/tmp/bots")
	err = cmd.Run()

	//remove generics
	cmd = exec.Command("sh", "prepScript.sh")
	cmd.Run()

	//add genes into bots
	cmd = exec.Command("cp", "-r", individualPath+"/genes ", "/tmp/bots/genes")
	err = cmd.Run()

	//tar it
	cmd = exec.Command("tar", "-pcvzf", "bots.tar.gz", "-C", "/tmp/", "bots/")
	err = cmd.Run()
	//check for errors
	if err != nil {
		log.Fatal(err)
		return
	}
}

func printStatus() {
	if debug {
		return
	}
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Printf("%s Status:\n\n", strings.Split(dataPath, "/")[1])
	fmt.Printf("ID\tCompleted\t Running\n")
	for i := FIRST_WORKER; i <= LAST_WORKER; i++ {
		fmt.Printf("%d\t%d\t\t|", i, workerCompleted[i])
		for n := 0; n < workerTasks[i]; n++ {
			fmt.Print("# ")
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Printf("Requested: %d | Completed: %d | Running: %d | Remaining:%d \n", requestedGames, completedGames, runningGames, remainingGames)
}

func IPToWorkerID(ip string) int {
	ip = strings.Split(ip, ":")[0]
	ip = strings.Split(ip, ".")[3]
	workerNum, err := strconv.Atoi(ip)
	if err != nil {
		log.Println("Atoi failed")
	}
	return workerNum
}

func GetWorkerIP(n int) string {
	return fmt.Sprintf("dota%d@10.10.10.%d", n, n)
}

func PrintUsage(progName string) {
	fmt.Println("Usage: " + progName + " <num games> [bot script folder path]")
}
