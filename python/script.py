import sys
import json
from enum import Enum
import helper

class State(Enum):
    Question = 1
    QuestionInstance = 2
    User = 3
    Score = 4
    PrevLocation = 5

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
    Done = 6


def main():
    state = json.loads(sys.argv[1])
    if state[State.PrevLocation.name]['IsFirst']:
        #state[State.Question.name] = questions['ClimateChange'][0]['questions'][0]
        state[State.Question.name], attempts = \
        helper.getFirstQuestion(state[State.Score.name],state[State.PrevLocation.name])
        state[State.QuestionInstance.name] = {
            QInst.status.name: Status.NewQuestion.name,
            QInst.answer.name: [],
            QInst.numAttempts.name: attempts
        }
    else:
        status = state[State.QuestionInstance.name][QInst.status.name]
        questionInstance = state[State.QuestionInstance.name]
        question = state[State.Question.name]

        if status == Status.Correct.name or status == Status.IncorrectNoAttempts.name:
            state[State.Question.name] = helper.getNextQuestion(
                question['Assignment'], question['level'], question['number'], state[State.Score.name])
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
#'{"Question":{"Assignment": "Climate Change","level": 1,"number": 2,"text": "Which variables are significant in the model? We will consider a variable signficant only if the p-value is below 0.05. (Select all that apply.)","options": ["MEI", "CO2", "CH4", "N20", "CFC.11", "CFC.12", "TSI", "Aerosols"],"correctAnswer": ["MEI", "CO2", "CFC.11", "CFC.12", "TSI", "Aerosols"],"explanation": "If you look at the model we created in the previous problem using summary(climatelm), all of the variables have at least one star except for CH4 and N2O. So MEI, CO2, CFC.11, CFC.12, TSI, and Aerosols are all significant.","attemptsOverall": 2, "Weight": 0},"QuestionInstance":{"status": "Correct","answer": ["MEI", "CO2", "CFC.11", "CFC.12", "TSI", "Aerosols"],"numAttempts": 1},"User":{"username": "omer"},"Score": {"Username": "6987787dd79cf0aecabdca8ddae95b4a1", "Assignment": "Climate Change", "Gender": "None", "Level_of_education": "None", "Enrollment_mode": "audit", "AgeCategory": "Null", "Ad1": 0, "Ad2": 0, "Ad3": 0, "Ad4": 0, "Sd1": 0, "Sd2": 0, "Sd3": 0, "Sd4": 0, "De1": 0, "De2": 0, "De3": 0, "De4": 0, "Cc1": 0, "Cc2": 0, "Cc3": 0, "Cc4": 0, "Rts1": 0, "Rts2": 0, "Rts3": 0, "Rts4": 0, "Score1_correct": 5, "Score1_attempts": 11, "Score2_correct": 2, "Score2_attempts": 3, "Score3_correct": 0, "Score3_attempts": 0, "Score4_correct": 0, "Score4_attempts": 0, "Next1": 1, "Next2": 1, "Next3": 1, "Next4": 1}, "IsFirst": "false"}'
