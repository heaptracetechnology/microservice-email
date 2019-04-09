# Mail as a microservice
An OMG service for Mail, it allows to send and receive emails.

[![Open Microservice Guide](https://img.shields.io/badge/OMG-enabled-brightgreen.svg?style=for-the-badge)](https://microservice.guide)
<!-- [![Build Status](https://travis-ci.com/heaptracetechnology/microservice-firebase.svg?branch=master)](https://travis-ci.com/heaptracetechnology/microservice-firebase)
[![codecov](https://codecov.io/gh/heaptracetechnology/microservice-firebase/branch/master/graph/badge.svg)](https://codecov.io/gh/heaptracetechnology/microservice-firebase)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-web.svg)](https://golangci.com) -->

## [OMG](hhttps://microservice.guide) CLI

### OMG

* omg validate
```
omg validate
```
* omg build
```
omg build
```
### Test Service

* Test the service by following OMG commands

### CLI

##### Send Email
```sh
$ omg run send -a from=<SENDER_MAIL_ADDRESS> -a password=<MAIL_ACCOUNT_PASSWORD> -a to=<RECEIVER_EMAIL_ADDRESS> -a subject=<EMAIL_SUBJECT> -a message=<EMAIL_MESSAGE_BODY> -e SMTP_HOST="smtp.gmail.com" -e SMTP_PORT="587"
```

## License
### [MIT](https://choosealicense.com/licenses/mit/)

## Docker
### Build
```
docker build -t microservice-mail .
```
### RUN
```
docker run -p 3000:3000 microservice-mail
```
