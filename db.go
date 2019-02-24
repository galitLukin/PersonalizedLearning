package main

import (
	"database/sql"
	"fmt"
	"time"
)

type userSchema struct {
	FName    string
	LName    string
	Email    string
	Userid   string
	Password string
}

type scores struct {
	Username           string
	Assignment         string
	Gender             string
	Level_of_education string
	Enrollment_mode    string
	AgeCategory        string
	Ad1                float32
	Ad2                float32
	Ad3                float32
	Ad4                float32
	Sd1                float32
	Sd2                float32
	Sd3                float32
	Sd4                float32
	De1                float32
	De2                float32
	De3                float32
	De4                float32
	Cc1                float32
	Cc2                float32
	Cc3                float32
	Cc4                float32
	Rts1               float32
	Rts2               float32
	Rts3               float32
	Rts4               float32
	Score1_correct     int
	Score1_attempts    int
	Score2_correct     int
	Score2_attempts    int
	Score3_correct     int
	Score3_attempts    int
	Score4_correct     int
	Score4_attempts    int
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

type grade struct {
	Username      string
	Assignment    string
	Level         int
	Number        int
	Correctness   int
	ScorePossible int
	Grade         float32
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
		err = rows.Scan(&s.Username, &s.Assignment, &s.Gender, &s.Level_of_education, &s.Enrollment_mode, &s.AgeCategory, &s.Ad1, &s.Ad2, &s.Ad3, &s.Ad4, &s.Sd1, &s.Sd2, &s.Sd3, &s.Sd4, &s.De1, &s.De2, &s.De3, &s.De4, &s.Cc1, &s.Cc2, &s.Cc3, &s.Cc4, &s.Rts1, &s.Rts2, &s.Rts3, &s.Rts4, &s.Score1_correct, &s.Score1_attempts, &s.Score2_correct, &s.Score2_attempts, &s.Score3_correct, &s.Score3_attempts, &s.Score4_correct, &s.Score4_attempts, &s.Next1, &s.Next2, &s.Next3, &s.Next4, &s.Grade)
		dbCheck(err)
	}
	if i == 0 {
		fmt.Println("User does not exist ...", user, assignment)
		//if rows is empty - this should not occur but if it does, it means we dont have past data on the user
		//so insert it with these default values
		q := fmt.Sprintf(`insert into test02.scores
			(username, assignment, gender, level_of_education, enrollment_mode, ageCategory, ad1, ad2, ad3, ad4,
			sd1, sd2, sd3, sd4, de1, de2, de3, de4, cc1, cc2, cc3, cc4, rts1, rts2, rts3, rts4,
			score1_correct, score1_attempts, score2_correct, score2_attempts, score3_correct, score3_attempts,
			score4_correct, score4_attempts, next1, next2, next3, next4, grade)
			values ("%s", "%s", "%s", "%s", "%s", "%s", "%f", "%f", "%f", "%f",
			   "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f",
			   "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%f");`,
			user, assignment, "None", "None", "audit", "Null", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0.0)

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
			err = rows.Scan(&s.Username, &s.Assignment, &s.Gender, &s.Level_of_education, &s.Enrollment_mode, &s.AgeCategory, &s.Ad1, &s.Ad2, &s.Ad3, &s.Ad4, &s.Sd1, &s.Sd2, &s.Sd3, &s.Sd4, &s.De1, &s.De2, &s.De3, &s.De4, &s.Cc1, &s.Cc2, &s.Cc3, &s.Cc4, &s.Rts1, &s.Rts2, &s.Rts3, &s.Rts4, &s.Score1_correct, &s.Score1_attempts, &s.Score2_correct, &s.Score2_attempts, &s.Score3_correct, &s.Score3_attempts, &s.Score4_correct, &s.Score4_attempts, &s.Next1, &s.Next2, &s.Next3, &s.Next4, &s.Grade)
			dbCheck(err)
		}
	}
	return s
}

