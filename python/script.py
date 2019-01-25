import sys
import json
from enum import Enum
import helper
import pprint

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
    pp = pprint.PrettyPrinter(indent=4)
    state = json.loads(sys.argv[1])
    pp.pprint(state)
    # This crashes here -> check how you are accessing the object with enums
    status = state[State.QuestionInstance.name][QInst.status.name]
    questionInstance = state[State.QuestionInstance.name]
    question = state[State.Question.name]
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

# example of json parameter that script recieves and returns
