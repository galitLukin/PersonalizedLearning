package main

import (
	"database/sql"
	"fmt"
	"time"
)

type scores struct {
	Username           string
	Assignment         string
	Gender             int
	Level_of_education int
	Enrollment_mode    int
	AgeCategory        int
	Next1              int
	Next2              int
	Next3              int
	Next4              int
	Grade              float32
}

type response struct {
	IsFirst     bool
	Level       int
	Number      int
	Attempt     int
	Correctness int
	Timestamp   string
}

type pastQ struct {
	Question    Question
	Correctness []bool
	Answer      []string
	AnswerTime  []string
	Attempts    []int
}

func dbCheck(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

/////////////////BEGINNING OF ASSIGNMENT////////////////

//get user data at the start of the assignment
func dbInitFetchUser(db *sql.DB, user string, assignment string) scores {
	var s scores
	fmt.Println("Getting user from scores  ...", user, assignment)
	q := fmt.Sprintf(`SELECT * FROM test02.scores
	WHERE username = "%s" AND assignment = "%s";`, user, assignment)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()
	i := 0
	for rows.Next() {
		i++
		err = rows.Scan(&s.Username, &s.Assignment, &s.Gender, &s.Level_of_education, &s.Enrollment_mode, &s.AgeCategory, &s.Next1, &s.Next2, &s.Next3, &s.Next4, &s.Grade)
		dbCheck(err)
	}
	if i == 0 {
		fmt.Println("User does not exist ...", user, assignment)
		//if rows is empty - this should not occur but if it does, it means we dont have past data on the user
		//so insert it with these default values
		q := fmt.Sprintf(`insert into test02.scores
			(username, assignment, gender, level_of_education, enrollment_mode, ageCategory,
 			next1, next2, next3, next4, grade)
			values ("%s", "%s", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%f");`,
			user, assignment, 0, 0, 0, 0, 1, 1, 1, 1, 0.0)

		stmt, err := db.Prepare(q)
		dbCheck(err)
		defer stmt.Close()

		_, err = stmt.Exec()
		dbCheck(err)

		q = fmt.Sprintf(`SELECT * FROM test02.scores
				WHERE username = "%s" AND assignment = "%s";`, user, assignment)
		rows, err := db.Query(q)
		dbCheck(err)
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&s.Username, &s.Assignment, &s.Gender, &s.Level_of_education, &s.Enrollment_mode, &s.AgeCategory, &s.Next1, &s.Next2, &s.Next3, &s.Next4, &s.Grade)
			dbCheck(err)
		}
	}
	return s
}

// get data on user's last location - if there was a previous session
func dbGetUserPrevLocation(db *sql.DB, qd QuestionData) response {
	fmt.Println("Getting users past location... ", qd.User.Uid, qd.User.AssignmentName)
	var r response
	r.IsFirst = true
	r.Level = 0
	q := fmt.Sprintf(`SELECT level, numb, attempt, correctness FROM test02.responses
   WHERE username = "%s" AND assignment = "%s" ORDER BY answer_timestamp DESC
   LIMIT 1;`, qd.User.Uid, qd.User.AssignmentName)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&r.Level, &r.Number, &r.Attempt, &r.Correctness)
		dbCheck(err)
	}
	return r
}

///////////////////////AFTER USER PRESSES SUBMIT AND BEFORE CALLING PYTHON SCRIPT//////////

//records response
//after user presses submit and before calling python script
//records only if user marked down an answer
func dbInsertResponse(db *sql.DB, qd QuestionData) {
	if len(qd.QuestionInstance.Answer[0]) > 0 && (qd.QuestionInstance.Status == "NewQuestion" || qd.QuestionInstance.Status == "IncorrectWithAttempts") {
		fmt.Println("Inserting user response  ...", qd.User.Uid, qd.User.AssignmentName)
		t := time.Now().UTC()
		tf := t.Format("20060102150405")
		q := fmt.Sprintf(`insert into test02.responses
		  (username, assignment, level, numb, attempt, correctness, score_possible, answer, answer_timestamp)
		  values ("%s", "%s", "%d", "%d", "%d", "%d", "%d", "%s", "%s");`,
			qd.User.Uid, qd.User.AssignmentName, qd.Question.Level, qd.Question.Number,
			qd.QuestionInstance.NumAttempts+1, 0, qd.Question.Weight, qd.QuestionInstance.Answer[0], tf)
		stmt, err := db.Prepare(q)
		dbCheck(err)
		defer stmt.Close()

		_, err = stmt.Exec()
		dbCheck(err)
	}
}

