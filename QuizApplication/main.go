package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

var (
	fileName *string
	timer    *int
	shuffle  *bool
	wg       sync.WaitGroup
)

func init() {
	fmt.Println("Welcome to the Quiz App")
	fileName = flag.String("csv", "problems.csv", "File name of input csv file which contains the problem set.")
	timer = flag.Int("timer", 30, "Set the timer to end the quiz.")
	shuffle = flag.Bool("shuffle", false, "Set shuffle flag to true to shuffle the problem set.")
	flag.Parse()
}

func main() {
	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("Error opening file:", *fileName)
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new CSV reader
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	problemSet := defineProblemSet(records, *shuffle)
	userResponse := make(map[int]string, len(problemSet))

	score := 0
	incorrect := 0
	timeUp := time.After(time.Second * time.Duration(*timer))
label:
	for i, problem := range problemSet {
		answerChan := make(chan string)
		fmt.Printf("Problem #%d: %s = ", i+1, problem.Question)
		go func() {
			var userAnswer string
			fmt.Scan(&userAnswer)
			answerChan <- userAnswer
		}()
		select {
		case <-timeUp:
			fmt.Println("\nYour time is up!!")
			break label
		case userAnswer, ok := <-answerChan:
			if ok {
				userResponse[i] = userAnswer
			} else {
				break label
			}
		}
	}

	for index := range userResponse {
		if checkAnswer(userResponse[index], problemSet[index].Answer) {
			score++
		} else {
			incorrect++
		}
	}

	fmt.Println("Your total score is: ", score)
	fmt.Println("You have answered", incorrect, "questions incorrectly")
}

func defineProblemSet(csvRows [][]string, shuffle bool) []Problem {
	lines := make([]Problem, len(csvRows))

	// Shuffle the questions if shuffle is set to true.
	perm := rand.Perm(len(lines))
	if !shuffle {
		sort.Ints(perm[:])
	}

	// Define the problems in the problem set.
	for row, data := range csvRows {
		lines[perm[row]] = Problem{
			Question: data[0],
			Answer:   strings.TrimSpace(data[1]),
		}
	}

	return lines
}

func checkAnswer(ans string, expected string) bool {
	if strings.EqualFold(strings.TrimSpace(ans), strings.TrimSpace(expected)) {
		return true
	}
	return false
}
