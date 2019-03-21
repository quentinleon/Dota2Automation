package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const TEMPLATE = "gene_template.lua"

func GetGeneFromFile(filename string) []float64 {
	var gene []float64
	rawdata, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	data := string(rawdata)
	lines := strings.Split(data, "\n")
	for i := 0; i < len(lines); i++ {
		numberStr := strings.Split(lines[i], "=")
		if len(numberStr) > 1 && !strings.Contains(numberStr[1], "{") {
			if strings.Contains(numberStr[1], ",") {
				numberStr[1] = numberStr[1][:len(numberStr[1])-1]
			}
			if strings.Contains(numberStr[1], " ") {
				numberStr[1] = numberStr[1][1:]
			}

			number, err := strconv.ParseFloat(numberStr[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			gene = append(gene, number)
		}
	}
	return (gene)
}

func WriteGeneToFile(gene []float64, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	template, err := ioutil.ReadFile(TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(template), "\n")
	for i, g := 0, 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "%f") {
			fmt.Fprintln(file, fmt.Sprintf(lines[i], gene[g]))
			g++
		} else {
			fmt.Fprintln(file, lines[i])
		}
	}
	file.Close()
}

func MutateGeneRand(gene []float64, mutationRate float64, maxMutationPercent float64) {
	src := rand.NewSource(time.Now().UnixNano())
	rando := rand.New(src)
	for i := 0; i < len(gene); i++ {
		rate := rando.Float64()
		if rate < mutationRate {
			amount := (rando.Float64() * (maxMutationPercent * 2)) - maxMutationPercent
			gene[i] *= (1 + amount)
		}
	}
}

func main() {
	gene := GetGeneFromFile("gene.lua")
	fmt.Println(gene)
	MutateGeneRand(gene, .1, .1)
	fmt.Println(gene)
	WriteGeneToFile(gene, "result.txt")
}
