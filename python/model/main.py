import answers
import users
import cleanData
import cluster
import pandas as pd
import julia

preprocess = 1
assignmentNames = ["ad","sd", "de"]
levelQuesAD = [6,5,5,5]
maxAttemptsAD = [3] * 6 + [1,3,2,2,2] + [1,1,5,3,3] + [2,3,3,2,2]
levelQuesSD = [8,3,4,3]
maxAttemptsSD = [3] * 8 + [2,1,1] + [2] * 6 + [1]
levelQuesDE = [5,5,6,4]
maxAttemptsDE = [3,2,3,3,2] + [2,2,3,1,4] + [3,3,1,5,5,3] + [3,5,5,1]
levelQuesCC = [2,3,2,2]
maxAttemptsCC = [5,2,1,2,2,5,4,4,5]
levelQuesRTS = [5,4,5,4]
maxAttemptsRTS = [3,3,2,2,2] + [3,2,2,2] + [5,3,2,1,2] + [5] * 4
levelQuesDFE = [6,4,5,6]
maxAttemptsDFE = [2,2,1,1,2,5] + [1,3,3,3] + [3,1,3,3,1] + [3,1,1,3,5,1]
assignmentLevels = [levelQuesAD, levelQuesSD, levelQuesDE]

if preprocess:
    #users
    users2018 = "../../../../../Desktop/Fall2018/PL/2018Data/userEnrolledData.csv"
    users2017 = "../../../../../Desktop/Fall2018/PL/2017Data/userEnrolledData.csv"
    userList = users.cleanUsers([users2018,users2017])

    #analytical detective
    ad2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/AnalyticalDetective/AnalyticalDetective{}_{}.csv"
    ad2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/AnalyticalDetective/AnalyticalDetective{}_{}.csv"
    anyticalDetective = answers.parseAndGroup(levelQuesAD,[ad2018,ad2017],"ad")

    #stock dynamics
    sd2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/StockDynamics/StockDynamics{}_{}.csv"
    sd2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/StockDynamics/StockDynamics{}_{}.csv"
    stockDynamics = answers.parseAndGroup(levelQuesSD,[sd2018,sd2017],"sd")

    #demographics and emloyment
    de2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/DemographicsEmployment/DemographicsEmployment{}_{}.csv"
    de2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/DemographicsEmployment/DemographicsEmployment{}_{}.csv"
    demographicsEmployment = answers.parseAndGroup(levelQuesDE,[de2018,de2017],"de")

    #climateChange
    cc2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/ClimateChange/ClimateChange{}_{}.csv"
    cc2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/ClimateChange/ClimateChange{}_{}.csv"
    climateChange = answers.parseAndGroup(levelQuesCC,[cc2018,cc2017],"cc")

    #reading test scores
    rts2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/ReadingTestScores/ReadingTestScores{}_{}.csv"
    rts2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/ReadingTestScores/ReadingTestScores{}_{}.csv"
    climateChange = answers.parseAndGroup(levelQuesRTS,[rts2018,rts2017],"rts")

    #detecting flu epedemics
    dfe2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/DetectingFluEpedemics/DetectingFluEpedemics{}_{}.csv"
    dfe2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/DetectingFluEpedemics/DetectingFluEpedemics{}_{}.csv"
    climateChange = answers.parseAndGroup(levelQuesDFE,[dfe2018,dfe2017],"dfe")

    #exam
    cols = ["username", "examScore", "courseYear"]
    levelExam17 = [7]
    exam2017 = "../../../../../Desktop/Fall2018/PL/2017Data/lrExam/{}.csv"
    exam2017 = answers.parseAndGroup(levelExam17,[exam2017],"exam2017")
    exam2017.loc[:,'examScore'] = exam2017.apply(lambda row: answers.calcScore(row, "exam2017", 1, 8),axis=1)
    exam2017.to_csv("data/temp1.csv")
    exam2017 = exam2017[cols]

    levelExam18 = [4]
    exam2018 = "../../../../../Desktop/Fall2018/PL/2018Data/lrExam/{}.csv"
    exam2018 = answers.parseAndGroup(levelExam18,[exam2018],"exam2018")
    exam2018.loc[:,'examScore'] = exam2018.apply(lambda row: answers.calcScore(row, "exam2018", 1, 5),axis=1)
    exam2018.to_csv("data/temp2.csv")
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
    data.to_csv("data/preprocessed.csv")
else:
    data = pd.read_csv("data/preprocessed.csv",header=0)

# current position in each level
# position = [2,3,2,1]
positions = [[[2,1,1,1],[2,2,1,1],[2,3,1,1]],
             [[2,2,1,1],[2,3,1,1],[2,2,2,1],[2,3,2,1]],
             [[2,2,2,1],[2,3,2,1],[2,2,2,2],[2,3,2,2]],
             [[2,3,2,2]]]
weight = [[1.9,1,0,0,0],
          [2.5,1.4,1,0,0],
          [0,1.4,1,1,0],
          [0,0,1,1,1]]
data1 = cluster.setYGroup(data,2)
dataLev = []
for p in positions[0]:
    dataFilled = answers.sortToPredict(data1, assignmentLevels, assignmentNames, levelQuesCC, 1, p, "cc", weight[0], 2)
    cleanData.convertTreatment(dataFilled)
    dataLev.append(dataFilled)
final = pd.concat(dataLev)
final.to_csv("data/climateChange1.csv")
data2 = cluster.setYGroup(data,3)
for currlevel in range(3,4):
    dataLev = []
    for p in positions[currlevel - 1]:
        dataFilled = answers.sortToPredict(data2, assignmentLevels, assignmentNames, levelQuesCC, currlevel, p, "cc", weight[currlevel - 1])
        cleanData.convertTreatment(dataFilled)
        dataLev.append(dataFilled)
    final = pd.concat(dataLev)
    final.to_csv("data/climateChange{}.csv".format(currlevel))
# prediction of y3 is exam
dataLev = []
for p in positions[3]:
    dataFilled = answers.sortToPredict(data2, assignmentLevels, assignmentNames, levelQuesCC, 4, p, "cc", weight[3])
    cleanData.convertTreatment(dataFilled)
    dataLev.append(dataFilled)
final = pd.concat(dataLev)
final.to_csv("data/climateChange4.csv")

j = julia.Julia()
x = j.include("PrescriptiveTree.jl")
