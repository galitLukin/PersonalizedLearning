package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	SelectedAnswers    = "selectedAnswers"
	Python             = "python3"
	PathToPythonScript = "./python/script.py"
)

// Question Data
type QuestionData struct {
	Question         Question
	QuestionInstance QuestionInstance
	User             User
}

// Question : Data for current question
type Question struct {
	Assignment      string   `json:"Assignment"`
	Level           int      `json:"level"`
	Number          int      `json:"number"`
	Text            string   `json:"text"`
	Options         []string `json:"options"`
	CorrectAnswer   []string `json:"correctAnswer"`
	Explanation     string   `json:"explanation"`
	AttemptsOverall int      `json:"attemptsOverall"`
	Weight          int      `json:"weight"`
}

// QuestionInstance : Meta data on current question
type QuestionInstance struct {
	Status      string   `json:"status"`
	Answer      []string `json:"answer"`
	NumAttempts int      `json:"numAttempts"`
}

// User : Current user doing the quiz
type User struct {
	Username string `json:"username"`
}

// Read Question from static JSON file
func getQuestion() QuestionData {
	jf, err := os.Open("static/json/simplequestion.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jf.Close()
	bv, _ := ioutil.ReadAll(jf)
	var cq QuestionData
	json.Unmarshal(bv, &cq)
	return cq
}

func getQuestionFromPythonScript(q QuestionData, s string) QuestionData {
	if s == "" {
		cmd := exec.Command(Python, PathToPythonScript)
		outb, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(outb))
		err = json.Unmarshal(outb, &q)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		cmd := exec.Command(Python, PathToPythonScript, s)
		outb, err := cmd.CombinedOutput()
		fmt.Println(string(outb))
		err = json.Unmarshal(outb, &q)
		if err != nil {
			fmt.Println(err)
		}
	}
	return q
}

func getNextQuizState(q QuestionData) QuestionData {
	if q.QuestionInstance.Answer == nil {
		return getQuestionFromPythonScript(q, "")
	}
	j, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}
	return getQuestionFromPythonScript(q, string(j))
}

func logQuestionData(q QuestionData) {
	fmt.Println("Question - Assignment: " + q.Question.Assignment)
	fmt.Println("Question - Level: " + strconv.Itoa(q.Question.Level))
	fmt.Println("Question - Number: " + strconv.Itoa(q.Question.Number))
	fmt.Println("Question - Text: " + q.Question.Text)
	fmt.Println("Question - Options: " + strings.Join(q.Question.Options, " "))
	fmt.Println("Question - CorrectAnswer: " + strings.Join(q.Question.CorrectAnswer, " "))
	fmt.Println("Question - Explanations: " + q.Question.Explanation)
	fmt.Println("Question - AttemptsOverall: " + strconv.Itoa(q.Question.AttemptsOverall))
	fmt.Println("Question - Weight: " + strconv.Itoa(q.Question.Weight))
	fmt.Println("Question Instance - Status: " + q.QuestionInstance.Status)
	fmt.Println("Question Instance - Answer: " + strings.Join(q.QuestionInstance.Answer, " "))
	fmt.Println("Question Instance - Num attempts: " + strconv.Itoa(q.QuestionInstance.NumAttempts))
	fmt.Println("User - Username: " + q.User.Username)
}
