INSERT INTO pl.weights (assignment, weight)
VALUES
    ("cc",14),
    ("rts",24),
    ("dfe",28);
SELECT * FROM pl.weights;

SELECT username, assignment FROM pl.scores
  WHERE username = "100" AND assignment = "Climate+Change";

insert into pl.responses
  (username, assignment, level, numb, attempt, correctness, score_possible, answer, answer_timestamp)
  values ("100", "cc", 1, 2, 1, 0, 2, "answer", "20060102150405");
SELECT * FROM pl.responses;

SELECT * FROM pl.scores WHERE username = "100" AND assignment = "Climate+Change";

update pl.responses SET correctness = 1
  WHERE username="100" AND assignment="cc" AND level=1 AND numb=2 AND attempt=1;

update pl.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + 1,
  score1_correct = score1_correct + 1 WHERE username = "100" AND assignment = "Climate+Change";
update pl.scores SET next1 = next1 + 1, score1_attempts = score1_attempts + 1
  WHERE username = "100" AND assignment = "Climate+Change";

update pl.scores SET cc1 = CASE
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
WHERE username = "100" AND assignment = "Climate+Change";

SELECT * FROM pl.scores;

update pl.scores SET
score1_correct = 0, score1_attempts = 0,
score2_correct = 0, score2_attempts = 0,
score3_correct = 0, score3_attempts = 0,
score4_correct = 0, score4_attempts = 0,
next1 = 1, next2 = 1, next3 = 1, next4 = 1;


CREATE TEMPORARY TABLE temp
SELECT username, assignment, level, numb, MAX(correctness) AS correctness, ANY_VALUE(score_possible) AS score_possible
FROM pl.responses
WHERE username = "100" AND assignment = "cc"
GROUP BY username, assignment, level, numb;

CREATE TEMPORARY TABLE temp1
SELECT assignment,
    SUM(correctness*score_possible)/SUM(score_possible) AS prescore
    FROM temp GROUP BY username, assignment;

SELECT prescore * weight AS score
    FROM temp1 INNER JOIN pl.weights ON temp1.assignment = pl.weights.assignment;

SELECT * FROM temp1;

DROP TEMPORARY TABLE temp;

DROP TEMPORARY TABLE temp1;
