# Email as a microservice
An OMG service for Email, it allows to send and receive emails.

[![Open Microservice Guide](https://img.shields.io/badge/OMG-enabled-brightgreen.svg?style=for-the-badge)](https://microservice.guide)
[![Build Status](https://travis-ci.com/heaptracetechnology/microservice-mail.svg?branch=master)](https://travis-ci.com/heaptracetechnology/microservice-mail)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-web.svg)](https://golangci.com)


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
$ omg run send -a from=<SENDER_MAIL_ADDRESS> -a to=<LIST_OF_RECEIVER_EMAIL_ADDRESS> -a subject=<EMAIL_SUBJECT> -a message=<EMAIL_MESSAGE_BODY> -e SMTP_HOST=<smtp.example.com> -e SMTP_PORT="465" -e PASSWORD=<APP_PASSWORD_FROM_SENDER_ACCOUNT>


```
##### Receive Email
```sh
$ omg subscribe receive mail -a username=<RECEIVER_MAIL_ADDRESS> -a label=<MAILBOX_LABEL> -e PASSWORD=<APP_PASSWORD_FROM_RECEIVER_ACCOUNT> -e IMAP_HOST=<imap.example.com> -e IMAP_PORT="993"
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
