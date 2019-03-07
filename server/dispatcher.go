package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

const IP_ADDR = "10.10.10.1"
const PORT = 1042
const MOUNT = "/Volumes/AOEU/dota/:/dota/"

var WORKERS = []string{"dota2@10.10.10.2", "dota3@10.10.10.3", "dota4@10.10.10.4", "dota5@10.10.10.5", "dota6@10.10.10.6"} //"dota7@10.10.10.7", "dota8@10.10.10.8", "dota9@10.10.10.9"}

var remainingGames = 0
var runningGames = 0
var completedGames = 0
var finished sync.WaitGroup

func main() {
	args := os.Args
	if len(args) < 2 {
		PrintUsage(args[0])
		return
	}
	var err error
	remainingGames, err = strconv.Atoi(args[1])
	if err != nil {
		PrintUsage(args[0])
		return
	}

	fmt.Printf("Running %d DotA games...\n\n", remainingGames)
	if remainingGames == 0 {
		return
	}
	//copy bot scripts to nfs and set up for bot games

	//start server
	finished.Add(1)
	go StartServer()

	//spin up all clients
	for index, worker := range WORKERS {
		fmt.Printf("Starting worker %d at %s\n", index+1, worker)
		//fmt.Printf("%s %s %s %s %s %s %s\n", "ssh", worker, "./Dota2Automation/client/client", IP_ADDR, strconv.Itoa(PORT), MOUNT, "&> /dev/null")
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
	http.HandleFunc("/", HandleUnknown)
	http.ListenAndServe(":1042", nil)
}

func HandleNew(w http.ResponseWriter, r *http.Request) {
	if remainingGames > 0 {
		remainingGames--
		runningGames++
		fmt.Fprintf(w, "yes")
		fmt.Printf("Worker %s starting new game.\n", r.RemoteAddr)
	} else {
		fmt.Fprintf(w, "no")
	}
}

func HandleDone(w http.ResponseWriter, r *http.Request) {
	runningGames--
	completedGames++
	fmt.Printf("Worker %s finished game.\n", r.RemoteAddr)
	//get json file
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	filename := "game" + strconv.Itoa(completedGames) + ".json"
	err = ioutil.WriteFile("testdata/"+filename, body, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if remainingGames == 0 && runningGames == 0 {
		finished.Done()
	}
}

func HandleUnknown(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Unknown request. GET: /new or POST: /done")
}

func PrintUsage(progName string) {
	fmt.Println("Usage: " + progName + " <num games> [bot script folder path]")
}
