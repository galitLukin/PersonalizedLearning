import json
import random
import julia

def isCorrect(answer, correctAnswer):
	if len(answer) is not len(correctAnswer):
		return False
	elif set(answer)^set(correctAnswer):
		return False
	else:
		return True

def getNextQuestion(assignment, level, number):
	with open('./python/LinearRegression.json') as f:
	    questions = json.load(f)

	#should recieve this from omer
	history = None

	j = julia.Julia()
	pred = j.include("./python/prediction.jl")
	treatment = j.func(assignment, level, history)
	assignment = assignment.replace(" ", "")
	if treatment == "A":
        #TODO: return correct json based on where the user is
		q = random.randint(0,4)
		return questions[assignment][level - 2]['questions'][q]
	elif treatment == "B":
		#TODO:return correct json based on where the user is
		q = random.randint(0,4)
		return questions[assignment][level - 1]['questions'][q]
	else:
		#TODO:return correct json based on where the user is
		#TODO: deal with finishing the assignment
		if level == 4:
			q = random.randint(0,4)
			return questions[assignment][level - 1]['questions'][q]
		else:
			q = random.randint(0,4)
			return questions[assignment][level]['questions'][q]
