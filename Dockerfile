FROM golang:alpine AS builder

COPY . /email
WORKDIR /email/cmd/server

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix nocgo -o /server .

FROM scratch
COPY --from=builder /server ./

CMD ["/server"]

EXPOSE 3000