// get data on user's last location - if there was a previous session
func dbGetUserPrevLocation(db *sql.DB, qd QuestionData) response {
	fmt.Println("Getting users past location... ", qd.User.Username, qd.Question.Assignment)
	var r response
	r.IsFirst = true
	r.Level = 0
	q := fmt.Sprintf(`SELECT level, numb, attempt, correctness FROM test02.responses
   WHERE username = "%s" AND assignment = "%s" ORDER BY answer_timestamp DESC
   LIMIT 1;`, qd.User.Username, qd.Question.Assignment)
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
	if qd.QuestionInstance.Answer != nil && (qd.QuestionInstance.Status == "NewQuestion" || qd.QuestionInstance.Status == "IncorrectWithAttempts") {
		fmt.Println("Inserting user response  ...", qd.User.Username, qd.Question.Assignment)
		t := time.Now()
		tf := t.Format("20060102150405")
		q := fmt.Sprintf(`insert into test02.responses
		  (username, assignment, level, numb, attempt, correctness, score_possible, answer, answer_timestamp)
		  values ("%s", "%s", "%d", "%d", "%d", "%d", "%d", "%s", "%s");`,
			qd.User.Username, qd.Question.Assignment, qd.Question.Level, qd.Question.Number,
			qd.QuestionInstance.NumAttempts+1, 0, qd.Question.Weight, qd.QuestionInstance.Answer, tf)
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
	fmt.Println("Getting user from scores  ...", qd.User.Username, qd.Question.Assignment)
	var s scores
	q := fmt.Sprintf(`SELECT * FROM test02.scores WHERE username = "%s" AND assignment = "%s";`,
		qd.User.Username, qd.Question.Assignment)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&s.Username, &s.Assignment, &s.Gender, &s.Level_of_education, &s.Enrollment_mode, &s.AgeCategory, &s.Ad1, &s.Ad2, &s.Ad3, &s.Ad4, &s.Sd1, &s.Sd2, &s.Sd3, &s.Sd4, &s.De1, &s.De2, &s.De3, &s.De4, &s.Cc1, &s.Cc2, &s.Cc3, &s.Cc4, &s.Rts1, &s.Rts2, &s.Rts3, &s.Rts4, &s.Score1_correct, &s.Score1_attempts, &s.Score2_correct, &s.Score2_attempts, &s.Score3_correct, &s.Score3_attempts, &s.Score4_correct, &s.Score4_attempts, &s.Next1, &s.Next2, &s.Next3, &s.Next4, &s.Grade)
		dbCheck(err)
	}
	return s
}

//////////////AFTER RECEIVING RESPONSE FROM PYTHON SCRIPT////////////

