import json
import networkx as nx
import re
import default

def isCorrect(answer, correctAnswer):
	if len(answer) is not len(correctAnswer):
		return False
	elif set(answer)^set(correctAnswer):
		return False
	else:
		return True

def getNextNode(history,lnr,curr):

	root_exp = nx.get_node_attributes(lnr,'label')[curr]
	root_exp = re.search('"(.*)"',root_exp).group(1)
	text = root_exp.split()
	feature = text[0]
	if feature == "Prescribe":
		return text[1]
	else:
		try:
			val = history[feature]
		except:
			return
		comp = text[1]
		threshold = text[2]
		direction = "right"
		if comp == "in":
			if val in threshold:
				direction = "left"
		else:
			if val < float(threshold):
				direction = "left"
		if direction == "left":
			children = list(map(int, list(lnr.successors(curr))))
			return str(min(children))
		else:
			children = list(map(int, list(lnr.successors(curr))))
			return str(max(children))
	return

def getNextQuestion(assignment, level, number):
	assignment = assignment.replace(" ", "")
	map={"ClimateChange": "cc", "ReadingTestScores": "rts", "DetectingFluEpedemics": "dfe"}
	mapQues={"ClimateChange":5,"ReadingTestScores":6,"DetectingFluEpedemics":7}
	asmt = map[assignment];
	with open('./python/LinearRegression.json', encoding='utf-8') as f:
	    questions = json.load(f)

	#should recieve this from omer
	historydb = ["None", "None", "audit", "Null", \
	1, 1, 1, 1, \
	1, 1, 1, 1, \
	1, 1, 0.8888888888888888, 1,\
	0, 0, 0, 0, \
	0, 0, 0, 0, \
	1, 0, 0, 0, \
	1, 0, 0, 0, \
	2, 1, 1, 1]

	historydb = []

	features = ["gender", "level_of_education", "enrollment_mode", "ageCategory", \
	"ad1", "ad2", "ad3", "ad4", "sd1", "sd2", "sd3", "sd4", "de1", "de2", "de3", "de4",\
	"cc1", "cc2", "cc3", "cc4", "rts1", "rts2", "rts3", "rts4",\
	"score1_correct", "score2_correct", "score3_correct", "score4_correct", \
	"score1_attempts", "score2_attempts", "score3_attempts", "score4_attempts",\
  	"next1", "next2", "next3", "next4"]
	try:
		history = dict(zip(features,historydb))
		for l in range(1,5):
			if history['score{}_attempts'.format(l)] > 0:
				history['score{}'.format(l)] = float(history['score{}_correct'.format(l)])/history['score{}_attempts'.format(l)]
			else:
				history['score{}'.format(l)] = 0
			#del history['score{}_correct'.format(l)]
			#del history['score{}_attempts'.format(l)]
	except:
		level,q = default.path(assignment,history,level)
		if level is not None and q is not None:
			return questions[assignment][level]['questions'][q]
		return


	lnr = nx.nx_pydot.read_dot('./python/model/{}/{}/pytree.dot'.format(level,asmt))
	treatment = '1'

	infLoop = 0
	while treatment not in ["A","B","C"]:
		try:
			treatment = getNextNode(history,lnr,treatment)
			infLoop += 1
			if infLoop > 10 or not treatment:
				level,q = default.path(assignment,history,level)
				if level is not None and q is not None:
					return questions[assignment][level]['questions'][q]
				return
		except:
			level,q = default.path(assignment,history,level)
			if level is not None and q is not None:
				return questions[assignment][level]['questions'][q]
			return

	lastQues = mapQues[assignment]
	prevLevelFull = False

	if treatment == "A":
		q = history["next{}".format(level - 1)] - 1
		if q < lastQues:
			return questions[assignment][level - 2]['questions'][q]
		prevLevelFull = True
	if treatment == "B" or prevLevelFull:
		q = history["next{}".format(level)] - 1
		if q < lastQues:
			return questions[assignment][level - 1]['questions'][q]
		prevLevelFull = True
	if treatment == "C" or prevLevelFull:
		while level < 4:
			q = history["next{}".format(level + 1)] - 1
			if q < lastQues:
				return questions[assignment][level]['questions'][q]
			level += 1
		return
