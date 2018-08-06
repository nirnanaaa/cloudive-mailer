# HTTP to SMTP Gateway

[![CircleCI](https://circleci.com/gh/nirnanaaa/cloudive-mailer.svg?style=svg)](https://circleci.com/gh/nirnanaaa/cloudive-mailer)

Cloudive mailer is offered in two run modes: a distributed, gateway-worker mode and a more minimalistic standalone mode. In Standalone mode all mail processing happens
within the same process and may overload your server.

## Distributed mode

Packages are available via docker hub:

```bash
# start up kafka

docker run \
    --rm -ti -p 9009:9009 \
    -e CLOUDIVE_HTTPD_ENABLED=true \ # prometheus metrics and /mail endpoint
    -e CLOUDIVE_HTTPD_BIND_ADDRESS=0.0.0.0:9092 \
    -e CLOUDIVE_KAFKA_BROKERS=kafka:9092 \
    -e CLOUDIVE_KAFKA_INBOUND_QUEUE=mail \ # where mails get accepted from within the internal network
    -e CLOUDIVE_KAFKA_OUTBOUND_QUEUE=mail-worker-queue \ # jobs are forwarded to an internal queue
    cloudive/mailer master

docker run \
    --rm -ti \
    -e CLOUDIVE_KAFKA_BROKERS=kafka:9092 \
    -e CLOUDIVE_KAFKA_INBOUND_QUEUE=mail-worker-queue \ # where mails get accepted from within the internal network
    -e CLOUDIVE_SMTP_ENABLED=true \
    -e CLOUDIVE_SMTP_HOSTNAME=smtp.office365.com \
    -e CLOUDIVE_SMTP_PORT=587 \
    -e CLOUDIVE_SMTP_USERNAME="someguy@somedomain.com" \
    -e CLOUDIVE_SMTP_PASSWORD="someguy" \
    -e CLOUDIVE_SMTP_FROM_NAME="someguy" \ # Default settings
    -e CLOUDIVE_SMTP_FROM_MAIL="someguy@somedomain.com" \ # Default settings
    cloudive/mailer worker
```


### Usage

```bash
curl -X POST \
  http://localhost:9009/mail \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
	"recipient": {
		"name": "Florian Kasper",
		"email": "someguy@somedomain.com"
	},
	"sender": {
		"name": "Florian Kasper",
		"email": "someguy@somedomain.com"
	},
	"subject": "test",
	"payload": "YWxwaGEgYmV0dGEgZ2FtbWE=",
	"attachments": [
		{
			"name": "favicon",
			"url": "https://google.com/favicon.ico"
		}
	]
}'
```
