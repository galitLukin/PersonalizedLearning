CREATE TABLE answered
SELECT DISTINCT username, assignment
FROM responses;

SELECT COUNT(*)
FROM answered;

CREATE TABLE needGrade
SELECT answered.username, answered.assignment 
FROM answered
INNER JOIN scores
ON answered.username = scores.username AND answered.assignment = scores.assignment
WHERE scores.grade = 0.0;

SELECT COUNT(*)
FROM needGrade;

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
SELECT username, assignment, MAX(level) as maxLev, SUM(score) as score, SUM(score_possible) as score_possible
FROM res2
GROUP BY username, assignment;

CREATE TABLE res4
SELECT username, assignment, score/(score_possible+1+2+3+4) as grade
FROM res3
WHERE maxLev = 1
UNION
SELECT username, assignment, score/(score_possible+2+3+4) as grade
FROM res3
WHERE maxLev = 2
UNION
SELECT username, assignment, score/(score_possible+3+4) as grade
FROM res3
WHERE maxLev = 3
UNION
SELECT username, assignment, score/(score_possible+4) as grade
FROM res3
WHERE maxLev = 4;

select count(*)
from res4;

create table toinsert
SELECT userMap.nickname, res4.assignment, res4.grade
FROM res4 INNER JOIN userMap
ON res4.username = userMap.username;

SELECT COUNT(*)
FROM toinsert;

SELECT userMap.nickname, res4.username, res4.assignment, res4.grade
FROM res4 LEFT OUTER JOIN userMap
ON res4.username = userMap.username;

DROP TABLE answered;
DROP TABLE needGrade;
DROP TABLE res1;
DROP TABLE res2;
DROP TABLE res3;
DROP TABLE res4;
DROP TABLE toinsert;
