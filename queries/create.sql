create table responses(
	username VARCHAR(100) NOT NULL,
	assignment VARCHAR(100) NOT NULL,
	level INT NOT NULL,
	numb INT NOT NULL,
  attempt INT NOT NULL,
	correctness INT NOT NULL,
  score_possible INT NOT NULL,
	answer VARCHAR(1000) NOT NULL,
  answer_timestamp TIMESTAMP,
  PRIMARY KEY ( username, assignment, level, numb, attempt)
);

create table scores(
	username VARCHAR(50) NOT NULL,
	gender VARCHAR(20) NOT NULL,
	level_of_education VARCHAR(20) NOT NULL,
	enrollment_mode VARCHAR(45) NOT NULL,
  ageCategory VARCHAR(45) NOT NULL,
	ad1 FLOAT NOT NULL,
	ad2 FLOAT NOT NULL,
  ad3 FLOAT NOT NULL,
  ad4 FLOAT NOT NULL,
	sd1 FLOAT NOT NULL,
	sd2 FLOAT NOT NULL,
  sd3 FLOAT NOT NULL,
  sd4 FLOAT NOT NULL,
	de1 FLOAT NOT NULL,
	de2 FLOAT NOT NULL,
  de3 FLOAT NOT NULL,
  de4 FLOAT NOT NULL,
	cc1 FLOAT NOT NULL,
	cc2 FLOAT NOT NULL,
  cc3 FLOAT NOT NULL,
  cc4 FLOAT NOT NULL,
	rts1 FLOAT NOT NULL,
	rts2 FLOAT NOT NULL,
  rts3 FLOAT NOT NULL,
  rts4 FLOAT NOT NULL,
	score1_correct INT NOT NULL,
  score1_attempts INT NOT NULL,
	score2_correct INT NOT NULL,
  score2_attempts INT NOT NULL,
	score3_correct INT NOT NULL,
  score3_attempts INT NOT NULL,
	score4_correct INT NOT NULL,
  score4_attempts INT NOT NULL,
	next1 FLOAT NOT NULL,
	next2 FLOAT NOT NULL,
  next3 FLOAT NOT NULL,
  next4 FLOAT NOT NULL,
  PRIMARY KEY ( username )
);

create table weights(
	assignment VARCHAR(100) NOT NULL,
	weight INT NOT NULL,
	PRIMARY KEY ( assignment )
);
