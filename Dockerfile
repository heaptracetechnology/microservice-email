FROM golang

RUN go get github.com/gorilla/mux

RUN go get github.com/emersion/go-imap

RUN go get github.com/emersion/go-message/mail

RUN go get github.com/emersion/go-imap/client

RUN go get github.com/cloudevents/sdk-go

WORKDIR /go/src/github.com/heaptracetechnology/microservice-mail

ADD . /go/src/github.com/heaptracetechnology/microservice-mail

RUN go install github.com/heaptracetechnology/microservice-mail

ENTRYPOINT microservice-mail

EXPOSE 3000