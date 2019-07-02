import pandas as pd
today = "0312m"
gradeReport = pd.read_csv("./python/model/data/exp/gradeReport/{}.csv".format(today),header=0)
userIDs = pd.read_csv("./python/model/data/exp/users/userids.csv",header=0)
gradeReport = gradeReport.loc[gradeReport['Experiment Group (Personalization in Linear Regression Assignment)'] == "Personalized"]
personalizedUsers = gradeReport.loc[:,['Student ID','Username']]
userIDs.columns = ['id', 'anonymized_id', 'username']
personalizedUsers.columns = ['id','nickname']
print(len(personalizedUsers))
userMap = pd.merge(userIDs, personalizedUsers, on='id', how='inner')
print(len(userMap))
userMap = userMap.loc[:,['username','nickname']]
userMap.to_csv("./python/model/data/exp/users/userMap.csv", index=False)

userMap = pd.merge(personalizedUsers, userIDs, on='id', how='left')
print(len(userMap))
userMap = userMap.loc[:,['username','nickname']]
userMap = userMap[userMap['username'].isnull()]
print(len(userMap))
userMap.to_csv("./python/model/data/exp/users/userMapMissing.csv", index=False)




