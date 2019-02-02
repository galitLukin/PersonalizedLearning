/*run offline before the start of the experiment for every user taking part in the experiment group
I will run this with the data from edx
part of the value will be constant and others will need to be updated but no rows will be added
TODO GALIT: write procedure to do this automatically from a csv*/
INSERT INTO scores (username, gender, level_of_education, enrollment_mode, ageCategory, ad1, ad2, ad3, ad4, sd1, sd2, sd3, sd4, de1, de2, de3, de4, cc1, cc2, cc3, cc4, rts1, rts2, rts3, rts4, score1_correct, score1_attempts, score2_correct, score2_attempts, score3_correct, score3_attempts, score4_correct, score4_attempts )
VALUES ();

/*run IF Question.level == 1 AND (QuestionInstance.status == "Correct" OR QuestionInstance.status == "IncorrectNoAttempts")
TODO OMER: insert username and make sure I put in QuestionInstance.numAttemps and 
QuestionInstance.status correctly*/
UPDATE scores SET score1_attempts = score1_attempts + QuestionInstance.numAttemps, 
score1_correct = CASE
   WHEN QuestionInstance.status= "Correct" THEN score1_correct + 1
   ELSE score1_correct
END
WHERE username = "";

/*run IF Question.level == 2 AND (QuestionInstance.status == "Correct" OR QuestionInstance.status == "IncorrectNoAttempts")
This is the same as the previous query but with changing the values of level 2 so same TODO*/
UPDATE scores SET score2_attempts = score2_attempts + QuestionInstance.numAttemps, 
score2_correct = CASE
   WHEN QuestionInstance.status= "Correct" THEN score2_correct + 1
   ELSE score2_correct
END
WHERE username = "";

/*run IF Question.level == 3 AND (QuestionInstance.status == "Correct" OR QuestionInstance.status == "IncorrectNoAttempts")
same for level 3*/
UPDATE scores SET score3_attempts = score3_attempts + QuestionInstance.numAttemps, 
score3_correct = CASE
   WHEN QuestionInstance.status= "Correct" THEN score3_correct + 1
   ELSE score3_correct
END
WHERE username = "";

/*run IF Question.level == 3 AND (QuestionInstance.status == "Correct" OR QuestionInstance.status == "IncorrectNoAttempts")
same for level 4*/
UPDATE scores SET score4_attempts = score4_attempts + QuestionInstance.numAttemps, 
score4_correct = CASE
   WHEN QuestionInstance.status= "Correct" THEN score4_correct + 1
   ELSE score4_correct
END
WHERE username = "";

/*run when assignment ends IF Question.Assignment == "Climate Change"
TODO GALIT: change the division so that its a float*/
UPDATE scores SET 
cc1 = CASE
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
   WHEN score1_attempts > 0 THEN score4_correct/score4_attempts
   ELSE 0
END;

/*run when assignment ends IF Question.Assignment == "Reading Test Scores"
TODO GALIT: change the division so that its a float
same as previous but for second assignment*/
UPDATE scores SET 
rts1 = CASE
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
   WHEN score1_attempts > 0 THEN score4_correct/score4_attempts
   ELSE 0
END;

/*run when assignment ends and after running previous two queries if relevant*/
UPDATE scores SET 
score1_correct = 0, score1_attempts = 0,
score2_correct = 0, score2_attempts = 0,
score3_correct = 0, score3_attempts = 0,
score4_correct = 0, score4_attempts = 0;
