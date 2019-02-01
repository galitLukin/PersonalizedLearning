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
	QuestionNumber  int      `json:"question_number"`
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

func (q QuestionData) getNextQuizState() QuestionData {
	j, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(Python, PathToPythonScript, string(j))
	outb, err := cmd.CombinedOutput()
	var cq QuestionData
	fmt.Println(string(outb))
	fmt.Println("***********")
	err = json.Unmarshal(outb, &cq)
	cq.logQuestionData()
	if err != nil {
		fmt.Println(err)
	}
	return cq
}

func (q QuestionData) logQuestionData() {
	fmt.Println("Question - Assignment: " + q.Question.Assignment)
	fmt.Println("Question - Level: " + strconv.Itoa(q.Question.Level))
	fmt.Println("Question - QuestionNumber: " + strconv.Itoa(q.Question.QuestionNumber))
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
