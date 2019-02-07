u := user{}
q := QuestionData{}

/////////////////BEGINNING OF ASSIGNMENT////////////////
//check if user is in users list - user_id is what is received from LTI
u, err0 := db.Query(`SELECT username FROM test02.scores WHERE username = "%s";`, user_id)

//if u is empty
r0 := fmt.Sprintf(`insert into test02.scores
  (username, gender, level_of_education, enrollment_mode, ageCategory, ad1, ad2, ad3, ad4,
  sd1, sd2, sd3, sd4, de1, de2, de3, de4, cc1, cc2, cc3, cc4, rts1, rts2, rts3, rts4,
  score1_correct, score1_attempts, score2_correct, score2_attempts, score3_correct, score3_attempts
  score4_correct, score4_attempts, next1, next2, next3, next4)
  values ("%s", "%s", "%s", "%s", "%s", "%f", "%f", "%f", "%f",
     "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f", "%f",
     "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d", "%d");`,
   user_id, "None", "None", "audit", "None", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
   0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1)

/////////////////AFTER PRESSING SUBMIT AND BEFORE CALLING PYTHON SCRIPT//////////

//run when user submits answer if q.QuestionInstance.answer is not empty
t := time.Now()
tf := t.Format("20060102150405")
r1 := fmt.Sprintf(`insert into test02.responses
  (username, assignment, level, numb, attempt, correctness, score_possible, answer, answer_timestamp)
  values ("%s", "%s", "%d", "%d", "%d", "%d", "%d", "%s", "%s");`,
  u.UserName, q.Question.Assignment, q.Question.level, q.Question.number,
  q.QuestionInstance.NumAttempts + 1, 0, q.Question.weight, q.QuestionInstance.answer, tf)

//run to get the user's row - this is the user's history that will be sent to python script
rows, err := db.Query(`SELECT * FROM test02.scores WHERE username = "%s";`, u.username)


//////////////AFTER RECEIVING RESPONSE FROM PYTHON SCRIPT////////////

//run IF q.QuestionInstance.status == "Correct" was return from script where q is the json returned from the script
r2 := fmt.Sprintf(`update test02.responses SET correctness = 1
  WHERE username="%s" AND assignment="%s" AND level="%d" AND numb="%d" AND attempt="%d";`,
   u.UserName, q.Question.Assignment, q.Question.level, q.Question.number, q.QuestionInstance.NumAttempts)

//run IF Question.level == 1 AND q.QuestionInstance.status == "Correct"
r3 := fmt.Sprintf(`update test02.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + "%d",
  score1_correct = score1_correct + 1 WHERE username = "%s";`, q.QuestionInstance.numAttemps,  u.UserName)

//run IF Question.level == 1 AND q.QuestionInstance.status == "IncorrectNoAttempts"
r4 := fmt.Sprintf(`update test02.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + "%d"
  WHERE username = "%s";`, q.QuestionInstance.numAttemps,  u.UserName)

//run IF Question.level == 2 AND q.QuestionInstance.status == "Correct"
r5 := fmt.Sprintf(`update test02.scores SET next2 = next2 + 1, score2_attempts = score2_attempts + "%d",
  score2_correct = score2_correct + 1 WHERE username = "%s";`, q.QuestionInstance.numAttemps,  u.UserName)

//run IF Question.level == 2 AND q.QuestionInstance.status == "IncorrectNoAttempts"
r6 := fmt.Sprintf(`update test02.scores SET next2 = next2 + 1, score2_attempts = score2_attempts + "%d"
  WHERE username = "%s";`, q.QuestionInstance.numAttemps,  u.UserName)

//run IF Question.level == 3 AND q.QuestionInstance.status == "Correct"
r7 := fmt.Sprintf(`update test02.scores SET next3 = next3 + 1, score3_attempts = score3_attempts + "%d",
  score3_correct = score3_correct + 1 WHERE username = "%s";`, q.QuestionInstance.numAttemps,  u.UserName)

//run IF Question.level == 3 AND q.QuestionInstance.status == "IncorrectNoAttempts"
r8 := fmt.Sprintf(`update test02.scores SET next3 = next3 + 1, score3_attempts = score3_attempts + "%d"
  WHERE username = "%s";`, q.QuestionInstance.numAttemps,  u.UserName)

//run IF Question.level == 4 AND q.QuestionInstance.status == "Correct"
r9 := fmt.Sprintf(`update test02.scores SET next2 = next4 + 1, score4_attempts = score4_attempts + "%d",
  score4_correct = score4_correct + 1 WHERE username = "%s";`, q.QuestionInstance.numAttemps,  u.UserName)

//run IF Question.level == 4 AND q.QuestionInstance.status == "IncorrectNoAttempts"
r10 := fmt.Sprintf(`update test02.scores SET next4 = next4 + 1, score4_attempts = score4_attempts + "%d"
  WHERE username = "%s";`, q.QuestionInstance.numAttemps,  u.UserName)


/////////////////END OF ASSIGNMENT//////////////////

//run when assignment ends IF q.Question.Assignment == "ClimateChange"
r11 := fmt.Sprintf(`update test02.scores SET cc1 = CASE
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
END;`)

//run when assignment ends IF q.Question.Assignment == "Reading Test Scores"
r12 := fmt.Sprintf(`update test02.scores SET rts1 = CASE
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
END;`)

//run when assignment ends and AFTER running previous three queries if they are relevant
r13 := fmt.Sprintf(`update test02.scores SET
score1_correct = "%d", score1_attempts = "%d",
score2_correct = "%d", score2_attempts = "%d",
score3_correct = "%d", score3_attempts = "%d",
score4_correct = "%d", score4_attempts = "%d",
next1 = "%d", next2 = "%d", next3 = "%d", next4 = "%d";`,
0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1)

//run when assignment ends to return the user's score - this is the user's score for the assignment
rows1, err1 := db.Query(`SELECT username, assignment, level, numb,
    MAX(correctness) AS correctness, FIRST_VALUE(score_possible) AS score_possible
    FROM test02.responses WHERE username = "%s" AND assignment = "%s"
    GROUP BY username, assignment, level, numb;)`, u.UserName, q.Question.Assignment)

rows2, err2 := db.Query(`SELECT assignment,
    SUM(correctness*score_possible)/SUM(score_possible) AS prescore
    FROM rows1 GROUP BY username, assignment;)`)

rows3, err3 := db.Query(`SELECT prescore * weight AS score
    FROM rows2 INNER JOIN test02.weights ON rows2.assignment = test02.weights.assignment;)`)
