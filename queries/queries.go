//////////////Tables///////////
//all 3 tables are in db under test02
//1.responses : this will save every response the user submits
/// Primary Key: username, assignment, level, numb, attempt
/// responses is used for logging and saving data for offline evaulation
//2.scores : this will save aggregated data on each user allowing us to save scores and track user location
/// Primary Key: username, assignment
/// scores is used for saving scores and location and the row of user,assignment is what is sent to the python script
//3.weights : mapping of assignment and the maximum score that the user can receive
/// Primary key: assignment
/// used at the end for calculating the user's score that is returned to edX

////////////////////QUERIES///////////////////

/////////////////BEGINNING OF ASSIGNMENT////////////////

//check if user, assignment exists - user_id and custom_component_display_name is what is received from LTI
rows, err := db.Query(`SELECT username, assignment FROM test02.scores
  WHERE username = "%s" AND assignment = "%s";`, user_id, custom_component_display_name)
if err != nil {
	log.Fatal(err)
}
defer rows.Close()
i := 0
for rows.Next(){
   i++
}
if i == 0 {
  //if rows0 is empty - this should not occur but if it does, it means we dont have past data on the user
  //so insert it with these default values
  r0 := fmt.Sprintf(`insert into test02.scores
    (username, assignment, gender, level_of_education, enrollment_mode, ageCategory, ad1, ad2, ad3, ad4,
    sd1, sd2, sd3, sd4, de1, de2, de3, de4, cc1, cc2, cc3, cc4, rts1, rts2, rts3, rts4,
    score1_correct, score1_attempts, score2_correct, score2_attempts, score3_correct, score3_attempts
    score4_correct, score4_attempts, next1, next2, next3, next4)
    values ("%s", "%s", "%s", "%s", "%s", "%s", "%f", "%f", "%f", "%f",
       "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f",
       "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d");`,
     user_id, custom_component_display_name, "None", "None", "audit", "None", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
     0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1)
}

///////////////////////AFTER USER PRESSES SUBMIT AND BEFORE CALLING PYTHON SCRIPT//////////

//run when user submits answer if q.QuestionInstance.answer is not empty
//save response to responses
t := time.Now()
tf := t.Format("20060102150405")
r := fmt.Sprintf(`insert into test02.responses
  (username, assignment, level, numb, attempt, correctness, score_possible, answer, answer_timestamp)
  values ("%s", "%s", "%d", "%d", "%d", "%d", "%d", "%s", "%s");`,
  u.UserName, q.Question.Assignment, q.Question.level, q.Question.number,
  q.QuestionInstance.NumAttempts + 1, 0, q.Question.weight, q.QuestionInstance.answer, tf)

//run to get the user's row - this is the user's history that will be sent to python script
rows, err := db.Query(`SELECT * FROM test02.scores WHERE username = "%s" AND assignment = "%s";`,
  user_id, custom_component_display_name)
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

//////////////AFTER RECEIVING RESPONSE FROM PYTHON SCRIPT////////////

//run IF q.QuestionInstance.status == "Correct" was returned from script where
//q is the json returned from the script
//if correct, this saves last response as the correct one
r := fmt.Sprintf(`update test02.responses SET correctness = 1
  WHERE username="%s" AND assignment="%s" AND level="%d" AND numb="%d" AND attempt="%d";`,
   u.UserName, q.Question.Assignment, q.Question.level, q.Question.number, q.QuestionInstance.NumAttempts)

//run IF Question.level == 1 AND q.QuestionInstance.status == "Correct"
//update scores table when user is done with the question
r := fmt.Sprintf(`update test02.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + "%d",
  score1_correct = score1_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
  q.QuestionInstance.numAttemps, user_id, custom_component_display_name)

//run IF Question.level == 1 AND q.QuestionInstance.status == "IncorrectNoAttempts"
//update scores table when user is done with the question
r := fmt.Sprintf(`update test02.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + "%d"
  WHERE username = "%s" AND assignment = "%s";`, q.QuestionInstance.numAttemps, user_id, custom_component_display_name)

//run IF Question.level == 2 AND q.QuestionInstance.status == "Correct"
//update scores table when user is done with the question
r := fmt.Sprintf(`update test02.scores SET next2 = next2 + 1, score2_attempts = score2_attempts + "%d",
  score2_correct = score2_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
  q.QuestionInstance.numAttemps, user_id, custom_component_display_name)

//run IF Question.level == 2 AND q.QuestionInstance.status == "IncorrectNoAttempts"
//update scores table when user is done with the question
r := fmt.Sprintf(`update test02.scores SET next2 = next2 + 1, score2_attempts = score2_attempts + "%d"
  WHERE username = "%s" AND assignment = "%s";`, q.QuestionInstance.numAttemps, user_id, custom_component_display_name)

