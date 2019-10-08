import json
import re
import default
import numpy as np
import datetime

def isCorrect(answer, correctAnswer, answerType):
	answer = answer.replace(' ','').split(",")
	correctAnswer = correctAnswer.replace(' ','').split(",")
	if answerType > 0 :
		if set(answer) == set(correctAnswer):
			return True
		else:
			return False
	else:
		if set(answer) == set(correctAnswer):
			return True
		for ca in correctAnswer:
			try:
				caf = float(ca)
				notInAnswer = True
				for a in answer:
					try:
						af = float(a)
						if abs(caf-af) <= 0.01*abs(caf):
							notInAnswer = False
					except:
						continue
				if notInAnswer:
					return False
			except:
				if ca not in answer:
					return False
		return True

def prepareFeatures(score, qi):
	featMap = {
		"attempt": qi["numAttempts"],
		"correctness": 1 if qi["status"] == "Correct" else 0,
		"t": qi["duration"],
		"enrollment_mode": score["enrollment_mode"]
	}
	return featMap


def standardizeFeat(feat, featMap, meu, sigma):
	return (featMap[feat] - meu)/sigma


def decide(userFeat, centroids, decisions):
	userFeat = np.array(userFeat)
	minDist = np.inf
	for i,clst in enumerate(centroids):
		dist = np.linalg.norm(userFeat-np.array(clst))
		if dist < minDist:
			minDist = dist
			bestCluster = i
	return decisions[bestCluster]


def getNextQuestion(assignment, level, number, score, qi):
	assignment = assignment.replace(" ", "")
	map={"Asmt1": "cc", "Asmt2": "rts", "Asmt3": "dfe"}
	mapQues={"Asmt1":5,"Asmt2":6,"Asmt3":6}
	asmt = map[assignment];
	lowerIt = lambda s: s[:1].lower() + s[1:] if s else ''
	for key in score:
		score[lowerIt(key)] = score.pop(key)

	with open('newLR.json', encoding='utf-8') as f:
	    questions = json.load(f)

	with open('clustering.json', encoding='utf-8') as fc:
	    clusteringData = json.load(fc)

	qid = questions[assignment][level - 1]['questions'][number - 1]["qid"]
	keyQid = str(qid)
	qClusteringData = clusteringData[str(qid)]
	featMap = prepareFeatures(score, qi)
	userFeat = []
	for feat in qClusteringData["features"]:
		userFeat.append(standardizeFeat(feat, featMap, qClusteringData["meu"][feat], qClusteringData["sigma"][feat]))
	treatment = decide(userFeat, clusteringData[keyQid]["centroids"], clusteringData[keyQid]["decisions"])

	lastQues = mapQues[assignment]
	prevLevelFull = False

	if treatment == 1:
		q = score["next{}".format(level - 1)] - 1
		if q < lastQues:
			return questions[assignment][level - 2]['questions'][q]
		prevLevelFull = True
	if treatment == 2 or prevLevelFull:
		q = score["next{}".format(level)] - 1
		if q < lastQues:
			return questions[assignment][level - 1]['questions'][q]
		prevLevelFull = True
	if treatment == 3 or prevLevelFull:
		if default.prequisiteSatisfied(assignment, level, number):
			while level < 4:
				q = score["next{}".format(level + 1)] - 1
				if q < lastQues:
					return questions[assignment][level]['questions'][q]
				level += 1
			return
		#finish prequisites on this level
		else:
			return questions[assignment][level-1]['questions'][number]

def getFirstQuestion(score, location, qi):
	assignment = score['Assignment'].replace(" ", "")
	level = location['Level']
	numb = location['Number']
	with open('./python/newLR.json', encoding='utf-8') as f:
		questions = json.load(f)
	if level == 0:
		return questions[assignment][0]['questions'][0], 0, datetime.datetime.utcnow().strftime("%m/%d/%Y %H:%M:%S")
	attemptsOverall = questions[assignment][level - 1]['questions'][numb - 1]['attemptsOverall']
	if location['Correctness'] == 1 or location['Attempt'] >= attemptsOverall:
		return getNextQuestion(assignment, level, numb, score, qi), 0, datetime.datetime.utcnow().strftime("%m/%d/%Y %H:%M:%S")
	return questions[assignment][level - 1]['questions'][numb - 1], location['Attempt'], qi["startTime"]
