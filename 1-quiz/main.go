package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {

	timeFlagPtr := flag.Int("time", 30, "timeout in seconds")
	flag.Parse()
	// read file problems.csv into a slice of strings
	file, err := os.Open("problems.csv")
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(file)
	stdinReader := bufio.NewReader(os.Stdin)

	// var t time.Duration = 30
	t := time.Duration(*timeFlagPtr)
	timeOut := time.After(t * time.Second)

	fmt.Printf("Welcome to the quiz!\nThe quiz time is %d seconds.\nPress Enter to start quiz\n", t)
	stdinReader.ReadString('\n')
	quiz(csvReader, stdinReader, timeOut)

}

func quiz(csvReader *csv.Reader, stdinReader *bufio.Reader, timeout <-chan time.Time) {
	questionCount := 0
	correctAnswersCount := 0
	userInputs := make(chan string)

quizLoop:
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break quizLoop
		}
		questionCount++
		question := record[0]
		correctAnswer := record[1]
		fmt.Printf("Question %d: %s\n", questionCount, question)

		go getUserInput(stdinReader, userInputs)
		select {
		case <-timeout:
			fmt.Println("You ran out of time!")
			for true {
				rest, _ := csvReader.ReadAll() // reads rest of lines
				questionCount += len(rest)
				break quizLoop
			}
		case input, _ := <-userInputs:
			userAnswer := strings.TrimSuffix(input, "\n")
			if correctAnswer == userAnswer {
				correctAnswersCount++
			}
		}
	}
	defer fmt.Printf("You scored %d out of %d.\n", correctAnswersCount, questionCount)
}

func getUserInput(stdinReader *bufio.Reader, userInputs chan<- string) {
	go func() {
		input, err := stdinReader.ReadString('\n')
		if err != nil {
			close(userInputs)
		}
		userInputs <- input
	}()
}