//run to get the user's history
//after user presses submit and before calling python script
//returns history only when user is moving to next question - ow null
func dbFetchUserInScores(db *sql.DB, qd QuestionData) scores {
	fmt.Println("Getting user from scores  ...", qd.User.Uid, qd.User.AssignmentName)
	var s scores
	q := fmt.Sprintf(`SELECT * FROM test02.scores WHERE username = "%s" AND assignment = "%s";`,
		qd.User.Uid, qd.User.AssignmentName)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&s.Username, &s.Assignment, &s.Gender, &s.Level_of_education, &s.Enrollment_mode, &s.AgeCategory, &s.Next1, &s.Next2, &s.Next3, &s.Next4, &s.Grade)
		dbCheck(err)
	}
	return s
}

//////////////AFTER RECEIVING RESPONSE FROM PYTHON SCRIPT////////////

//if correct, this saves last response as the correct one - or does nothing
func dbUpdateResponse(db *sql.DB, qd QuestionData) {
	fmt.Println("Updating user response  ...", qd.User.Uid, qd.User.AssignmentName)
	if qd.QuestionInstance.Status == "Correct" {
		q := fmt.Sprintf(`update test02.responses SET correctness = 1
		  WHERE username="%s" AND assignment="%s" AND level="%d" AND numb="%d" AND attempt="%d";`,
			qd.User.Uid, qd.User.AssignmentName, qd.Question.Level, qd.Question.Number, qd.QuestionInstance.NumAttempts)
		stmt, err := db.Prepare(q)
		dbCheck(err)
		defer stmt.Close()

		_, err = stmt.Exec()
		dbCheck(err)
	}

}

//update scores table when user is done with the question
func dbUpdateScores(db *sql.DB, qd QuestionData) {
	var q string
	fmt.Println("Updating user scores  ...", qd.User.Uid, qd.User.AssignmentName)
	if qd.QuestionInstance.Status == "Correct" || qd.QuestionInstance.Status == "IncorrectNoAttempts" {
		if qd.Question.Level == 1 {
			q = fmt.Sprintf(`update test02.scores SET next1 = next1 + 1 WHERE username = "%s" AND assignment = "%s";`,
				qd.User.Uid, qd.User.AssignmentName)
		} else if qd.Question.Level == 2 {
			q = fmt.Sprintf(`update test02.scores SET next2 = next2 + 1 WHERE username = "%s" AND assignment = "%s";`,
				qd.User.Uid, qd.User.AssignmentName)
		} else if qd.Question.Level == 3 {
			q = fmt.Sprintf(`update test02.scores SET next3 = next3 + 1 WHERE username = "%s" AND assignment = "%s";`,
				qd.User.Uid, qd.User.AssignmentName)
		} else {
			q = fmt.Sprintf(`update test02.scores SET next4 = next4 + 1 WHERE username = "%s" AND assignment = "%s";`,
				qd.User.Uid, qd.User.AssignmentName)
		}
		stmt, err := db.Prepare(q)
		dbCheck(err)
		defer stmt.Close()

		_, err = stmt.Exec()
		dbCheck(err)
	}
}

// run when user after question is received from python script
// updates only if the question is finished (complete/incorrectWithNoAttempts)
func dbUpdateFinishedQuestion(db *sql.DB, qd QuestionData) {
	dbUpdateResponse(db, qd)
	dbUpdateScores(db, qd)
}

