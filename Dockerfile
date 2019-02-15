FROM golang:1.8-onbuild
RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip3 install --upgrade pip
RUN pip3 install graphviz
RUN pip3 install pydot
RUN pip3 install networkx
RUN go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
RUN go get github.com/jordic/lti/oauth
#RUN go get github.com/satori/go.uuid && go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
CMD go run *.go
MAINTAINER arieg419@gmail.com
