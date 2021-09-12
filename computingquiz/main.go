package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var scanner = bufio.NewScanner(os.Stdin)

func readQuizFromCSV(fileName string) (map[string]string, error) {
	f, err := os.Open(fmt.Sprintf("%s", fileName))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	quiz := make(map[string]string)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(record) < 2 {
			continue
		}
		quiz[record[0]] = record[1]
	}

	return quiz, nil
}

func askQuestions(quiz map[string]string, timeLimit int) int {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	done := make(chan string)

	noOfCorrectAns := 0
	go func() {
		for question, answer := range quiz {
			fmt.Println(fmt.Sprintf("Q: %s", question))
			fmt.Print("A: ")

			var givenAns string
			scanner.Scan()
			givenAns = scanner.Text()

			if givenAns == answer {
				noOfCorrectAns++
			}
		}

		done <- "done"
	}()

	select {
	case <-done:
	case <-timer.C:
		return noOfCorrectAns
	}

	return noOfCorrectAns
}

func main() {
	fileNamePtr := flag.String("f", "computingquiz", "File name that contains the quiz questions")
	timeLimitPtr := flag.Int("t", 30, "Time Limit to  answer questions")
	flag.Parse()

	quiz, _ := readQuizFromCSV(fmt.Sprintf("%s.csv", *fileNamePtr))

	fmt.Print("Tap any key to start")
	scanner.Scan()

	fmt.Println("Started...")
	noOfCorrectAns := askQuestions(quiz, *timeLimitPtr)

	fmt.Printf("\nYou got %d questions correct\n", noOfCorrectAns)
}
