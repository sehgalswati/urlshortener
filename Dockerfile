FROM golang:latest

ADD . /go/src/github.com/sehgalswati/urlshortener/

WORKDIR /go/src/github.com/sehgalswati/urlshortener/

RUN go get && go build

RUN rm Dockerfile

RUN cp ./Docker/Dockerfile .

CMD tar cvzf - .
