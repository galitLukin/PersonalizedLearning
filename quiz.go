package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	SelectedAnswers           = "selectedAnswers"
	Python                    = "python3"
	PathToPythonScript        = "./python/script.py"
	PathToPastQuestionsScript = "./python/pastQuestions.py"
	QuizReferrer              = "http://3.16.157.40/quiz"
)

// Question Data
type QuestionData struct {
	Question         Question
	QuestionInstance QuestionInstance
	User             user
	Score            scores
	PrevLocation     response
}

// Question : Data for current question
type Question struct {
	Assignment      string   `json:"Assignment"`
	Level           int      `json:"level"`
	Number          int      `json:"number"`
	Text            string   `json:"text"`
	Options         []string `json:"options"`
	CorrectAnswer   string   `json:"correctAnswer"`
	Explanation     string   `json:"explanation"`
	AttemptsOverall int      `json:"attemptsOverall"`
	Weight          int      `json:"weight"`
	AnswerType      int      `json:"answerType"`
	Qid							int		 	 `json:"qid"`
}

// QuestionInstance : Meta data on current question
type QuestionInstance struct {
	Status      string   `json:"status"`
	Answer      []string `json:"answer"`
	NumAttempts int      `json:"numAttempts"`
	StartTime 	time     `json:"startTime"`
	EndTime 		time     `json:"endTime"`
	Duration 		int      `json:"duration"`
}

// Read Question from static JSON file
// debugging purposes
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
	cmd := exec.Command(Python, PathToPythonScript, s)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(outb))
	err = json.Unmarshal(outb, &q)
	if err != nil {
		fmt.Println(err)
	}
	return q
}

func getNextQuizState(q QuestionData) QuestionData {
	j, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}
	return getQuestionFromPythonScript(q, string(j))
}

func getAllPastQuestions(qd QuestionData, pqs []pastQ) []pastQ {
	fmt.Println("Getting past questions  ...", qd.User.Uid, qd.User.AssignmentName)
	j, err := json.Marshal(pqs)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(Python, PathToPastQuestionsScript, string(j))
	outb, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	var qs []Question
	err = json.Unmarshal(outb, &qs)
	if err != nil {
		fmt.Println(err)
	}
	for i, q := range qs {
		pqs[i].Question = q
	}
	return pqs
}

func logQuestion(q Question) {
	fmt.Println("Question - Assignment: " + q.Assignment)
	fmt.Println("Question - Level: " + strconv.Itoa(q.Level))
	fmt.Println("Question - Number: " + strconv.Itoa(q.Number))
	fmt.Println("Question - Text: " + q.Text)
	fmt.Println("Question - Options: " + strings.Join(q.Options, " "))
	fmt.Println("Question - CorrectAnswer: " + q.CorrectAnswer)
	fmt.Println("Question - Explanations: " + q.Explanation)
	fmt.Println("Question - AttemptsOverall: " + strconv.Itoa(q.AttemptsOverall))
	fmt.Println("Question - Weight: " + strconv.Itoa(q.Weight))
	fmt.Println("Question - AnswerType: " + strconv.Itoa(q.AnswerType))
}

func logQuestionData(q QuestionData) {
	fmt.Println("Question - Assignment: " + q.User.AssignmentName)
	fmt.Println("Question - Level: " + strconv.Itoa(q.Question.Level))
	fmt.Println("Question - Number: " + strconv.Itoa(q.Question.Number))
	fmt.Println("Question - Text: " + q.Question.Text)
	fmt.Println("Question - Options: " + strings.Join(q.Question.Options, " "))
	fmt.Println("Question - CorrectAnswer: " + q.Question.CorrectAnswer)
	fmt.Println("Question - Explanations: " + q.Question.Explanation)
	fmt.Println("Question - AttemptsOverall: " + strconv.Itoa(q.Question.AttemptsOverall))
	fmt.Println("Question - Weight: " + strconv.Itoa(q.Question.Weight))
	fmt.Println("Question - AnswerType: " + strconv.Itoa(q.Question.AnswerType))
	fmt.Println("Question Instance - Status: " + q.QuestionInstance.Status)
	fmt.Println("Question Instance - Answer: " + strings.Join(q.QuestionInstance.Answer, " "))
	fmt.Println("Question Instance - Num attempts: " + strconv.Itoa(q.QuestionInstance.NumAttempts))
	fmt.Println("User - Username: " + q.User.Uid)
}
