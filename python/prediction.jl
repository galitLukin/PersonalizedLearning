using DataFrames
using Suppressor
using Compat: @warn
@suppress begin
    @warn(using OptimalTrees)
end

function func(assignment, level, history)
    #fix this path to be per assignment
    lnr = OptimalTrees.readjson("./python/model/$level/tree.json")
    #the following function is temporary until i get history from omer
    history = tempFunc(assignment, level)
    treatment_prediction, outcome_predictions = OptimalTrees.predict(lnr, history)
    return treatment_prediction[1]
end

#teporary function until history is given
function tempFunc(assignment, level)
    if assignment == "Climate Change"
        df = readtable("./python/model/data/cc$level.csv", header=true, makefactors=true)
    end
    X = df[4:end-3]
    return X[1,:]
end
