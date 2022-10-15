package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func quiz(f *os.File, resultsChan chan<- int) {
	r := csv.NewReader(f)
	for {
		line, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		question, answer := line[0], line[1]
		fmt.Println(question)

		var userInput string
		fmt.Scanln(&userInput)

		if userInput == answer {
			resultsChan <- 1
		} else {
			resultsChan <- 0
		}
	}
	resultsChan <- -1
}

func timer(quizDuration time.Duration, resultsChan chan<- int) {
	quizEndTime := time.Now().Add(quizDuration)
	for {

		if time.Now().After(quizEndTime) {
			resultsChan <- -1
		}
	}
}

func main() {
	f, err := os.Open("C:\\coding\\gophersizes\\quiz\\problems.csv")

	quizDuration := time.Duration(5) * time.Second

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println("Press Enter to Start Quiz")
	fmt.Scanln()

	resultsChan := make(chan int)

	go timer(quizDuration, resultsChan)
	go quiz(f, resultsChan)

	totalAnswered := 0
	totalCorrect := 0
	finished := false

	for {
		result := <-resultsChan
		if result == -1 {
			finished = true
		} else {
			totalAnswered++
			totalCorrect += result
		}
		if finished == true {
			break
		}
	}

	fmt.Printf("You got %d correct out of %d", totalCorrect, totalAnswered)

}
