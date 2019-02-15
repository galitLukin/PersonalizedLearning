import pandas as pd
import numpy as np
import math
import json

# download all of 2019 data in ./python/model/data/exp/ and replace with 2018 data1
# do it for 3 assignments from unit 1 and for users enrolled and user userids
# run script that output scores csv
# use mysql bench to insert it into db

def location(row):
	if pd.notnull(row['country']):
		return row['country']
	else:
		return row['location']


def ageCategories(yob):
	if yob is None or ( isinstance(yob,float) and math.isnan(yob) ):
		return "Null"
	cutoffs = [1970,1985,1990,1995]
	if int(yob) <= cutoffs[0]:
		return "<={}".format(cutoffs[0])
	for i in range(1,len(cutoffs)):
		if int(yob) > cutoffs[i-1] and int(yob) <= cutoffs[i]:
			return "{}-{}".format(cutoffs[i-1],cutoffs[i])
	return ">{}".format(cutoffs[len(cutoffs) - 1])


def cleanUserData(df):
    cols = ["username","id","location","year_of_birth","gender","level_of_education","enrollment_mode","country"]
    df = df.loc[:,cols]
    #need to normalize location if we decide to use it
    df.loc[:,'location'] = df.apply(lambda x: location(x),axis=1)
    df.gender.replace(np.nan, "None", inplace=True)
    df.level_of_education.replace(np.nan, "None", inplace=True)
    df.enrollment_mode.replace(np.nan, "None", inplace=True)
    df.year_of_birth.replace("None", np.nan, inplace=True)
    #df.["ageCategory"] = np.nan
    df.loc[:,'ageCategory'] = df['year_of_birth'].apply(lambda row: ageCategories(row))
    df = df.drop(['country', 'year_of_birth'], axis=1)
    return df


def insertAssignment(df):
    #df["assignment"] = np.nan
    df.at[:, 'assignment'] = pd.Series([["Climate+Change", "Reading+Test+Scores", "Detecting+Flu+Epidemics+via+Search+Engine+Query+Data"]] * len(df.index))
    s = df.apply(lambda x: pd.Series(x['assignment']), axis=1).stack().reset_index(level=1, drop=True)
    s.name = 'assignment'
    df = df.drop('assignment', axis=1).join(s)
    df['assignment'] = pd.Series(df['assignment'], dtype=object)
    return df


def CustomParser(data):
	return json.loads(data)


def ParseAndCombine(newCols,mapCols,levelQues,data,qtype):
	final  = pd.DataFrame()
	for i in range(1,len(levelQues) + 1):
		for j in range(1,levelQues[i-1]+1):
			df = pd.read_csv(data.format(i,j), converters={'state':CustomParser},header=0)
			for c in newCols:
				df[c] = None
			for index, row in df.iterrows():
				d = row['state']
				for c in newCols:
					if c in d.keys():
						df.at[index,c] = d[c]
					else:
						df.at[index,c] = None
			for c in mapCols:
				df[c] = None
			df = df[pd.notnull(df['correct_map'])]
			for index, row in df.iterrows():
				for c in mapCols:
					df.at[index,c] = []
				for d in row['correct_map'].values():
					for c in mapCols:
						if c in d.keys():
							df.at[index,c].append(d[c])
						else:
							df.at[index,c]= None
			df = df.drop(['state'], axis=1)
			df = df.drop(['correct_map'], axis=1)
			df['qId'] = "{}{}.{}".format(qtype,i,j)
			df['level'] = i
			df['correctness'] = df.correctness.apply(lambda x: [1.0 if clause == "correct" else 0.0 for clause in x])
			if qtype == "cc":
				if i == 3 and j == 2:
					df['correctness'] = df.correctness.apply(lambda x: x[1])
				elif i == 4 and j == 1:
					df['correctness'] = df.correctness.apply(lambda x: x[0])
				elif i == 2 and j == 2:
					df['correctness'] = df.correctness.apply(lambda x: x[0])
				elif i == 2 and j == 3:
					df['correctness'] = df.correctness.apply(lambda x: x[1])
				else:
					df['correctness'] = df.correctness.apply(lambda x: sum(x)/float(len(x)))
			elif qtype == "rts":
				if i == 1 and j == 4:
					df['correctness'] = df.correctness.apply(lambda x: x[1])
				elif i == 1 and j == 5:
					df['correctness'] = df.correctness.apply(lambda x: x[0])
				elif i == 2 and j == 3:
					df['correctness'] = df.correctness.apply(lambda x: x[0])
				elif i == 2 and j == 4:
					df['correctness'] = df.correctness.apply(lambda x: x[1])
				else:
					df['correctness'] = df.correctness.apply(lambda x: sum(x)/float(len(x)))
			elif qtype == "dfe":
				if i == 3 and j == 3:
					df['correctness'] = df.correctness.apply(lambda x: x[1])
				elif i == 3 and j == 4:
					df['correctness'] = df.correctness.apply(lambda x: x[0])
				elif i == 4 and j == 2:
					df['correctness'] = df.correctness.apply(lambda x: x[1])
				elif i == 4 and j == 3:
					df['correctness'] = df.correctness.apply(lambda x: x[0])
				else:
					df['correctness'] = df.correctness.apply(lambda x: sum(x)/float(len(x)))
			else:
				df['correctness'] = df.correctness.apply(lambda x: sum(x)/float(len(x)))
			cols = ["qId","level","username","attempts","correctness"]
			final = pd.concat([final, df[cols]])
	return final


