package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
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
