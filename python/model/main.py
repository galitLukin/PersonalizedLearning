import answers
import users
import cleanData
import cluster
import pandas as pd
import julia

#users
users2018 = "../../../../../Desktop/Fall2018/PL/2018Data/userEnrolledData.csv"
users2017 = "../../../../../Desktop/Fall2018/PL/2017Data/userEnrolledData.csv"
userList = users.cleanUsers([users2018,users2017])

assignmentNames = ["ad","sd", "de"]

#analytical detective
levelQuesAD = [6,5,5,5]
maxAttemptsAD = [3] * 6 + [1,3,2,2,2] + [1,1,5,3,3] + [2,3,3,2,2]
ad2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/AnalyticalDetective/AnalyticalDetective{}_{}.csv"
ad2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/AnalyticalDetective/AnalyticalDetective{}_{}.csv"
anyticalDetective = answers.parseAndGroup(levelQuesAD,[ad2018,ad2017],"ad")

#stock dynamics
levelQuesSD = [8,3,4,3]
maxAttemptsSD = [3] * 8 + [2,1,1] + [2] * 6 + [1]
sd2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/StockDynamics/StockDynamics{}_{}.csv"
sd2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/StockDynamics/StockDynamics{}_{}.csv"
stockDynamics = answers.parseAndGroup(levelQuesSD,[sd2018,sd2017],"sd")

#demographics and emloyment
levelQuesDE = [5,5,6,4]
maxAttemptsDE = [3,2,3,3,2] + [2,2,3,1,4] + [3,3,1,5,5,3] + [3,5,5,1]
de2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/DemographicsEmployment/DemographicsEmployment{}_{}.csv"
de2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/DemographicsEmployment/DemographicsEmployment{}_{}.csv"
demographicsEmployment = answers.parseAndGroup(levelQuesDE,[de2018,de2017],"de")

assignmentLevels = [levelQuesAD, levelQuesSD, levelQuesDE]

#climateChange
#existing 2,3,3,1
levelQuesCC = [2,3,2,2]
maxAttemptsCC = [5,2,1,2,2,5,4,4,5]
cc2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/ClimateChange/ClimateChange{}_{}.csv"
cc2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/ClimateChange/ClimateChange{}_{}.csv"
climateChange = answers.parseAndGroup(levelQuesCC,[cc2018,cc2017],"cc")

#exam
cols = ["username", "examScore", "courseYear"]
levelExam17 = [7]
exam2017 = "../../../../../Desktop/Fall2018/PL/2017Data/lrExam/{}.csv"
exam2017 = answers.parseAndGroup(levelExam17,[exam2017],"exam2017")
exam2017.loc[:,'examScore'] = exam2017.apply(lambda row: -5 * answers.calcScore(row, "exam2017", 1, 8),axis=1)
exam2017 = exam2017[cols]

levelExam18 = [4]
exam2018 = "../../../../../Desktop/Fall2018/PL/2018Data/lrExam/{}.csv"
exam2018 = answers.parseAndGroup(levelExam18,[exam2018],"exam2018")
exam2018.loc[:,'examScore'] = exam2018.apply(lambda row: answers.calcScore(row, "exam2018", 1, 5),axis=1)
exam2018 = exam2018[cols]
exam = pd.concat([exam2017, exam2018])

#combine and clean data
data = pd.merge(userList, climateChange, on=['username','courseYear'], how='inner')
data = pd.merge(data, anyticalDetective, on=['username','courseYear'], how='left')
data = pd.merge(data, stockDynamics, on=['username','courseYear'], how='left')
data = pd.merge(data, demographicsEmployment, on=['username','courseYear'], how='left')
data = pd.merge(data, exam, on=['username','courseYear'], how='left')
data = cleanData.fillMissingData(data, assignmentNames)
data = cleanData.fixAttemptsForIncorrect(data, levelQuesAD, "ad", maxAttemptsAD)
data = cleanData.fixAttemptsForIncorrect(data, levelQuesSD, "sd", maxAttemptsSD)
data = cleanData.fixAttemptsForIncorrect(data, levelQuesDE, "de", maxAttemptsDE)
data = cleanData.fixAttemptsForIncorrect(data, levelQuesCC, "cc", maxAttemptsCC)
cleanData.makeAgeCategorical(data)
data.drop_duplicates()
cluster.clusterStudents(data, assignmentNames)

# current position in each level
# position = [2,3,2,1]
positions = [[[2,1,1,1],[2,2,1,1],[2,3,1,1]],
             [[2,2,1,1],[2,3,1,1],[2,2,2,1],[2,3,2,1]],
             [[2,2,2,1],[2,3,2,1],[2,2,2,2],[2,3,2,2]],
             [[2,3,2,2]]]
weight = [1,1,1,1,1]
data1 = cluster.setYGroup(data,2)
dataLev = []
for p in positions[0]:
    dataFilled = answers.sortToPredict(data1, assignmentLevels, assignmentNames, levelQuesCC, 1, p, "cc", weight, 2)
    cleanData.convertTreatment(dataFilled)
    dataLev.append(dataFilled)
final = pd.concat(dataLev)
final.to_csv("data/climateChange1.csv")
data2 = cluster.setYGroup(data,3)
for currlevel in range(2,4):
    dataLev = []
    for p in positions[currlevel - 1]:
        dataFilled = answers.sortToPredict(data2, assignmentLevels, assignmentNames, levelQuesCC, currlevel, p, "cc", weight)
        cleanData.convertTreatment(dataFilled)
        dataLev.append(dataFilled)
    final = pd.concat(dataLev)
    final.to_csv("data/climateChange{}.csv".format(currlevel))
# prediction of y3 is exam
dataLev = []
for p in positions[3]:
    dataFilled = answers.sortToPredict(data2, assignmentLevels, assignmentNames, levelQuesCC, 4, p, "cc", weight)
    cleanData.convertTreatment(dataFilled)
    dataLev.append(dataFilled)
final = pd.concat(dataLev)
final.to_csv("data/climateChange4.csv".format(4))

j = julia.Julia()
x = j.include("PrescriptiveTree.jl")
