import sys
import json
from enum import Enum
import helper

class State(Enum):
    Question = 1
    QuestionInstance = 2
    User = 3


class QInst(Enum):
    status = 1
    answer = 2
    numAttempts = 3


class Status(Enum):
    Correct = 1
    IncorrectNoAttempts = 2
    IncorrectWithAttempts = 3
    Incomplete = 4
    NewQuestion = 5


def main():
    if len(sys.argv) < 2:
        state = {}
        # history = ???
        with open('./python/LinearRegression.json', encoding='utf-8') as f:
            questions = json.load(f)
        # TODO:  this should be based on where the user is in which assignment
        # so first call to script should be with scores row
        # all other calls to script should be with Json+scores row or just Json - ask omer about this
        state[State.Question.name] = questions['ClimateChange'][0]['questions'][0]
        # state[State.Question.name] = helper.getFirstQuestion(questions, history)
        state[State.QuestionInstance.name] = {
            QInst.status.name: Status.NewQuestion.name,
            QInst.answer.name: [],
            QInst.numAttempts.name: 0
        }
        #TODO: insert user_id that is received
        state['User'] = {
            'username': "omer"
        }
    else:
        state = json.loads(sys.argv[1])

        status = state[State.QuestionInstance.name][QInst.status.name]
        questionInstance = state[State.QuestionInstance.name]
        question = state[State.Question.name]

        if status == Status.Correct.name or status == Status.IncorrectNoAttempts.name:
            state[State.Question.name] = helper.getNextQuestion(
                question['Assignment'], question['level'], question['number'])
            state[State.QuestionInstance.name] = {
                QInst.status.name: Status.NewQuestion.name,
                QInst.answer.name: [],
                QInst.numAttempts.name: 0
            }
        else:
            # user is in process of answering
            if not questionInstance[QInst.answer.name][0]:
                # user did not answer question
                state[State.QuestionInstance.name][QInst.status.name] = Status.Incomplete.name
            else:
                # user answered question
                state[State.QuestionInstance.name][QInst.numAttempts.name] = questionInstance[QInst.numAttempts.name] + 1
                if helper.isCorrect(questionInstance[QInst.answer.name], question['correctAnswer']):
                    state[State.QuestionInstance.name][QInst.status.name] = Status.Correct.name
                elif state[State.QuestionInstance.name][QInst.numAttempts.name] < question['attemptsOverall']:
                    state[State.QuestionInstance.name][QInst.status.name] = Status.IncorrectWithAttempts.name
                else:
                    state[State.QuestionInstance.name][QInst.status.name] = Status.IncorrectNoAttempts.name

    print(json.dumps(state))


if __name__ == "__main__":
    main()

# example of json parameter that script recieves and returns
#'{"Question":{"Assignment": "Climate Change","level": 1,"number": 2,"text": "Which variables are significant in the model? We will consider a variable signficant only if the p-value is below 0.05. (Select all that apply.)","options": ["MEI", "CO2", "CH4", "N20", "CFC.11", "CFC.12", "TSI", "Aerosols"],"correctAnswer": ["MEI", "CO2", "CFC.11", "CFC.12", "TSI", "Aerosols"],"explanation": "If you look at the model we created in the previous problem using summary(climatelm), all of the variables have at least one star except for CH4 and N2O. So MEI, CO2, CFC.11, CFC.12, TSI, and Aerosols are all significant.","attemptsOverall": 2, "Weight": 0},"QuestionInstance":{"status": "Correct","answer": ["MEI", "CO2", "CFC.11", "CFC.12", "TSI", "Aerosols"],"numAttempts": 1},"User":{"username": "omer"}}'
