import sys
import json
import datetime
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
    s = '{"Question":{"Assignment": "Asmt1","level": 1,"number": 2,"text": "Which variables are significant in the model? We will consider a variable signficant only if the p-value is below 0.05. (Select all that apply.)","options": ["MEI", "CO2", "CH4", "N20", "CFC.11", "CFC.12", "TSI", "Aerosols"],"correctAnswer": "MEI,CO2,CFC.11,CFC.12,TSI,Aerosols","explanation": "If you look at the model we created in the previous problem using summary(climatelm), all of the variables have at least one star except for CH4 and N2O. So MEI, CO2, CFC.11, CFC.12, TSI, and Aerosols are all significant.","attemptsOverall": 2, "weight": 0,"answerType":2,"qid": 10102},"QuestionInstance":{"status": "IncorrectNoAttempts","answer": "CO2,MEI,CFC.11,CFC.12,TSI,Aerosols","numAttempts": 2,"startTime": "10/08/2019 10:49:01","endTime": "10/08/2019 10:49:11","duration": 60000},"User":{"Uid": "omer", "AssignmentName": "Asmt1"},"Score": {"Username": "6987787dd79cf0aecabdca8ddae95b4a1", "Assignment": "Asmt1", "Gender": "None", "Level_of_education": "None", "Enrollment_mode": "audit", "AgeCategory": "Null", "Next1": 1, "Next2": 1, "Next3": 1, "Next4": 1, "Grade": 0.0}, "PrevLocation": {"IsFirst" : false, "Level" : 1, "Number" : 2, "Attempt" : 0, "Correctness" : 0, "Timestamp" : ""}}'
    state = json.loads(s)
    #state = json.loads(sys.argv[1])
    if state[State.PrevLocation.name]['IsFirst']:
        state[State.Question.name], attempts, startTime = \
        helper.getFirstQuestion(state[State.Score.name],state[State.PrevLocation.name],state[State.QuestionInstance.name])
        state[State.QuestionInstance.name] = {
            QInst.status.name: Status.NewQuestion.name if state[State.Question.name] else Status.Done.name,
            QInst.answer.name: [],
            QInst.numAttempts.name: attempts,
            QInst.startTime.name: startTime,
            QInst.endTime.name: None,
            QInst.duration.name: None
        }
    else:
        status = state[State.QuestionInstance.name][QInst.status.name]
        questionInstance = state[State.QuestionInstance.name]
        question = state[State.Question.name]

        if status == Status.Correct.name or status == Status.IncorrectNoAttempts.name:
            state[State.Question.name] = helper.getNextQuestion(
                question['Assignment'], question['level'], question['number'], state[State.Score.name], state[State.QuestionInstance.name])
            state[State.QuestionInstance.name] = {
                QInst.status.name: Status.NewQuestion.name if state[State.Question.name] else Status.Done.name,
                QInst.answer.name: [],
                QInst.numAttempts.name: 0,
                QInst.startTime.name: datetime.datetime.utcnow().strftime("%m/%d/%Y %H:%M:%S"),
                QInst.endTime.name: None,
                QInst.duration.name: None
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
                    now = datetime.datetime.utcnow()
                    state[State.QuestionInstance.name][QInst.status.name] = Status.Correct.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []
                    state[State.QuestionInstance.name][QInst.endTime.name] = now.strftime("%m/%d/%Y %H:%M:%S")
                    d = now - datetime.datetime.strptime(state[State.QuestionInstance.name][QInst.startTime.name], "%m/%d/%Y %H:%M:%S")
                    state[State.QuestionInstance.name][QInst.duration.name] = d.total_seconds() if d.total_seconds() <= 1800 else 1800
                elif state[State.QuestionInstance.name][QInst.numAttempts.name] < question['attemptsOverall']:
                    state[State.QuestionInstance.name][QInst.status.name] = Status.IncorrectWithAttempts.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []
                else:
                    now = datetime.datetime.utcnow()
                    state[State.QuestionInstance.name][QInst.status.name] = Status.IncorrectNoAttempts.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []
                    state[State.QuestionInstance.name][QInst.endTime.name] = now.strftime("%m/%d/%Y %H:%M:%S")
                    d = now - datetime.datetime.strptime(state[State.QuestionInstance.name][QInst.startTime.name], "%m/%d/%Y %H:%M:%S")
                    state[State.QuestionInstance.name][QInst.duration.name] = d.total_seconds() if d.total_seconds() <= 1800 else 1800

    print(json.dumps(state))


if __name__ == "__main__":
    main()

# example of json parameter that script recieves and returns
#'{"Question":{"Assignment": "Asmt1","level": 1,"number": 2,"text": "Which variables are significant in the model? We will consider a variable signficant only if the p-value is below 0.05. (Select all that apply.)","options": ["MEI", "CO2", "CH4", "N20", "CFC.11", "CFC.12", "TSI", "Aerosols"],"correctAnswer": ["MEI,CO2,CFC.11,CFC.12,TSI,Aerosols"],"explanation": "If you look at the model we created in the previous problem using summary(climatelm), all of the variables have at least one star except for CH4 and N2O. So MEI, CO2, CFC.11, CFC.12, TSI, and Aerosols are all significant.","attemptsOverall": 2, "Weight": 0,"AnswerType":2,"qid": 10102},"QuestionInstance":{"status": "NewQuestion","answer": "CO2,MEI,CFC.11,CFC.12,TSI,Aerosols","numAttempts": 1,"startTime": "10/08/2019 10:49:01","endTime": "10/08/2019 10:49:11","duration": 10.0},"User":{"Uid": "omer", "AssignmentName": "Asmt1"},"Score": {"Username": "6987787dd79cf0aecabdca8ddae95b4a1", "Assignment": "Asmt1", "Gender": "None", "Level_of_education": "None", "Enrollment_mode": "audit", "AgeCategory": "Null", "Next1": 1, "Next2": 1, "Next3": 1, "Next4": 1, "Grade": 0.0}, "PrevLocation": {"IsFirst" : false, "Level" : 1, "Number" : 2, "Attempt" : 0, "Correctness" : 0, "Timestamp" : ""}}'
