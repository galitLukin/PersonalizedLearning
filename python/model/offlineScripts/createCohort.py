import pandas as pd
today = "0312m"

gradeReport = pd.read_csv("./python/model/data/exp/gradeReport/{}.csv".format(today),header=0)
gradeReport = gradeReport.loc[gradeReport['Experiment Group (Personalization in Linear Regression Assignment)'] == "Personalized"]
print(len(gradeReport))
df = gradeReport.loc[gradeReport['Cohort Name'] != "Personalized"]
df.to_csv("./python/model/data/exp/cohorts/{}.csv".format(today), index=False)

cohort = pd.read_csv("./python/model/data/exp/cohorts/all.csv",header=0)
df = df.loc[:,'Username']
df = df.drop_duplicates().to_frame()
frames = [cohort, df]
result = pd.concat(frames)
result = result.drop_duplicates()
result.to_csv("./python/model/data/exp/cohorts/all.csv", index=False)
print(len(result))