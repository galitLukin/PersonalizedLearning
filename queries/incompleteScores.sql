CREATE TABLE answered
SELECT DISTINCT username, assignment
FROM responses;

CREATE TABLE needGrade
SELECT answered.username, answered.assignment 
FROM answered
LEFT OUTER JOIN scores
ON answered.username = scores.username AND answered.assignment = scores.assignment
WHERE scores.grade = 0.0;

CREATE TABLE res1
SELECT responses.username, responses.assignment, level, numb,
	    MAX(correctness) AS correctness, ANY_VALUE(score_possible) AS score_possible
FROM responses
INNER JOIN needGrade ON responses.username = needGrade.username AND responses.assignment = needGrade.assignment
GROUP BY responses.username, responses.assignment, level, numb;

CREATE TABLE res2
SELECT username, assignment, level, numb, correctness*score_possible as score, score_possible
FROM res1;

CREATE TABLE res3
SELECT username, assignment, SUM(score) as score, SUM(score_possible) as score_possible
FROM res2
GROUP BY username, assignment;

SELECT username, assignment, score/score_possible as grade
FROM res3;

DROP TABLE answered;
DROP TABLE needGrade;
DROP TABLE res1;
DROP TABLE res2;
DROP TABLE res3;
