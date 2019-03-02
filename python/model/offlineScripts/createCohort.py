import pandas as pd
start = "0226"
today = "0301"
gradeReport = pd.read_csv("./python/model/data/exp/gradeReport/{}.csv".format(today),header=0)
gradeReport = gradeReport.loc[gradeReport['Experiment Group (Personalization in Linear Regression Assignment)'] == "Personalized"]
personalizedUsers = gradeReport.loc[:,'Username']
personalizedUsers = personalizedUsers.drop_duplicates().to_frame()
cohort = pd.read_csv("./python/model/data/exp/cohorts/{}.csv".format(start),header=0)
df = cohort.merge(personalizedUsers, on=['Username'], how='left', indicator=True)
df = df.loc[df['_merge'] == 'left_only']
df.to_csv("./python/model/data/exp/cohorts/{}.csv".format(today), index=False)