def calcScore(row, assignment, level, position):
	correct = 0
	attempts = 0
	for j in range(1, position):
		correct = correct + row["{}{}.{}_correct".format(assignment,level,j)]
		attempts = attempts + row["{}{}.{}_attempts".format(assignment,level,j)]
	return float(correct)/attempts if attempts > 0 else 0


def groupPerUser(df, levelQues, qtype):
    answers = df.groupby(['username']).agg(lambda x: list(x)).reset_index()
    for i in range(1,len(levelQues) + 1):
        for j in range(1,levelQues[i-1]+1):
            answers["{}{}.{}_attempts".format(qtype,i,j)] = 0.0
            answers["{}{}.{}_correct".format(qtype,i,j)] = 0.0
    for index, row in answers.iterrows():
        questions = row['qId']
        level = row['level']
        attempts = row['attempts']
        correct = row['correctness']
        for i in range(len(questions)):
            a = questions[i] + "_attempts"
            c = questions[i] + "_correct"
            answers.at[index,a] = attempts[i]
            answers.at[index,c] = correct[i]
    for level in range(1,5):
        answers.loc[:,qtype+str(level)] = answers.apply(lambda row: calcScore(row, qtype, level, levelQues[level-1]+1),axis=1)
    for c in answers.columns:
        if "_attempts" in c or "_correct" in c:
            answers = answers.drop(columns=c)
    answers = answers.drop(columns=['qId', 'level','attempts','correctness'])
    return answers


def parseAndGroup(levelQues,data,qtype):
	newCols = ["correct_map"," input_state","last_submission_time","attempts","score","done","student_answers","seed"]
	mapCols = ["hint","hintmode","correctness","msg","answervariable","npoints","queuestate"]
	dfs = []
	df = ParseAndCombine(newCols, mapCols, levelQues, data, qtype)
	return groupPerUser(df, levelQues, qtype)


userIDs = pd.read_csv("./python/model/data/exp/users/userids.csv",header=0)
userData = pd.read_csv("./python/model/data/exp/users/users.csv",header=0)
userList = cleanUserData(userData)
userIDs.columns = ['id', 'anonymized_id', 'user_id']
userList = pd.merge(userList, userIDs, on='id', how='inner')
userList = insertAssignment(userList)
cols = ["username", "user_id", "assignment", "gender","level_of_education","enrollment_mode","ageCategory"]
users = userList[cols]

levelQuesAD = [6,5,5,5]
maxAttemptsAD = [3] * 6 + [1,3,2,2,2] + [1,1,5,3,3] + [2,3,3,2,2]
levelQuesSD = [8,3,4,3]
maxAttemptsSD = [3] * 8 + [2,1,1] + [2] * 6 + [1]
levelQuesDE = [5,5,6,4]
maxAttemptsDE = [3,2,3,3,2] + [2,2,3,1,4] + [3,3,1,5,5,3] + [3,5,5,1]

#analytical detective
ad2019 = "./python/model/data/exp/AnalyticalDetective/AnalyticalDetective{}_{}.csv"
#TODO: dont use answer - create new one here
anyticalDetective = parseAndGroup(levelQuesAD,ad2019,"ad")

#stock dynamics
sd2019 = "./python/model/data/exp/StockDynamics/StockDynamics{}_{}.csv"
stockDynamics = parseAndGroup(levelQuesSD,sd2019,"sd")

#demographics and emloyment
de2019 = "./python/model/data/exp/DemographicsEmployment/DemographicsEmployment{}_{}.csv"
demographicsEmployment = parseAndGroup(levelQuesDE,de2019,"de")

data = pd.merge(users, anyticalDetective, on='username', how='left')
data = pd.merge(data, stockDynamics, on='username', how='left')
data = pd.merge(data, demographicsEmployment, on='username', how='left')
cols = list(data.columns.values)
for c in cols:
    data[c].fillna(0.0, inplace=True)
data.drop_duplicates()

cols = ['score1_correct','score2_correct','score3_correct','score4_correct',\
'score1_attempts','score2_attempts','score3_attempts','score4_attempts']
for c in cols:
    data.loc[:,c] = [0]*len(data.index)
cols = ['next1','next2','next3','next4']
for c in cols:
    data.loc[:,c] = [1]*len(data.index)
data.loc[:,'grade'] = ""*len(data.index)
data = data.drop('username', axis=1)
data.to_csv("./python/model/data/exp/scores.csv", index=False)
