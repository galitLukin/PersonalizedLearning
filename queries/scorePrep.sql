CREATE TABLE allUsers
select username, assignment, grade
from scores;

CREATE TABLE newUsers
SELECT allUsers.username, allUsers.assignment, allUsers.grade
FROM allUsers LEFT OUTER JOIN userHasGrade 
ON allUsers.username =  userHasGrade.username AND allUsers.assignment = userHasGrade.assignment 
AND allUsers.grade = userHasGrade.grade;

#these are the users-assignment that need the score to be put in edX - save to usersNeedGrades
SELECT userMap.nickname, newUsers.assignment, newUsers.grade
FROM newUsers INNER JOIN userMap
ON newUsers.username = userMap.username;

CREATE TABLE withNewUsers
SELECT username, assignment, grade, 0 as isNew
FROM userHasGrade
UNION ALL
SELECT username, assignment, grade, 1 as isNew
FROM newUsers;

CREATE TABLE updatedUsers
SELECT username, assignment, MAX(isNew) as isNew
FROM withNewUsers
GROUP BY username, assignment;

#this replaces all of userHasGrade - save into userHasGrade
select withNewUsers.username, withNewUsers.assignment, withNewUsers.grade
FROM withNewUsers
INNER JOIN updatedUsers ON
withNewUsers.username = updatedUsers.username AND
withNewUsers.assignment = updatedUsers.assignment AND
withNewUsers.isNew = updatedUsers.isNew;

DROP TABLE allUsers;
DROP TABLE newUsers;
DROP TABLE withNewUsers;
DROP TABLE updatedUsers;