/////////////////END OF ASSIGNMENT//////////////////

// // end of assignment happens when the user closes the tab or finishes the assignment
// func dbCalculateScores(db *sql.DB, qd QuestionData) {
// 	//save user's score in each level in all three of the user's rows (one row per assignment)
// 	var q string
// 	var ss scores
// 	if qd.User.AssignmentName == "Asmt1" {
// 		q = fmt.Sprintf(`update test02.scores SET cc1 = CASE
// 			WHEN score1_attempts > 0 THEN score1_correct/score1_attempts
// 			   ELSE 0
// 			END,
// 			cc2 = CASE
// 			   WHEN score2_attempts > 0 THEN score2_correct/score2_attempts
// 			   ELSE 0
// 			END,
// 			cc3 = CASE
// 			   WHEN score3_attempts > 0 THEN score3_correct/score3_attempts
// 			   ELSE 0
// 			END,
// 			cc4 = CASE
// 			   WHEN score4_attempts > 0 THEN score4_correct/score4_attempts
// 			   ELSE 0
// 			END
// 			WHERE username = "%s" AND assignment = "%s";`, qd.User.Uid, qd.User.AssignmentName)
//
// 		stmt, err := db.Prepare(q)
// 		dbCheck(err)
// 		defer stmt.Close()
//
// 		_, err = stmt.Exec()
// 		dbCheck(err)
//
// 		q = fmt.Sprintf(`SELECT cc1, cc2, cc3, cc4
// 			    FROM test02.scores WHERE username = "%s" AND assignment = "%s";`, qd.User.Uid, qd.User.AssignmentName)
//
// 		rows, err := db.Query(q)
// 		dbCheck(err)
// 		defer rows.Close()
//
// 		// the primary key of the table is username, assignment, so there will only be one row
// 		for rows.Next() {
// 			err := rows.Scan(&ss.Cc1, &ss.Cc2, &ss.Cc3, &ss.Cc4)
// 			dbCheck(err)
// 			q := fmt.Sprintf(`update test02.scores SET cc1 = "%f", cc2 = "%f", cc3 = "%f", cc4 = "%f"
// 			    WHERE username = "%s";`, ss.Cc1, ss.Cc2, ss.Cc3, ss.Cc4, qd.User.Uid)
//
// 			stmt, err := db.Prepare(q)
// 			dbCheck(err)
// 			defer stmt.Close()
//
// 			_, err = stmt.Exec()
// 			dbCheck(err)
//
// 		}
// 	} else if qd.User.AssignmentName == "Asmt2" {
// 		q = fmt.Sprintf(`update test02.scores SET rts1 = CASE
// 				WHEN score1_attempts > 0 THEN score1_correct/score1_attempts
// 				   ELSE 0
// 				END,
// 				rts2 = CASE
// 				   WHEN score2_attempts > 0 THEN score2_correct/score2_attempts
// 				   ELSE 0
// 				END,
// 				rts3 = CASE
// 				   WHEN score3_attempts > 0 THEN score3_correct/score3_attempts
// 				   ELSE 0
// 				END,
// 				rts4 = CASE
// 				   WHEN score4_attempts > 0 THEN score4_correct/score4_attempts
// 				   ELSE 0
// 				END
// 				WHERE username = "%s" AND assignment = "%s";`, qd.User.Uid, qd.User.AssignmentName)
//
// 		stmt, err := db.Prepare(q)
// 		dbCheck(err)
// 		defer stmt.Close()
//
// 		_, err = stmt.Exec()
// 		dbCheck(err)
//
// 		q = fmt.Sprintf(`SELECT rts1, rts2, rts3, rts4
// 				    FROM test02.scores WHERE username = "%s" AND assignment = "%s";`, qd.User.Uid, qd.User.AssignmentName)
//
// 		rows, err := db.Query(q)
// 		dbCheck(err)
// 		defer rows.Close()
//
// 		// the primary key of the table is username, assignment, so there will only be one row
// 		for rows.Next() {
// 			err := rows.Scan(&ss.Rts1, &ss.Rts2, &ss.Rts3, &ss.Rts4)
// 			dbCheck(err)
// 			q := fmt.Sprintf(`update test02.scores SET rts1 = "%f", rts2 = "%f", rts3 = "%f", rts4 = "%f"
// 				    WHERE username = "%s";`, ss.Rts1, ss.Rts2, ss.Rts3, ss.Rts4, qd.User.Uid)
//
// 			stmt, err := db.Prepare(q)
// 			dbCheck(err)
// 			defer stmt.Close()
//
// 			_, err = stmt.Exec()
// 			dbCheck(err)
// 		}
// 	}
//}

