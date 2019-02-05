import json
import networkx as nx
import re

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
	#TODO: deal with if feature not in feature space - return rand or smthing smart
	text = root_exp.split()
	feature = text[0]
	if feature == "Prescribe":
		return text[1]
	else:
		val = history[feature]
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

def getNextQuestion(assignment, level, number):
	map={"Climate Change": "cc", "Reading Test Scores": "rts", "Detecting Flu Epedemics": "dfe"}
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

	features = ["gender", "level_of_education", "enrollment_mode", "ageCategory", \
	"ad1", "ad2", "ad3", "ad4", "sd1", "sd2", "sd3", "sd4", "de1", "de2", "de3", "de4",\
	"cc1", "cc2", "cc3", "cc4", "rts1", "rts2", "rts3", "rts4",\
	"score1_correct", "score2_correct", "score3_correct", "score4_correct", \
	"score1_attempts", "score2_attempts", "score3_attempts", "score4_attempts",\
  	"next1", "next2", "next3", "next4"]
	history = dict(zip(features,historydb))

	for l in range(1,5):
		if history['score{}_attempts'.format(l)] > 0:
			history['score{}'.format(l)] = float(history['score{}_correct'.format(l)])/history['score{}_attempts'.format(l)]
		else:
			history['score{}'.format(l)] = 0
		del history['score{}_correct'.format(l)]
		del history['score{}_attempts'.format(l)]

	lnr = nx.nx_pydot.read_dot('./python/model/{}/{}/pytree.dot'.format(level,asmt))
	treatment = '1'

	#TODO deal with possible infinite loop
	while treatment not in ["A","B","C"]:
		treatment = getNextNode(history,lnr,treatment)

	assignment = assignment.replace(" ", "")
	if treatment == "A":
		q = history["next{}".format(level - 1)] - 1
		return questions[assignment][level - 2]['questions'][q]
	elif treatment == "B":
		q = history["next{}".format(level)] - 1
		return questions[assignment][level - 1]['questions'][q]
	else:
		if level == 4:
			#TODO: deal with finishing the assignment
			q = history["next{}".format(level)]
			return questions[assignment][level - 1]['questions'][q] - 1
		else:
			q = history["next{}".format(level + 1)] - 1
			return questions[assignment][level]['questions'][q]
