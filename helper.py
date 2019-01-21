import json
import random

def isCorrect(answer, correctAnswer):
	if len(answer) is not len(correctAnswer):
		return False
	elif set(answer)^set(correctAnswer):
		return False
	else:
		return True

def getNextQuestion():
    with open('LinearRegression.json') as f:
        questions = json.load(f)
    level = random.randint(0,3)
    q = random.randint(0,4)
    return questions['ClimateChange'][level]['questions'][q]
    #first get history of user
	# j = julia.Julia()
	# pred = j.include("prediction.jl")
    # treatment = pred.func(state, history)
    # if treatment == "A":
    #     #return correct json
    #     return
    # elif treatment == "B":
    #     #return correct json
    #     return
    # else:
    #     #return correct json
    #     return
