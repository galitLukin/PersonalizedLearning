FROM golang:1.8-onbuild
RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip3 install --upgrade pip
RUN pip3 install graphviz
RUN pip3 install pydot
RUN pip3 install networkx
RUN pip3 install Flask==0.10.1
RUN pip3 install Chalice==1.3.0
RUN pip3 install oauth2==1.9.0.post1
RUN pip3 install oauthlib==2.0.6
RUN pip3 install pyflakes==1.2.3
RUN pip3 install pytest==2.9.2
RUN pip3 install pytest-cache==1.0
RUN pip3 install pytest-cov==2.3.0
RUN pip3 install pytest-flakes==1.0.1
RUN pip3 install pytest-pep8==1.0.6
RUN pip3 install httplib2==0.9.2
RUN pip3 install six==1.11.0
RUN go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
#RUN go get github.com/jordic/lti/oauth
#RUN go get github.com/satori/go.uuid && go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
CMD go run *.go
MAINTAINER arieg419@gmail.com
