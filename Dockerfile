FROM golang:1.8-onbuild
RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip3 install --upgrade pip
RUN pip3 install git+https://github.com/mitodl/pylti.git
RUN pip3 install graphviz
RUN pip3 install pydot
RUN pip3 install networkx
RUN pip3 install Flask==0.12.3
RUN pip3 install Flask-WTF==0.10.2
RUN pip3 install Jinja2==2.7.3
RUN pip3 install MarkupSafe==0.23
RUN pip3 install Pygments==1.6
RUN pip3 install WTForms==2.0.1
RUN pip3 install Werkzeug==0.9.6
RUN pip3 install httplib2==0.9
RUN pip3 install itsdangerous==0.24
RUN pip3 install oauth==1.0.1
RUN pip3 install oauth2==1.9.0.post1
RUN pip3 install uWSGI==2.0.13
RUN pip3 install urwid==1.2.1
RUN pip3 install Sphinx==1.2.3
RUN pip3 install gunicorn>=19.5.0
RUN pip3 install pylti>=0.1.3
RUN go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
#RUN go get github.com/jordic/lti/oauth
#RUN go get github.com/satori/go.uuid && go get golang.org/x/crypto/bcrypt && go get github.com/go-sql-driver/mysql
CMD go run *.go
MAINTAINER arieg419@gmail.com
