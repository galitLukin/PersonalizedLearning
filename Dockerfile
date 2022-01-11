FROM golang:1.13
WORKDIR /PersonalizedLearning
COPY . .
RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip3 install --upgrade pip
RUN pip3 install graphviz
RUN pip3 install pydot
RUN pip3 install pydotplus
RUN pip3 install networkx==1.11
RUN go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
RUN go get github.com/satori/go.uuid && go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
CMD go run *.go
MAINTAINER lukingalit@gmail.com
