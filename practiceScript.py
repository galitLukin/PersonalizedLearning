import json
import random

# this script only returns a random json

def getNextQuestion():
    with open('LinearRegression.json') as f:
        questions = json.load(f)
        level = random.randint(0,3)
        q = random.randint(0,4)
        print(level,q)
        return questions['ClimateChange'][level]['questions'][q]

def main():
    response = {}
    response["question"] = getNextQuestion()
    response["questionInstance"] = { "status": "new question", "answer": [], "numAttempts": 0 },
    response["user"] = { "username": "galit" }
    print(json.dumps(response))
    return json.dumps(response)

if __name__ == "__main__":
    main()
