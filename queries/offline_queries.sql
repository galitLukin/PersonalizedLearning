/*run these queries offline before the start of the experiment*/

INSERT INTO weights (assignment, weight)
VALUES
    ("Climate Change",14),
    ("Reading Test Scores",24),
    ("Detecting Flu Epedemics",28);

/* add the last columns to the file where all values will be zeros for cc, rtx, scores
and ones for next*/
LOAD DATA INFILE 'location of file'
INTO TABLE test02.scores
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;
