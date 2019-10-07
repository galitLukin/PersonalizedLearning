import json
import networkx as nx
import re
import default

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

def getNextQuestion(assignment, level, number, score, status):
	assignment = assignment.replace(" ", "")
	map={"Asmt1": "cc", "Asmt2": "rts", "Asmt3": "dfe"}
	mapQues={"Asmt1":5,"Asmt2":6,"Asmt3":6}
	asmt = map[assignment];
	lowerIt = lambda s: s[:1].lower() + s[1:] if s else ''
	for key in score:
		score[lowerIt(key)] = score.pop(key)

	with open('./python/LinearRegression.json', encoding='utf-8') as f:
	    questions = json.load(f)
	try:
		for l in range(1,5):
			if score['score{}_attempts'.format(l)] > 0:
				score['score{}'.format(l)] = float(score['score{}_correct'.format(l)])/score['score{}_attempts'.format(l)]
			else:
				score['score{}'.format(l)] = 0
	except:
		level,q = default.basicPath(assignment, level, number)
		if level is not None and q is not None:
			return questions[assignment][level]['questions'][q]
		return

	lnr = nx.nx_pydot.read_dot('./python/model/{}/{}/pytree.dot'.format(level,asmt))
	treatment = '1'

	infLoop = 0
	while treatment not in ["A","B","C"]:
		try:
			treatment = getNextNode(score,lnr,treatment)
			infLoop += 1
			if infLoop > 10 or not treatment:
				level,q = default.path(assignment,score,level)
				if level is not -1:
					return questions[assignment][level]['questions'][q]
				return
		except:
			level,q = default.path(assignment,score,level)
			if level is not -1:
				return questions[assignment][level]['questions'][q]
			return

	lastQues = mapQues[assignment]
	prevLevelFull = False

	if treatment == "A":
		q = score["next{}".format(level - 1)] - 1
		if q < lastQues:
			return questions[assignment][level - 2]['questions'][q]
		prevLevelFull = True
	if treatment == "B" or prevLevelFull:
		q = score["next{}".format(level)] - 1
		if q < lastQues:
			return questions[assignment][level - 1]['questions'][q]
		prevLevelFull = True
	if treatment == "C" or prevLevelFull:
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

def getFirstQuestion(score, location):
	assignment = score['Assignment'].replace(" ", "")
	level = location['Level']
	numb = location['Number']
	with open('./python/LinearRegression.json', encoding='utf-8') as f:
		questions = json.load(f)
	if level == 0:
		return questions[assignment][0]['questions'][0], 0
	attemptsOverall = questions[assignment][level - 1]['questions'][numb - 1]['attemptsOverall']
	if location['Correctness'] == 1 or location['Attempt'] >= attemptsOverall:
		return getNextQuestion(assignment, level, numb, score), 0
	return questions[assignment][level - 1]['questions'][numb - 1], location['Attempt']
