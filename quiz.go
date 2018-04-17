package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
	"encoding/csv"
	"io"
	"flag"
	"time"
)

const DEFAULT_QUIZ string = "./problems.csv"
const DELIMITER byte = '\n'

type QuestionAndAnswer struct {
	question string
	answer string
}

func main () {
	csvPath := flag.String("path", "./problems.csv", "a path")
	timeoutSeconds := flag.Int("timeout", 30, "timeout for quiz")
	flag.Parse()
	fmt.Println(*csvPath)
	file := openCsvFile(*csvPath)
	parsedCsv := parseCsvFileWithCSV(file)
	file.Close()
	numberOfQuestions, numberOfCorrectQuestions := askQuestions(parsedCsv, *timeoutSeconds)
	printResults(numberOfQuestions, numberOfCorrectQuestions)

}

func openCsvFile(fileLocation string) (*os.File) {
	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(fileLocation + " : is not a valid file")
	}
	return file
}

func parseCsvFile(file *os.File) ([]QuestionAndAnswer) {
	var parsedCsv []QuestionAndAnswer
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString(DELIMITER)

		if err != nil {
			break
		}

		parsedCsv = append(parsedCsv, parseQuizLine(line))
	}
	return parsedCsv
}

func parseCsvFileWithCSV(file *os.File) ([]QuestionAndAnswer) {
	var parsedCsv []QuestionAndAnswer
	reader := bufio.NewReader(file)
	r := csv.NewReader(reader)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		parsedCsv = append(parsedCsv, QuestionAndAnswer{question:record[0], answer:record[1]})
	}
	return parsedCsv
}


func parseQuizLine(line string) (QuestionAndAnswer) {
	index := strings.Index(line, ",")
	return QuestionAndAnswer{
		question:line[:index],
		answer:line[index+1:len(line) -1],
	}
}

func askQuestions(questions []QuestionAndAnswer, timeoutSeconds int) (numberOfQuestions int, numberOfCorrectQuestions int){
	numberOfCorrectQuestions  = 0
	for index, questionAndAnswer := range questions {
		var timeout bool = false
		ch := make (chan string)
		// Don't know how to make this less ugly...
		go func() {
			ch <- askQuestion(index, questionAndAnswer.question)
		}()
		var answer string
		select {
			case answer = <- ch : {
				if (answer == questionAndAnswer.answer) {
					numberOfCorrectQuestions ++
				}
			}
			case <-time.After(time.Second * time.Duration(timeoutSeconds)): timeout = true
		}

		if timeout{
			break
		}

	}
	return len (questions), numberOfCorrectQuestions
}

func askQuestion (index int, question string) (answer string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Question #", index+1, question)
	answer, _ = reader.ReadString('\n')
	return answer[:len(answer)-1];
}

func printResults(numberOfQuestions int, numberOfCorrectAnswers int) {
	fmt.Println("Total number of questions:", numberOfQuestions)
	fmt.Println("Total number of correct answers:", numberOfCorrectAnswers)
}