
println(time())
using DataFrames
using Suppressor
using Compat: @warn
@suppress begin
    @warn(using OptimalTrees)
end

function func(assignment, level, history)
    map=Dict([("Climate Change", "cc"), ("Reading Test Scores", "rts"), ("Detecting Flu Epedemics", "dfe")]);
    asmt = map[assignment];
    #fix this path to be per assignment
    lnr = OptimalTrees.readjson("./python/model/$level/$asmt/tree.json");
    #the following function is temporary until i get history from omer
    #history = tempFunc(asmt, level)
    history = DataFrame(Gender = "None", level_of_education = "None", enrollment_mode = "audit", ageCategory = "Null", ad1 = 1, ad2 = 1, ad3 = 1, ad4 = 1, sd1 = 1, sd2 = 1, sd3 = 1, sd4 = 1, de1 = 1, de2 = 1, de3 = 0.8888888888888888, de4 = 1, score1 = 1, score2 = 0, score3 = 0, score4 = 0);
    treatment_prediction, outcome_predictions = OptimalTrees.predict(lnr, history);
    return treatment_prediction[1];
end

println(time())
func("Climate Change", 1, nothing)
println(time())
