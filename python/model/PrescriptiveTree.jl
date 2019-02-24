using DataFrames
using MLDataUtils
using Suppressor
using Compat: @warn
@suppress begin
    @warn(using OptimalTrees)
end
using FileIO
using JLD2
using GraphViz
using Base.Test

function splitData(level, asmt)
  df = readtable("data/$asmt$level-train.csv", header=true, makefactors=true)
  fix1, p1 = false,1.2
  fix2, p2 = true,1.1
  fix3, p3 = true,0.95
  removescore4 = false
  if fix1
      for i=1:size(df)[1]
          if df[i,:Ygroup] == "A"
              df[i,:y1] = df[i,:y1]*p1
              df[i,:y] = df[i,:y]*p1
          end
      end
  end
  if fix2
      for i=1:size(df)[1]
          if df[i,:Ygroup] == "B"
              df[i,:y2] = df[i,:y2]*p2
              df[i,:y] = df[i,:y]*p2
          end
      end
  end
  if fix3
      for i=1:size(df)[1]
          if df[i,:Ygroup] == "C"
              df[i,:y3] = df[i,:y3]*p3
              df[i,:y] = df[i,:y]*p3
          end
      end
  end
  if level == 1
      X = df[3:end-2]
  else
      if removescore4
          X = df[3:end-4]
      else
          X = df[3:end-3]
      end
  end
  Y = df[1]
  T = df[2]
  if level == 1
      outcomes = df[:, [:y2, :y3]]
  else
      outcomes = df[:, [:y1, :y2, :y3]]
  end
  return stratifiedobs((X, outcomes, Y, T), p=0.75)
end

function trainTree(X, Y, T, level, depth, meu, asmt)
  lnr = OptimalTrees.OptimalTreePrescriptionMinimizer(
    ls_num_tree_restarts=100,
    treatment_minbucket=50,
    prescription_factor=meu
    )

  grid = OptimalTrees.GridSearch(
      lnr,
      Dict(:max_depth => depth) #, :prescription_factor => 0.:0.2:0.8,)
  )

  OptimalTrees.fit!(grid, X, T, Y)#,sample_weight=sample_weight) #, validation_criterion=:prediction_accuracy)
  #@show grid.best_score, grid.best_params
  lnr = grid.best_lnr
  #@show lnr

  plotname = "$level/$asmt/tree.dot"
  OptimalTrees.writedot(plotname, lnr)
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

depth=[[4,5],[4,5],[4,5],[4,5]]
meu=[0.55,0.55,0.55,0.55]
for asmt in ["rts"]#, "rts", "dfe"]
    for level in 4:4
        (train_X, train_outcomes, train_Y, train_T), (test_X, test_outcomes, test_Y, test_T) = splitData(level, asmt)
        lnr = trainTree(train_X, train_Y, train_T, level, depth[level], meu[level], asmt)
        treatment_accuracy, r2, accuracy= evaluate(lnr, test_X, test_outcomes, level)
        println(asmt)
        println(level)
        println(treatment_accuracy)
        println(accuracy)
    end
end
