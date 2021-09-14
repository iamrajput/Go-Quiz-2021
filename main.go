package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "Time limit for the quiz")
	shuffle := flag.Bool("shuffle", false, "Shuffle the questions (default 'false')")

	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll() //read all the line in the csv
	if err != nil {
		exit("failed to open the csv file")
	}
	problems := parseLine(lines)

	//if user want to shuffle the problems
	if *shuffle {
		problems = randomize(problems)
	}

	//fmt.Println(problems)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	fmt.Printf("***Thank you for participating in our quiz you have %d*** \n", *timeLimit)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s= ", i+1, p.que)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s \n", &answer)
			answerCh <- answer //sending the anser to the channel
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d. \n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.ans {
				correct++
				fmt.Println("***Correct***")
			} else {
				fmt.Println("===InCorrect===")
			}

		}
	}
	fmt.Printf("You scored %d out of %d. \n", correct, len(problems))

}

func parseLine(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			que: line[0],
			ans: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	que string
	ans string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func randomize(problems []problem) []problem {
	source := rand.NewSource(time.Now().UnixNano()) // newSource need int64
	r := rand.New(source)

	for i := range problems {
		newPosition := r.Intn(len(problems) - 1)
		problems[i], problems[newPosition] = problems[newPosition], problems[i] //swap deck
	}
	return problems
}
