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
        state[State.Question.name], attempts = \
        helper.getFirstQuestion(state[State.Score.name],state[State.PrevLocation.name])
        state[State.QuestionInstance.name] = {
            QInst.status.name: Status.NewQuestion.name if state[State.Question.name] else Status.Done.name,
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
                QInst.status.name: Status.NewQuestion.name if state[State.Question.name] else Status.Done.name,
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
                if helper.isCorrect(questionInstance[QInst.answer.name][0], question['correctAnswer'], question['answerType']):
                    state[State.QuestionInstance.name][QInst.status.name] = Status.Correct.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []
                elif state[State.QuestionInstance.name][QInst.numAttempts.name] < question['attemptsOverall']:
                    state[State.QuestionInstance.name][QInst.status.name] = Status.IncorrectWithAttempts.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []
                else:
                    state[State.QuestionInstance.name][QInst.status.name] = Status.IncorrectNoAttempts.name
                    state[State.QuestionInstance.name][QInst.answer.name] = []

    print(json.dumps(state))


if __name__ == "__main__":
    main()

# example of json parameter that script recieves and returns
#'{"Question":{"Assignment": "Asmt1","level": 1,"number": 2,"text": "Which variables are significant in the model? We will consider a variable signficant only if the p-value is below 0.05. (Select all that apply.)","options": ["MEI", "CO2", "CH4", "N20", "CFC.11", "CFC.12", "TSI", "Aerosols"],"correctAnswer": ["MEI,CO2,CFC.11,CFC.12,TSI,Aerosols"],"explanation": "If you look at the model we created in the previous problem using summary(climatelm), all of the variables have at least one star except for CH4 and N2O. So MEI, CO2, CFC.11, CFC.12, TSI, and Aerosols are all significant.","attemptsOverall": 2, "Weight": 0,"AnswerType":2},"QuestionInstance":{"status": "NewQuestion","answer": "CO2,MEI,CFC.11,CFC.12,TSI,Aerosols","numAttempts": 1},"User":{"Uid": "omer", AssignmentName: "Asmt1"},"Score": {"Username": "6987787dd79cf0aecabdca8ddae95b4a1", "Assignment": "Asmt1", "Gender": "None", "Level_of_education": "None", "Enrollment_mode": "audit", "AgeCategory": "Null", "Ad1": 0, "Ad2": 0, "Ad3": 0, "Ad4": 0, "Sd1": 0, "Sd2": 0, "Sd3": 0, "Sd4": 0, "De1": 0, "De2": 0, "De3": 0, "De4": 0, "Cc1": 0, "Cc2": 0, "Cc3": 0, "Cc4": 0, "Rts1": 0, "Rts2": 0, "Rts3": 0, "Rts4": 0, "Score1_correct": 5, "Score1_attempts": 11, "Score2_correct": 2, "Score2_attempts": 3, "Score3_correct": 0, "Score3_attempts": 0, "Score4_correct": 0, "Score4_attempts": 0, "Next1": 1, "Next2": 1, "Next3": 1, "Next4": 1, "Grade": 0.0}, "PrevLocation": {"IsFirst" : false, "Level" : 1, "Number" : 2, "Attempt" : 0, "Correctness" : 0, "Timestamp" : ""}}'



#'{"Question":{"Assignment":"","level":0,"number":0,"text":"","options":null,"correctAnswer":"","explanation":"","attemptsOverall":0,"weight":0,"answerType":0},"QuestionInstance":{"status":"","answer":null,"numAttempts":0},"User":{"Uid":"2","AssignmentName":"Asmt1"},"Score":{"Username":"2","Assignment":"Asmt1","Gender":"None","Level_of_education":"None","Enrollment_mode":"audit","AgeCategory":"Null","Ad1":0,"Ad2":0,"Ad3":0,"Ad4":0,"Sd1":0,"Sd2":0,"Sd3":0,"Sd4":0,"De1":0,"De2":0,"De3":0,"De4":0,"Cc1":0,"Cc2":0,"Cc3":0,"Cc4":0,"Rts1":0,"Rts2":0,"Rts3":0,"Rts4":0,"Score1_correct":0,"Score1_attempts":4,"Score2_correct":0,"Score2_attempts":4,"Score3_correct":0,"Score3_attempts":2,"Score4_correct":0,"Score4_attempts":2,"Next1":3,"Next2":4,"Next3":2,"Next4":2,"Grade":0},"PrevLocation":{"IsFirst":true,"Level":1,"Number":1,"Attempt":1,"Correctness":0,"Timestamp":""}}'
