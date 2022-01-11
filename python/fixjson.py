##making cluster json so that that values of the centroids are not standardized
#import json
#
#
# newdata = {}
# with open('oldclustering.json') as json_file:
#     data = json.load(json_file)
#     for qid,vals in data.items():
#         if qid == "10405":
#             newdata[qid] = newdata["10404"]
#         else:
#             newdata[qid] = {}
#             newdata[qid]["features"] = vals["features"]
#             newdata[qid]["decisions"] = vals["decisions"]
#             newdata[qid]["centroids"] = [[0.0 for feat in vals["features"]] for d in vals["decisions"]]
#             for i,f in enumerate(vals["features"]):
#                 meu = vals["meu"][f]
#                 sigma = vals["sigma"][f]
#                 for j,c in enumerate(vals["centroids"]):
#                     newdata[qid]["centroids"][j][i] = c[i] * sigma + meu
#                     if f == "t" and c[i] * sigma + meu > 1800:
#                         newdata[qid]["centroids"][j][i] = 1800
#     print(json.dumps(newdata))

####creating asner sheets for TAS
# import pandas as pd
#
# df = pd.DataFrame()
# with open('newLR.json') as json_file:
#      data = json.load(json_file)
#      for a,val in data.items():
#          i = 0
#          df = pd.DataFrame(index = range(30), columns = ["ID","Question","Options","Answer","Explanation"])
#          for l in [1,2,3,4]:
#              valLev = val[l-1]
#              questions = valLev["questions"]
#              for q in questions:
#                  df.at[i,:] = [q["qid"],q["text"],str(q["options"]),q["correctAnswer"],q["explanation"]]
#                  i = i + 1
#          df.to_csv("{}.csv".format(a),index = False)
