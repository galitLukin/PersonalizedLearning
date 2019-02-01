FROM golang:1.8-onbuild
RUN go get github.com/satori/go.uuid && go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
CMD go run *.go
MAINTAINER arieg419@gmail.com