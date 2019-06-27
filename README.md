# _Email_ OMG Microservice

[![Open Microservice Guide](https://img.shields.io/badge/OMG%20Enabled-👍-green.svg?)](https://microservice.guide)
[![Build Status](https://travis-ci.com/heaptracetechnology/microservice-mail.svg?branch=master)](https://travis-ci.com/heaptracetechnology/microservice-mail)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-web.svg)](https://golangci.com)

An OMG service for Email, it allows to send and receive emails.

## Direct usage in [Storyscript](https://storyscript.io/):

##### Send Email
```coffee
>>> email send from:'abc@example.com' to:'xyz@example.com' subject:'subjectForMail' message:'messageBody'
{"success": true/false,"message": "success/failure message","statuscode": "SMTPStatusCode"}
```

Curious to [learn more](https://docs.storyscript.io/)?

✨🍰✨

## Usage with [OMG CLI](https://www.npmjs.com/package/omg)

##### Send Email
```shell
$ omg run send -a from=<SENDER_MAIL_ADDRESS> -a to=<LIST_OF_RECEIVER_EMAIL_ADDRESS> -a subject=<EMAIL_SUBJECT> -a message=<EMAIL_MESSAGE_BODY> -e SMTP_HOST=<smtp.example.com> -e SMTP_PORT="465" -e PASSWORD=<APP_PASSWORD_FROM_SENDER_ACCOUNT>
```
##### Receive Email
```shell
$ omg subscribe receive mail -a username=<RECEIVER_MAIL_ADDRESS> -a label=<MAILBOX_LABEL> -e PASSWORD=<APP_PASSWORD_FROM_RECEIVER_ACCOUNT> -e IMAP_HOST=<imap.example.com> -e IMAP_PORT="993"
```

**Note**: the OMG CLI requires [Docker](https://docs.docker.com/install/) to be installed.

## License
[MIT License](https://github.com/omg-services/email/blob/master/LICENSE).

