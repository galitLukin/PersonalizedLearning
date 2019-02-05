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
	history = {"Gender" : "None", "level_of_education" : "None", "enrollment_mode" : "audit",\
	"ageCategory" : "Null", "ad1" : 1, "ad2" : 1, "ad3" : 1, "ad4" : 1,\
	"sd1" : 1, "sd2" : 1, "sd3" : 1, "sd4" : 1, \
	"de1" : 1, "de2" : 1, "de3" : 0.8888888888888888, "de4" : 1,\
	"score1" : 1, "score2" : 0, "score3" : 0, "score4" : 0,\
	"next1": 2, "next2": 1, "next3": 1, "next4": 1}

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
