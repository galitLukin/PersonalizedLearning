import sys
import json

def main():
    pastQuestions = json.loads(sys.argv[1])
    with open('./python/LinearRegression.json', encoding='utf-8') as f:
	    questions = json.load(f)
    pastQuestionsData = []
    a = pastQuestions[0]['Question']['Assignment'].replace(" ","")
    for pq in pastQuestions:
        q = pq['Question']
        pastQuestionsData.append(questions[a][q['level']-1]['questions'][q['number']-1])
    print(json.dumps(pastQuestionsData))

if __name__ == "__main__":
    main()
