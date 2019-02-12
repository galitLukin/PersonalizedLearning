package main

import (
	"database/sql"
	"fmt"
)

type userSchema struct {
	FName    string
	LName    string
	Email    string
	Userid   string
	Password string
}

type scoresSchema struct {
	Username	string
	Assignment	string
	Gender	string
	Level_of_education	string
	Enrollment_mode	string
	AgeCategory	string
	Ad1	float32
	Ad2	float32
	Ad3	float32
	Ad4	float32
	Sd1	float32
	Sd2	float32
	Sd3	float32
	Sd4	float32
	De1	float32
	De2	float32
	De3	float32
	De4	float32
	Cc1	float32
	Cc2	float32
	Cc3	float32
	Cc4	float32
	Rts1	float32
	Rts2	float32
	Rts3	float32
	Rts4	float32
	Score1_correct int
	Score1_attempts int
	Score2_correct int
	Score2_attempts int
	Score3_correct int
	Score3_attempts int
	Score4_correct int
	Score4_attempts int
	Next1 int
	Next2 int
	Next3 int
	Next4 int
}

func dbGetUsers(db *sql.DB) string {
	rows, err := db.Query(`SELECT * FROM users;`)
	dbCheck(err)
	defer rows.Close()

	var s string

	// query
	for rows.Next() {
		u := user{}
		err = rows.Scan(&u.UserName, &u.First, &u.Last, &u.UserName, &u.Password)
		dbCheck(err)
		s += fmt.Sprintf(`email: "%s" firstName: "%s", lastName: "%s", passWord: "%s"`, u.UserName, u.First, u.Last, u.Password)
		s += "\n"
	}
	return s
}

func dbGetUser(db *sql.DB, email string) user {
	q := fmt.Sprintf(`SELECT * FROM users WHERE email="%s";`, email)
	fmt.Println(q)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()

	// data to be used in query
	var u user
	for rows.Next() {
		err = rows.Scan(&u.UserName, &u.First, &u.Last, &u.UserName, &u.Password)
		dbCheck(err)
		s := fmt.Sprintf(`email: "%s" firstName: "%s", lastName: "%s", passWord: "%s"`, u.UserName, u.First, u.Last, u.Password)
		fmt.Printf(`RETRIEVED USER: %#v`, s)
	}
	return u
}

