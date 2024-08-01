package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

// function to exit in case of an error
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv in the formart of 'question, answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file : %s\n", *csvFilename))
	}
	defer file.Close()
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided csv file")
	}
}
