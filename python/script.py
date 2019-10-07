import sys
import json
import datetime as dt
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
    startTime = 4
    endTime = 5
    duration = 6

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
        state[State.Question.name], attempts = \
        helper.getFirstQuestion(state[State.Score.name],state[State.PrevLocation.name])
        state[State.QuestionInstance.name] = {
            QInst.status.name: Status.NewQuestion.name if state[State.Question.name] else Status.Done.name,
            QInst.answer.name: [],
            QInst.numAttempts.name: attempts,
            QInst.startTime.name: QInst.startTime.name: state[State.QuestionInstance.name][QInst.startTime.name] if state[State.QuestionInstance.name][QInst.startTime.name] else dt.datetime.now()
        }
    else:
        status = state[State.QuestionInstance.name][QInst.status.name]
        questionInstance = state[State.QuestionInstance.name]
        question = state[State.Question.name]

        if status == Status.Correct.name or status == Status.IncorrectNoAttempts.name:
            state[State.Question.name] = helper.getNextQuestion(
                question['Assignment'], question['level'], question['number'], state[State.Score.name], status)
            state[State.QuestionInstance.name] = {
                QInst.status.name: Status.NewQuestion.name if state[State.Question.name] else Status.Done.name,
                QInst.answer.name: [],
                QInst.numAttempts.name: 0
                QInst.startTime.name: dt.datetime.now()
            }
        else:
            # user is in process of answering
            if not questionInstance[QInst.answer.name][0]:
                # user did not answer question
                state[State.QuestionInstance.name][QInst.status.name] = Status.Incomplete.name
            else:
                # user answered question
                state[State.QuestionInstance.name][QInst.numAttempts.name] = questionInstance[QInst.numAttempts.name] + 1
                if helper.isCorrect(questionInstance[QInst.answer.name][0], question['correctAnswer'], question['answerType']):
                    state[State.QuestionInstance.name][QInst.status.name] = Status.Correct.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []
                    state[State.QuestionInstance.name][QInst.endTime.name] = dt.datetime.now()
                    d = state[State.QuestionInstance.name][QInst.endTime.name] - state[State.QuestionInstance.name][QInst.startTime.name]
                    state[State.QuestionInstance.name][QInst.duration.name] = d.total_seconds() if d.total_seconds() <= 1800 else 1800
                elif state[State.QuestionInstance.name][QInst.numAttempts.name] < question['attemptsOverall']:
                    state[State.QuestionInstance.name][QInst.status.name] = Status.IncorrectWithAttempts.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []
                else:
                    state[State.QuestionInstance.name][QInst.status.name] = Status.IncorrectNoAttempts.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []
                    state[State.QuestionInstance.name][QInst.endTime.name] = dt.datetime.now()
                    d = state[State.QuestionInstance.name][QInst.endTime.name] - state[State.QuestionInstance.name][QInst.startTime.name]
                    state[State.QuestionInstance.name][QInst.duration.name] = d.total_seconds() if d.total_seconds() <= 1800 else 1800

    print(json.dumps(state))


if __name__ == "__main__":
    main()

# example of json parameter that script recieves and returns
#'{"Question":{"Assignment": "Asmt1","level": 1,"number": 2,"text": "Which variables are significant in the model? We will consider a variable signficant only if the p-value is below 0.05. (Select all that apply.)","options": ["MEI", "CO2", "CH4", "N20", "CFC.11", "CFC.12", "TSI", "Aerosols"],"correctAnswer": ["MEI,CO2,CFC.11,CFC.12,TSI,Aerosols"],"explanation": "If you look at the model we created in the previous problem using summary(climatelm), all of the variables have at least one star except for CH4 and N2O. So MEI, CO2, CFC.11, CFC.12, TSI, and Aerosols are all significant.","attemptsOverall": 2, "Weight": 0,"AnswerType":2},"QuestionInstance":{"status": "NewQuestion","answer": "CO2,MEI,CFC.11,CFC.12,TSI,Aerosols","numAttempts": 1},"User":{"Uid": "omer", AssignmentName: "Asmt1"},"Score": {"Username": "6987787dd79cf0aecabdca8ddae95b4a1", "Assignment": "Asmt1", "Gender": "None", "Level_of_education": "None", "Enrollment_mode": "audit", "AgeCategory": "Null", "Next1": 1, "Next2": 1, "Next3": 1, "Next4": 1, "Grade": 0.0}, "PrevLocation": {"IsFirst" : false, "Level" : 1, "Number" : 2, "Attempt" : 0, "Correctness" : 0, "Timestamp" : ""}}'



#'{"Question":{"Assignment":"","level":0,"number":0,"text":"","options":null,"correctAnswer":"","explanation":"","attemptsOverall":0,"weight":0,"answerType":0},"QuestionInstance":{"status":"","answer":null,"numAttempts":0},"User":{"Uid":"2","AssignmentName":"Asmt1"},"Score":{"Username":"2","Assignment":"Asmt1","Gender":"None","Level_of_education":"None","Enrollment_mode":"audit","AgeCategory":"Null","Next1":3,"Next2":4,"Next3":2,"Next4":2,"Grade":0},"PrevLocation":{"IsFirst":true,"Level":1,"Number":1,"Attempt":1,"Correctness":0,"Timestamp":""}}'
