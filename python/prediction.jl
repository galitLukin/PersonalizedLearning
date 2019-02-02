using DataFrames
using Suppressor
using Compat: @warn
@suppress begin
    @warn(using OptimalTrees)
end

function func(assignment, level, history)
    map=Dict([("Climate Change", "cc"), ("Reading Test Scores", "rts"), ("Detecting Flu Epedemics", "dfe")])
    asmt = map[assignment]
    #fix this path to be per assignment
    lnr = OptimalTrees.readjson("./python/model/$level/$asmt/tree.json")
    #the following function is temporary until i get history from omer
    history = tempFunc(asmt, level)
    treatment_prediction, outcome_predictions = OptimalTrees.predict(lnr, history)
    return treatment_prediction[1]
end

#teporary function until history is given
function tempFunc(assignment, level)
    df = readtable("./python/model/data/$assignment$level.csv", header=true, makefactors=true)
    X = df[4:end-3]
    return X[1,:]
end
