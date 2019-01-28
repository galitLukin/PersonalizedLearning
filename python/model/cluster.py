import pandas as pd
import numpy as np
from sklearn.cluster import KMeans
from sklearn.metrics import silhouette_score
import matplotlib.pyplot as plt
from sklearn.decomposition import PCA

def Grade(row, start, length):
	s = 0.0
	d = 0
	for i in range(start, start + length,2):
		s = s + row[i+1]/row[i] if row[i] > 0 else 0
		d = d + 1
	return s/d

def clusterStudents(ques, assignmentNames):

	df = ques.drop(['username', 'location', 'year_of_birth', 'gender', 'level_of_education',
       'enrollment_mode', 'ageCategory'], axis=1)

	for asmt in assignmentNames:
		asmtCols = [col for col in ques.columns if asmt in col]
		asmtDF = ques[asmtCols]
		for i in range(1,5):
			cols = [c for c in asmtDF.columns if '{}.'.format(i) in c]
			asmtLev = asmtDF[cols]
			df.loc[:,'score{}_{}'.format(i,asmt)] = asmtLev.apply(lambda row: Grade(row, 0, len(cols)),axis=1)
		if asmt == "cc":
			continue
		else:
			df.loc[:,'score{}'.format(asmt)] = (df.loc[:,'score1_{}'.format(asmt)] + df.loc[:,'score2_{}'.format(asmt)] + df.loc[:,'score3_{}'.format(asmt)] + df.loc[:,'score4_{}'.format(asmt)])/4.0
			df = df.drop(['score1_{}'.format(asmt), 'score2_{}'.format(asmt), 'score3_{}'.format(asmt), 'score4_{}'.format(asmt)], axis=1)

	cols = [c for c in df.columns if 'score' in c]
	df=df[cols]

	kmeans = KMeans(n_clusters=2)
	kmeans.fit(df)
	y_kmeans = kmeans.predict(df)
	ques['cluster'] = y_kmeans
	pca = PCA(n_components=2)
	principalComponents = pca.fit_transform(df)
	principalDf = pd.DataFrame(data = principalComponents
             , columns = ['principal component 1', 'principal component 2'])
	# plt.scatter(principalDf.iloc[:, 0], principalDf.iloc[:, 1], c=y_kmeans, s=50, cmap='viridis')
	# plt.show()
	# print(silhouette_score(df, y_kmeans))
	return

def setYGroup(data,groups):
	data["Ygroup"] = np.nan
	clusters = data.cluster.unique()
	for c in clusters:
		data.loc[data['cluster'] == c,'Ygroup'] = np.random.randint(1, groups+1, data[data['cluster'] == c].shape[0])
	if groups == 2:
		data.Ygroup.replace(2, 3, inplace=True)
		data.Ygroup.replace(1, 2, inplace=True)
	return data
