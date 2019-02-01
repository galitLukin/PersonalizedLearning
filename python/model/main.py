import answers
import users
import cleanData
import cluster
import pandas as pd
import julia

preprocess = 0
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
    ccCols = list(climateChange)

    #reading test scores
    rts2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/ReadingTestScores/ReadingTestScores{}_{}.csv"
    rts2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/ReadingTestScores/ReadingTestScores{}_{}.csv"
    testScores = answers.parseAndGroup(levelQuesRTS,[rts2018,rts2017],"rts")
    tsCols = list(testScores)

    climateChange = pd.merge(climateChange, testScores, on=['username','courseYear'], how='left')
    climateChange.loc[:,'rtsAvg'] = climateChange.apply(lambda row: [answers.calcScore(row, "rts", i, levelQuesRTS[i-1]) for i in range(1,5)],axis=1)
    climateChange.loc[:,'rtsAvg'] = climateChange.rtsAvg.apply(lambda row: float(sum(row))/len(row) if row else 0)
    climateChange.loc[:,'rtsAvg'] = climateChange.rtsAvg.apply(lambda row: row if row else 0)
    climateChange = climateChange.loc[:, ccCols + ['rtsAvg']]

    #detecting flu epedemics
    dfe2018 = "../../../../../Desktop/Fall2018/PL/2018Data/Assignment/DetectingFluEpedemics/DetectingFluEpedemics{}_{}.csv"
    dfe2017 = "../../../../../Desktop/Fall2018/PL/2017Data/Assignment/DetectingFluEpedemics/DetectingFluEpedemics{}_{}.csv"
    fluEpedemics = answers.parseAndGroup(levelQuesDFE,[dfe2018,dfe2017],"dfe")

    testScores = pd.merge(testScores, fluEpedemics, on=['username','courseYear'], how='left')
    testScores.loc[:,'dfeAvg'] = testScores.apply(lambda row: [answers.calcScore(row, "dfe", i, levelQuesDFE[i-1]) for i in range(1,5)],axis=1)
    testScores.loc[:,'dfeAvg'] = testScores.dfeAvg.apply(lambda row: float(sum(row))/len(row) if row else 0)
    testScores.loc[:,'dfeAvg'] = testScores.dfeAvg.apply(lambda row: row if row else 0)
    testScores = testScores.loc[:, tsCols + ['dfeAvg']]

    #exam
    cols = ["username", "examScore", "courseYear"]
    levelExam17 = [7]
    exam2017 = "../../../../../Desktop/Fall2018/PL/2017Data/lrExam/{}.csv"
    exam2017 = answers.parseAndGroup(levelExam17,[exam2017],"exam2017")
    exam2017.loc[:,'examScore'] = exam2017.apply(lambda row: answers.calcScore(row, "exam2017", 1, 8),axis=1)
    exam2017 = exam2017[cols]

    levelExam18 = [4]
    exam2018 = "../../../../../Desktop/Fall2018/PL/2018Data/lrExam/{}.csv"
    exam2018 = answers.parseAndGroup(levelExam18,[exam2018],"exam2018")
    exam2018.loc[:,'examScore'] = exam2018.apply(lambda row: answers.calcScore(row, "exam2018", 1, 5),axis=1)
    exam2018 = exam2018[cols]
    exam = pd.concat([exam2017, exam2018])

    personalizedAssignmentData = [climateChange, testScores, fluEpedemics]
    personalizedAssignmentNames = ["cc", "rts", "dfe"]
    personalizedMaxAttempts = [maxAttemptsCC, maxAttemptsRTS, maxAttemptsDFE]
    personalizedLevelQues = [levelQuesCC, levelQuesRTS, levelQuesDFE]

    assignmentData = [anyticalDetective,stockDynamics,demographicsEmployment]
    assignmentNames = ["ad","sd", "de"]
    maxAttempts = [maxAttemptsAD, maxAttemptsSD, maxAttemptsDE]
    levelQues = [levelQuesAD, levelQuesSD, levelQuesDE]
    #combine and clean data
    for i in range(len(personalizedAssignmentNames)):
        data = pd.merge(userList, personalizedAssignmentData[i], on=['username','courseYear'], how='inner')
        for j in range(len(assignmentNames)):
            data = pd.merge(data, assignmentData[j], on=['username','courseYear'], how='left')
        data = pd.merge(data, exam, on=['username','courseYear'], how='left')
        data = cleanData.fillMissingData(data, assignmentNames)
        for j in range(len(assignmentNames)):
            data = cleanData.fixAttemptsForIncorrect(data, levelQues[j], assignmentNames[j], maxAttempts[j])
        data = cleanData.fixAttemptsForIncorrect(data, personalizedLevelQues[i], personalizedAssignmentNames[i], personalizedMaxAttempts[i])
        cleanData.makeAgeCategorical(data)
        data.drop_duplicates()
        cluster.clusterStudents(data, assignmentNames)
        data.to_csv("data/preprocessed_{}.csv".format(personalizedAssignmentNames[i]))

