package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "csv file with format 'question,answer'")
	limit := flag.Int64("timer", 3, "time limit per question, by default 3 second")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Println("Failed to open CSV File")
		os.Exit(1)
	}

	readCSV := csv.NewReader(file)
	lines, err := readCSV.ReadAll()
	if err != nil {
		exitProgram("Failed to parse CSV File")
	}

	fmt.Println("Quiz Game")

	problems := parseLines(lines)
	correct := 0
	chanAnswer := make(chan string)
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s \n", i+1, problem.question)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			chanAnswer <- answer
		}()
		select {
		case ans := <-chanAnswer:
			if ans == problem.answer {
				correct++
			}
		case <-time.After(time.Second * time.Duration(*limit)):
			exitProgram(fmt.Sprintf("Times up, your score %d of %d.\n", correct, len(problems)))
		}

	}

	exitProgram(fmt.Sprintf("Your score %d of %d.\n", correct, len(problems)))
}

type problem struct {
	question string
	answer   string
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))
	for i, line := range lines {
		result[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return result
}

func exitProgram(message string) {
	fmt.Println(message)
	os.Exit(1)
}
