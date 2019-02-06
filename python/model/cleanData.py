import pandas as pd
import numpy as np
import math


def fillMissingData(df, assinmentNames):
	df = df.drop(columns=['id'])
	cols = list(df.columns.values)
	for c in cols:
		for assignment in assinmentNames:
			if assignment in c:
				df[c].fillna(0, inplace=True)
	if "cc1" in cols:
		df['cc1'].fillna(0, inplace=True)
		df['cc2'].fillna(0, inplace=True)
		df['cc3'].fillna(0, inplace=True)
		df['cc4'].fillna(0, inplace=True)
		if "rts1" in cols:
			df['rts1'].fillna(0, inplace=True)
			df['rts2'].fillna(0, inplace=True)
			df['rts3'].fillna(0, inplace=True)
			df['rts4'].fillna(0, inplace=True)
	df.location.replace(np.nan, "None", inplace=True)
	df.loc[:,'location'] = df['location'].apply(lambda x: x.upper())
	df.location.replace("BROOKLYN, NEW YORK", "US", inplace=True)
	df.location.replace("ALABAMA", "US", inplace=True)
	df.location.replace("MALAGA SPAIN", "ES", inplace=True)
	df.location.replace("CHICAGO", "US", inplace=True)
	df.location.replace("PORTUGAL", "PT", inplace=True)
	df.location.replace("PUNE,MAHARASTRA,INDIA", "IN", inplace=True)
	df.location.replace("BOGOTA", "CO", inplace=True)
	df.location.replace("PORTUGAL", "PT", inplace=True)
	df.location.replace("CHENNAI", "IN", inplace=True)
	df.location.replace("CANADA", "CA", inplace=True)
	df.location.replace("LIEGE, BELGIUM", "CA", inplace=True)
	df.location.replace("MONTCLAIR,NJ,USA", "US", inplace=True)
	df.location.replace("SUITA, OSAKA, JAPAN", "JP", inplace=True)
	df.location.replace("LIEGE, BELGIUM", "BE", inplace=True)
	df.location.replace("ANTWERP, BELGIUM", "BE", inplace=True)
	df.location.replace("DELHI, INDIA", "IN", inplace=True)
	df.location.replace("FREMONT, CALIFORNIA, USA", "US", inplace=True)
	df.location.replace("AVON MN", "US", inplace=True)
	df.location.replace("ILLINOIS", "US", inplace=True)
	df.location.replace("CHICAGO", "US", inplace=True)
	df.location.replace("USA", "US", inplace=True)
	df.location.replace("MAHARASTRA", "IN", inplace=True)
	df.gender.replace(0, "None", inplace=True)
	df.level_of_education.replace(np.nan, "None", inplace=True)
	df.enrollment_mode.replace(np.nan, "None", inplace=True)
	df.year_of_birth.replace("None", np.nan, inplace=True)
	df.examScore.replace(np.nan, 0.0, inplace=True)
	return df


def fixAttemptsForIncorrect(df, levelQues, qtype, maxAttempts = []):
	maxAttemptsDict = {}
	keys = []
	for i in range(1,len(levelQues) + 1):
		for j in range(1,levelQues[i-1]+1):
			keys.append("{}{}.{}".format(qtype,i,j))
	for i in range(len(keys)):
		maxAttemptsDict[keys[i]] = maxAttempts[i]
	for i in range(1,len(levelQues) + 1):
		for j in range(1,levelQues[i-1]+1):
			df.loc[:,"{}{}.{}_attempts".format(qtype,i,j)] = df.apply(lambda row: maxAttemptsDict["{}{}.{}".format(qtype,i,j)] \
						if row["{}{}.{}_correct".format(qtype,i,j)] == 0 \
						else row["{}{}.{}_attempts".format(qtype,i,j)],axis=1)
	return df

def treatmentMap(treatment):
	if treatment == 1:
		return "A"
	elif treatment == 2:
		return "B"
	return "C"


def convertTreatment(df):
	df.loc[:,'Ygroup'] = df['Ygroup'].apply(lambda row: treatmentMap(row))


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


def makeAgeCategorical(df):
	df["ageCategory"] = np.nan
	df.loc[:,'ageCategory'] = df['year_of_birth'].apply(lambda row: ageCategories(row))


def addCols(data, asmt):
	if asmt == "cc":
		for i in range(1,5):
			data.insert(15+i, "cc{}".format(i), [0]*len(data.index))
		for i in range(1,5):
			data.insert(19+i, "rts{}".format(i), [0]*len(data.index))
	elif asmt == "rts":
		for i in range(1,5):
			data.insert(19+i, "rts{}".format(i), [0]*len(data.index))
	return data
