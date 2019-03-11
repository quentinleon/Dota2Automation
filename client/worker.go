package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/korovkin/limiter"
)

const usage = "usage: ./client ip port"
const maxCont = 4

var ip string
var port string
var homeDir string
var volume = "/dota:/dota"
var mountPath string
var cont = 0
var done = false

func url() string {
	return strings.Join([]string{"http://", ip, ":", port}, "")
}

func urlPath(path string) string {
	return strings.Join([]string{url(), path}, "/")
}

func populateBotDir(body []byte) {
	os.RemoveAll(homeDir + "/dota/game/dota/scripts/vscripts/bots")
	if err := ioutil.WriteFile("/tmp/bots.tar.gz", body, 0644); err != nil {
		panic(err)
	}
	cmd := exec.Command("tar", "xpf", "/tmp/bots.tar.gz", "-C", homeDir+"/dota/game/dota/scripts/vscripts/")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

}

func aquireBots() {
	resp, err := http.Get(urlPath("bots"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if string(body) == "no" {
		return
	}
	populateBotDir(body)
}

func requestStart() bool {
	resp, err := http.Get(urlPath("new"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	fmt.Println(string(body))
	if string(body) == "yes" {
		return true
	}
	done = true
	return false
}

func postDone(s string) {
	req, err := http.NewRequest("POST", urlPath("done"), bytes.NewBuffer([]byte(s)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp.Body.Close()
}

func runGame() {
	if requestStart() {
		out, err := exec.Command("/usr/local/bin/docker", "run", "-v", mountPath, "arwn/dota").Output()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//postDone(string(cmd)[:len(string(cmd))-1])
		/* post json file */
		postDone(string(out))
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println(usage)
		return
	}
	ip = os.Args[1]
	port = os.Args[2]

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir = usr.HomeDir
	mountPath = homeDir + volume

	aquireBots()

	limit := limiter.NewConcurrencyLimiter(4)
	for !done {
		limit.Execute(runGame)
		time.Sleep(time.Second * 3) /* because other computers are slow */
	}
	limit.Wait()
}