func dbCreateUser(db *sql.DB, newUser user) string {
	q := fmt.Sprintf(`insert into test02.users (fName, lName, email, password) values ("%s", "%s", "%s", "%s");`, newUser.First, newUser.Last, newUser.UserName, newUser.Password)
	stmt, err := db.Prepare(q)
	dbCheck(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	dbCheck(err)

	n, err := r.RowsAffected()
	dbCheck(err)

	return fmt.Sprintf("%s%d", "INSERTED RECORD ", n)
}

func dbDeleteUser(db *sql.DB) string {
	stmt, err := db.Prepare(`DELETE FROM test02.users WHERE fName="Omer";`)
	dbCheck(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	dbCheck(err)

	n, err := r.RowsAffected()
	dbCheck(err)

	return fmt.Sprintf("%s%d", "DELETED RECORD ", n)
}

func dbCheck(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//check if user, assignment exists - otherwise, create a row with default values
//beginning of assignment
func dbUserInScores(db *sql.DB, user string, assignment string) string {
	q := fmt.Sprintf(`SELECT username, assignment FROM test02.scores
	  WHERE username = "%s" AND assignment = "%s";`, user, assignment)
	fmt.Println(q)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()

	i := 0
	for rows.Next(){
	   i++
 		 return fmt.Printf(`RETRIEVED USER: %s`, user)
	}
	if i == 0 {
	  //if rows is empty - this should not occur but if it does, it means we dont have past data on the user
	  //so insert it with these default values
	  q := fmt.Sprintf(`insert into test02.scores
	    (username, assignment, gender, level_of_education, enrollment_mode, ageCategory, ad1, ad2, ad3, ad4,
	    sd1, sd2, sd3, sd4, de1, de2, de3, de4, cc1, cc2, cc3, cc4, rts1, rts2, rts3, rts4,
	    score1_correct, score1_attempts, score2_correct, score2_attempts, score3_correct, score3_attempts
	    score4_correct, score4_attempts, next1, next2, next3, next4)
	    values ("%s", "%s", "%s", "%s", "%s", "%s", "%f", "%f", "%f", "%f",
	       "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f",
	       "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d");`,
	     user, assignment, "None", "None", "audit", "Null", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
	     0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1)

		stmt, err := db.Prepare(q)
	 	dbCheck(err)
	 	defer stmt.Close()

	 	r, err := stmt.Exec()
	 	dbCheck(err)

	 	n, err := r.RowsAffected()
	 	dbCheck(err)

	 	return fmt.Sprintf("%s%d", "INSERTED RECORD ", n)
	}
}

//records response
//after user presses submit and before calling python script
func dbInsertResponse(db *sql.DB, qd QuestionData) string {
	//run when user submits answer if q.QuestionInstance.answer is not empty
	//save response to responses
	t := time.Now()
	tf := t.Format("20060102150405")
	q := fmt.Sprintf(`insert into test02.responses
	  (username, assignment, level, numb, attempt, correctness, score_possible, answer, answer_timestamp)
	  values ("%s", "%s", "%d", "%d", "%d", "%d", "%d", "%s", "%s");`,
	  qd.User.Username, qd.Question.Assignment, qd.Question.Level, qd.Question.Number,
	  qd.QuestionInstance.NumAttempts, 0, qd.Question.Weight, qd.QuestionInstance.Answer, tf)
	stmt, err := db.Prepare(q)
	dbCheck(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	dbCheck(err)

	n, err := r.RowsAffected()
	dbCheck(err)

	return fmt.Sprintf("%s%d", "INSERTED RECORD ", n)
}

//run to get the user's row - this returns one row and is the user's history that will be sent to python script
//after user presses submit and before calling python script
func dbUserInScores(db *sql.DB, user string, assignment string) string {
	q := fmt.Sprintf(`SELECT * FROM test02.scores WHERE username = "%s" AND assignment = "%s";`,
	  user, assignment)
	fmt.Println(q)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()

	var s scoresSchema
	for rows.Next() {
		err = rows.Scan(&s.Username, &s.Assignment, &s.Gender, &s.Level_of_education, &s.Enrollment_mode, &s.AgeCategory, &s.Ad1, &s.Ad2, &s.Ad3, &s.Ad4, &s.Sd1, &s.Sd2, &s.Sd3, &s.Sd4, &s.De1, &s.De2, &s.De3, &s.De4, &s.Cc1, &s.Cc2, &s.Cc3, &s.Cc4, &s.Rts1, &s.Rts2, &s.Rts3, &s.Rts4, &s.Score1_correct, &s.Score1_attempts, &s.Score2_correct, &s.Score2_attempts, &s.Score3_correct, &s.Score3_attempts, &s.Score4_correct, &s.Score4_attempts, &s.Next1, &s.Next2, &s.Next3, &s.Next4)
		dbCheck(err)
	}
	return s
}

//////////////AFTER RECEIVING RESPONSE FROM PYTHON SCRIPT////////////

//run IF q.QuestionInstance.status == "Correct" was returned from script where
//q is the json returned from the script
//if correct, this saves last response as the correct one
func dbUpdateResponse(db *sql.DB, qd QuestionData) string {
	q := fmt.Sprintf(`update test02.responses SET correctness = 1
	  WHERE username="%s" AND assignment="%s" AND level="%d" AND numb="%d" AND attempt="%d";`,
	   qd.User.Username, qd.Question.Assignment, qd.Question.Level, qd.Question.Number,qd.QuestionInstance.NumAttempts)
	stmt, err := db.Prepare(q)
	dbCheck(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	dbCheck(err)

	n, err := r.RowsAffected()
	dbCheck(err)

	return fmt.Sprintf("%s%d", "UPDATED RECORD ", n)
}

//run IF QuestionInstance.status == "IcorrectNoAttempts" or QuestionInstance.status == "Correct"
//update scores table when user is done with the question
func dbUpdateScores(db *sql.DB, qd QuestionData) string {
	if qd.QuestionInstance.Status == "Correct"{
		if qd.Question.Level == 1{
				q := fmt.Sprintf(`update test02.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + "%d",
			  score1_correct = score1_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
			  qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else if qd.Question.Level == 2{
				q := fmt.Sprintf(`update test02.scores SET next2 = next2 + 1, score2_attempts = score2_attempts + "%d",
			  score2_correct = score2_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
			  qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else if qd.Question.Level == 3{
				q := fmt.Sprintf(`update test02.scores SET next3 = next3 + 1, score3_attempts = score3_attempts + "%d",
			  score3_correct = score3_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
			  qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else{
				q := fmt.Sprintf(`update test02.scores SET next2 = next4 + 1, score4_attempts = score4_attempts + "%d",
			  score4_correct = score4_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
			  qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		}
	} else{
		if qd.Question.Level == 1{
				q := fmt.Sprintf(`update test02.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + "%d"
					WHERE username = "%s" AND assignment = "%s";`, qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else if qd.Question.Level == 2{
				q := fmt.Sprintf(`update test02.scores SET next2 = next2 + 1, score2_attempts = score2_attempts + "%d"
					WHERE username = "%s" AND assignment = "%s";`, qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else if qd.Question.Level == 3{
				q := fmt.Sprintf(`update test02.scores SET next3 = next3 + 1, score3_attempts = score3_attempts + "%d"
					WHERE username = "%s" AND assignment = "%s";`, qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		} else{
				q := fmt.Sprintf(`update test02.scores SET next4 = next4 + 1, score4_attempts = score4_attempts + "%d"
					WHERE username = "%s" AND assignment = "%s";`, qd.QuestionInstance.NumAttempts, qd.User.Username, qd.Question.Assignment)
		}
	}

	stmt, err := db.Prepare(q)
	dbCheck(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	dbCheck(err)

	n, err := r.RowsAffected()
	dbCheck(err)

	return fmt.Sprintf("%s%d", "UPDATED RECORD ", n)
}

/////////////////END OF ASSIGNMENT//////////////////
// end of assignment happens when the user closes the tab or

func dbCalculateAssignmentScores(db *sql.DB, qd QuestionData) string {
	//save user's score in each level in all three of the user's rows (one row per assignment)
	if qd.Question.Assignment == "Climate Change"{
			q := fmt.Sprintf(`update test02.scores SET cc1 = CASE
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

			r, err := stmt.Exec()
			dbCheck(err)

			n, err := r.RowsAffected()
			dbCheck(err)

			q = db.Sprintf(`SELECT cc1, cc2, cc3, cc4
			    FROM test02.scores WHERE username = "%s" AND assignment = "%s";`, qd.User.Username, qd.Question.Assignment)

			rows, err := db.Query(q)
			dbCheck(err)
			defer rows.Close()

			// the primary key of the table is username, assignment, so there will only be one row
			for rows.Next() {
				err := rows.Scan(&cc1, &cc2, &cc3, &cc4)
				dbCheck(err)
			  q := fmt.Sprintf(`update test02.scores SET cc1 = "%f", cc2 = "%f", cc3 = "%f", cc4 = "%f"
			    WHERE username = "%s";`, cc1, cc2, cc3, cc4, qd.User.Username)

				stmt, err := db.Prepare(q)
			 	dbCheck(err)
			 	defer stmt.Close()

			 	r, err := stmt.Exec()
			 	dbCheck(err)

			 	n, err := r.RowsAffected()
			 	dbCheck(err)
			}
			return fmt.Sprintf("%s%d", "UPDATED RECORD ", n)

		} else if q.Question.Assignment == "Reading Test Scores"{
				q := fmt.Sprintf(`update test02.scores SET rts1 = CASE
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

				r, err := stmt.Exec()
				dbCheck(err)

				n, err := r.RowsAffected()
				dbCheck(err)

				q = db.Sprintf(`SELECT rts1, rts2, rts3, rts4
				    FROM test02.scores WHERE username = "%s" AND assignment = "%s";`, qd.User.Username, qd.Question.Assignment)

				rows, err := db.Query(q)
				dbCheck(err)
				defer rows.Close()

				// the primary key of the table is username, assignment, so there will only be one row
				for rows.Next() {
					err := rows.Scan(&cc1, &cc2, &cc3, &cc4)
					dbCheck(err)
					q := fmt.Sprintf(`update test02.scores SET rts1 = "%f", rts2 = "%f", rts3 = "%f", rts4 = "%f"
				    WHERE username = "%s";`, rts1, rts2, rts3, rts4, qd.User.Username)

					stmt, err := db.Prepare(q)
					dbCheck(err)
					defer stmt.Close()

					r, err := stmt.Exec()
					dbCheck(err)

					n, err := r.RowsAffected()
					dbCheck(err)
				}
				return fmt.Sprintf("%s%d", "UPDATED RECORD ", n)
		}
}
