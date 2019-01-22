import sys
import json
from enum import Enum
import helper

class State(Enum):
    question = 1
    questionInstance = 2
    user = 3

class QInst(Enum):
    status = 1
    answer = 2
    numAttempts = 3

class Status(Enum):
    correct = 1
    incorrectNoAttempts = 2
    incorrectWithAttempts = 3
    incomplete = 4
    newQuestion = 5

def main():
    state = json.loads(sys.argv[1])
    status = state[State.questionInstance.name][QInst.status.name]
    questionInstance = state[State.questionInstance.name]
    question = state[State.question.name]
    if status == Status.correct.name or status == Status.incorrectNoAttempts.name:
        state[State.question.name] = helper.getNextQuestion(question['Assignment'])
        state[State.questionInstance.name] = {
            QInst.status.name: Status.newQuestion.name,
            QInst.answer.name: [],
            QInst.numAttempts.name: 0
        }
    else:
        #user is in process of answering
        if questionInstance[QInst.answer.name]: #user answered question
            state[State.questionInstance.name][QInst.numAttempts.name] += 1
            if helper.isCorrect(questionInstance[QInst.answer.name],question['correctAnswer']):
                state[State.questionInstance.name][QInst.status.name] = Status.correct.name
            elif questionInstance[QInst.numAttempts.name] < question['attemptsOverall']:
                state[State.questionInstance.name][QInst.status.name] = Status.incorrectWithAttempts.name
            else:
                state[State.questionInstance.name][QInst.status.name] = Status.incorrectNoAttempts.name
        else:
            #user did not answer question
            state[State.questionInstance.name][QInst.status.name] = Status.incomplete.name
    return json.dumps(state)

if __name__ == "__main__":
    main()

#example of json parameter that script recieves
#'{"question":{"Assignment": "ClimateChange","level": 1,"number": 2,"text": "Which variables are significant in the model? We will consider a variable signficant only if the p-value is below 0.05. (Select all that apply.)","options": ["MEI", "CO2", "CH4", "N20", "CFC.11", "CFC.12", "TSI", "Aerosols"],"correctAnswer": ["MEI", "CO2", "CFC.11", "CFC.12", "TSI", "Aerosols"],"explanation": "If you look at the model we created in the previous problem using summary(climatelm), all of the variables have at least one star except for CH4 and N2O. So MEI, CO2, CFC.11, CFC.12, TSI, and Aerosols are all significant.","attemptsOverall": 2},"questionInstance":{"status": "newQuestion","answer": ["MEI", "CO2"],"numAttempts": 0},"user":{"username": "omer"}}'
