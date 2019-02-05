FROM golang:1.8-onbuild
RUN apt-get update && apt-get install -y python3 python3-pip
# RUN apt-get update && apt-get install -y julia=0.6.4
RUN pip3 install --upgrade pip
RUN pip3 install graphviz
RUN pip3 install pydot
RUN git clone https://github.com/networkx/networkx.git
# RUN pip3 install --no-use-pep517 pandas
# RUN pip3 install --no-use-pep517 julia
RUN go get github.com/satori/go.uuid && go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
CMD go run *.go
MAINTAINER arieg419@gmail.com