//runs when user finishes assignment - returns grade for edX
func dbCalculateGrade(db *sql.DB, qd QuestionData) float32 {
	fmt.Println("Getting user from responses  ...", qd.User.Uid, qd.User.AssignmentName)
	var g struct {
		Username      string
		Assignment    string
		Level         int
		Number        int
		Correctness   int
		ScorePossible int
		Grade         float32
	}
	score := 0
	scorePossible := 0
	q := fmt.Sprintf(`SELECT username, assignment, level, numb,
	    MAX(correctness) AS correctness, ANY_VALUE(score_possible) AS score_possible
	    FROM test02.responses WHERE username = "%s" AND assignment = "%s"
	    GROUP BY username, assignment, level, numb;`, qd.User.Uid, qd.User.AssignmentName)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&g.Username, &g.Assignment, &g.Level, &g.Number, &g.Correctness, &g.ScorePossible)
		dbCheck(err)
		score += g.Correctness * g.ScorePossible
		scorePossible += g.ScorePossible
	}

	if scorePossible > 0 {
		g.Grade = float32(score) / float32(scorePossible)
	} else {
		g.Grade = 0.0
	}
	g.Grade = float32(int(g.Grade*100)) / 100

	q = fmt.Sprintf(`update test02.scores SET grade = "%f"
				WHERE username = "%s" AND assignment = "%s";`, g.Grade, qd.User.Uid, qd.User.AssignmentName)

	stmt, err := db.Prepare(q)
	dbCheck(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	dbCheck(err)

	return g.Grade
}

// run if user finished assignment
func dbAssignmentDone(db *sql.DB, qd QuestionData) float32 {
	//dbCalculateScores(db, qd)
	return dbCalculateGrade(db, qd)
}

//run to get the user's history for past questions
//this is called when user wants to see their past questions
func dbFetchUserInResponses(db *sql.DB, qd QuestionData) []pastQ {
	fmt.Println("Getting user from responses  ...", qd.User.Uid, qd.User.AssignmentName)
	q := fmt.Sprintf(`SELECT level, numb, correctness, answer, attempt, answer_timestamp
			FROM test02.responses
			WHERE username = "%s" AND assignment = "%s"
			ORDER BY answer_timestamp ASC;`, qd.User.Uid, qd.User.AssignmentName)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()
	var userResponses []pastQ
	var pq pastQ
	var ques Question
	var level int
	var number int
	var correctness bool
	var answer string
	var attempt int
	var answertime string
	ques.Assignment = qd.User.AssignmentName
	for rows.Next() {
		err = rows.Scan(&ques.Level, &ques.Number, &correctness, &answer, &attempt, &answertime)
		dbCheck(err)
		if ques.Level == level && ques.Number == number {
			pq.Correctness = append(pq.Correctness, correctness)
			pq.Answer = append(pq.Answer, answer)
			pq.AnswerTime = append(pq.AnswerTime, answertime)
			pq.Attempts = append(pq.Attempts, attempt)
		} else {
			level = ques.Level
			number = ques.Number
			pq.Correctness = []bool {correctness}
			pq.Answer = []string {answer}
			pq.AnswerTime = []string {answertime}
			pq.Attempts =[]int {attempt}
		}
		pq.Question = ques
		dbCheck(err)
		userResponses = append(userResponses, pq)
	}
	return userResponses
}
