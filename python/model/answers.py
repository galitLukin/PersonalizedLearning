import json
import pandas as pd
import numpy as np

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


def groupPerUser(df, levelQues, qtype):
	answers = df.groupby(['username','courseYear']).agg(lambda x: list(x)).reset_index()
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
	answers = answers.drop(columns=['qId', 'level','attempts','correctness'])
	return answers


def parseAndGroup(levelQues,data,qtype):
	newCols = ["correct_map"," input_state","last_submission_time","attempts","score","done","student_answers","seed"]
	mapCols = ["hint","hintmode","correctness","msg","answervariable","npoints","queuestate"]
	dfs = []
	for course in data:
		d = ParseAndCombine(newCols, mapCols, levelQues, course, qtype)
		if "18" in course:
			d.loc[:,'courseYear'] = [18]*len(d.index)
		else:
			d.loc[:,'courseYear'] = [17]*len(d.index)
		dfs.append(d)
	df = pd.concat(dfs)
	return groupPerUser(df, levelQues, qtype)


def calcScore(row, assignment, level, position):
	correct = 0
	attempts = 0
	for j in range(1, position):
		correct = correct + row["{}{}.{}_correct".format(assignment,level,j)]
		attempts = attempts + row["{}{}.{}_attempts".format(assignment,level,j)]
	return float(correct)/attempts if attempts > 0 else 0

def countCorrect(row, assignment, level, position):
	correct = 0
	for j in range(1, position):
		correct = correct + row["{}{}.{}_correct".format(assignment,level,j)]
	return correct

def countAttempts(row, assignment, level, position):
	attempts = 0
	for j in range(1, position):
		attempts = attempts + row["{}{}.{}_attempts".format(assignment,level,j)]
	return attempts

def calcY(row, assignment, level, position):
	correct = row["{}{}.{}_correct".format(assignment,level,position)]
	attempts = row["{}{}.{}_attempts".format(assignment,level,position)]
	return float(correct)/attempts if attempts > 0 else 0

def sortToPredict(data, levelQues, assignmentName, personalizedLevelQues, currLevel, position, assignment, w, groups=3):
	cols = ['location', 'gender', 'level_of_education', 'enrollment_mode', 'Ygroup', 'ageCategory', 'examScore']
	if assignment == "rts" or assignment == "dfe":
		cols = cols + ["cc1","cc2","cc3","cc4"]
		if assignment == "dfe":
			cols = cols + ["rts1","rts2","rts3","rts4"]
	df = data.loc[:,cols]
	for k in range(len(assignmentName)):
		for i in range(1,len(levelQues[k]) + 1):
			df.loc[:,"{}{}".format(assignmentName[k],i)] = np.zeros(df.shape[0])
			for j in range(1,levelQues[k][i-1]+1):
				df.loc[:,"{}{}.{}".format(assignmentName[k],i,j)] = data.apply(lambda row: float(row["{}{}.{}_correct".format(assignmentName[k],i,j)])/row["{}{}.{}_attempts".format(assignmentName[k],i,j)], axis=1)
				df.loc[:,"{}{}".format(assignmentName[k],i)] = df.loc[:,"{}{}".format(assignmentName[k],i)] + df.loc[:,"{}{}.{}".format(assignmentName[k],i,j)]
			df.loc[:,"{}{}".format(assignmentName[k],i)] = df.loc[:,"{}{}".format(assignmentName[k],i)].apply(lambda x: float(x)/levelQues[k][i-1])
	for i in range(1,5):
		#df.loc[:,'score'+str(i)+'_correct'] = data.apply(lambda row: countCorrect(row, assignment, i, position[i-1]),axis=1)
		#df.loc[:,'score'+str(i)+'_attempts'] = data.apply(lambda row: countAttempts(row, assignment, i, position[i-1]),axis=1)
		df.loc[:,'score'+str(i)] = data.apply(lambda row: calcScore(row, assignment, i, position[i-1]),axis=1)
	#for i in range(1,5):
	#	df.loc[:,'next'+str(i)] = [position[i-1]]*len(df.index)
	if groups == 3:
		df.loc[:,'y1'] = data.apply(lambda row: -(currLevel - 1)*w[currLevel-2] * calcY(row, assignment, currLevel - 1, position[currLevel - 2]),axis=1)
		df.loc[:,'y2'] = data.apply(lambda row: -(currLevel)*w[currLevel-1] * calcY(row, assignment, currLevel, position[currLevel - 1]),axis=1)
		if currLevel == 4:
			if assignment == "dfe":
				df.loc[:,'y3'] = data.apply(lambda row: -5*w[currLevel]*row['examScore'],axis=1)
			elif assignment == "rts":
				df.loc[:,'y3'] = data.apply(lambda row: -5*w[currLevel]*row['dfeAvg'],axis=1)
			elif assignment == "cc":
				df.loc[:,'y3'] = data.apply(lambda row: -5*w[currLevel]*row['rtsAvg'],axis=1)
		else:
			df.loc[:,'y3'] = data.apply(lambda row: -(currLevel + 1)*w[currLevel] * calcY(row, assignment, currLevel + 1, position[currLevel]),axis=1)
		df.loc[:,'y'] = np.nan
		df.loc[df['Ygroup'] == 1,'y'] = df.loc[df['Ygroup'] == 1,'y1']
		df.loc[df['Ygroup'] == 2,'y'] = df.loc[df['Ygroup'] == 2,'y2']
		df.loc[df['Ygroup'] == 3,'y'] = df.loc[df['Ygroup'] == 3,'y3']
	else:
		df.loc[:,'y2'] = data.apply(lambda row: -(currLevel)*w[currLevel-1] * calcY(row, assignment, currLevel, position[currLevel - 1]),axis=1)
		df.loc[:,'y3'] = data.apply(lambda row: -(currLevel + 1)*w[currLevel] * calcY(row, assignment, currLevel + 1, position[currLevel]),axis=1)
		df.loc[:,'y'] = np.nan
		df.loc[df['Ygroup'] == 2,'y'] = df.loc[df['Ygroup'] == 2,'y2']
		df.loc[df['Ygroup'] == 3,'y'] = df.loc[df['Ygroup'] == 3,'y3']
	cols = list(df)
	cols = cols[-1:] + cols[4:5] + cols[0:4] + cols[5:-1]
	df = df.loc[:,cols]
	subcols = ['y','Ygroup','gender','level_of_education','enrollment_mode','ageCategory']
	for asmt in assignmentName:
		for i in range(1,5):
			subcols.append(asmt+str(i))
	if assignment == "rts" or assignment == "dfe":
		subcols = subcols + ["cc1","cc2","cc3","cc4"]
		if assignment == "dfe":
			subcols = subcols + ["rts1","rts2","rts3","rts4"]
	for i in range(1,5):
		#subcols.append('score'+str(i)+"_correct")
		#subcols.append('score'+str(i)+"_attempts")
		subcols.append('score'+str(i))
	#for i in range(1,5):
	#	subcols.append('next'+str(i))
	if groups == 3:
		subcols.append('y1')
	subcols.append('y2')
	subcols.append('y3')
	df = df.loc[:,subcols]
	return df
