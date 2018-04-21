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
	"math/rand"
)

const DEFAULT_QUIZ string = "./problems-in-order.csv"
const DELIMITER byte = '\n'

type QuestionAndAnswer struct {
	question string
	answer string
}

func main () {
	csvPath := flag.String("path", DEFAULT_QUIZ, "a path")
	timeoutSeconds := flag.Int("timeout", 30, "timeout for quiz")
	shuffleQuestions := flag.Bool("shuffle", true, "Shuffle Questions")
	flag.Parse()
	fmt.Println(*csvPath)
	file := openCsvFile(*csvPath)
	parsedCsv := parseCsvFileWithCSV(file)
	file.Close()
	if (*shuffleQuestions) {
		parsedCsv = shuffleQuestionsAndAnswers(parsedCsv)
	}
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

func shuffleQuestionsAndAnswers(input []QuestionAndAnswer) (output []QuestionAndAnswer) {
	rand.Seed(42)
	var usedIndexes [] int
	for i:=0; i < len(input); i++{
		indexToBePlaced := rand.Intn(len(input))
		for isInArray(indexToBePlaced, usedIndexes) {
			indexToBePlaced = rand.Intn(len(input))
		}
		usedIndexes = append(usedIndexes, indexToBePlaced)
		output = append(output, input[indexToBePlaced])
	}
	return output
}

func isInArray(value int, array []int) (isInArray bool){
	for _,arrayValue := range(array) {
		if value == arrayValue {
			return true
		}


	}
	return false
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
				if (isAnswerCorrect(answer, questionAndAnswer.answer)) {
					numberOfCorrectQuestions ++
				}
			}
		case <-time.After(time.Second * time.Duration(timeoutSeconds)):
			timeout = true
		}

		if timeout {
			break
		}

	}
	return len(questions), numberOfCorrectQuestions
}

func isAnswerCorrect(providedAnswer string, expectedAnswer string) bool {
	return strings.TrimSpace(strings.ToLower(providedAnswer)) == strings.TrimSpace(strings.ToLower(expectedAnswer))
}

func askQuestion (index int, question string) (answer string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Question #", index+1, question)
	answer, _ = reader.ReadString(DELIMITER)
	return answer[:len(answer)-1];
}

func printResults(numberOfQuestions int, numberOfCorrectAnswers int) {
	fmt.Println("Total number of questions:", numberOfQuestions)
	fmt.Println("Total number of correct answers:", numberOfCorrectAnswers)
}