//run IF Question.level == 3 AND q.QuestionInstance.status == "Correct"
//update scores table when user is done with the question
r := fmt.Sprintf(`update test02.scores SET next3 = next3 + 1, score3_attempts = score3_attempts + "%d",
  score3_correct = score3_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
  q.QuestionInstance.numAttemps, user_id, custom_component_display_name)

//run IF Question.level == 3 AND q.QuestionInstance.status == "IncorrectNoAttempts"
//update scores table when user is done with the question
r := fmt.Sprintf(`update test02.scores SET next3 = next3 + 1, score3_attempts = score3_attempts + "%d"
  WHERE username = "%s" AND assignment = "%s";`, q.QuestionInstance.numAttemps, user_id, custom_component_display_name)

//run IF Question.level == 4 AND q.QuestionInstance.status == "Correct"
//update scores table when user is done with the question
r := fmt.Sprintf(`update test02.scores SET next2 = next4 + 1, score4_attempts = score4_attempts + "%d",
  score4_correct = score4_correct + 1 WHERE username = "%s" AND assignment = "%s";`,
  q.QuestionInstance.numAttemps, user_id, custom_component_display_name)

//run IF Question.level == 4 AND q.QuestionInstance.status == "IncorrectNoAttempts"
//update scores table when user is done with the question
r := fmt.Sprintf(`update test02.scores SET next4 = next4 + 1, score4_attempts = score4_attempts + "%d"
  WHERE username = "%s" AND assignment = "%s";`, q.QuestionInstance.numAttemps,  user_id, custom_component_display_name)


/////////////////END OF ASSIGNMENT//////////////////
// end of assignment happens when the user closes the tab or
// when the python script returns a specific Json (we have to corrdinate which)

//run when assignment ends IF q.Question.Assignment == "ClimateChange"
//save user's score in each level in all three of the user's rows (one row per assignment)
r := fmt.Sprintf(`update test02.scores SET cc1 = CASE
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
WHERE username = "%s" AND assignment = "%s";`, user_id, "Climate+Change")

rows, err = db.Query(`SELECT cc1, cc2, cc3, cc4
    FROM test02.scores WHERE username = "%s" AND assignment = "%s";`, user_id, "Climate+Change")
if err != nil {
	log.Fatal(err)
}
defer rows.Close()
// the primary key of the table is username, assignment, so there will only be one row in row11
for rows.Next() {
	err := rows11.Scan(&cc1, &cc2, &cc3, &cc4)
	if err != nil {
		log.Fatal(err)
	}
  r := fmt.Sprintf(`update test02.scores SET cc1 = "%f", cc2 = "%f", cc3 = "%f", cc4 = "%f"
    WHERE username = "%s";`, cc1, cc2, cc3, cc4, user_id)
}
err = rows.Err()
if err != nil {
	log.Fatal(err)
}

//run when assignment ends IF q.Question.Assignment == "Reading Test Scores"
//save user's score in each level in all three of the user's rows (one row per assignment)
r := fmt.Sprintf(`update test02.scores SET rts1 = CASE
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
WHERE username = "%s" AND assignment = "%s";`, user_id, "Reading+Test+Scores")

row, err = db.Query(`SELECT rts1, rts2, rts3, rts4
    FROM test02.scores WHERE username = "%s" AND assignment = "%s";`, user_id, "Reading+Test+Scores")
if err != nil {
	log.Fatal(err)
}
defer rows.Close()
// the primary key of the table is username, assignment, so there will only be one row in rows12
for rows.Next() {
	err := rows.Scan(&rts1, &rts2, &rts3, &rts4)
	if err != nil {
		log.Fatal(err)
	}
  r := fmt.Sprintf(`update test02.scores SET rts1 = "%f", rts2 = "%f", rts3 = "%f", rts4 = "%f"
    WHERE username = "%s";`, rts1, rts2, rts3, rts4, user_id)
}
err = rows.Err()
if err != nil {
	log.Fatal(err)
}

//don't need to do this previous query for the third assignment since it is the last

//run when assignment ends to return the user's score - this is the user's score for the assignment
rows, err := db.Query(`SELECT username, assignment, level, numb,
    MAX(correctness) AS correctness, FIRST_VALUE(score_possible) AS score_possible
    FROM test02.responses WHERE username = "%s" AND assignment = "%s"
    GROUP BY username, assignment, level, numb;`, u.UserName, q.Question.Assignment)
if err != nil {
    	log.Fatal(err)
    }
defer rows.Close()

rows, err := db.Query(`SELECT assignment,
    SUM(correctness*score_possible)/SUM(score_possible) AS prescore
    FROM rows GROUP BY username, assignment;`)
if err != nil {
    	log.Fatal(err)
    }
defer rows.Close()

rows, err := db.Query(`SELECT prescore * weight AS score
    FROM rows INNER JOIN test02.weights ON rows2.assignment = test02.weights.assignment;`)
if err != nil {
  	log.Fatal(err)
  }
defer rows.Close()

// the previous query will return one row will one column score
// this will be the value to return if the user finished the assignment (the python script returned the finishing value)
// IF the user leaves the page without finishing the assignment, we will penalize his score by running the following and returning that instead of the previous value
rows, err := db.Query(`SELECT 0.5 * score FROM rows;`)
if err != nil {
  	log.Fatal(err)
  }
defer rows.Close()