assignmentNames = ["ad","sd", "de"]
personalizedAssignmentNames = ["cc","rts","dfe"]
dataSets = []
dataSets.append(pd.read_csv("data/preprocessed_cc.csv",header=0))
dataSets.append(pd.read_csv("data/preprocessed_rts.csv",header=0))
dataSets.append(pd.read_csv("data/preprocessed_dfe.csv",header=0))
levelQues = [levelQuesCC, levelQuesRTS, levelQuesDFE]
[6,4,5,6]
# current position in each level
positionsSet = [
             [[[2,1,1,1],[2,2,1,1],[2,3,1,1]],
             [[2,2,1,1],[2,3,1,1],[2,2,2,1],[2,3,2,1]],
             [[2,2,2,1],[2,3,2,1],[2,2,2,2],[2,3,2,2]],
             [[2,3,2,2]]],
             [[[5,1,1,1],[5,4,1,1],[4,4,1,1],[4,2,1,1]],
             [[5,4,1,1],[5,3,1,1],[4,4,3,1],[4,3,4,1]],
             [[5,4,5,1],[5,4,4,1],[5,3,5,4],[5,3,4,3]],
             [[5,4,5,4],[5,4,4,4],[5,4,5,3],[5,4,4,3]]],
             [[[6,1,1,1],[6,4,1,1],[5,4,1,1],[5,3,1,1]],
             [[6,4,1,1],[6,4,5,1],[6,3,4,1],[5,3,5,1]],
             [[6,4,5,1],[6,4,4,1],[6,4,5,6],[6,3,4,5]],
             [[6,4,5,6],[6,4,4,6],[6,4,5,5],[6,3,4,5]]]
            ]
weight = [[1.9,1,0,0,0],
          [2.5,1.4,1,0,0],
          [0,1.4,1,1,0],
          [0,0,1.4,1.2,1]]
i = 0
for data in dataSets:
    positions = positionsSet[i]
    data1 = cluster.setYGroup(data,2)
    dataLev = []
    for p in positions[0]:
        dataFilled = answers.sortToPredict(data1, assignmentLevels, assignmentNames, levelQues[i], 1, p, personalizedAssignmentNames[i], weight[0], 2)
        cleanData.convertTreatment(dataFilled)
        dataLev.append(dataFilled)
    final = pd.concat(dataLev)
    final.to_csv("data/{}1.csv".format(personalizedAssignmentNames[i]))
    data2 = cluster.setYGroup(data,3)
    for currlevel in range(2,4):
        dataLev = []
        for p in positions[currlevel - 1]:
            dataFilled = answers.sortToPredict(data2, assignmentLevels, assignmentNames, levelQues[i], currlevel, p, personalizedAssignmentNames[i], weight[currlevel - 1])
            cleanData.convertTreatment(dataFilled)
            dataLev.append(dataFilled)
        final = pd.concat(dataLev)
        final.to_csv("data/{}{}.csv".format(personalizedAssignmentNames[i], currlevel))
    # prediction of y3 is exam
    dataLev = []
    for p in positions[3]:
        dataFilled = answers.sortToPredict(data2, assignmentLevels, assignmentNames, levelQues[i], 4, p, personalizedAssignmentNames[i], weight[3])
        cleanData.convertTreatment(dataFilled)
        dataLev.append(dataFilled)
    final = pd.concat(dataLev)
    final.to_csv("data/{}4.csv".format(personalizedAssignmentNames[i]))
    i += 1

# j = julia.Julia()
# x = j.include("PrescriptiveTree.jl")
