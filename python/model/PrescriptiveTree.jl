using DataFrames
using MLDataUtils
using OptimalTrees
using FileIO
using JLD2
using GraphViz
using Base.Test

function splitData(level)
  df1 = readtable("data/cc$level.csv", header=true, makefactors=true)
  df2 = readtable("data/rts$level.csv", header=true, makefactors=true)
  df3 = readtable("data/dfe$level.csv", header=true, makefactors=true)
  df = [df1; df2; df3]
  X = df[4:end-3]
  Y = df[2]
  T = df[3]
  if level == 1
      outcomes = df[:, [:y2, :y3]]
  else
      outcomes = df[:, [:y1, :y2, :y3]]
  end
  # if level == 4
  #     X = delete!(X, :enrollment_mode)
  # end
  return stratifiedobs((X, outcomes, Y, T), p=0.75)
end

function trainTree(X, Y, T, level, depth, meu)
  lnr = OptimalTrees.OptimalTreePrescriptionMinimizer(
    ls_num_tree_restarts=100,
    treatment_minbucket=50,
    prescription_factor=meu
    #regression_sparsity=:all,
    #regression_weighted_betas=true,
    #regression_lambda=0.1,
    )

  grid = OptimalTrees.GridSearch(
      lnr,
      Dict(:max_depth => depth) #, :prescription_factor => 0.:0.2:0.8,)
  )

#   best_treatment = map(i -> minOutcome(outcomes[i, :], level), 1:size(X, 1))
#   treatments = ["A","B","C"]
#   d=Dict([(i,count(x->x==i,  best_treatment = map(i -> minOutcome(outcomes[i, :], level), 1:size(X, 1))
# )) for i in treatments])
#
#   if level == 1
#       sample_weight = Dict(
#           "B" => 99.0,#1/(d["B"]/length(T)),
#           "C" => 1.0#1/(d["C"]/length(T))
#         )
#   elseif level == 2
#       sample_weight = Dict(
#           "A" => 100.0,#1/(d["A"]/length(T)),
#           "B" => 80.0,#1/(d["B"]/length(T)),
#           "C" => 1.0#1/(d["C"]/length(T))
#         )
#   elseif level == 3
#       sample_weight = Dict(
#           "A" => 100.0,#1/(d["A"]/length(T)),
#           "B" => 20.0,#1/(d["B"]/length(T)),
#           "C" => 1.0#1/(d["C"]/length(T))
#         )
#   else
#       sample_weight = Dict(
#           "A" => 1.0,#1/(d["A"]/length(T)),
#           "B" => 4.0,#1/(d["B"]/length(T)),
#           "C" => 1.0#1/(d["C"]/length(T))
#         )
#   end

  OptimalTrees.fit!(grid, X, T, Y)#,sample_weight=sample_weight) #, validation_criterion=:prediction_accuracy)
  @show grid.best_score, grid.best_params
  lnr = grid.best_lnr
  @show lnr

  plotname = "$level/model2/tree.dot"
  OptimalTrees.writedot(plotname, lnr)
  @save "$level/model2/mytree.jld2" lnr
  OptimalTrees.writejson("$level/model2/tree.json", lnr)
  run(`dot -Tpng $plotname -o $(replace(plotname, ".dot", ".png"))`)
  return lnr
end

function minOutcome(outcome,level)
  r = rand(1)[1]
  if level == 1
      if outcome[1,1] <= outcome[1,2]
          if outcome[1,1] == outcome[1,2]
              if r > 0.5
                  return "C"
              end
          else
              return "B"
          end
      else
          return "C"
      end
      return "C"
  else
      if outcome[1,1] <= outcome[1,2] && outcome[1,1] <= outcome[1,3]
          if outcome[1,1] == outcome[1,2] && outcome[1,1] == outcome[1,3]
              if r < 1/3
                  return "C"
              end
              if r < 2/3
                  return "B"
              end
          elseif outcome[1,1] == outcome[1,2]
              if r < 0.5
                  return "B"
              end
          elseif outcome[1,1] == outcome[1,3]
              if r < 0.5
                  return "C"
              end
          else
              return "A"
          end
          return "A"
      elseif outcome[1,2] <= outcome[1,1] && outcome[1,2] <= outcome[1,3]
          if outcome[1,2] == outcome[1,1]
              if r < 0.5
                  return "A"
              end
          elseif outcome[1,2] == outcome[1,3]
              if r < 0.5
                  return "C"
              end
          else
              return "B"
          end
          return "B"
      else
        return "C"
      end
  end
end

function evaluate(lnr, X, outcomes, level)
  best_treatment = map(i -> minOutcome(outcomes[i, :], level), 1:size(X, 1))
  #check best treatment
  # treatments = ["A","B","C"]
  # d=Dict([(i,count(x->x==i,best_treatment)) for i in treatments])
  # println((d["A"]/length(best_treatment)))
  # println((d["B"]/length(best_treatment)))
  # println((d["C"]/length(best_treatment)))
  #####
  best_outcome = map(i -> minimum(vec(convert(Array, outcomes[i, :]))), 1:size(X, 1))
  treatment_predictions, outcome_predictions = OptimalTrees.predict(lnr, X)
  treatment_accuracy = mean(best_treatment .== treatment_predictions)
  # r2 = 1 - sum(abs2, best_outcome .- outcome_predictions) /
  #                        sum(abs2, best_outcome .- mean(best_outcome))
  accuracy = -1 * mean(outcome_predictions .- best_outcome)
  treatment_accuracy, r2, accuracy
end

depth=[[4,5],[4,5],[5,6],[5,6]]
meu=[0.55,0.55,0.55,0.55]
for level in 4:4
    (train_X, train_outcomes, train_Y, train_T), (test_X, test_outcomes, test_Y, test_T) = splitData(level)
    lnr = trainTree(train_X, train_Y, train_T, level, depth[level], meu[level])
    treatment_accuracy, r2, accuracy= evaluate(lnr, test_X, test_outcomes, level)
    println(treatment_accuracy)
    println(accuracy)
end

#0.878, -0.079
#0.585, -0.15
#0.702, -0.254
#0.4368, -0.81
