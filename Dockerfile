FROM golang:alpine AS builder

COPY . /email
WORKDIR /email

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix nocgo -o /email .

FROM scratch
COPY --from=builder /email ./

CMD ["/email"]

EXPOSE 3000
