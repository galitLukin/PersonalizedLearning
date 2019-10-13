def static(assignment, history, level):
    Asmt1 = { 1:2, 2:3, 3:2, 4:2 }
    Asmt2 = { 1:5, 2:4 ,3:5, 4:4 }
    Asmt3 = { 1:6, 2:4, 3:5, 4:6 }

    defaultPath = None
    if assignment == "Asmt1":
        defaultPath = Asmt1
    elif assignment == "Asmt2":
        defaultPath = Asmt2
    else:
        defaultPath = Asmt3

    if history["Next{}".format(level)] < defaultPath[level]:
        return level - 1, history["Next{}".format(level)] - 1
    else:
        while level < 4:
            if history["Next{}".format(level+1)] < defaultPath[level+1]:
                return level, history["Next{}".format(level+1)] - 1
            level += 1
    return -1,0

def prequisiteSatisfied(assignment, level, number):
    mandatory = None
    if assignment == "Asmt1":
        mandatory = { 1:1, 2:2, 3:2, 4:2 }
    elif assignment == "Asmt2":
        mandatory = { 1:4, 2:2 ,3:1, 4:1 }
    else:
        mandatory = { 1:3, 2:3, 3:3, 4:4 }

    if number < mandatory[level]:
        return False
    return True
