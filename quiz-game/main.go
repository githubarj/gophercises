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

// struct to store quiz problems with answers
type problem struct {
	question string
	answer   string
}

// function to create a slice that contains the lines as structs
func parseLines(lines [][]string) []problem {
	parsed := make([]problem, len(lines))

	for index, line := range lines {
		parsed[index] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return parsed
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
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

	// parsing the csv lines into a slice of the problem struct i created
	problems := parseLines(lines)

	correct := 0
	for index, p := range problems {
		fmt.Printf("Question %d : %s \n", index+1, p.question)
		answer := getUserInput()
		if answer == p.answer {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}
