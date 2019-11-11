FROM golang

RUN go get github.com/gorilla/mux

RUN go get github.com/emersion/go-imap

RUN go get github.com/emersion/go-message/mail

RUN go get github.com/emersion/go-imap/client

RUN go get github.com/cloudevents/sdk-go

WORKDIR /go/src/github.com/oms-services/email

ADD . /go/src/github.com/oms-services/email

RUN go install github.com/oms-services/email

ENTRYPOINT email

EXPOSE 3000