ClimateChange = { 1:2, 2:3, 3:2, 4:2 }
ReadingTestScores = { 1:5, 2:4 ,3:5, 4:4 }
DetectingFluEpidemics = { 1:6, 2:4, 3:5, 4:6 }

def path(assignment, history, level):
    map={"ClimateChange": ClimateChange,
    "ReadingTestScores": ReadingTestScores,
    "DetectingFluEpedemics": DetectingFluEpidemics}

    defaultPath = map[assignment]
    if history["next{}".format(level)] < defaultPath[level]:
        return level - 1, history["next{}".format(level)] - 1
    else:
        while level < 4:
            if history["next{}".format(level+1)] < defaultPath[level+1]:
                return level, history["next{}".format(level+1)] - 1
            level += 1
    return