//if correct, this saves last response as the correct one - ow does nothing
func dbUpdateResponse(db *sql.DB, qd QuestionData) {
	fmt.Println("Updating user response  ...", qd.User.Username, qd.Question.Assignment)
	if qd.QuestionInstance.Status == "Correct" {
		q := fmt.Sprintf(`update test02.responses SET correctness = 1
		  WHERE username="%s" AND assignment="%s" AND level="%d" AND numb="%d" AND attempt="%d";`,
			qd.User.Username, qd.Question.Assignment, qd.Question.Level, qd.Question.Number, qd.QuestionInstance.NumAttempts)
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
	fmt.Println("Updating user scores  ...", qd.User.Username, qd.Question.Assignment)
	if qd.QuestionInstance.Status == "Correct" {
		if qd.Question.Level == 1 {
			q = fmt.Sprintf(`update test02.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + "%d",
			  score1_correct = score1_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
				qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else if qd.Question.Level == 2 {
			q = fmt.Sprintf(`update test02.scores SET next2 = next2 + 1, score2_attempts = score2_attempts + "%d",
			  score2_correct = score2_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
				qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else if qd.Question.Level == 3 {
			q = fmt.Sprintf(`update test02.scores SET next3 = next3 + 1, score3_attempts = score3_attempts + "%d",
			  score3_correct = score3_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
				qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else {
			q = fmt.Sprintf(`update test02.scores SET next4 = next4 + 1, score4_attempts = score4_attempts + "%d",
			  score4_correct = score4_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
				qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		}
		stmt, err := db.Prepare(q)
		dbCheck(err)
		defer stmt.Close()

		_, err = stmt.Exec()
		dbCheck(err)

	} else if qd.QuestionInstance.Status == "IncorrectNoAttempts" {
		if qd.Question.Level == 1 {
			q = fmt.Sprintf(`update test02.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + "%d"
					WHERE username = "%s" AND assignment = "%s";`, qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else if qd.Question.Level == 2 {
			q = fmt.Sprintf(`update test02.scores SET next2 = next2 + 1, score2_attempts = score2_attempts + "%d"
					WHERE username = "%s" AND assignment = "%s";`, qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else if qd.Question.Level == 3 {
			q = fmt.Sprintf(`update test02.scores SET next3 = next3 + 1, score3_attempts = score3_attempts + "%d"
					WHERE username = "%s" AND assignment = "%s";`, qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else {
			q = fmt.Sprintf(`update test02.scores SET next4 = next4 + 1, score4_attempts = score4_attempts + "%d"
					WHERE username = "%s" AND assignment = "%s";`, qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
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

// end of assignment happens when the user closes the tab or finishes the assignment
func dbCalculateScores(db *sql.DB, qd QuestionData) {
	//save user's score in each level in all three of the user's rows (one row per assignment)
	var q string
	var ss scores
	if qd.Question.Assignment == "Climate Change" {
		q = fmt.Sprintf(`update test02.scores SET cc1 = CASE
			WHEN score1_attempts > 0 THEN score1_correct/score1_attempts
			   ELSE 0
			END,
			cc2 = CASE
			   WHEN score2_attempts > 0 THEN score2_correct/score2_attempts
			   ELSE 0
			END,
			cc3 = CASE
			   WHEN score3_attempts > 0 THEN score3_correct/score3_attempts
			   ELSE 0
			END,
			cc4 = CASE
			   WHEN score4_attempts > 0 THEN score4_correct/score4_attempts
			   ELSE 0
			END
			WHERE username = "%s" AND assignment = "%s";`, qd.User.Username, qd.Question.Assignment)

		stmt, err := db.Prepare(q)
		dbCheck(err)
		defer stmt.Close()

		_, err = stmt.Exec()
		dbCheck(err)

		q = fmt.Sprintf(`SELECT cc1, cc2, cc3, cc4
			    FROM test02.scores WHERE username = "%s" AND assignment = "%s";`, qd.User.Username, qd.Question.Assignment)

		rows, err := db.Query(q)
		dbCheck(err)
		defer rows.Close()

		// the primary key of the table is username, assignment, so there will only be one row
		for rows.Next() {
			err := rows.Scan(&ss.Cc1, &ss.Cc2, &ss.Cc3, &ss.Cc4)
			dbCheck(err)
			q := fmt.Sprintf(`update test02.scores SET cc1 = "%f", cc2 = "%f", cc3 = "%f", cc4 = "%f"
			    WHERE username = "%s";`, ss.Cc1, ss.Cc2, ss.Cc3, ss.Cc4, qd.User.Username)

			stmt, err := db.Prepare(q)
			dbCheck(err)
			defer stmt.Close()

			_, err = stmt.Exec()
			dbCheck(err)

		}
	} else if qd.Question.Assignment == "Reading Test Scores" {
		q = fmt.Sprintf(`update test02.scores SET rts1 = CASE
				WHEN score1_attempts > 0 THEN score1_correct/score1_attempts
				   ELSE 0
				END,
				rts2 = CASE
				   WHEN score2_attempts > 0 THEN score2_correct/score2_attempts
				   ELSE 0
				END,
				rts3 = CASE
				   WHEN score3_attempts > 0 THEN score3_correct/score3_attempts
				   ELSE 0
				END,
				rts4 = CASE
				   WHEN score4_attempts > 0 THEN score4_correct/score4_attempts
				   ELSE 0
				END
				WHERE username = "%s" AND assignment = "%s";`, qd.User.Username, qd.Question.Assignment)

		stmt, err := db.Prepare(q)
		dbCheck(err)
		defer stmt.Close()

		_, err = stmt.Exec()
		dbCheck(err)

		q = fmt.Sprintf(`SELECT rts1, rts2, rts3, rts4
				    FROM test02.scores WHERE username = "%s" AND assignment = "%s";`, qd.User.Username, qd.Question.Assignment)

		rows, err := db.Query(q)
		dbCheck(err)
		defer rows.Close()

		// the primary key of the table is username, assignment, so there will only be one row
		for rows.Next() {
			err := rows.Scan(&ss.Rts1, &ss.Rts2, &ss.Rts3, &ss.Rts4)
			dbCheck(err)
			q := fmt.Sprintf(`update test02.scores SET rts1 = "%f", rts2 = "%f", rts3 = "%f", rts4 = "%f"
				    WHERE username = "%s";`, ss.Rts1, ss.Rts2, ss.Rts3, ss.Rts4, qd.User.Username)

			stmt, err := db.Prepare(q)
			dbCheck(err)
			defer stmt.Close()

			_, err = stmt.Exec()
			dbCheck(err)
		}
	}
}

//runs when user finishes assignment - returns grade for edX
func dbCalculateGrade(db *sql.DB, qd QuestionData) float32 {
	var g grade
	score := 0
	scorePossible := 0
	q := fmt.Sprintf(`SELECT username, assignment, level, numb,
	    MAX(correctness) AS correctness, ANY_VALUE(score_possible) AS score_possible
	    FROM test02.responses WHERE username = "%s" AND assignment = "%s"
	    GROUP BY username, assignment, level, numb;`, qd.User.Username, qd.Question.Assignment)
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
	g.Grade = float32(int(g.Grade * 100)) / 100

	q = fmt.Sprintf(`update test02.scores SET grade = "%f"
				WHERE username = "%s" AND assignment = "%s";`, g.Grade, qd.User.Username, qd.Question.Assignment)

	stmt, err := db.Prepare(q)
	dbCheck(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	dbCheck(err)

	return g.Grade
}

// run if user finished assignment
func dbAssignmentDone(db *sql.DB, qd QuestionData) float32 {
	dbCalculateScores(db, qd)
	return dbCalculateGrade(db, qd)
}
