import pandas as pd
import numpy as np

def location(row):
	if pd.notnull(row['country']):
		return row['country']
	else:
		return row['location']

def cleanUsers(data):
	cols = ["id","username","location","year_of_birth","gender","level_of_education","enrollment_mode","country"]
	dfs = []
	for course in data:
		df = pd.read_csv(course,header=0)
		df = df[cols]
		df['location'] = df.apply(lambda x: location(x),axis=1)
		df[cols[0:7]]
		if "18" in course:
			df.loc[:,'courseYear'] = np.ones(df.shape[0]) * 18
		elif "17" in course:
			df.loc[:,'courseYear'] = np.ones(df.shape[0]) * 17
		dfs.append(df)
	df = pd.concat(dfs)
	return df
