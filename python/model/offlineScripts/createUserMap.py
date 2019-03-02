import pandas as pd
gradeReport = pd.read_csv("./python/model/data/exp/users/grade_report.csv",header=0)
userIDs = pd.read_csv("./python/model/data/exp/users/userids.csv",header=0)
gradeReport = gradeReport.loc[gradeReport['Experiment Group (Personalization in Linear Regression Assignment)'] == "Personalized"]
personalizedUsers = gradeReport.loc[:,['Student ID','Username']]
userIDs.columns = ['id', 'anonymized_id', 'username']
personalizedUsers.columns = ['id','nickname']
userMap = pd.merge(userIDs, personalizedUsers, on='id', how='inner')
userMap = userMap.loc[:,['username','nickname']]
userMap.to_csv("./python/model/data/exp/userMap.csv", index=False)
