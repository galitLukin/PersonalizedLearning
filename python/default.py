ClimateChange = { 1:2, 2:3, 3:2, 4:2 }
ReadingTestScores = { 1:5, 2:4 ,3:5, 4:4 }
DetectingFluEpidemics = { 1:6, 2:4, 3:5, 4:6 }

def path(assignment, history, level):
    defaultPath = None
    if assignment == "ClimateChange":
        defaultPath = ClimateChange
    elif assignment == "ReadingTestScores":
        defaultPath = ReadingTestScores
    else:
        defaultPath = DetectingFluEpidemics

    if history["next{}".format(level)] < defaultPath[level]:
        return level - 1, history["next{}".format(level)] - 1
    else:
        while level < 4:
            if history["next{}".format(level+1)] < defaultPath[level+1]:
                return level, history["next{}".format(level+1)] - 1
            level += 1
    return -1,0


def basicPath(assignment, level, number):
    defaultPath = None
    if assignment == "ClimateChange":
        defaultPath = ClimateChange
    elif assignment == "ReadingTestScores":
        defaultPath = ReadingTestScores
    else:
        defaultPath = DetectingFluEpidemics

    if number < defaultPath[level]:
        return level - 1, number
    else:
        if level == 4:
            return
        else:
            return level, 0


def prequisiteSatisfied(assignment, level, number):
    mandatory = None
    if assignment == "ClimateChange":
        mandatory = { 1:1, 2:2, 3:2, 4:2 }
    elif assignment == "ReadingTestScores":
        mandatory = { 1:4, 2:2 ,3:1, 4:1 }
    else:
        mandatory = { 1:1, 2:1, 3:1, 4:1 }

    if number < mandatory[level]:
        return False
    return True
