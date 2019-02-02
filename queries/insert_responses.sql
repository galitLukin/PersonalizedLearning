/*run when user submits answer
TODO OMER: need to add values to query - the values arr from the same json that is sent to the python script:
username = you have it
assignment = Question.Assignment
level = Question.level
number = Question.number
attempt = QuestionInstance.NumAttempts + 1
correctness = 0
score_possible = Question.weight
answer = QuestionInstance.answer
*/
INSERT INTO responses (username, assignment, level, number, attempt, correctness, score_possible, answer, timestamp)
VALUES ();

/*run IF QuestionInstance.status == "Correct"
TODO OMER: insert values of where clause - the same values as above but this time, it is from the json 
returned from the python script*/
UPDATE responses SET correctness = 1
WHERE username = ""
AND assignment = ""
AND level = 0
AND number = 0 
AND attempt = 0
