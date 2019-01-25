using DataFrames
using MLDataUtils
using OptimalTrees
using FileIO
using JLD2
using GraphViz
using Base.Test

function func(state, history)
    #upload lnr from file f depending on state - discuss with dimitris
    lnr = OptimalTrees.readjson(f)
    #get X of person
    treatment_prediction, outcome_predictions = OptimalTrees.predict(lnr, history)
    return treatment_prediction[1]
end
