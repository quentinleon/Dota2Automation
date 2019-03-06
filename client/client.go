package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/korovkin/limiter"
)

const usage = "usage: ./client ip port"
const maxCont = 4

var ip string
var port string
var volume = "~/dota:/dota"
var cont = 0
var done = false

func url() string {
	return strings.Join([]string{"http://", ip, ":", port}, "")
}

func urlPath(path string) string {
	return strings.Join([]string{url(), path}, "/")
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
		cmd, err := exec.Command("/usr/local/bin/docker", "run", "-v", volume, "dota").Output()
		if err != nil {
			fmt.Println(err)
		}
		postDone(string(cmd)[:len(string(cmd))-1])
		/* post json file */
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println(usage)
		return
	}
	ip = os.Args[1]
	port = os.Args[2]

	limit := limiter.NewConcurrencyLimiter(4)
	for !done {
		limit.Execute(runGame)
		time.Sleep(time.Second * 3) /* because other computers are slow */
	}
	limit.Wait()
}
