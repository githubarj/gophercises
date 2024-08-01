package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
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

func shuffleProblems(problems []problem) {
	rand.Seed(time.Now().UnixNano())
	for i := len(problems) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		problems[i], problems[j] = problems[j], problems[i]
	}
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv in the formart of 'question, answer'")
	shuffleQuiz := flag.Bool("shuffle", false, "shuffle the quiz order")
	timeLimit := flag.Int("limit", 30, "set the time limit for the quiz in seconds")
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

	if *shuffleQuiz {
		shuffleProblems(problems)
	}

	fmt.Println("Press enter to start the quiz. You have", *timeLimit, "seconds")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	problemChan := make(chan *problem)
	answerChan := make(chan string)

	go func() {
		for _, p := range problems {
			problemChan <- &p
			<-answerChan
		}
		close(problemChan)
	}()

	go func() {
		for p := range problemChan {
			fmt.Printf("%s", p.question)
			answerChan <- getUserInput()
		}
	}()

quizLoop:
	for i := 0; i < len(problems); i++ {
		select {
		case <-timer.C:
			fmt.Println("\n Time is up!")
			break quizLoop
		case p := <-problemChan:
			answer := <-answerChan
			if answer == p.answer {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}
