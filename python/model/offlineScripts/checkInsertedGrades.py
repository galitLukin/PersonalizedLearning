import pandas as pd
today = "0308"

gradeReport = pd.read_csv("./python/model/data/exp/gradeReport/{}.csv".format(today),header=0)
gradeReport = gradeReport.loc[gradeReport['Experiment Group (Personalization in Linear Regression Assignment)'] == "Personalized"]
gradeReport = gradeReport.loc[:,['Student ID','Homework 2: Assignment 2']]
gradeReport.columns = ['id','hw2grade']
gradeReport = gradeReport.loc[gradeReport['hw2grade'] != "Not Attempted"]
userIDs = pd.read_csv("./python/model/data/exp/users/userids.csv",header=0)
userIDs.columns = ['id', 'anonymized_id', 'username']
userMap = pd.merge(userIDs, gradeReport, on='id', how='inner')
userMap = userMap.loc[:,['username','hw2grade']]

print(len(userMap))

scores = pd.read_csv("./python/model/data/exp/users/userHasGrade.csv".format(today),header=0)
scores = scores.groupby('username')['grade'].mean().to_frame()
print(len(scores))

res = pd.merge(userMap, scores, on='username', how='outer')
res = res[res['grade'].isnull() | res['hw2grade'].isnull()]
print(len(res))

res.to_csv("./python/model/data/exp/users/missingScores.csv", index=False)
