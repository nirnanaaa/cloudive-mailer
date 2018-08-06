# Distributed Mailing

[![CircleCI](https://circleci.com/gh/nirnanaaa/cloudive-mailer.svg?style=svg)](https://circleci.com/gh/nirnanaaa/cloudive-mailer)

## Installation

Packages are available via docker hub:

```bash
# start up kafka

docker run \
    --rm -ti -p 9009:9009 \
    -e CLOUDIVE_KAFKA_BROKERS=kafka:9092 \
    cloudive/mailer master

docker run \
    --rm -ti -e CLOUDIVE_KAFKA_BROKERS=kafka:9092 \
    -e CLOUDIVE_HTTPD_ENABLED=true \
    -e CLOUDIVE_HTTPD_BIND_ADDRESS="0.0.0.0:9009" \
    -e CLOUDIVE_SMTP_ENABLED=true \
    -e CLOUDIVE_SMTP_HOSTNAME=smtp.office365.com \
    -e CLOUDIVE_SMTP_PORT=587 \
    -e CLOUDIVE_SMTP_USERNAME="someguy@somedomain.com" \
    -e CLOUDIVE_SMTP_PASSWORD="someguy" \
    -e CLOUDIVE_SMTP_FROM_NAME="someguy" \
    -e CLOUDIVE_SMTP_FROM_MAIL="someguy@somedomain.com" \
    cloudive/mailer worker
```


## Usage

```bash
curl -X POST \
  http://localhost:9009/mail \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
	"trace_id": "1234",
	"recipient": {
		"name": "Florian Kasper",
		"email": "someguy@somedomain.com",
		"tracking_id": "1234"
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
